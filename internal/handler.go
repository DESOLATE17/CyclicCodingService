package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Entry
}

func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{
		logger: logger.WithField("component", "handler"),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.PUT("/", h.DataLink)

	return r
}

func (h *Handler) DataLink(context *gin.Context) {

}
