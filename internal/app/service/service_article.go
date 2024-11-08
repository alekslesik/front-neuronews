package service

import (
	"math"
	"net/http"
	"strconv"

	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/internal/pkg/grabber"
	"github.com/alekslesik/neuro-news/internal/pkg/template"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type ArticleService interface {
	GetHomePaginateData(string) (*template.TemplateData, error)

	GrabNewArticle() (*model.Article, error)
	GetArticleByURL(url string) (model.Article, error)
	InsertArticleImage(i *model.Image, a *model.Article) error
	GetArticleTemplateData(url string) (*template.TemplateData, error)
	GetCategoryArticlesData(page, url string) (*template.TemplateData, error)

	RenderTemplate(w http.ResponseWriter, r *http.Request, name string, td *template.TemplateData) error
}

type articleService struct {
	ar model.ArticleModel
	t  *template.Template
	l  *logger.Logger
	g  *grabber.Grabber
}

// GetHomePaginateData return template data for home page
func (as *articleService) GetHomePaginateData(page string) (*template.TemplateData, error) {
	const op = "service.GetHomeTemplateData()"
	var err error

	as.t.TemplateData.TemplateDataArticle.CarouselArticles, err = as.getHomeCarouselArticles()
	if err != nil {
		as.l.Error().Msgf("%s: get home carouser articles data error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.TrendingArticlesTop, err = as.getHomeTrendingArticlesTop()
	if err != nil {
		as.l.Error().Msgf("%s: get home trending articles data top error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.TrendingArticlesBottom, err = as.getHomeTrendingArticlesBottom()
	if err != nil {
		as.l.Error().Msgf("%s: get home trending articles data bottom error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.NewsArticles, err = as.getHomeNewsArticles()
	if err != nil {
		as.l.Error().Msgf("%s: get home news articles data error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.SportArticles, err = as.getHomeSportArticles()
	if err != nil {
		as.l.Error().Msgf("%s: get home sport news articles data error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.VideoArticles, err = as.getHomeVideoArticles()
	if err != nil {
		as.l.Error().Msgf("%s: get home video articles data error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.PaginationArticles, err = as.getHomePaginationArticles(page)
	if err != nil {
		as.l.Error().Msgf("%s: get pagination articles data on home page error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataPage, err = as.GetPaginationTemplateData(page)
	if err != nil {
		as.l.Error().Msgf("%s: get pagination page data on home page error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.RandomArticles, err = as.getRandomArticles(5)
	if err != nil {
		as.l.Error().Msgf("%s: get random page data on home page error > %s", op, err)
		return nil, err
	}

	return &as.t.TemplateData, nil
}

// GetCategoryArticlesData return template data for category page
func (as *articleService) GetCategoryArticlesData(url, page string) (*template.TemplateData, error) {
	const op = "service.GetCategoryListData()"
	var err error

	as.t.TemplateData.TemplateDataArticle.PaginationArticles, err = as.getCategoryPaginationArticles(url, page)
	if err != nil {
		as.l.Error().Msgf("%s: get pagination articles data on category page error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataPage, err = as.GetPaginationTemplateData(page)
	if err != nil {
		as.l.Error().Msgf("%s: get pagination page data on category page error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.RandomArticles, err = as.getRandomArticles(5)
	if err != nil {
		as.l.Error().Msgf("%s: get random page data on category page error > %s", op, err)
		return nil, err
	}

	return &as.t.TemplateData, nil
}

func (as *articleService) getHomeCarouselArticles() ([]model.Article, error) {
	return as.ar.SelectHomeCarouselArticles()
}

func (as *articleService) getHomeTrendingArticlesTop() ([]model.Article, error) {
	return as.ar.SelectHomeTrendingArticlesTop()
}

func (as *articleService) getHomeTrendingArticlesBottom() ([]model.Article, error) {
	return as.ar.SelectHomeTrendingArticlesBottom()
}

func (as *articleService) getHomeNewsArticles() ([]model.Article, error) {
	return as.ar.SelectHomeNewsArticles()
}

func (as *articleService) getHomeSportArticles() ([]model.Article, error) {
	return as.ar.SelectHomeSportArticles()
}

func (as *articleService) getHomeVideoArticles() ([]model.Article, error) {
	return as.ar.SelectHomeVideoArticles()
}

// getHomePaginationArticles return []model.Article for pagination part
func (as *articleService) getHomePaginationArticles(page string) ([]model.Article, error) {
	op := "service.GetHomePaginationArticles"
	limit := 15
	var offset int

	if page == "" {
		offset = 0
	} else {
		p, err := strconv.Atoi(page)
		if err != nil {
			as.l.Error().Msgf("%s: atoi convert pagination page error > %s", op, err)
			return nil, err
		}
		offset = (p - 1) * limit
	}

	return as.ar.SelectHomePaginationArticles(limit, offset)
}

// getRandomArticles return 5 random articles
func (as *articleService) getRandomArticles(limit int) ([]model.Article, error) {
	return as.ar.GetRandomArticles(limit)
}

// getCategoryPaginationArticles return []model.Article with pagination on category page
func (as *articleService) getCategoryPaginationArticles(url, page string) ([]model.Article, error) {
	op := "service.GetHomePaginationArticles"
	limit := 15
	var offset int

	if page == "" {
		offset = 0
	} else {
		p, err := strconv.Atoi(page)
		if err != nil {
			as.l.Error().Msgf("%s: atoi convert pagination page error > %s", op, err)
			return nil, err
		}
		offset = (p - 1) * limit
	}

	return as.ar.SelectCategoryArticles(url, limit, offset)
}

// GetPaginationTemplateData return pagination page data
func (as *articleService) GetPaginationTemplateData(page string) (*template.TemplateDataPage, error) {
	const op = "service.GetPaginationTemplateData()"

	articlesOnPage := 15
	articlesCount, err := as.ar.CountArticles()
	if err != nil {
		as.l.Error().Msgf("%s: render template error > %s", op, err)
		return nil, err
	}

	totalPages := math.Ceil(float64(articlesCount) / float64(articlesOnPage))

	var currentPage int

	if page == "" {
		currentPage = 1
	} else {
		currentPage, err = strconv.Atoi(page)
		if err != nil {
			as.l.Error().Msgf("%s: convert page string pagination number to int > %s", op, err)
			return nil, err
		}
		if currentPage > int(totalPages) {
			as.l.Error().Msgf("%s: pagination page out of range > %s", op, err)
			return nil, err
		}
	}

	data := &template.TemplateDataPage{
		TotalPaginationPages:  int(totalPages),
		CurrentPaginationPage: currentPage,
	}

	return data, nil
}

// GetArticleTemplateData return template data for article page by article URL
func (as *articleService) GetArticleTemplateData(url string) (*template.TemplateData, error) {
	const op = "service.GetArticleByURL()"
	var err error

	as.t.TemplateData.TemplateDataArticle.Article, err = as.GetArticleByURL(url)
	if err != nil {
		as.l.Error().Msgf("%s: get article template data error > %s", op, err)
		return nil, err
	}

	as.t.TemplateData.TemplateDataArticle.RandomArticles, err = as.getRandomArticles(10)
	if err != nil {
		as.l.Error().Msgf("%s: get random page data on article page error > %s", op, err)
		return nil, err
	}

	return &as.t.TemplateData, nil
}

func (as *articleService) GetArticleByURL(url string) (model.Article, error) {
	return as.ar.SelectArticleByURL(url)
}

// GrabNewArticle grab new article from news site
func (as *articleService) GrabNewArticle() (*model.Article, error) {
	const op = "service.GetNewArticle()"

	// get article model without image
	a, err := as.g.GrabArticle()
	if err != nil {
		as.l.Error().Msgf("%s: grab article error > %s", op, err)
		return nil, err
	}

	return a, err
}

// InsertArticleImage insert article to DB
func (as *articleService) InsertArticleImage(image *model.Image, article *model.Article) error {
	const op = "service.InsertArticleImage()"

	err := as.ar.InsertArticleImage(image, article)
	if err != nil {
		as.l.Error().Msgf("%s: insert article error > %s", op, err)
		return err
	}
	return nil
}

// RenderTemplate render page with received data
func (as *articleService) RenderTemplate(w http.ResponseWriter, r *http.Request, n string, td *template.TemplateData) error {
	const op = "service.RenderTemplate()"

	err := as.t.Render(w, r, n, td)
	if err != nil {
		as.l.Error().Msgf("%s: render template error > %s", op, err)
		return err
	}

	return nil
}
