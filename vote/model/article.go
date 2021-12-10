package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mrfyo/example-redis/util"
)

const (
	ArticleName           = "article"
	ArticleRecordKey      = "record:article"
	ArticlePublishZsetKey = "publish:article"
	ArticleScoreZsetKey   = "score:article"
	voteBase              = 432
)

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Poster  string `json:"poster"`
	Time    int64  `json:"time"`
	Votes   int    `json:"votes"`
	Content string `json:"content"`
}

func (Article) TableName() string {
	return "article"
}

func (art *Article) KeyName() string {
	return util.KeyGenerate(art.TableName(), strconv.Itoa(art.ID))
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

func (article *Article) Published() (err error) {
	_, err = redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := article.KeyName()
		p.ZAdd(ctx, ArticlePublishZsetKey, &redis.Z{
			Score:  float64(article.Time),
			Member: key,
		})

		p.ZAdd(ctx, ArticleScoreZsetKey, &redis.Z{
			Score:  article.GetScore(),
			Member: key,
		})

		return nil
	})

	return
}

// VoteBy 用户给文章投票
func (article *Article) VoteBy(user *User) (err error) {
	key := util.KeyGenerate("vote", article.KeyName())

	ok, err := redisDB.SIsMember(ctx, key, user.KeyName()).Result()
	if ok || err != nil {
		return fmt.Errorf("has been voted by %s", user.Username)
	}

	redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.SAdd(ctx, key, user.KeyName())

		article.Votes++
		key = article.KeyName()
		p.HIncrBy(ctx, key, "votes", 1)

		p.ZIncrBy(ctx, ArticleScoreZsetKey, voteBase, key)

		return nil
	})

	return
}

// GetScore 计算文章评分
func (article *Article) GetScore() float64 {
	score := article.Time/1000 + int64(voteBase*article.Votes)
	return float64(score)
}

func ArticleOfMap(m map[string]string) *Article {
	article := Article{
		Title:   m["title"],
		Link:    m["link"],
		Poster:  m["poster"],
		Content: m["content"],
	}
	article.ID, _ = strconv.Atoi(m["id"])
	article.Time, _ = strconv.ParseInt(m["time"], 10, 64)
	article.Votes, _ = strconv.Atoi(m["votes"])

	return &article
}

//
// DAO: Article
//

func CreateArticle(article *Article) (err error) {
	article.ID, err = nextID(article.TableName())
	if err != nil {
		return
	}
	article.Link = fmt.Sprintf("http://127.0.0.1:8080/api/articles/%d", article.ID)
	article.Time = time.Now().Unix()

	key := article.KeyName()

	if err := redisDB.HSet(ctx, key, article.ToMap()).Err(); err != nil {
		log.Println(err)
	}

	_, err = redisDB.ZAdd(ctx, ArticleRecordKey, &redis.Z{
		Score:  float64(article.Time),
		Member: key,
	}).Result()

	return
}

func RemoveArticle(art *Article) (err error) {
	redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := art.KeyName()
		p.Del(ctx, key)
		p.ZRem(ctx, ArticleRecordKey, key)
		return nil
	})

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

func GetArticleById(ID int) (art *Article, err error) {

	key := util.KeyGenerate(ArticleName, strconv.Itoa(ID))
	cmd := redisDB.HGetAll(ctx, key)

	if err = cmd.Err(); err != nil {
		log.Println(err)
		return
	}

	art = ArticleOfMap(cmd.Val())

	return
}

// GetAllArticleByScore 获取分数最高的 n 篇文章
func GetAllArticleByScore(num int) (articles []*Article) {

	keys, err := redisDB.ZRevRangeByScore(ctx, ArticleScoreZsetKey, &redis.ZRangeBy{
		Count: int64(num),
	}).Result()

	if err != nil {
		return
	}

	ids, err := util.BatchExtraID(keys)
	if err != nil {
		return
	}

	return GetAllArticleIDs(ids)
}

func GetAllArticleByPage(offset, limit int) (artilces []*Article) {

	keys, err := redisDB.ZRevRangeByScore(ctx, ArticleRecordKey, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(offset),
		Count:  int64(limit),
	}).Result()

	if err != nil {
		log.Println(err)
		return
	}

	ids, err := util.BatchExtraID(keys)
	if err != nil {
		return
	}

	return GetAllArticleIDs(ids)

}

func GetAllArticleIDs(ids []int) (articles []*Article) {

	var cmds []*redis.StringStringMapCmd
	redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, id := range ids {
			key := util.KeyGenerate(ArticleName, strconv.Itoa(id))
			cmds = append(cmds, redisDB.HGetAll(ctx, key))
		}
		return nil
	})

	for _, cmd := range cmds {
		if cmd.Err() != nil {
			continue
		}
		articles = append(articles, ArticleOfMap(cmd.Val()))
	}

	return
}

func CountArticle() (count int64) {
	i, err := redisDB.ZCount(ctx, ArticleRecordKey, "-inf", "+inf").Result()
	if err != nil {
		return 0
	}
	return i
}
