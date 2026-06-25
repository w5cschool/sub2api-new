package admin

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// toResponsePagination converts pagination.PaginationResult to response.PaginationResult
func toResponsePagination(p *pagination.PaginationResult) *response.PaginationResult {
	if p == nil {
		return nil
	}
	return &response.PaginationResult{
		Total:    p.Total,
		Page:     p.Page,
		PageSize: p.PageSize,
		Pages:    p.Pages,
	}
}

// SubscriptionHandler handles admin subscription management
type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

// NewSubscriptionHandler creates a new admin subscription handler
func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// AssignSubscriptionRequest represents assign subscription request
type AssignSubscriptionRequest struct {
	UserID       int64    `json:"user_id" binding:"required"`
	GroupID      int64    `json:"group_id" binding:"required"`
	ValidityDays int      `json:"validity_days" binding:"omitempty,max=36500"` // max 100 years
	Notes        string   `json:"notes"`
	PriceUSD     *float64 `json:"price_usd" binding:"omitempty,min=0"`
}

// BulkAssignSubscriptionRequest represents bulk assign subscription request
type BulkAssignSubscriptionRequest struct {
	UserIDs      []int64  `json:"user_ids" binding:"required,min=1"`
	GroupID      int64    `json:"group_id" binding:"required"`
	ValidityDays int      `json:"validity_days" binding:"omitempty,max=36500"` // max 100 years
	Notes        string   `json:"notes"`
	PriceUSD     *float64 `json:"price_usd" binding:"omitempty,min=0"`
}

// AdjustSubscriptionRequest represents adjust subscription request (extend or shorten)
type AdjustSubscriptionRequest struct {
	Days int `json:"days" binding:"required,min=-36500,max=36500"` // negative to shorten, positive to extend
}

// List handles listing all subscriptions with pagination and filters
// GET /api/v1/admin/subscriptions
func (h *SubscriptionHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)

	// Parse optional filters
	var userID, groupID *int64
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		if id, err := strconv.ParseInt(groupIDStr, 10, 64); err == nil {
			groupID = &id
		}
	}
	status := c.Query("status")
	platform := c.Query("platform")

	// Parse sorting parameters
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	subscriptions, pagination, err := h.subscriptionService.List(c.Request.Context(), page, pageSize, userID, groupID, status, platform, sortBy, sortOrder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.AdminUserSubscription, 0, len(subscriptions))
	for i := range subscriptions {
		out = append(out, *dto.UserSubscriptionFromServiceAdmin(&subscriptions[i]))
	}
	response.PaginatedWithResult(c, out, toResponsePagination(pagination))
}

// ListRecords handles listing admin subscription assignment records.
// GET /api/v1/admin/subscription-records
func (h *SubscriptionHandler) ListRecords(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	filters, ok := parseSubscriptionRecordFilters(c)
	if !ok {
		return
	}

	records, pagination, err := h.subscriptionService.ListSubscriptionRecords(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.SubscriptionRecord, 0, len(records))
	for i := range records {
		out = append(out, *dto.SubscriptionRecordFromService(&records[i]))
	}
	response.PaginatedWithResult(c, out, toResponsePagination(pagination))
}

// RecordStats handles amount statistics for admin subscription assignment records.
// GET /api/v1/admin/subscription-records/stats
func (h *SubscriptionHandler) RecordStats(c *gin.Context) {
	filters, ok := parseSubscriptionRecordFilters(c)
	if !ok {
		return
	}

	stats, err := h.subscriptionService.GetSubscriptionRecordStats(c.Request.Context(), filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.SubscriptionRecordStatsFromService(stats))
}

// ExportRecords handles exporting admin subscription records as CSV.
// GET /api/v1/admin/subscription-records/export
func (h *SubscriptionHandler) ExportRecords(c *gin.Context) {
	filters, ok := parseSubscriptionRecordFilters(c)
	if !ok {
		return
	}

	records, err := h.subscriptionService.ExportSubscriptionRecords(c.Request.Context(), filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	loc := subscriptionRecordCSVLocation(c)
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"record_id",
		"operation",
		"user_id",
		"user_email",
		"group_id",
		"group_name",
		"subscription_id",
		"price_usd",
		"validity_days",
		"starts_at",
		"expires_at",
		"assigned_by_id",
		"assigned_by_email",
		"recorded_at",
		"notes",
	}); err != nil {
		response.InternalError(c, "Failed to export subscription records: "+err.Error())
		return
	}

	for i := range records {
		record := records[i]
		if err := writer.Write([]string{
			strconv.FormatInt(record.ID, 10),
			record.Operation,
			strconv.FormatInt(record.UserID, 10),
			subscriptionRecordUserEmail(&record),
			strconv.FormatInt(record.GroupID, 10),
			subscriptionRecordGroupName(&record),
			subscriptionRecordInt64PtrString(record.SubscriptionID),
			strconv.FormatFloat(record.PriceUSD, 'f', 2, 64),
			strconv.Itoa(record.ValidityDays),
			formatSubscriptionRecordCSVTime(record.StartsAt, loc),
			formatSubscriptionRecordCSVTime(record.ExpiresAt, loc),
			subscriptionRecordInt64PtrString(record.AssignedBy),
			subscriptionRecordAssignedByEmail(&record),
			formatSubscriptionRecordCSVTime(record.CreatedAt, loc),
			record.Notes,
		}); err != nil {
			response.InternalError(c, "Failed to export subscription records: "+err.Error())
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		response.InternalError(c, "Failed to export subscription records: "+err.Error())
		return
	}

	filename := subscriptionRecordExportFilename(filters, loc)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", buf.Bytes())
}

func subscriptionRecordCSVLocation(c *gin.Context) *time.Location {
	if tz := c.Query("timezone"); tz != "" {
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}
	return time.UTC
}

func subscriptionRecordExportFilename(filters service.SubscriptionRecordFilters, loc *time.Location) string {
	start := "all"
	if filters.StartTime != nil {
		start = filters.StartTime.In(loc).Format("20060102")
	}
	end := "all"
	if filters.EndTime != nil {
		end = filters.EndTime.In(loc).Format("20060102")
	}
	return fmt.Sprintf("subscription-records-%s-to-%s.csv", start, end)
}

func formatSubscriptionRecordCSVTime(t time.Time, loc *time.Location) string {
	if t.IsZero() {
		return ""
	}
	return t.In(loc).Format("2006-01-02 15:04:05")
}

func subscriptionRecordInt64PtrString(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}

func subscriptionRecordUserEmail(record *service.SubscriptionRecord) string {
	if record == nil {
		return ""
	}
	if record.User != nil {
		return record.User.Email
	}
	if record.Subscription != nil && record.Subscription.User != nil {
		return record.Subscription.User.Email
	}
	return ""
}

func subscriptionRecordGroupName(record *service.SubscriptionRecord) string {
	if record == nil {
		return ""
	}
	if record.Group != nil {
		return record.Group.Name
	}
	if record.Subscription != nil && record.Subscription.Group != nil {
		return record.Subscription.Group.Name
	}
	return ""
}

func subscriptionRecordAssignedByEmail(record *service.SubscriptionRecord) string {
	if record == nil {
		return ""
	}
	if record.AssignedByUser != nil {
		return record.AssignedByUser.Email
	}
	if record.Subscription != nil && record.Subscription.AssignedByUser != nil {
		return record.Subscription.AssignedByUser.Email
	}
	return ""
}

func parseSubscriptionRecordFilters(c *gin.Context) (service.SubscriptionRecordFilters, bool) {
	var filters service.SubscriptionRecordFilters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid user_id")
			return filters, false
		}
		filters.UserID = &id
	}
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		id, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid group_id")
			return filters, false
		}
		filters.GroupID = &id
	}
	if startStr := c.Query("start_time"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_time")
			return filters, false
		}
		filters.StartTime = &t
	}
	if endStr := c.Query("end_time"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_time")
			return filters, false
		}
		filters.EndTime = &t
	}
	return filters, true
}

// GetByID handles getting a subscription by ID
// GET /api/v1/admin/subscriptions/:id
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid subscription ID")
		return
	}

	subscription, err := h.subscriptionService.GetByID(c.Request.Context(), subscriptionID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserSubscriptionFromServiceAdmin(subscription))
}

// GetProgress handles getting subscription usage progress
// GET /api/v1/admin/subscriptions/:id/progress
func (h *SubscriptionHandler) GetProgress(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid subscription ID")
		return
	}

	progress, err := h.subscriptionService.GetSubscriptionProgress(c.Request.Context(), subscriptionID)
	if err != nil {
		response.NotFound(c, "Subscription not found")
		return
	}

	response.Success(c, progress)
}

// Assign handles assigning a subscription to a user
// POST /api/v1/admin/subscriptions/assign
func (h *SubscriptionHandler) Assign(c *gin.Context) {
	var req AssignSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Get admin user ID from context
	adminID := getAdminIDFromContext(c)

	subscription, err := h.subscriptionService.AssignSubscription(c.Request.Context(), &service.AssignSubscriptionInput{
		UserID:       req.UserID,
		GroupID:      req.GroupID,
		ValidityDays: req.ValidityDays,
		AssignedBy:   adminID,
		Notes:        req.Notes,
		PriceUSD:     req.PriceUSD,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.UserSubscriptionFromServiceAdmin(subscription))
}

// BulkAssign handles bulk assigning subscriptions to multiple users
// POST /api/v1/admin/subscriptions/bulk-assign
func (h *SubscriptionHandler) BulkAssign(c *gin.Context) {
	var req BulkAssignSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Get admin user ID from context
	adminID := getAdminIDFromContext(c)

	result, err := h.subscriptionService.BulkAssignSubscription(c.Request.Context(), &service.BulkAssignSubscriptionInput{
		UserIDs:      req.UserIDs,
		GroupID:      req.GroupID,
		ValidityDays: req.ValidityDays,
		AssignedBy:   adminID,
		Notes:        req.Notes,
		PriceUSD:     req.PriceUSD,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.BulkAssignResultFromService(result))
}

// Extend handles adjusting a subscription (extend or shorten)
// POST /api/v1/admin/subscriptions/:id/extend
func (h *SubscriptionHandler) Extend(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid subscription ID")
		return
	}

	var req AdjustSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	adminID := getAdminIDFromContext(c)
	idempotencyPayload := struct {
		SubscriptionID int64                     `json:"subscription_id"`
		Body           AdjustSubscriptionRequest `json:"body"`
	}{
		SubscriptionID: subscriptionID,
		Body:           req,
	}
	executeAdminIdempotentJSON(c, "admin.subscriptions.extend", idempotencyPayload, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		subscription, execErr := h.subscriptionService.AdminAdjustSubscription(ctx, subscriptionID, req.Days, adminID)
		if execErr != nil {
			return nil, execErr
		}
		return dto.UserSubscriptionFromServiceAdmin(subscription), nil
	})
}

// ResetSubscriptionQuotaRequest represents the reset quota request
type ResetSubscriptionQuotaRequest struct {
	Daily   bool `json:"daily"`
	Weekly  bool `json:"weekly"`
	Monthly bool `json:"monthly"`
}

// ResetQuota resets daily, weekly, and/or monthly usage for a subscription.
// POST /api/v1/admin/subscriptions/:id/reset-quota
func (h *SubscriptionHandler) ResetQuota(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid subscription ID")
		return
	}
	var req ResetSubscriptionQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if !req.Daily && !req.Weekly && !req.Monthly {
		response.BadRequest(c, "At least one of 'daily', 'weekly', or 'monthly' must be true")
		return
	}
	adminID := getAdminIDFromContext(c)
	sub, err := h.subscriptionService.AdminResetQuotaWithRecord(c.Request.Context(), subscriptionID, req.Daily, req.Weekly, req.Monthly, adminID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.UserSubscriptionFromServiceAdmin(sub))
}

// Revoke handles revoking a subscription
// DELETE /api/v1/admin/subscriptions/:id
func (h *SubscriptionHandler) Revoke(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid subscription ID")
		return
	}

	adminID := getAdminIDFromContext(c)
	err = h.subscriptionService.AdminRevokeSubscription(c.Request.Context(), subscriptionID, adminID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Subscription revoked successfully"})
}

// ListByGroup handles listing subscriptions for a specific group
// GET /api/v1/admin/groups/:id/subscriptions
func (h *SubscriptionHandler) ListByGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid group ID")
		return
	}

	page, pageSize := response.ParsePagination(c)

	subscriptions, pagination, err := h.subscriptionService.ListGroupSubscriptions(c.Request.Context(), groupID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.AdminUserSubscription, 0, len(subscriptions))
	for i := range subscriptions {
		out = append(out, *dto.UserSubscriptionFromServiceAdmin(&subscriptions[i]))
	}
	response.PaginatedWithResult(c, out, toResponsePagination(pagination))
}

// ListByUser handles listing subscriptions for a specific user
// GET /api/v1/admin/users/:id/subscriptions
func (h *SubscriptionHandler) ListByUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	subscriptions, err := h.subscriptionService.ListUserSubscriptions(c.Request.Context(), userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]dto.AdminUserSubscription, 0, len(subscriptions))
	for i := range subscriptions {
		out = append(out, *dto.UserSubscriptionFromServiceAdmin(&subscriptions[i]))
	}
	response.Success(c, out)
}

// Helper function to get admin ID from context
func getAdminIDFromContext(c *gin.Context) int64 {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		return 0
	}
	return subject.UserID
}
