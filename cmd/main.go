package main

import (
	"github.com/Verce11o/Hezzl-Go/internal/app"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	"github.com/Verce11o/Hezzl-Go/lib/logger"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg)

	app.Run(log, cfg)
}
