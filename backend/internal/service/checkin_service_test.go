package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type fakeCheckInStore struct {
	matches                map[uint]*model.Match
	checkIns               map[string]*model.CheckIn
	players                map[uint]model.Player
	tags                   map[uint]model.Tag
	eligiblePlayersByMatch map[uint]map[uint]struct{}
	nextCheckInID          uint
	nextPlayerRatingID     uint
	failCreateCheckIn      error
	failUpdateCheckIn      error
	failReplaceRatings     error
	failReplaceTags        error
}

type fakeCheckInRepository struct {
	store *fakeCheckInStore
}

func (f *fakeCheckInRepository) FindMatchByID(id uint) (*model.Match, error) {
	match, ok := f.store.matches[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *match
	return &copy, nil
}

func (f *fakeCheckInRepository) FindCheckInByUserAndMatch(userID, matchID uint) (*model.CheckIn, error) {
	checkIn, ok := f.store.checkIns[checkInKey(userID, matchID)]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return cloneCheckIn(checkIn), nil
}

func (f *fakeCheckInRepository) GetEligiblePlayerIDs(matchID uint, playerIDs []uint) (map[uint]struct{}, error) {
	result := map[uint]struct{}{}
	eligible := f.store.eligiblePlayersByMatch[matchID]
	for _, id := range playerIDs {
		if _, ok := eligible[id]; ok {
			result[id] = struct{}{}
		}
	}
	return result, nil
}

func (f *fakeCheckInRepository) GetActiveTagsByIDs(tagIDs []uint) ([]model.Tag, error) {
	tags := make([]model.Tag, 0, len(tagIDs))
	for _, id := range tagIDs {
		tag, ok := f.store.tags[id]
		if ok && tag.IsActive {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func (f *fakeCheckInRepository) WithTransaction(fn func(repo repository.CheckInRepository) error) error {
	cloned := cloneStore(f.store)
	txRepo := &fakeCheckInRepository{store: cloned}
	if err := fn(txRepo); err != nil {
		return err
	}
	f.store = cloned
	return nil
}

func (f *fakeCheckInRepository) CreateCheckIn(checkIn *model.CheckIn) error {
	if f.store.failCreateCheckIn != nil {
		return f.store.failCreateCheckIn
	}
	f.store.nextCheckInID++
	created := *checkIn
	created.ID = f.store.nextCheckInID
	created.CreatedAt = time.Now().UTC()
	created.UpdatedAt = created.CreatedAt
	f.store.checkIns[checkInKey(checkIn.UserID, checkIn.MatchID)] = &created
	checkIn.ID = created.ID
	checkIn.CreatedAt = created.CreatedAt
	checkIn.UpdatedAt = created.UpdatedAt
	return nil
}

func (f *fakeCheckInRepository) UpdateCheckIn(checkIn *model.CheckIn) error {
	if f.store.failUpdateCheckIn != nil {
		return f.store.failUpdateCheckIn
	}
	existing, ok := f.store.checkIns[checkInKey(checkIn.UserID, checkIn.MatchID)]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	updated := *checkIn
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = time.Now().UTC()
	updated.Tags = cloneTags(existing.Tags)
	updated.PlayerRatings = clonePlayerRatings(existing.PlayerRatings)
	f.store.checkIns[checkInKey(checkIn.UserID, checkIn.MatchID)] = &updated
	checkIn.UpdatedAt = updated.UpdatedAt
	checkIn.CreatedAt = updated.CreatedAt
	return nil
}

func (f *fakeCheckInRepository) ReplacePlayerRatings(checkInID uint, ratings []model.PlayerRating) error {
	if f.store.failReplaceRatings != nil {
		return f.store.failReplaceRatings
	}
	checkIn, ok := f.findCheckInByID(checkInID)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	newRatings := make([]model.PlayerRating, 0, len(ratings))
	for _, rating := range ratings {
		f.store.nextPlayerRatingID++
		copied := rating
		copied.ID = f.store.nextPlayerRatingID
		copied.Player = f.store.players[rating.PlayerID]
		newRatings = append(newRatings, copied)
	}
	checkIn.PlayerRatings = newRatings
	return nil
}

func (f *fakeCheckInRepository) ReplaceCheckInTags(checkInID uint, tagIDs []uint) error {
	if f.store.failReplaceTags != nil {
		return f.store.failReplaceTags
	}
	checkIn, ok := f.findCheckInByID(checkInID)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	newTags := make([]model.Tag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		newTags = append(newTags, f.store.tags[tagID])
	}
	checkIn.Tags = newTags
	return nil
}

func (f *fakeCheckInRepository) findCheckInByID(id uint) (*model.CheckIn, bool) {
	for _, checkIn := range f.store.checkIns {
		if checkIn.ID == id {
			return checkIn, true
		}
	}
	return nil, false
}

func TestCheckInServiceGetMyCheckInReturnsNilWhenMissing(t *testing.T) {
	svc := NewCheckInService(newFakeCheckInRepository())

	result, err := svc.GetMyCheckIn(1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}
}

func TestCheckInServiceCreateCheckInSuccess(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)

	result, err := svc.CreateCheckIn(1, 10, validCheckInRequest())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.ID == 0 {
		t.Fatalf("expected created check-in, got %#v", result)
	}
	if len(result.Tags) != 2 || len(result.PlayerRatings) != 2 {
		t.Fatalf("expected tags and player ratings, got %#v", result)
	}
}

func TestCheckInServiceCreateRejectsDuplicateCreate(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)
	req := validCheckInRequest()

	if _, err := svc.CreateCheckIn(1, 10, req); err != nil {
		t.Fatalf("expected initial create to succeed, got %v", err)
	}

	_, err := svc.CreateCheckIn(1, 10, req)
	if !errors.Is(err, ErrCheckInAlreadyExists) {
		t.Fatalf("expected ErrCheckInAlreadyExists, got %v", err)
	}
}

func TestCheckInServiceCreateRejectsNonFinishedMatch(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)

	_, err := svc.CreateCheckIn(2, 10, validCheckInRequest())
	assertValidationError(t, err)
}

func TestCheckInServiceCreateRejectsInvalidPayloadCases(t *testing.T) {
	cases := []struct {
		name string
		req  dto.UpsertCheckInRequestDTO
	}{
		{
			name: "invalid score",
			req: func() dto.UpsertCheckInRequestDTO {
				req := validCheckInRequest()
				req.MatchRating = 11
				return req
			}(),
		},
		{
			name: "short review too long",
			req: func() dto.UpsertCheckInRequestDTO {
				req := validCheckInRequest()
				long := makeString(281)
				req.ShortReview = &long
				return req
			}(),
		},
		{
			name: "invalid tags",
			req: func() dto.UpsertCheckInRequestDTO {
				req := validCheckInRequest()
				req.Tags = []uint{1, 99}
				return req
			}(),
		},
		{
			name: "duplicate player entries",
			req: func() dto.UpsertCheckInRequestDTO {
				req := validCheckInRequest()
				req.PlayerRatings = []dto.PlayerRatingInputDTO{
					{PlayerID: 101, Rating: 7},
					{PlayerID: 101, Rating: 8},
				}
				return req
			}(),
		},
		{
			name: "player not in match",
			req: func() dto.UpsertCheckInRequestDTO {
				req := validCheckInRequest()
				req.PlayerRatings = []dto.PlayerRatingInputDTO{{PlayerID: 999, Rating: 8}}
				return req
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeCheckInRepository()
			svc := NewCheckInService(repo)

			_, err := svc.CreateCheckIn(1, 10, tc.req)
			assertValidationError(t, err)
		})
	}
}

func TestCheckInServiceAllowsMoreThanFiveRosterPlayers(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)
	req := validCheckInRequest()
	req.PlayerRatings = []dto.PlayerRatingInputDTO{
		{PlayerID: 101, Rating: 7},
		{PlayerID: 102, Rating: 8},
		{PlayerID: 103, Rating: 9},
		{PlayerID: 104, Rating: 6},
		{PlayerID: 105, Rating: 8},
		{PlayerID: 106, Rating: 7},
	}

	result, err := svc.CreateCheckIn(1, 10, req)
	if err != nil {
		t.Fatalf("expected create to allow full roster rating, got %v", err)
	}
	if len(result.PlayerRatings) != 6 {
		t.Fatalf("expected six player ratings, got %#v", result.PlayerRatings)
	}
}

func TestCheckInServiceUpdateReplacesChildren(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)
	req := validCheckInRequest()

	created, err := svc.CreateCheckIn(1, 10, req)
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	updatedReq := validCheckInRequest()
	updatedReq.Tags = []uint{2}
	updatedReq.PlayerRatings = []dto.PlayerRatingInputDTO{{PlayerID: 102, Rating: 9}}
	updatedReq.MatchRating = 9

	updated, err := svc.UpdateCheckIn(1, 10, updatedReq)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	if updated.ID != created.ID {
		t.Fatalf("expected same check-in id, got %d vs %d", updated.ID, created.ID)
	}
	if len(updated.Tags) != 1 || updated.Tags[0].ID != 2 {
		t.Fatalf("expected replaced tags, got %#v", updated.Tags)
	}
	if len(updated.PlayerRatings) != 1 || updated.PlayerRatings[0].Player.ID != 102 {
		t.Fatalf("expected replaced player ratings, got %#v", updated.PlayerRatings)
	}
}

func TestCheckInServiceUpdateRejectsMissingCheckIn(t *testing.T) {
	svc := NewCheckInService(newFakeCheckInRepository())

	_, err := svc.UpdateCheckIn(1, 10, validCheckInRequest())
	if !errors.Is(err, ErrCheckInMissing) {
		t.Fatalf("expected ErrCheckInMissing, got %v", err)
	}
}

func TestCheckInServiceCreateRollsBackOnChildWriteFailure(t *testing.T) {
	repo := newFakeCheckInRepository()
	repo.store.failReplaceTags = errors.New("boom")
	svc := NewCheckInService(repo)

	_, err := svc.CreateCheckIn(1, 10, validCheckInRequest())
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(repo.store.checkIns) != 0 {
		t.Fatalf("expected rollback, found %d check-ins", len(repo.store.checkIns))
	}
}

func TestCheckInServiceUpdateRollsBackOnChildWriteFailure(t *testing.T) {
	repo := newFakeCheckInRepository()
	svc := NewCheckInService(repo)

	created, err := svc.CreateCheckIn(1, 10, validCheckInRequest())
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	repo.store.failReplaceRatings = errors.New("boom")
	updateReq := validCheckInRequest()
	updateReq.MatchRating = 10
	updateReq.Tags = []uint{2}
	updateReq.PlayerRatings = []dto.PlayerRatingInputDTO{{PlayerID: 102, Rating: 9}}

	_, err = svc.UpdateCheckIn(1, 10, updateReq)
	if err == nil {
		t.Fatalf("expected update error")
	}

	stored := repo.store.checkIns[checkInKey(10, 1)]
	if stored.MatchRating != created.MatchRating || len(stored.Tags) != len(created.Tags) || len(stored.PlayerRatings) != len(created.PlayerRatings) {
		t.Fatalf("expected rollback to preserve previous state, got %#v", stored)
	}
}

func newFakeCheckInRepository() *fakeCheckInRepository {
	homeTeam := model.Team{ID: 1, Name: "Arsenal", Slug: "arsenal"}
	awayTeam := model.Team{ID: 2, Name: "Chelsea", Slug: "chelsea"}
	store := &fakeCheckInStore{
		matches: map[uint]*model.Match{
			1: {ID: 1, Status: model.MatchStatusFinished},
			2: {ID: 2, Status: model.MatchStatusScheduled},
		},
		checkIns: map[string]*model.CheckIn{},
		players: map[uint]model.Player{
			101: {ID: 101, Name: "Bukayo Saka", Slug: "bukayo-saka", TeamID: 1, Team: homeTeam},
			102: {ID: 102, Name: "Cole Palmer", Slug: "cole-palmer", TeamID: 2, Team: awayTeam},
			103: {ID: 103, Name: "Martin Odegaard", Slug: "martin-odegaard", TeamID: 1, Team: homeTeam},
			104: {ID: 104, Name: "Declan Rice", Slug: "declan-rice", TeamID: 1, Team: homeTeam},
			105: {ID: 105, Name: "Nicolas Jackson", Slug: "nicolas-jackson", TeamID: 2, Team: awayTeam},
			106: {ID: 106, Name: "Enzo Fernandez", Slug: "enzo-fernandez", TeamID: 2, Team: awayTeam},
		},
		tags: map[uint]model.Tag{
			1: {ID: 1, Name: "Thriller", Slug: "thriller", IsActive: true},
			2: {ID: 2, Name: "Tense", Slug: "tense", IsActive: true},
			3: {ID: 3, Name: "Inactive", Slug: "inactive", IsActive: false},
		},
		eligiblePlayersByMatch: map[uint]map[uint]struct{}{
			1: {101: {}, 102: {}, 103: {}, 104: {}, 105: {}, 106: {}},
		},
		nextCheckInID:      100,
		nextPlayerRatingID: 200,
	}
	return &fakeCheckInRepository{store: store}
}

func validCheckInRequest() dto.UpsertCheckInRequestDTO {
	review := "Great match"
	note1 := "Clutch"
	note2 := "Sharp"
	return dto.UpsertCheckInRequestDTO{
		WatchedType:    "FULL",
		SupporterSide:  "NEUTRAL",
		MatchRating:    8,
		HomeTeamRating: 7,
		AwayTeamRating: 8,
		ShortReview:    &review,
		WatchedAt:      time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC),
		Tags:           []uint{1, 2},
		PlayerRatings: []dto.PlayerRatingInputDTO{
			{PlayerID: 101, Rating: 8, Note: &note1},
			{PlayerID: 102, Rating: 9, Note: &note2},
		},
	}
}

func assertValidationError(t *testing.T, err error) {
	t.Helper()
	var validationErr *CheckInValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func makeString(length int) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = 'a'
	}
	return string(buf)
}

func checkInKey(userID, matchID uint) string {
	return fmt.Sprintf("%d:%d", userID, matchID)
}

func cloneStore(store *fakeCheckInStore) *fakeCheckInStore {
	cloned := &fakeCheckInStore{
		matches:                map[uint]*model.Match{},
		checkIns:               map[string]*model.CheckIn{},
		players:                map[uint]model.Player{},
		tags:                   map[uint]model.Tag{},
		eligiblePlayersByMatch: map[uint]map[uint]struct{}{},
		nextCheckInID:          store.nextCheckInID,
		nextPlayerRatingID:     store.nextPlayerRatingID,
		failCreateCheckIn:      store.failCreateCheckIn,
		failUpdateCheckIn:      store.failUpdateCheckIn,
		failReplaceRatings:     store.failReplaceRatings,
		failReplaceTags:        store.failReplaceTags,
	}

	for id, match := range store.matches {
		copy := *match
		cloned.matches[id] = &copy
	}
	for key, checkIn := range store.checkIns {
		cloned.checkIns[key] = cloneCheckIn(checkIn)
	}
	for id, player := range store.players {
		cloned.players[id] = player
	}
	for id, tag := range store.tags {
		cloned.tags[id] = tag
	}
	for matchID, players := range store.eligiblePlayersByMatch {
		copied := map[uint]struct{}{}
		for playerID := range players {
			copied[playerID] = struct{}{}
		}
		cloned.eligiblePlayersByMatch[matchID] = copied
	}
	return cloned
}

func cloneCheckIn(checkIn *model.CheckIn) *model.CheckIn {
	if checkIn == nil {
		return nil
	}
	copy := *checkIn
	copy.Tags = cloneTags(checkIn.Tags)
	copy.PlayerRatings = clonePlayerRatings(checkIn.PlayerRatings)
	return &copy
}

func cloneTags(tags []model.Tag) []model.Tag {
	copied := make([]model.Tag, len(tags))
	copy(copied, tags)
	return copied
}

func clonePlayerRatings(ratings []model.PlayerRating) []model.PlayerRating {
	copied := make([]model.PlayerRating, len(ratings))
	for i, rating := range ratings {
		copied[i] = rating
	}
	return copied
}
