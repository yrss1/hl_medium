package app

import (
	"fmt"
	"medium/internal/handler"
	"medium/internal/repository"
	"medium/internal/service/todo"
	"medium/pkg/server"
)

func Run() {
	dbSource := "postgres://postgres:postgres@localhost:5432/e-shop?sslmode=disable"

	repositories, err := repository.New(repository.WithPostgresStore(dbSource))
	if err != nil {
		panic(err.Error())
		return
	}
	todoService, err := todo.New(todo.WithTaskRepository(repositories.Task))
	if err != nil {
		panic(err)
		return
	}
	handlers, err := handler.New(handler.Dependencies{TodoService: todoService}, handler.WithHTTPHandler())
	if err != nil {
		panic(err)
		return
	}
	servers, err := server.New(server.WithHTTPServer(handlers.HTTP, "8080"))
	if err != nil {
		panic(err)
		return
	}
	if err = servers.Run(); err != nil {
		fmt.Println("run error")
		panic(err)
		return
	}
	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
}
