package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mrfyo/example-redis/result"
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
	_, err = redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := KeyGenerate("vote", article.KeyName())

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

//
// DAO: Article
//

func CreateArticle(article *Article) (err error) {
	article.ID, err = NextID(article.TableName())
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

func GetArticleById(ID int) (art *Article, err error) {

	key := KeyGenerate(ArticleName, strconv.Itoa(ID))
	cmd := redisDB.HGetAll(ctx, key)

	if err = cmd.Err(); err != nil {
		log.Println(err)
		return
	}

	m := cmd.Val()

	art = &Article{
		ID:      ID,
		Title:   m["title"],
		Link:    m["link"],
		Poster:  m["poster"],
		Content: m["content"],
	}
	art.Time, _ = strconv.ParseInt(m["time"], 10, 64)
	art.Votes, _ = strconv.Atoi(m["votes"])
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

	ids, err := BatchExtraID(keys)
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

	ids, err := BatchExtraID(keys)
	if err != nil {
		return
	}

	return GetAllArticleIDs(ids)

}

func GetAllArticleIDs(ids []int) (articles []*Article) {
	for _, id := range ids {
		article, err := GetArticleById(id)
		if err != nil {
			continue
		}
		articles = append(articles, article)
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

//
// API: Article
//

func AddArticleHandler(c *gin.Context) {
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		result.Fail(c, 1, "Form Error")
		return
	}

	if AnyEmptyStr(article.Title, article.Content, article.Poster) {
		result.Fail(c, 2, "Param Not Empty")
		return
	}

	if err := CreateArticle(&article); err != nil {
		result.Fail(c, 10, "Create Fail.")
		return
	}

	result.Success(c, nil)
}

func DeleteArticleHandler(c *gin.Context) {

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ID <= 0 {
		result.Fail(c, 2, "Path Param Error: ID")
		return
	}

	article, err := GetArticleById(ID)
	if err != nil {
		result.Fail(c, 10, "Article Not Exist")
		return
	}

	if err := RemoveArticle(article); err != nil {
		result.Fail(c, 20, "Delete Fail")
	}
}

func TopListArticleHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		result.Fail(c, 2, "Query Param Error: limit")
		return
	}

	articles := GetAllArticleByScore(limit)

	list := make([]map[string]interface{}, 0, len(articles))

	for _, article := range articles {
		list = append(list, article.ToMap())
	}

	result.Success(c, gin.H{
		"total": len(articles),
		"items": list,
	})
}

func ListArticleHandler(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		result.Fail(c, 2, "Query Param Error: offset")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit <= 0 {
		result.Fail(c, 2, "Query Param Error: limit")
		return
	}

	articles := GetAllArticleByPage(offset, limit)

	total := CountArticle()

	result.Success(c, gin.H{
		"total": total,
		"items": articles,
	})
}

func VoteArticleHandler(c *gin.Context) {
	form := struct {
		UserId    int `json:"userId"`
		ArticleId int `json:"articleId"`
	}{}

	if err := c.ShouldBindJSON(&form); err != nil {
		result.Fail(c, 1, "Form Error")
		return
	}

	if form.UserId <= 0 {
		result.Fail(c, 2, "Param Error: userId")
		return
	}

	if form.ArticleId <= 0 {
		result.Fail(c, 2, "Param Error: articleId")
		return
	}

	user, err := GetUser(form.UserId)
	if err != nil {
		result.Fail(c, 10, "Not Exist: User")
		return
	}

	article, err := GetArticleById(form.ArticleId)
	if err != nil {
		result.Fail(c, 10, "Not Exist: article")
		return
	}

	err = article.VoteBy(user)
	if err != nil {
		result.Fail(c, 30, "voted fail")
		return
	}

	result.Success(c, nil)
}
