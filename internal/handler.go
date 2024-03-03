package internal

import (
	"cyclic_code/pkg/coder"
	"cyclic_code/pkg/decoder"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Handler struct {
	logger  *logrus.Entry
	coder   coder.Coder
	decoder decoder.Decoder
}

func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{
		logger:  logger.WithField("component", "handler"),
		coder:   coder.NewCoderImpl(11, 1, logger),
		decoder: decoder.NewDecoderImpl(logger),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.PUT("/", h.DataLink)

	return r
}

func (h *Handler) DataLink(context *gin.Context) {
	segment, err := io.ReadAll(context.Request.Body)
	if err != nil {
		h.logger.Error(err)
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	h.logger.Info("segment: ", segment)

	encodedSegment := h.coder.EncodeSegment(segment)
	h.logger.Info("encoded segment: ", encodedSegment)

	decodedSegment := h.decoder.DecodeSegment(encodedSegment)
	h.logger.Info("decoded segment: ", decodedSegment)

	context.Data(http.StatusOK, "application/json", decodedSegment)
}
