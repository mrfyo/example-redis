package main

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Poster string `json:"poster"`
	Time   uint64 `json:"time"`
	Votes  int    `json:"votes"`
}


func CreateArticle(atricle *Article) {
	
}
