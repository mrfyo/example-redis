package main

import (
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	ArticleName           = "article"
	ArticlePublishZsetKey = "publish:article"
	ArticleScoreZsetKey   = "score:article"
	voteBase              = 432
)

type Article struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	Poster string `json:"poster"`
	Time   int64  `json:"time"`
	Votes  int    `json:"votes"`
}

func (Article) TableName() string {
	return "article"
}

func (art *Article) KeyName() string {
	return KeyGenerate(art.TableName(), strconv.Itoa(art.ID))
}

func (art *Article) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     art.ID,
		"title":  art.Title,
		"link":   art.Link,
		"poster": art.Poster,
		"time":   art.Time,
		"votes":  art.Votes,
	}
}

func CreateArticle(article *Article) (err error) {
	article.ID, err = NextID(article.TableName())
	if err != nil {
		return
	}
	key := article.KeyName()

	intCmd := redisDB.HSet(ctx, key, article.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}

	return
}

func UpdateArticle(article *Article) (err error) {
	key := article.KeyName()
	intCmd := redisDB.HSet(ctx, key, article.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func RemoveArticle(art *Article) (err error) {
	key := art.KeyName()

	intCmd := redisDB.Del(ctx, key)
	if err = intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func GetArticle(ID int) (art *Article, err error) {

	key := KeyGenerate(ArticleName, strconv.Itoa(ID))
	cmd := redisDB.HGetAll(ctx, key)

	if err = cmd.Err(); err != nil {
		log.Println(err)
		return
	}

	m := cmd.Val()

	art = &Article{
		ID:     ID,
		Title:  m["title"],
		Link:   m["link"],
		Poster: m["poster"],
	}
	art.Time, _ = strconv.ParseInt(m["time"], 10, 64)
	art.Votes, _ = strconv.Atoi(m["votes"])
	return
}

func PublishArticle(article *Article) (err error) {

	_, err = redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := article.KeyName()
		p.ZAdd(ctx, ArticlePublishZsetKey, &redis.Z{
			Score:  float64(article.Time),
			Member: key,
		})

		p.ZAdd(ctx, ArticleScoreZsetKey, &redis.Z{
			Score:  getArticleScore(article),
			Member: key,
		})

		return nil
	})

	return
}

// VoteArticle 用户给文章投票
func VoteArticle(user *User, article *Article) (err error) {
	_, err = redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := KeyGenerate("vote", article.KeyName())

		p.SAdd(ctx, key, user.KeyName())

		article.Votes++
		key = article.KeyName()
		p.HSet(ctx, key, "votes", article.Votes)

		p.ZAdd(ctx, ArticleScoreZsetKey, &redis.Z{
			Score:  getArticleScore(article),
			Member: key,
		})

		return nil
	})

	return
}

func getArticleScore(article *Article) float64 {
	score := article.Time/1000 + int64(voteBase*article.Votes)
	return float64(score)
}

// TopScoreArticle 获取分数最高的 n 篇文章
func TopScoreArticle(num int) (articles []*Article) {

	cmd := redisDB.ZRevRangeByScore(ctx, ArticleScoreZsetKey, &redis.ZRangeBy{
		Count: int64(num),
	})

	if err := cmd.Err(); err != nil {
		return
	}

	members := cmd.Val()

	for _, member := range members {
		if ID, err := ExtraID(member); err != nil {
			continue
		} else {
			article, err := GetArticle(ID)
			if err != nil {
				articles = append(articles, article)
			}
		}
	}
	return
}
