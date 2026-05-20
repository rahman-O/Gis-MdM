package http

import (
	"errors"
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

func (h *Handler) JSON(c *gin.Context) {
	body, err := h.svc.JSON(c.Request.Context(), c.Param("id"), parseQuery(c))
	if err != nil {
		if errors.Is(err, qrapp.ErrNotFound) {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(body))
}

func (h *Handler) PNG(c *gin.Context) {
	png, err := h.svc.PNG(c.Request.Context(), c.Param("id"), parseQuery(c))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", png)
}
