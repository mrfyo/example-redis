package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrfyo/example-redis/model"
	"github.com/mrfyo/example-redis/result"
	"github.com/mrfyo/example-redis/util"
)

func addArticleHandler(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		result.Fail(c, 1, "Form Error")
		return
	}

	if util.AnyEmptyStr(article.Title, article.Content, article.Poster) {
		result.Fail(c, 2, "Param Not Empty")
		return
	}

	if err := model.CreateArticle(&article); err != nil {
		result.Fail(c, 10, "Create Fail.")
		return
	}

	result.Success(c, nil)
}

func deleteArticleHandler(c *gin.Context) {

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ID <= 0 {
		result.Fail(c, 2, "Path Param Error: ID")
		return
	}

	article, err := model.GetArticleById(ID)
	if err != nil {
		result.Fail(c, 10, "Article Not Exist")
		return
	}

	if err := model.RemoveArticle(article); err != nil {
		result.Fail(c, 20, "Delete Fail")
	}

	result.Success(c, nil)
}

func topListArticleHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		result.Fail(c, 2, "Query Param Error: limit")
		return
	}

	articles := model.GetAllArticleByScore(limit)

	list := make([]map[string]interface{}, 0, len(articles))

	for _, article := range articles {
		list = append(list, article.ToMap())
	}

	result.Success(c, gin.H{
		"total": len(articles),
		"items": list,
	})
}

func listArticleHandler(c *gin.Context) {
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

	articles := model.GetAllArticleByPage(offset, limit)

	total := model.CountArticle()

	result.Success(c, gin.H{
		"total": total,
		"items": articles,
	})
}

func publishArticleHandler(c *gin.Context) {

	form := struct {
		UserId    int `json:"userId"`
		ArticleId int `json:"articleId"`
	}{}

	if err := c.ShouldBindJSON(&form); err != nil {
		result.Fail(c, 1, "Form Error")
		return
	}

	article, err := model.GetArticleById(form.ArticleId)
	if err != nil {
		result.Fail(c, 10, "Article Not Exist")
		return
	}

	if err := article.Published(); err != nil {
		result.Fail(c, 30, err.Error())
		return
	}

	result.Success(c, nil)

}

func voteArticleHandler(c *gin.Context) {
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

	user, err := model.GetUser(form.UserId)
	if err != nil {
		result.Fail(c, 10, "Not Exist: User")
		return
	}

	article, err := model.GetArticleById(form.ArticleId)
	if err != nil {
		result.Fail(c, 10, "Not Exist: article")
		return
	}

	err = article.VoteBy(user)
	if err != nil {
		result.Fail(c, 30, err.Error())
		return
	}

	result.Success(c, nil)
}
