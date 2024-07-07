package handler

import (
	"github.com/gin-gonic/gin"
	"medium/internal/config"
	"medium/internal/handler/http"
	"medium/internal/service/todo"
	"medium/pkg/server/router"
)

type Dependencies struct {
	Configs     config.Configs
	TodoService *todo.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
}
type Configuration func(h *Handler) error

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
		HTTP:         router.New(),
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()
		taskHandler := http.NewTaskHandler(h.dependencies.TodoService)

		api := h.HTTP.Group("/api/todo-list/")
		{
			taskHandler.Routes(api)
		}
		return
	}
}
