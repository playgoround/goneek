package types

import "time"

type Article struct {
	title    string
	date     time.Time
	contents map[string]interface{}
}

func NewArticle(
	title string,
	date time.Time,
	contents map[string]interface{},
) Article {
	return Article{title, date, contents}
}
func (article Article) GetTitle() string                    { return article.title }
func (article Article) GetContents() map[string]interface{} { return article.contents }
func (article Article) GetDate() time.Time                  { return article.date }

type ArticleInfo struct {
	title string
	url   string
	html  []byte
}

func NewArticleInfo(
	title string,
	url string,
	html []byte,
) ArticleInfo {
	return ArticleInfo{title, url, html}
}
func (article ArticleInfo) GetTitle() string { return article.title }
func (article ArticleInfo) GetUrl() string   { return article.url }
func (article ArticleInfo) GetHtml() []byte  { return article.html }
