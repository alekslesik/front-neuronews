package router

import (
	"net/http"

	"github.com/alekslesik/neuro-news/internal/app/handler"
)

type Router struct {
	h *handler.AppHandler
}

func New(handler *handler.AppHandler) *Router {
	return &Router{
		h: handler,
	}
}

func (r *Router) Route() http.Handler {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	mux := http.NewServeMux()

	// dynamic pages
	mux.HandleFunc("/{$}", r.h.ArticleHandler.Home)

	mux.Handle("/news/", http.RedirectHandler("/", http.StatusSeeOther))
	mux.HandleFunc("/news/{category}/{article}/", r.h.ArticleHandler.Article)
	mux.HandleFunc("/news/{category}/", r.h.ArticleHandler.Category)

	// static pages
	mux.HandleFunc("/about", r.h.CommonHandler.GetAboutPage)
	mux.HandleFunc("/contact", r.h.CommonHandler.GetContactPage)
	mux.HandleFunc("/advertisement", r.h.CommonHandler.GetAdvertisementPage)
	mux.HandleFunc("/privacy", r.h.CommonHandler.GetPrivacyPage)

	// file server for static files
	fileServer := http.FileServer(http.Dir("./website/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
