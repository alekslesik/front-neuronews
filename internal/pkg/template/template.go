package template

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"path/filepath"
	"time"

	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/pkg/logger"
	"github.com/justinas/nosurf"
)

const (
	UserID   = "userID"
	UserName = "userName"
)

type Cache map[string]*template.Template

type Template struct {
	cache        Cache
	log          *logger.Logger
	TemplateData TemplateData
}

type TemplateData struct {
	// AuthenticatedUser *models.User
	CurrentYear         int
	UserName            string
	Flash               string
	CSRFToken           string
	TemplateDataArticle TemplateDataArticle
	TemplateDataPage    *TemplateDataPage
	TemplateDataImages  TemplateDataImages

	// Form              *forms.Form
	// File              *models.File
	// Files             []*models.File
}

type TemplateDataPage struct {
	TotalPaginationPages  int
	CurrentPaginationPage int
	Category              string
}

type TemplateDataArticle struct {
	CarouselArticles       []model.Article
	TrendingArticlesTop    []model.Article
	TrendingArticlesBottom []model.Article
	NewsArticles           []model.Article
	SportArticles          []model.Article
	VideoArticles          []model.Article
	AllArticles            []model.Article
	PaginationArticles     []model.Article
	RandomArticles         []model.Article
	PaginationData         string
	Article                model.Article
}

type TemplateDataImages struct {
	RandomImages []model.Image
}

// New return instance of template
func New(log *logger.Logger) *Template {
	return &Template{
		cache:        make(Cache),
		log:          log,
		TemplateData: TemplateData{},
	}
}

// AddCache add new cache of files in dir to template cache
func (t *Template) AddCache(dir string) (*Template, error) {
	const op = "template.AddCache()"

	cache, err := t.newCache(dir)
	if err != nil {
		t.log.Error().Msgf("%s: add cache to template error > %s", op, err)
		return nil, err
	}

	t.cache = cache
	return t, nil
}

func (t *Template) newCache(dir string) (Cache, error) {
	const op = "template.newCache()"

	// init new map keeping cache
	cache := map[string]*template.Template{}

	// use func Glob to get all filepathes slice with '.page.html' ext
	entries, err := filepath.Glob(filepath.Join(dir, "**/*.page.html"))
	if err != nil {
		t.log.Error().Msgf("%s: error glob *.page.html in dir %v > %s", op, dir, err)
		return nil, err
	}

	for _, e := range entries {
		// get filename from filepath
		name := filepath.Base(e)

		// The template.FuncMap must be registered with the template set before
		// call the ParseFiles() method. This means we have to use template.New
		// create an empty template set, use the Funcs() method t
		ts, err := template.New(name).Funcs(functions).ParseFiles(e)
		if err != nil {
			t.log.Error().Msgf("%s: template create error > %s", op, err)
			return nil, err
		}

		// use ParseGlob to add all frame patterns (base.layout.html)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			t.log.Error().Msgf("%s: error glob *.layout.html to template > %s", op, err)
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			t.log.Err(err).Msgf("%s: error glob *.partial.html to template > %s", op, err)
			return nil, err
		}

		// add received patterns set to cache, using page name
		// (ext home.page.html) as a key for our map
		cache[name] = ts
	}

	return cache, nil
}

// Render add template data and render
func (t *Template) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) error {
	const op = "templates.Render()"

	// extract pattern depending "name"
	ts, ok := t.cache[name]
	if !ok {
		t.log.Error().Msgf("%s: pattern %s not exist", op, name)
		return fmt.Errorf("%s: pattern %s not exist", op, name)
	}

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write template to the buffer, instead straight to http.ResponseWriter
	err := ts.Execute(buf, AddDefaultData(td, r))
	if err != nil {
		t.log.Error().Msgf("%s: template %v not executed > %s", op, ts, err)
		return err
	}

	// write buffer to http.ResponseWriter
	buf.WriteTo(w)
	return nil
}

// AddDefaultData Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	// Add current time.
	td.CurrentYear = time.Now().Year()
	// Add flash message.
	// Check if user is authenticate.
	// td.AuthenticatedUser = AuthenticatedUser(r)
	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)
	// Add User Name to template
	// td.UserName = app.UserName

	return td
}

// // Return userID ID from session
// func AuthenticatedUser(r *http.Request) *models.User {
// 	user, ok := r.Context().Value(UserID).(*models.User)
// 	if !ok {
// 		return nil
// 	}
// 	return user
// }
