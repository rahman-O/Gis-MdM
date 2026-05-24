package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// RolloutHandlers serves assignment and rollout endpoints (018).
type RolloutHandlers struct {
	assign  *profileapp.AssignmentService
	rollout *profileapp.RolloutStatusService
	enable  *profileapp.EnableService
	draft   *profileapp.DraftService
}

func NewRolloutHandlers(
	assign *profileapp.AssignmentService,
	rollout *profileapp.RolloutStatusService,
	enable *profileapp.EnableService,
	draft *profileapp.DraftService,
) *RolloutHandlers {
	return &RolloutHandlers{assign: assign, rollout: rollout, enable: enable, draft: draft}
}

func (h *RolloutHandlers) RegisterOnProfile(g *gin.RouterGroup) {
	g.GET("/:id/versions", h.ListVersions)
	g.POST("/:id/versions/:versionId/fork-draft", h.ForkDraft)
	g.GET("/:id/assignments", h.ListAssignments)
	g.GET("/:id/assignments/impact", h.AssignmentImpact)
	g.PUT("/:id/assignments", h.PutAssignment)
	g.DELETE("/:id/assignments/:assignmentId", h.DeleteAssignment)
	g.GET("/:id/rollout/devices", h.ListRolloutDevices)
	g.POST("/:id/rollout/recompute", h.RecomputeRollout)
	g.POST("/:id/disable", h.Disable)
	g.POST("/:id/enable", h.Enable)
}

func profileIDParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return 0, false
	}
	return id, true
}

func (h *RolloutHandlers) ListVersions(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	data, err := h.draft.ListVersions(c.Request.Context(), p, profileID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) ForkDraft(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	versionID, err := strconv.Atoi(c.Param("versionId"))
	if err != nil || versionID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	meta, err := h.draft.ForkDraftFromVersion(c.Request.Context(), p, profileID, versionID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, meta)
}

func (h *RolloutHandlers) ListAssignments(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	data, err := h.assign.List(c.Request.Context(), p, profileID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) AssignmentImpact(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	treeNodeID, _ := strconv.Atoi(c.Query("treeNodeId"))
	data, err := h.assign.Impact(c.Request.Context(), p, profileID, treeNodeID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) PutAssignment(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	var req domain.PutAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := h.assign.Put(c.Request.Context(), p, profileID, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) DeleteAssignment(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	assignmentID, err := strconv.Atoi(c.Param("assignmentId"))
	if err != nil || assignmentID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	if err := h.assign.Delete(c.Request.Context(), p, profileID, assignmentID); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (h *RolloutHandlers) ListRolloutDevices(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	q := domain.RolloutDevicesQuery{Status: c.Query("status")}
	if v := c.Query("treeNodeId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			q.TreeNodeID = &n
		}
	}
	q.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	q.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	data, err := h.rollout.ListDevices(c.Request.Context(), p, profileID, q)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) RecomputeRollout(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	if err := h.rollout.RecomputeProfile(c.Request.Context(), p, profileID); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, gin.H{"recomputed": true})
}

func (h *RolloutHandlers) Disable(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	data, err := h.enable.SetEnabled(c.Request.Context(), p, profileID, false)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *RolloutHandlers) Enable(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, ok := profileIDParam(c)
	if !ok {
		return
	}
	data, err := h.enable.SetEnabled(c.Request.Context(), p, profileID, true)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}
