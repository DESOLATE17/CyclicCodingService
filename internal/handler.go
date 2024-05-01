package internal

import (
	"bytes"
	_ "cyclic_code/docs"
	"cyclic_code/pkg/coder"
	"cyclic_code/pkg/decoder"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.PUT("/code", h.DataLink)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// DataLink godoc
// @Summary      Канальный уровень
// @Description  Эмитирует взаимодействие с удаленным сетевым узлом через канал с помехами, теряет сообщение с вероятностью 1% и использует циклический код[7, 4]
// @Tags         DataLink
// @Accept       json
// @Produce      json
// @Param        user  body  models.Data  true  "сегмент данных"
// @Success      200
// @Failure      400  {object}  error
// @Router       /code [post]
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

	context.Status(http.StatusOK)

	// transport layer url
	url := "http://localhost:8081/transfer"

	req, err := http.NewRequest("POST", url, bytes.NewReader(decodedSegment))
	if err != nil {
		h.logger.Error("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		h.logger.Error("Error sending request:", err)
		return
	}
	defer res.Body.Close()
}
