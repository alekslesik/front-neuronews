package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alekslesik/neuro-news/pkg/logger"
)

// UserHandler handle requests related with users
type UserHandler struct {
	AppHandler *AppHandler
	l          *logger.Logger
}

// NewUserHandler create new instance of UserHandler
func NewUserHandler(appHandler *AppHandler, l *logger.Logger) *UserHandler {
	return &UserHandler{
		AppHandler: appHandler,
		l:          l,
	}
}

// GetUser return user by ID
func (a *ArticleHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	userID := 1
	user, err := a.AppHandler.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send user like a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
