package app

import (
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/usecase/repo/file"
	"github.com/SETTER2000/prove/internal/usecase/repo/memory"
	"github.com/SETTER2000/prove/internal/usecase/repo/sql"
	"github.com/SETTER2000/prove/scripts"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/SETTER2000/prove/internal/controller/http/v1"
	"github.com/SETTER2000/prove/internal/server"
	"github.com/SETTER2000/prove/internal/usecase"
	"github.com/SETTER2000/prove/pkg/log/logger"
	"github.com/xlab/closer"
)

var (
	OSString      string
	archString    string
	versionString = "N/A" // version app
	dateString    = "N/A" // date build
	commitString  = "N/A" // id commit
)

func Run() {
	closer.Bind(cleanup)
	// logging
	l := logger.GetLogger()
	l.Info("logger initialized")

	// seed
	rand.Seed(time.Now().UnixNano())

	var repo usecase.ProveRepo

	if !scripts.CheckEnvironFlag("DATABASE_DSN", config.GetConfig().ConnectDB) {
		if config.GetConfig().FileStorage == "" {
			l.Warn("In memory storage!!!")
			repo = memory.New()
		} else {
			l.Info("File storage - is work...")
			repo = file.New()
		}
	} else {
		// DB
		db, err := sql.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "db connection not created: %e\n", err)
		}
		l.Info("DB SQL - is work...")
		repo = sql.New(db)
	}

	fmt.Printf("%s:\n OS/Arch: \t%s/%s\n Version: \t%s\n Build: \t%s\n Commit: \t%s\n Author: \t%s\n", config.GetConfig().Name, OSString,
		archString, versionString, dateString, commitString, config.GetConfig().Author)

	// создаем слой usecase передав объект подключения к хранилищу данных repo
	apiUseCase := usecase.New(repo)
	if err := apiUseCase.ReadService(); err != nil {
		l.Error(fmt.Errorf("app - Read - shorturlUseCase.ReadService: %w", err))
	}

	handlers := v1.NewServerHandler(apiUseCase)

	httpServer := server.New(handlers.InitRoutes(),
		server.Host(),
		// опция подключения HTTPS
		server.EnableHTTPS(),
	)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		if err := apiUseCase.SaveService(); err != nil {
			l.Error(fmt.Errorf("app - Save - proveUseCase.SaveService: %w", err))
		}
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	closer.Hold()

	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func cleanup() {
	fmt.Println("Hang on! I'm closing some DBs, wiping some trails..")
	time.Sleep(1 * time.Second)
	fmt.Println("  Done...")
}
