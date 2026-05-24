package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/profiles/* endpoints.
type Handler struct {
	draft         *profileapp.DraftService
	publish       *profileapp.PublishService
	hub           *profileapp.HubService
	versionDelete *profileapp.VersionDeleteService
}

func NewHandler(draft *profileapp.DraftService, publish *profileapp.PublishService, hub *profileapp.HubService, versionDelete *profileapp.VersionDeleteService) *Handler {
	return &Handler{draft: draft, publish: publish, hub: hub, versionDelete: versionDelete}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("", h.List)
	g.POST("", h.Create)
	if h.hub != nil {
		RegisterHubRoutes(g, h.hub)
	}
	g.GET("/:id", h.GetMeta)
	g.GET("/:id/versions/:versionId", h.GetVersion)
	g.PUT("/:id/versions/:versionId", h.SaveVersion)
	g.GET("/:id/impact", h.Impact)
	g.POST("/:id/versions/:versionId/publish", h.Publish)
	if h.versionDelete != nil {
		g.DELETE("/:id/versions/:versionId", h.DeleteVersion)
	}
}

func principal(c *gin.Context) (*platformauth.Principal, bool) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		c.Status(403)
		return nil, false
	}
	return p, true
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, profileapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, profileapp.ErrDuplicateProfile):
		response.DuplicateEntity(c, "error.duplicate.profile")
	case errors.Is(err, profileapp.ErrProfileNotFound), errors.Is(err, profileapp.ErrVersionNotFound):
		response.ErrorEnvelope(c, "error.notfound.profile")
	case errors.Is(err, profileapp.ErrNotDraftVersion):
		response.ErrorEnvelope(c, "error.profile.version.notdraft")
	case errors.Is(err, profileapp.ErrConfirmImpactRequired):
		response.ErrorEnvelope(c, "error.profile.publish.confirm_required")
	case errors.Is(err, profileapp.ErrProfileDisabled):
		response.ErrorEnvelope(c, "error.profile.disabled")
	case errors.Is(err, profileapp.ErrVersionNotPublished):
		response.ErrorEnvelope(c, "error.profile.version.notPublished")
	case errors.Is(err, profileapp.ErrAssignmentConfirmRequired):
		response.ErrorEnvelope(c, "error.profile.assignment.confirmRequired")
	case errors.Is(err, profileapp.ErrAssignmentNotFound):
		response.ErrorEnvelope(c, "error.profile.assignment.nodeNotFound")
	case errors.Is(err, profileapp.ErrVersionDeleteActivePublished):
		response.ErrorEnvelope(c, "error.profile.version.delete.activePublished")
	case errors.Is(err, profileapp.ErrVersionDeleteAssigned):
		response.ErrorEnvelope(c, "error.profile.version.delete.assigned")
	case errors.Is(err, profileapp.ErrVersionDeleteDevicesTarget):
		response.ErrorEnvelope(c, "error.profile.version.delete.devicesTarget")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

func parseIDs(c *gin.Context) (profileID, versionID int, ok bool) {
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil || profileID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return 0, 0, false
	}
	versionID, err = strconv.Atoi(c.Param("versionId"))
	if err != nil || versionID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return 0, 0, false
	}
	return profileID, versionID, true
}

func (h *Handler) List(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var data []domain.ProfileListItem
	var err error
	if h.hub != nil {
		data, err = h.hub.List(c.Request.Context(), p)
	} else {
		data, err = h.draft.List(c.Request.Context(), p)
	}
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Create(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	meta, err := h.draft.Create(c.Request.Context(), p, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, meta)
}

func (h *Handler) GetMeta(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil || profileID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	meta, err := h.draft.GetMeta(c.Request.Context(), p, profileID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, meta)
}

func (h *Handler) GetVersion(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, versionID, ok := parseIDs(c)
	if !ok {
		return
	}
	resp, err := h.draft.GetVersion(c.Request.Context(), p, profileID, versionID)
	if err != nil {
		mapErr(c, err)
		return
	}
	data := cfgdomain.ConfigurationResponseMap(&resp.Payload)
	data["profileId"] = resp.ProfileID
	data["versionId"] = resp.VersionID
	data["versionNumber"] = resp.VersionNumber
	data["versionStatus"] = resp.Status
	response.OK(c, data)
}

func (h *Handler) SaveVersion(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, versionID, ok := parseIDs(c)
	if !ok {
		return
	}
	raw, err := c.GetRawData()
	if err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	payload, err := cfgdomain.ParseConfigurationBody(raw)
	if err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	id := profileID
	payload.ID = &id
	if err := h.draft.SaveDraft(c.Request.Context(), p, profileID, versionID, payload); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, gin.H{"profileId": profileID, "versionId": versionID})
}

func (h *Handler) Impact(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil || profileID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := h.publish.Impact(c.Request.Context(), p, profileID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) DeleteVersion(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	if h.versionDelete == nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	profileID, versionID, ok := parseIDs(c)
	if !ok {
		return
	}
	_, err := h.versionDelete.Delete(c.Request.Context(), p, profileID, versionID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, gin.H{"profileId": profileID, "versionId": versionID})
}

func (h *Handler) Publish(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, versionID, ok := parseIDs(c)
	if !ok {
		return
	}
	var req domain.PublishRequest
	_ = c.ShouldBindJSON(&req)
	data, err := h.publish.Publish(c.Request.Context(), p, profileID, versionID, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}
