package main

import (
	"github.com/tv2169145/store_users-api/app"
	"github.com/tv2169145/store_users-api/logger"
)

func main() {
	logger.Info("about to start the application...")
	app.StartApplication()
}
