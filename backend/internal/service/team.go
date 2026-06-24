package service

import (
	"context"
	"sort"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	TeamRoleMember = "member"
	TeamRoleLeader = "leader"
)

var (
	ErrTeamNotFound              = infraerrors.NotFound("TEAM_NOT_FOUND", "team not found")
	ErrTeamNameRequired          = infraerrors.BadRequest("TEAM_NAME_REQUIRED", "team name is required")
	ErrTeamNameExists            = infraerrors.Conflict("TEAM_NAME_EXISTS", "team name already exists")
	ErrTeamMemberInvalid         = infraerrors.BadRequest("TEAM_MEMBER_INVALID", "team member is invalid")
	ErrTeamMemberAlreadyAssigned = infraerrors.Conflict("TEAM_MEMBER_ALREADY_ASSIGNED", "user already belongs to another team")
	ErrTeamLeaderRequired        = infraerrors.Forbidden("TEAM_LEADER_REQUIRED", "team leader permission required")
	ErrTeamMembershipNotFound    = infraerrors.NotFound("TEAM_MEMBERSHIP_NOT_FOUND", "team membership not found")
)

type Team struct {
	ID          int64
	Name        string
	Description *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Members     []TeamMember
}

type TeamMember struct {
	ID        int64
	TeamID    int64
	UserID    int64
	Role      string
	JoinedAt  time.Time
	UpdatedAt time.Time
	Team      *Team
	User      *User
}

type TeamFilters struct {
	Search string
	Status string
}

type TeamInput struct {
	Name        string
	Description *string
	Status      string
	Members     []TeamMemberInput
}

type TeamMemberInput struct {
	UserID int64
	Role   string
}

type TeamMemberUsageSummary struct {
	UserID          int64
	Email           string
	Username        string
	Role            string
	TodayActualCost float64
	TotalActualCost float64
	TodayRequests   int64
	TotalRequests   int64
	TodayTokens     int64
	TotalTokens     int64
}

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	Update(ctx context.Context, team *Team) error
	GetByID(ctx context.Context, id int64) (*Team, error)
	List(ctx context.Context, params pagination.PaginationParams, filters TeamFilters) ([]Team, *pagination.PaginationResult, error)
	SoftDelete(ctx context.Context, id int64) error
	ReplaceMembers(ctx context.Context, teamID int64, members []TeamMemberInput) error
	ListMembers(ctx context.Context, teamID int64) ([]TeamMember, error)
	GetMembershipByUserID(ctx context.Context, userID int64) (*TeamMember, error)
}

type TeamService struct {
	teamRepo  TeamRepository
	userRepo  UserRepository
	usageRepo UsageLogRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository, usageRepo UsageLogRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo, usageRepo: usageRepo}
}

func (s *TeamService) List(ctx context.Context, page, pageSize int, filters TeamFilters, sortBy, sortOrder string) ([]Team, *pagination.PaginationResult, error) {
	params := pagination.PaginationParams{Page: page, PageSize: pageSize, SortBy: sortBy, SortOrder: sortOrder}
	return s.teamRepo.List(ctx, params, filters)
}

func (s *TeamService) Create(ctx context.Context, input TeamInput) (*Team, error) {
	team := &Team{}
	if err := applyTeamFields(team, input); err != nil {
		return nil, err
	}
	members, err := s.normalizeMembers(ctx, input.Members, 0)
	if err != nil {
		return nil, err
	}
	team.Members = inputsToTeamMembers(members)
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	if len(team.Members) > 0 {
		if err := s.teamRepo.ReplaceMembers(ctx, team.ID, teamMembersToInputs(team.Members)); err != nil {
			return nil, err
		}
	}
	return s.teamRepo.GetByID(ctx, team.ID)
}

func (s *TeamService) GetByID(ctx context.Context, id int64) (*Team, error) {
	if id <= 0 {
		return nil, ErrTeamNotFound
	}
	return s.teamRepo.GetByID(ctx, id)
}

func (s *TeamService) Update(ctx context.Context, id int64, input TeamInput) (*Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := applyTeamFields(team, input); err != nil {
		return nil, err
	}
	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}
	return s.teamRepo.GetByID(ctx, id)
}

func (s *TeamService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrTeamNotFound
	}
	return s.teamRepo.SoftDelete(ctx, id)
}

func (s *TeamService) ReplaceMembers(ctx context.Context, teamID int64, members []TeamMemberInput) (*Team, error) {
	if teamID <= 0 {
		return nil, ErrTeamNotFound
	}
	if _, err := s.teamRepo.GetByID(ctx, teamID); err != nil {
		return nil, err
	}
	normalized, err := s.normalizeMembers(ctx, members, teamID)
	if err != nil {
		return nil, err
	}
	if err := s.teamRepo.ReplaceMembers(ctx, teamID, normalized); err != nil {
		return nil, err
	}
	return s.teamRepo.GetByID(ctx, teamID)
}

func (s *TeamService) GetMyTeam(ctx context.Context, userID int64) (*TeamMember, error) {
	if userID <= 0 {
		return nil, ErrTeamMembershipNotFound
	}
	return s.teamRepo.GetMembershipByUserID(ctx, userID)
}

func (s *TeamService) ListLeaderMembers(ctx context.Context, leaderUserID int64) ([]TeamMember, error) {
	membership, err := s.requireLeader(ctx, leaderUserID)
	if err != nil {
		return nil, err
	}
	return s.teamRepo.ListMembers(ctx, membership.TeamID)
}

func (s *TeamService) LeaderVisibleUserIDs(ctx context.Context, leaderUserID int64, memberID int64) ([]int64, error) {
	members, err := s.ListLeaderMembers(ctx, leaderUserID)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(members))
	for _, member := range members {
		if memberID > 0 {
			if member.UserID == memberID {
				return []int64{member.UserID}, nil
			}
			continue
		}
		ids = append(ids, member.UserID)
	}
	if memberID > 0 {
		return nil, ErrTeamMembershipNotFound
	}
	return ids, nil
}

func (s *TeamService) GetLeaderMembersUsageSummary(ctx context.Context, leaderUserID int64) ([]TeamMemberUsageSummary, error) {
	members, err := s.ListLeaderMembers(ctx, leaderUserID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	if len(userIDs) == 0 {
		return []TeamMemberUsageSummary{}, nil
	}
	stats, err := s.usageRepo.GetBatchUserUsageStats(ctx, userIDs, time.Unix(0, 0).UTC(), time.Now())
	if err != nil {
		return nil, err
	}
	out := make([]TeamMemberUsageSummary, 0, len(members))
	for _, member := range members {
		item := TeamMemberUsageSummary{
			UserID: member.UserID,
			Role:   member.Role,
		}
		if member.User != nil {
			item.Email = member.User.Email
			item.Username = member.User.Username
		}
		if s := stats[member.UserID]; s != nil {
			item.TodayActualCost = s.TodayActualCost
			item.TotalActualCost = s.TotalActualCost
			item.TodayRequests = s.TodayRequests
			item.TotalRequests = s.TotalRequests
			item.TodayTokens = s.TodayTokens
			item.TotalTokens = s.TotalTokens
		}
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].TotalActualCost == out[j].TotalActualCost {
			return out[i].UserID < out[j].UserID
		}
		return out[i].TotalActualCost > out[j].TotalActualCost
	})
	return out, nil
}

func applyTeamFields(team *Team, input TeamInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return ErrTeamNameRequired
	}
	team.Name = name
	team.Description = normalizeOptionalString(input.Description)
	team.Status = strings.TrimSpace(input.Status)
	if team.Status == "" {
		team.Status = StatusActive
	}
	return nil
}

func (s *TeamService) normalizeMembers(ctx context.Context, members []TeamMemberInput, currentTeamID int64) ([]TeamMemberInput, error) {
	if len(members) == 0 {
		return nil, nil
	}
	seen := make(map[int64]struct{}, len(members))
	out := make([]TeamMemberInput, 0, len(members))
	for _, member := range members {
		if member.UserID <= 0 {
			return nil, ErrTeamMemberInvalid
		}
		if _, ok := seen[member.UserID]; ok {
			continue
		}
		seen[member.UserID] = struct{}{}
		role := normalizeTeamRole(member.Role)
		if _, err := s.userRepo.GetByID(ctx, member.UserID); err != nil {
			return nil, err
		}
		existing, err := s.teamRepo.GetMembershipByUserID(ctx, member.UserID)
		if err != nil && !infraerrors.IsNotFound(err) {
			return nil, err
		}
		if existing != nil && existing.TeamID != currentTeamID {
			return nil, ErrTeamMemberAlreadyAssigned
		}
		out = append(out, TeamMemberInput{UserID: member.UserID, Role: role})
	}
	return out, nil
}

func (s *TeamService) requireLeader(ctx context.Context, userID int64) (*TeamMember, error) {
	membership, err := s.teamRepo.GetMembershipByUserID(ctx, userID)
	if err != nil {
		return nil, ErrTeamLeaderRequired
	}
	if membership.Role != TeamRoleLeader {
		return nil, ErrTeamLeaderRequired
	}
	return membership, nil
}

func normalizeTeamRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case TeamRoleLeader:
		return TeamRoleLeader
	default:
		return TeamRoleMember
	}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func inputsToTeamMembers(inputs []TeamMemberInput) []TeamMember {
	out := make([]TeamMember, 0, len(inputs))
	for _, input := range inputs {
		out = append(out, TeamMember{UserID: input.UserID, Role: normalizeTeamRole(input.Role)})
	}
	return out
}

func teamMembersToInputs(members []TeamMember) []TeamMemberInput {
	out := make([]TeamMemberInput, 0, len(members))
	for _, member := range members {
		out = append(out, TeamMemberInput{UserID: member.UserID, Role: normalizeTeamRole(member.Role)})
	}
	return out
}
