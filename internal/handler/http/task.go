package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"medium/internal/domain/task"
	"medium/internal/service/todo"
	"medium/pkg/server/response"
	"medium/pkg/store"
	"net/http"
	"strings"
)

type TaskHandler struct {
	todoService *todo.Service
}

func NewTaskHandler(s *todo.Service) *TaskHandler {
	return &TaskHandler{todoService: s}
}
func (h *TaskHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/tasks")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT(":id/done", h.mark)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)
	}
}
func (h *TaskHandler) list(c *gin.Context) {
	status := c.DefaultQuery("status", "active")

	if status != "active" && status != "done" {
		response.BadRequest(c, errors.New("invalid status"), status)
		return
	}

	res, err := h.todoService.ListTasks(c, status)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, res)
}
func (h *TaskHandler) add(c *gin.Context) {
	req := task.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.todoService.CreateTask(c, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "failed to parse:"):
			response.BadRequest(c, err, req)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
	response.OK(c, res)
}

func (h *TaskHandler) get(c *gin.Context) {
	id := c.Param("id")
	res, err := h.todoService.GetTask(c, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

func (h *TaskHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := task.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := h.todoService.UpdateTask(c, id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.todoService.DeleteTask(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) mark(c *gin.Context) {
	id := c.Param("id")
	if err := h.todoService.MarkTask(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
	c.Status(http.StatusNoContent)
}
