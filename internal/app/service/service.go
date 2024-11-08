package service

import (
	"github.com/alekslesik/neuro-news/internal/pkg/grabber"
	"github.com/alekslesik/neuro-news/internal/pkg/template"

	"github.com/alekslesik/neuro-news/internal/app/repository"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type Services struct {
	r *repository.MySQLRepository
	l *logger.Logger
	t *template.Template
	g *grabber.Grabber
}

func New(r *repository.MySQLRepository, l *logger.Logger, t *template.Template, g *grabber.Grabber) *Services {
	return &Services{
		r: r,
		t: t,
		l: l,
		g: g,
	}
}

func (s *Services) GetArticleService() ArticleService {
	return &articleService{
		ar: s.r.GetArticleRepository(),
		t:  s.t,
		l:  s.l,
		g:  s.g,
	}
}

func (s *Services) GetUserService() UserService {
	return &userService{
		ur: s.r.GetUserRepository(),
	}
}

func (s *Services) GetImageService() ImageService {
	return &imageService{
		ir: s.r.GetImageRepository(),
		l:  s.l,
		g:  s.g,
	}
}
