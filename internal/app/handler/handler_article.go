package handler

import (
	"net/http"
	"strings"

	"github.com/alekslesik/neuro-news/internal/pkg/template"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

// ArticleHandler handle requests related with articles
type ArticleHandler struct {
	AppHandler *AppHandler
	l          *logger.Logger
}

// NewArticleHandler create new instance of ArticleHandler.
func NewArticleHandler(appHandler *AppHandler, l *logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		AppHandler: appHandler,
		l:          l,
	}
}

// Home GET handler for home page with pagination: GET /?PAGEN_1
func (a *ArticleHandler) Home(w http.ResponseWriter, r *http.Request) {
	const (
		op       = "GetHomeArticles()"
		tmplFile = "home.page.html"
	)

	var (
		td  *template.TemplateData
		err error
	)

	page := r.URL.Query().Get("PAGEN_1")

	td, err = a.AppHandler.articleService.GetHomePaginateData(page)
	if err != nil {
		a.l.Error().Msgf("%s: GetHomePaginateData error > %s", op, err)
	}

	td, err = a.AppHandler.imageService.AddRandomImages(td)
	if err != nil {
		a.l.Error().Msgf("%s: AddRandomImages error > %s", op, err)
	}

	err = a.AppHandler.articleService.RenderTemplate(w, r, tmplFile, td)
	if err != nil {
		a.l.Error().Msgf("%s: RenderTemplate home page error > %s", op, err)
	}
}

// Category GET handler for category page with pagination: GET /{category}?PAGEN_1
func (a *ArticleHandler) Category(w http.ResponseWriter, r *http.Request) {
	const (
		op       = "Category()"
		tmplFile = "category.page.html"
	)

	var (
		td  *template.TemplateData
		err error
	)

	// get article URL
	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) < 3 {
		http.NotFound(w, r)
		return
	}

	url := strings.Trim(strings.Trim(urlParts[len(urlParts)-2], "\""), " ")
	page := r.URL.Query().Get("PAGEN_1")

	td, err = a.AppHandler.articleService.GetCategoryArticlesData(url, page)
	if err != nil {
		a.l.Error().Msgf("%s: GetArticle error > %s", op, err)
	}

	td.TemplateDataPage.Category = url

	td, err = a.AppHandler.imageService.AddRandomImages(td)
	if err != nil {
		a.l.Error().Msgf("%s: AddRandomImages error > %s", op, err)
	}

	err = a.AppHandler.articleService.RenderTemplate(w, r, tmplFile, td)
	if err != nil {
		a.l.Error().Msgf("%s: RenderTemplate article page error > %s", op, err)
	}
}

// Article GET handler for article page: GET /{category}/{article}
func (a *ArticleHandler) Article(w http.ResponseWriter, r *http.Request) {
	const (
		op       = "GetArticle()"
		tmplFile = "article.page.html"
	)

	var (
		td  *template.TemplateData
		err error
	)

	// get article URL
	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) < 3 {
		http.NotFound(w, r)
		return
	}

	url := strings.Trim(strings.Trim(urlParts[len(urlParts)-2], "\""), " ")

	td, err = a.AppHandler.articleService.GetArticleTemplateData(url)
	if err != nil {
		a.l.Error().Msgf("%s: GetArticle error > %s", op, err)
	}

	td, err = a.AppHandler.imageService.AddRandomImages(td)
	if err != nil {
		a.l.Error().Msgf("%s: AddRandomImages error > %s", op, err)
	}

	err = a.AppHandler.articleService.RenderTemplate(w, r, tmplFile, td)
	if err != nil {
		a.l.Error().Msgf("%s: RenderTemplate article page error > %s", op, err)
	}
}
