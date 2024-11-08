package repository

import (
	"database/sql"

	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type MySQLUserRepository struct {
	db *sql.DB
	l  *logger.Logger
}

func (r *MySQLUserRepository) GetUserByID(id int) (*model.User, error) {
	return nil, nil
}
