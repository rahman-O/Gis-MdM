package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	notifapp "github.com/gis-mdm/server-backend-go/internal/modules/notifications/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

func checkSig(signature, value string) bool {
	return sharedcrypto.CheckRequestSignature(signature, value)
}

type Handler struct {
	svc            *notifapp.Service
	secureEnroll   bool
	hashSecret     string
}

func NewHandler(svc *notifapp.Service, secureEnrollment bool, hashSecret string) *Handler {
	return &Handler{svc: svc, secureEnroll: secureEnrollment, hashSecret: hashSecret}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/device/:deviceNumber", h.GetDeviceMessages)
}

// RegisterPolling mounts /rest/notification/polling/:deviceNumber on the engine root.
func (h *Handler) RegisterPolling(engine *gin.Engine) {
	engine.GET("/rest/notification/polling/:deviceNumber", h.LongPoll)
}

// GetDeviceMessages godoc
// @Summary Get device notifications
// @Tags Notifications
// @Produce json
// @Param deviceNumber path string true "Device number"
// @Success 200 {object} response.Envelope
// @Router /notifications/device/{deviceNumber} [get]
func (h *Handler) GetDeviceMessages(c *gin.Context) {
	deviceNumber := c.Param("deviceNumber")
	msgs, err := h.svc.GetPending(c.Request.Context(), deviceNumber)
	if err != nil {
		if errors.Is(err, notifapp.ErrDeviceNotFound) {
			response.ErrorEnvelope(c, "error.notfound.device")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, msgs)
}

func (h *Handler) LongPoll(c *gin.Context) {
	deviceNumber := c.Param("deviceNumber")
	if h.secureEnroll {
		sig := c.GetHeader("X-Request-Signature")
		if !checkSig(sig, h.hashSecret+deviceNumber) {
			response.PermissionDenied(c)
			return
		}
	}
	msgs, err := h.svc.PollPending(c.Request.Context(), deviceNumber)
	if err != nil {
		if errors.Is(err, notifapp.ErrDeviceNotFound) {
			response.ErrorEnvelope(c, "error.notfound.device")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if len(msgs) == 0 {
		c.Status(http.StatusOK)
		return
	}
	response.OK(c, msgs)
}
