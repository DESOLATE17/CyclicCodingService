package main

import (
	"cyclic_code/internal"
	"github.com/sirupsen/logrus"
	"time"
)

// @title Канальный уровень
// @version 1.0
// @description Данный уровень эмитирует взаимодействие с удаленным сетевым узлом через канал с помехами

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	handler := internal.NewHandler(logger)
	r := handler.InitRouter()
	r.Run("0.0.0.0:8080")
}
