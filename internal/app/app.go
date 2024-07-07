package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"medium/internal/config"
	"medium/internal/handler"
	"medium/internal/repository"
	"medium/internal/service/todo"
	"medium/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	configs, err := config.New()
	if err != nil {
		fmt.Printf("ERR_INIT_CONFIGS: %v", err)
		return
	}

	repositories, err := repository.New(repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		fmt.Printf("ERR_INIT_REPOSITORIES: %v", err)
		return
	}

	todoService, err := todo.New(todo.WithTaskRepository(repositories.Task))
	if err != nil {
		fmt.Printf("ERR_INIT_TODO_SERVICE: %v", err)
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			TodoService: todoService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		fmt.Printf("ERR_INIT_HANDLERS: %v", err)
		return
	}

	servers, err := server.New(server.WithHTTPServer(handlers.HTTP, configs.APP.Port))
	if err != nil {
		fmt.Printf("ERR_RUN_SERVERS: %v", err)
		return
	}
	gin.SetMode(gin.ReleaseMode)
	if err = servers.Run(); err != nil {
		fmt.Println("run error")
		panic(err)
		return
	}
	fmt.Println("http server started on http://localhost:" + configs.APP.Port)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err = servers.Stop(ctx); err != nil {
		panic(err)
	}

	fmt.Println("running cleanup tasks...")

	fmt.Println("server was successful shutdown.")

	//fmt.Println("Press any key to exit...")
	//var input string
	//fmt.Scanln(&input)
}
