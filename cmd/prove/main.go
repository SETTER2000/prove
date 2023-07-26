package main

import (
	"github.com/SETTER2000/prove/internal/app"
	"github.com/SETTER2000/prove/pkg/log/logger"
)

func main() {
	logger.New()
	app.Run()
}
