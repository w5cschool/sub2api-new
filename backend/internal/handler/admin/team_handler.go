package admin

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

type TeamMemberRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}

type CreateTeamRequest struct {
	Name        string              `json:"name" binding:"required"`
	Description *string             `json:"description"`
	Status      string              `json:"status"`
	Members     []TeamMemberRequest `json:"members"`
}

type UpdateTeamRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

type ReplaceTeamMembersRequest struct {
	Members []TeamMemberRequest `json:"members"`
}

func (h *TeamHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	filters := service.TeamFilters{
		Search: strings.TrimSpace(c.Query("search")),
		Status: strings.TrimSpace(c.Query("status")),
	}
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	teams, result, err := h.teamService.List(c.Request.Context(), page, pageSize, filters, sortBy, sortOrder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, dto.TeamsFromService(teams), result.Total, page, pageSize)
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	id, ok := parseTeamID(c)
	if !ok {
		return
	}
	team, err := h.teamService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamFromService(team))
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	team, err := h.teamService.Create(c.Request.Context(), service.TeamInput{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Members:     teamMemberRequestsToInputs(req.Members),
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamFromService(team))
}

func (h *TeamHandler) Update(c *gin.Context) {
	id, ok := parseTeamID(c)
	if !ok {
		return
	}
	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	team, err := h.teamService.Update(c.Request.Context(), id, service.TeamInput{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamFromService(team))
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id, ok := parseTeamID(c)
	if !ok {
		return
	}
	if err := h.teamService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "team deleted"})
}

func (h *TeamHandler) ReplaceMembers(c *gin.Context) {
	id, ok := parseTeamID(c)
	if !ok {
		return
	}
	var req ReplaceTeamMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	team, err := h.teamService.ReplaceMembers(c.Request.Context(), id, teamMemberRequestsToInputs(req.Members))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamFromService(team))
}

func parseTeamID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid team ID")
		return 0, false
	}
	return id, true
}

func teamMemberRequestsToInputs(members []TeamMemberRequest) []service.TeamMemberInput {
	out := make([]service.TeamMemberInput, 0, len(members))
	for _, member := range members {
		out = append(out, service.TeamMemberInput{UserID: member.UserID, Role: member.Role})
	}
	return out
}
