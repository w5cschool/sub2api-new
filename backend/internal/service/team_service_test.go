package service

import (
	"context"
	"sort"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type teamUserRepoStub struct{}

func (teamUserRepoStub) Create(context.Context, *User) error { return nil }
func (teamUserRepoStub) GetByID(_ context.Context, id int64) (*User, error) {
	return &User{ID: id, Email: "user@example.com"}, nil
}
func (teamUserRepoStub) GetByIDIncludeDeleted(_ context.Context, id int64) (*User, error) {
	return &User{ID: id, Email: "user@example.com"}, nil
}
func (teamUserRepoStub) GetByEmail(context.Context, string) (*User, error) { return &User{}, nil }
func (teamUserRepoStub) GetFirstAdmin(context.Context) (*User, error)      { return &User{}, nil }
func (teamUserRepoStub) Update(context.Context, *User) error               { return nil }
func (teamUserRepoStub) Delete(context.Context, int64) error               { return nil }
func (teamUserRepoStub) GetUserAvatar(context.Context, int64) (*UserAvatar, error) {
	return nil, nil
}
func (teamUserRepoStub) UpsertUserAvatar(context.Context, int64, UpsertUserAvatarInput) (*UserAvatar, error) {
	return nil, nil
}
func (teamUserRepoStub) DeleteUserAvatar(context.Context, int64) error { return nil }
func (teamUserRepoStub) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (teamUserRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (teamUserRepoStub) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) {
	return map[int64]*time.Time{}, nil
}
func (teamUserRepoStub) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error) {
	return nil, nil
}
func (teamUserRepoStub) UpdateUserLastActiveAt(context.Context, int64, time.Time) error {
	return nil
}
func (teamUserRepoStub) UpdateBalance(context.Context, int64, float64) error { return nil }
func (teamUserRepoStub) DeductBalance(context.Context, int64, float64) error { return nil }
func (teamUserRepoStub) UpdateConcurrency(context.Context, int64, int) error { return nil }
func (teamUserRepoStub) BatchSetConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
}
func (teamUserRepoStub) BatchAddConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
}
func (teamUserRepoStub) BatchUpdateLimits(context.Context, []int64, *int, *int) (int, error) {
	return 0, nil
}
func (teamUserRepoStub) ExistsByEmail(context.Context, string) (bool, error) { return false, nil }
func (teamUserRepoStub) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	return 0, nil
}
func (teamUserRepoStub) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	return nil
}
func (teamUserRepoStub) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	return nil
}
func (teamUserRepoStub) ListUserAuthIdentities(context.Context, int64) ([]UserAuthIdentityRecord, error) {
	return nil, nil
}
func (teamUserRepoStub) UnbindUserAuthProvider(context.Context, int64, string) error {
	return nil
}
func (teamUserRepoStub) UpdateTotpSecret(context.Context, int64, *string) error { return nil }
func (teamUserRepoStub) EnableTotp(context.Context, int64) error                { return nil }
func (teamUserRepoStub) DisableTotp(context.Context, int64) error               { return nil }

type fakeTeamRepo struct {
	nextID      int64
	teams       map[int64]*Team
	members     map[int64][]TeamMember
	memberByUID map[int64]*TeamMember
}

func newFakeTeamRepo() *fakeTeamRepo {
	return &fakeTeamRepo{
		nextID:      1,
		teams:       map[int64]*Team{},
		members:     map[int64][]TeamMember{},
		memberByUID: map[int64]*TeamMember{},
	}
}

func (r *fakeTeamRepo) Create(_ context.Context, team *Team) error {
	clone := *team
	clone.ID = r.nextID
	r.nextID++
	now := time.Now()
	clone.CreatedAt = now
	clone.UpdatedAt = now
	clone.Members = nil
	r.teams[clone.ID] = &clone
	team.ID = clone.ID
	team.CreatedAt = clone.CreatedAt
	team.UpdatedAt = clone.UpdatedAt
	return nil
}

func (r *fakeTeamRepo) Update(_ context.Context, team *Team) error {
	existing, ok := r.teams[team.ID]
	if !ok {
		return ErrTeamNotFound
	}
	existing.Name = team.Name
	existing.Description = team.Description
	existing.Status = team.Status
	existing.UpdatedAt = time.Now()
	return nil
}

func (r *fakeTeamRepo) GetByID(_ context.Context, id int64) (*Team, error) {
	team, ok := r.teams[id]
	if !ok || team.DeletedAt != nil {
		return nil, ErrTeamNotFound
	}
	out := *team
	out.Members = r.copyMembers(id)
	return &out, nil
}

func (r *fakeTeamRepo) List(_ context.Context, params pagination.PaginationParams, _ TeamFilters) ([]Team, *pagination.PaginationResult, error) {
	out := make([]Team, 0, len(r.teams))
	for _, team := range r.teams {
		if team.DeletedAt == nil {
			clone := *team
			clone.Members = r.copyMembers(team.ID)
			out = append(out, clone)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, &pagination.PaginationResult{Total: int64(len(out)), Page: params.Page, PageSize: params.PageSize, Pages: 1}, nil
}

func (r *fakeTeamRepo) SoftDelete(_ context.Context, id int64) error {
	team, ok := r.teams[id]
	if !ok {
		return ErrTeamNotFound
	}
	now := time.Now()
	team.DeletedAt = &now
	for _, member := range r.members[id] {
		delete(r.memberByUID, member.UserID)
	}
	delete(r.members, id)
	return nil
}

func (r *fakeTeamRepo) ReplaceMembers(_ context.Context, teamID int64, inputs []TeamMemberInput) error {
	if _, ok := r.teams[teamID]; !ok {
		return ErrTeamNotFound
	}
	for _, member := range r.members[teamID] {
		delete(r.memberByUID, member.UserID)
	}
	members := make([]TeamMember, 0, len(inputs))
	for i, input := range inputs {
		member := TeamMember{
			ID:        int64(i + 1),
			TeamID:    teamID,
			UserID:    input.UserID,
			Role:      normalizeTeamRole(input.Role),
			JoinedAt:  time.Now(),
			UpdatedAt: time.Now(),
			User:      &User{ID: input.UserID, Email: "user@example.com"},
		}
		members = append(members, member)
		cloned := member
		r.memberByUID[input.UserID] = &cloned
	}
	r.members[teamID] = members
	return nil
}

func (r *fakeTeamRepo) ListMembers(_ context.Context, teamID int64) ([]TeamMember, error) {
	return r.copyMembers(teamID), nil
}

func (r *fakeTeamRepo) GetMembershipByUserID(_ context.Context, userID int64) (*TeamMember, error) {
	member, ok := r.memberByUID[userID]
	if !ok {
		return nil, ErrTeamMembershipNotFound
	}
	out := *member
	if team, ok := r.teams[member.TeamID]; ok {
		teamCopy := *team
		out.Team = &teamCopy
	}
	return &out, nil
}

func (r *fakeTeamRepo) copyMembers(teamID int64) []TeamMember {
	source := r.members[teamID]
	out := make([]TeamMember, len(source))
	copy(out, source)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Role == out[j].Role {
			return out[i].UserID < out[j].UserID
		}
		return out[i].Role > out[j].Role
	})
	return out
}

func TestTeamService_CreateNormalizesMembersAndAllowsMultipleLeaders(t *testing.T) {
	repo := newFakeTeamRepo()
	svc := NewTeamService(repo, teamUserRepoStub{}, nil)

	team, err := svc.Create(context.Background(), TeamInput{
		Name:   " Engineering ",
		Status: "",
		Members: []TeamMemberInput{
			{UserID: 1, Role: "leader"},
			{UserID: 2, Role: "leader"},
			{UserID: 3, Role: "invalid"},
			{UserID: 3, Role: "leader"},
		},
	})

	require.NoError(t, err)
	require.Equal(t, "Engineering", team.Name)
	require.Equal(t, StatusActive, team.Status)
	require.Len(t, team.Members, 3)
	require.Equal(t, TeamRoleLeader, repo.memberByUID[1].Role)
	require.Equal(t, TeamRoleLeader, repo.memberByUID[2].Role)
	require.Equal(t, TeamRoleMember, repo.memberByUID[3].Role)
}

func TestTeamService_RejectsUserAlreadyAssignedToDifferentTeam(t *testing.T) {
	repo := newFakeTeamRepo()
	repo.teams[10] = &Team{ID: 10, Name: "Other", Status: StatusActive}
	repo.memberByUID[42] = &TeamMember{TeamID: 10, UserID: 42, Role: TeamRoleMember}
	svc := NewTeamService(repo, teamUserRepoStub{}, nil)

	_, err := svc.Create(context.Background(), TeamInput{
		Name:    "New Team",
		Members: []TeamMemberInput{{UserID: 42, Role: TeamRoleLeader}},
	})

	require.Error(t, err)
	require.True(t, infraerrors.IsConflict(err))
}

func TestTeamService_DeleteReleasesMemberships(t *testing.T) {
	repo := newFakeTeamRepo()
	svc := NewTeamService(repo, teamUserRepoStub{}, nil)
	team, err := svc.Create(context.Background(), TeamInput{
		Name:    "Support",
		Members: []TeamMemberInput{{UserID: 7, Role: TeamRoleLeader}},
	})
	require.NoError(t, err)

	require.NoError(t, svc.Delete(context.Background(), team.ID))

	_, err = svc.GetMyTeam(context.Background(), 7)
	require.Error(t, err)
	require.True(t, infraerrors.IsNotFound(err))
}

func TestTeamService_LeaderVisibleUserIDs(t *testing.T) {
	repo := newFakeTeamRepo()
	svc := NewTeamService(repo, teamUserRepoStub{}, nil)
	_, err := svc.Create(context.Background(), TeamInput{
		Name: "Sales",
		Members: []TeamMemberInput{
			{UserID: 1, Role: TeamRoleLeader},
			{UserID: 2, Role: TeamRoleMember},
		},
	})
	require.NoError(t, err)

	ids, err := svc.LeaderVisibleUserIDs(context.Background(), 1, 0)
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{1, 2}, ids)

	ids, err = svc.LeaderVisibleUserIDs(context.Background(), 1, 2)
	require.NoError(t, err)
	require.Equal(t, []int64{2}, ids)

	_, err = svc.LeaderVisibleUserIDs(context.Background(), 1, 999)
	require.Error(t, err)
	require.True(t, infraerrors.IsNotFound(err))

	_, err = svc.LeaderVisibleUserIDs(context.Background(), 2, 0)
	require.Error(t, err)
	require.True(t, infraerrors.IsForbidden(err))
}
