package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// TeamHandler handles team-scoped user endpoints.
type TeamHandler struct {
	teamService  *service.TeamService
	usageService *service.UsageService
}

func NewTeamHandler(teamService *service.TeamService, usageService *service.UsageService) *TeamHandler {
	return &TeamHandler{teamService: teamService, usageService: usageService}
}

// Me returns the current user's team membership and role.
// GET /api/v1/team/me
func (h *TeamHandler) Me(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	membership, err := h.teamService.GetMyTeam(c.Request.Context(), subject.UserID)
	if err != nil {
		if infraerrors.IsNotFound(err) {
			response.Success(c, dto.TeamMe{CanViewUsage: false})
			return
		}
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamMe{
		Team:         dto.TeamFromService(membership.Team),
		Membership:   dto.TeamMemberFromService(membership),
		CanViewUsage: membership.Role == service.TeamRoleLeader,
	})
}

// Members lists all members in the leader's current team.
// GET /api/v1/team/members
func (h *TeamHandler) Members(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	members, err := h.teamService.ListLeaderMembers(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamMembersFromService(members))
}

// Usage lists usage logs visible to a team leader.
// GET /api/v1/team/usage
func (h *TeamHandler) Usage(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	memberID, err := parseOptionalTeamMemberID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userIDs, err := h.teamService.LeaderVisibleUserIDs(c.Request.Context(), subject.UserID, memberID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	filters, err := parseTeamUsageListFilters(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	filters.UserIDs = userIDs

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}
	records, result, err := h.usageService.ListWithFilters(c.Request.Context(), params, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.UsageLog, 0, len(records))
	for i := range records {
		out = append(out, *dto.UsageLogFromService(&records[i]))
	}
	response.Paginated(c, out, result.Total, page, pageSize)
}

// Stats returns aggregate usage stats visible to a team leader.
// GET /api/v1/team/usage/stats
func (h *TeamHandler) Stats(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	memberID, err := parseOptionalTeamMemberID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userIDs, err := h.teamService.LeaderVisibleUserIDs(c.Request.Context(), subject.UserID, memberID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	filters, err := parseTeamUsageStatsFilters(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	filters.UserIDs = userIDs
	stats, err := h.usageService.GetStatsWithFilters(c.Request.Context(), filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, stats)
}

// MembersSummary returns per-member today/total usage for the leader's team.
// GET /api/v1/team/usage/members-summary
func (h *TeamHandler) MembersSummary(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	items, err := h.teamService.GetLeaderMembersUsageSummary(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.TeamMemberUsageSummaryFromService(items))
}

func parseOptionalTeamMemberID(c *gin.Context) (int64, error) {
	raw := strings.TrimSpace(c.Query("member_id"))
	if raw == "" {
		return 0, nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, invalidQueryError("Invalid member_id")
	}
	return id, nil
}

func parseTeamUsageListFilters(c *gin.Context) (usagestats.UsageLogFilters, error) {
	filters, err := parseTeamUsageCommonFilters(c)
	if err != nil {
		return usagestats.UsageLogFilters{}, err
	}
	if exactTotalRaw := strings.TrimSpace(c.Query("exact_total")); exactTotalRaw != "" {
		parsed, err := strconv.ParseBool(exactTotalRaw)
		if err != nil {
			return usagestats.UsageLogFilters{}, invalidQueryError("Invalid exact_total value, use true or false")
		}
		filters.ExactTotal = parsed
	}
	userTZ := c.Query("timezone")
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", startDateStr, userTZ)
		if err != nil {
			return usagestats.UsageLogFilters{}, invalidQueryError("Invalid start_date format, use YYYY-MM-DD")
		}
		filters.StartTime = &t
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", endDateStr, userTZ)
		if err != nil {
			return usagestats.UsageLogFilters{}, invalidQueryError("Invalid end_date format, use YYYY-MM-DD")
		}
		t = t.AddDate(0, 0, 1)
		filters.EndTime = &t
	}
	return filters, nil
}

func parseTeamUsageStatsFilters(c *gin.Context) (usagestats.UsageLogFilters, error) {
	filters, err := parseTeamUsageCommonFilters(c)
	if err != nil {
		return usagestats.UsageLogFilters{}, err
	}
	userTZ := c.Query("timezone")
	now := timezone.NowInUserLocation(userTZ)
	var startTime, endTime time.Time
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr != "" && endDateStr != "" {
		startTime, err = timezone.ParseInUserLocation("2006-01-02", startDateStr, userTZ)
		if err != nil {
			return usagestats.UsageLogFilters{}, invalidQueryError("Invalid start_date format, use YYYY-MM-DD")
		}
		endTime, err = timezone.ParseInUserLocation("2006-01-02", endDateStr, userTZ)
		if err != nil {
			return usagestats.UsageLogFilters{}, invalidQueryError("Invalid end_date format, use YYYY-MM-DD")
		}
		endTime = endTime.AddDate(0, 0, 1)
	} else {
		switch c.DefaultQuery("period", "today") {
		case "today":
			startTime = timezone.StartOfDayInUserLocation(now, userTZ)
		case "week":
			startTime = now.AddDate(0, 0, -7)
		case "month":
			startTime = now.AddDate(0, -1, 0)
		default:
			startTime = timezone.StartOfDayInUserLocation(now, userTZ)
		}
		endTime = now
	}
	filters.StartTime = &startTime
	filters.EndTime = &endTime
	return filters, nil
}

func parseTeamUsageCommonFilters(c *gin.Context) (usagestats.UsageLogFilters, error) {
	var filters usagestats.UsageLogFilters
	if apiKeyIDStr := strings.TrimSpace(c.Query("api_key_id")); apiKeyIDStr != "" {
		id, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
		if err != nil || id < 0 {
			return filters, invalidQueryError("Invalid api_key_id")
		}
		filters.APIKeyID = id
	}
	filters.Model = c.Query("model")
	filters.BillingMode = strings.TrimSpace(c.Query("billing_mode"))
	if requestTypeStr := strings.TrimSpace(c.Query("request_type")); requestTypeStr != "" {
		parsed, err := service.ParseUsageRequestType(requestTypeStr)
		if err != nil {
			return filters, err
		}
		value := int16(parsed)
		filters.RequestType = &value
	} else if streamStr := strings.TrimSpace(c.Query("stream")); streamStr != "" {
		val, err := strconv.ParseBool(streamStr)
		if err != nil {
			return filters, invalidQueryError("Invalid stream value, use true or false")
		}
		filters.Stream = &val
	}
	if billingTypeStr := strings.TrimSpace(c.Query("billing_type")); billingTypeStr != "" {
		val, err := strconv.ParseInt(billingTypeStr, 10, 8)
		if err != nil {
			return filters, invalidQueryError("Invalid billing_type")
		}
		bt := int8(val)
		filters.BillingType = &bt
	}
	return filters, nil
}

type invalidQueryError string

func (e invalidQueryError) Error() string {
	return string(e)
}
