package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type Team struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
	Members     []TeamMember `json:"members,omitempty"`
}

type TeamMember struct {
	ID        int64     `json:"id"`
	TeamID    int64     `json:"team_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

type TeamMe struct {
	Team         *Team       `json:"team"`
	Membership   *TeamMember `json:"membership"`
	CanViewUsage bool        `json:"can_view_usage"`
}

type TeamMemberUsageSummary struct {
	UserID          int64   `json:"user_id"`
	Email           string  `json:"email"`
	Username        string  `json:"username"`
	Role            string  `json:"role"`
	TodayActualCost float64 `json:"today_actual_cost"`
	TotalActualCost float64 `json:"total_actual_cost"`
	TodayRequests   int64   `json:"today_requests"`
	TotalRequests   int64   `json:"total_requests"`
	TodayTokens     int64   `json:"today_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
}

func TeamFromService(team *service.Team) *Team {
	if team == nil {
		return nil
	}
	out := &Team{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Status:      team.Status,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
		DeletedAt:   team.DeletedAt,
	}
	if len(team.Members) > 0 {
		out.Members = TeamMembersFromService(team.Members)
	}
	return out
}

func TeamsFromService(teams []service.Team) []Team {
	out := make([]Team, 0, len(teams))
	for i := range teams {
		if team := TeamFromService(&teams[i]); team != nil {
			out = append(out, *team)
		}
	}
	return out
}

func TeamMemberFromService(member *service.TeamMember) *TeamMember {
	if member == nil {
		return nil
	}
	out := &TeamMember{
		ID:        member.ID,
		TeamID:    member.TeamID,
		UserID:    member.UserID,
		Role:      member.Role,
		JoinedAt:  member.JoinedAt,
		UpdatedAt: member.UpdatedAt,
		User:      UserFromServiceShallow(member.User),
	}
	return out
}

func TeamMembersFromService(members []service.TeamMember) []TeamMember {
	out := make([]TeamMember, 0, len(members))
	for i := range members {
		if member := TeamMemberFromService(&members[i]); member != nil {
			out = append(out, *member)
		}
	}
	return out
}

func TeamMemberUsageSummaryFromService(items []service.TeamMemberUsageSummary) []TeamMemberUsageSummary {
	out := make([]TeamMemberUsageSummary, 0, len(items))
	for _, item := range items {
		out = append(out, TeamMemberUsageSummary{
			UserID:          item.UserID,
			Email:           item.Email,
			Username:        item.Username,
			Role:            item.Role,
			TodayActualCost: item.TodayActualCost,
			TotalActualCost: item.TotalActualCost,
			TodayRequests:   item.TodayRequests,
			TotalRequests:   item.TotalRequests,
			TodayTokens:     item.TodayTokens,
			TotalTokens:     item.TotalTokens,
		})
	}
	return out
}
