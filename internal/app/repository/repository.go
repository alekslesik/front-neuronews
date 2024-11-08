package repository

import (
	"database/sql"

	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type MySQLRepository struct {
	db *sql.DB
	l  *logger.Logger
}

func New(db *sql.DB, l *logger.Logger) *MySQLRepository {
	return &MySQLRepository{
		db: db,
		l:  l,
	}
}

func (r *MySQLRepository) GetArticleRepository() model.ArticleModel {
	return &MySQLArticleRepository{db: r.db, l: r.l}
}

func (r *MySQLRepository) GetUserRepository() model.UserModel {
	return &MySQLUserRepository{db: r.db, l: r.l}
}

func (r *MySQLRepository) GetImageRepository() model.ImageModel {
	return &MySQLImageRepository{db: r.db, l: r.l}
}
