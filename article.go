package main

import (
	"log"
	"strconv"
)

const ArticleName = "article"

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