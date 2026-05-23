package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	qrapp "github.com/gis-mdm/server-backend-go/internal/modules/qrcode/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
)

type Handler struct {
	svc *qrapp.Service
}

func NewHandler(svc *qrapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/json/:id", h.JSON)
	g.GET("/:id", h.PNG)
}

func parseQuery(c *gin.Context) domain.QRQuery {
	size, _ := strconv.Atoi(c.Query("size"))
	return domain.QRQuery{
		DeviceID:       c.Query("deviceId"),
		CreateOnDemand: c.Query("create"),
		UseID:          c.Query("useId"),
		Groups:         c.QueryArray("group"),
		Size:           size,
	}
}

func writeQRError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, qrapp.ErrNotFound):
		c.String(http.StatusNotFound, "configuration not found for this QR key")
	case errors.Is(err, qrapp.ErrMainAppURLMissing):
		c.String(http.StatusBadRequest, "main application has no download URL: upload an APK for the Main App version or set launcher URL override")
	default:
		c.String(http.StatusInternalServerError, "failed to generate QR provisioning data")
	}
}

func (h *Handler) JSON(c *gin.Context) {
	body, err := h.svc.JSON(c.Request.Context(), c.Param("id"), parseQuery(c))
	if err != nil {
		slog.Warn("qr json failed", "key", c.Param("id"), "err", err)
		writeQRError(c, err)
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(body))
}

func (h *Handler) PNG(c *gin.Context) {
	png, err := h.svc.PNG(c.Request.Context(), c.Param("id"), parseQuery(c))
	if err != nil {
		slog.Warn("qr png failed", "key", c.Param("id"), "err", err)
		writeQRError(c, err)
		return
	}
	c.Data(http.StatusOK, "image/png", png)
}
