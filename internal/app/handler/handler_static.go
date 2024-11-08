package handler

import (
	"net/http"

	"github.com/alekslesik/neuro-news/internal/pkg/template"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

// CommonHandler handle requests related with articles
type CommonHandler struct {
	AppHandler *AppHandler
	l          *logger.Logger
}

// NewCommonHandler create new instance of CommonHandler.
func NewCommonHandler(appHandler *AppHandler, l *logger.Logger) *CommonHandler {
	return &CommonHandler{
		AppHandler: appHandler,
		l:          l,
	}
}

func (a *CommonHandler) GetAboutPage(w http.ResponseWriter, r *http.Request) {
	// const op = "GetAboutPage()"

	a.AppHandler.templates.Render(w, r, "about.page.html", &template.TemplateData{})

}

func (a *CommonHandler) GetContactPage(w http.ResponseWriter, r *http.Request) {
	// const op = "GetContactPage()"

	a.AppHandler.templates.Render(w, r, "contact.page.html", &template.TemplateData{})

}

func (a *CommonHandler) GetAdvertisementPage(w http.ResponseWriter, r *http.Request) {
	// const op = "GetContactPage()"

	a.AppHandler.templates.Render(w, r, "advertisement.page.html", &template.TemplateData{})

}

func (a *CommonHandler) GetPrivacyPage(w http.ResponseWriter, r *http.Request) {
	// const op = "GetContactPage()"

	a.AppHandler.templates.Render(w, r, "privacy.page.html", &template.TemplateData{})

}
