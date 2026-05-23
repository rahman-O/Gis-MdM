package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	syncapp "github.com/gis-mdm/server-backend-go/internal/modules/sync/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

type Handler struct {
	svc        *syncapp.Service
	hashSecret string
}

func NewHandler(svc *syncapp.Service, hashSecret string) *Handler {
	return &Handler{svc: svc, hashSecret: hashSecret}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/configuration/:deviceId", h.GetConfiguration)
	g.POST("/configuration/:deviceId", h.PostConfiguration)
	g.POST("/info", h.PostInfo)
	g.POST("/applicationSettings/:deviceId", h.PostApplicationSettings)
}

func (h *Handler) writeSync(c *gin.Context, resp *domain.SyncResponse) {
	if sig := sharedcrypto.SignSyncResponse(h.hashSecret, resp); sig != "" {
		c.Header("X-Response-Signature", sig)
	}
	c.Header("X-IP-Address", c.ClientIP())
	response.OK(c, resp)
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, syncapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, syncapp.ErrDeviceNotFound):
		response.ErrorEnvelope(c, "error.notfound.device")
	case errors.Is(err, syncapp.ErrDeviceExists):
		response.ErrorEnvelope(c, "error.duplicate.device")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// GetConfiguration godoc
// @Summary Get device configuration
// @Tags Sync
// @Produce json
// @Param deviceId path string true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /public/sync/configuration/{deviceId} [get]
func (h *Handler) GetConfiguration(c *gin.Context) {
	deviceID := c.Param("deviceId")
	resp, err := h.svc.GetConfiguration(c.Request.Context(), deviceID, c.GetHeader("X-Request-Signature"), c.GetHeader("X-CPU-Arch"))
	if err != nil {
		mapErr(c, err)
		return
	}
	h.writeSync(c, resp)
}

func (h *Handler) PostConfiguration(c *gin.Context) {
	deviceID := c.Param("deviceId")
	var opts domain.DeviceCreateOptions
	_ = c.ShouldBindJSON(&opts)
	resp, err := h.svc.EnrollConfiguration(c.Request.Context(), deviceID, opts, c.GetHeader("X-Request-Signature"), c.GetHeader("X-CPU-Arch"))
	if err != nil {
		mapErr(c, err)
		return
	}
	h.writeSync(c, resp)
}

func (h *Handler) PostInfo(c *gin.Context) {
	var info domain.DeviceInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.UpdateInfo(c.Request.Context(), info); err != nil {
		mapErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Status: "OK"})
}

func (h *Handler) PostApplicationSettings(c *gin.Context) {
	deviceID := c.Param("deviceId")
	var settings []domain.SyncApplicationSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.SaveApplicationSettings(c.Request.Context(), deviceID, settings); err != nil {
		mapErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.Envelope{Status: "OK"})
}
