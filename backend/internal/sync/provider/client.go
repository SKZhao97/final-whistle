package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"final-whistle/backend/internal/config"
)

const footballDataBaseURL = "https://api.football-data.org/v4"

type FootballDataClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	limiter    *Limiter
}

func NewFootballDataClient(cfg config.SyncConfig) *FootballDataClient {
	return &FootballDataClient{
		baseURL: footballDataBaseURL,
		token:   cfg.FootballDataAPIToken,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		limiter: NewLimiter(cfg.SafeRateLimitPerMinute, time.Minute),
	}
}

func (c *FootballDataClient) ListCompetitionTeams(ctx context.Context, competitionCode string) ([]Team, error) {
	if competitionCode == "" {
		return nil, fmt.Errorf("competition code is required")
	}

	var response struct {
		Teams []struct {
			ID        int     `json:"id"`
			Name      string  `json:"name"`
			ShortName *string `json:"shortName"`
			Crest     *string `json:"crest"`
		} `json:"teams"`
	}
	if err := c.getJSON(ctx, fmt.Sprintf("/competitions/%s/teams", url.PathEscape(competitionCode)), nil, &response); err != nil {
		return nil, err
	}

	teams := make([]Team, 0, len(response.Teams))
	for _, item := range response.Teams {
		teams = append(teams, Team{
			ExternalID: fmt.Sprintf("%d", item.ID),
			Name:       item.Name,
			ShortName:  item.ShortName,
			CrestURL:   item.Crest,
		})
	}
	return teams, nil
}

func (c *FootballDataClient) ListTeamPlayers(ctx context.Context, teamExternalID string) ([]Player, error) {
	if teamExternalID == "" {
		return nil, fmt.Errorf("team external id is required")
	}

	var response struct {
		Squad []struct {
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Position *string `json:"position"`
		} `json:"squad"`
	}
	if err := c.getJSON(ctx, fmt.Sprintf("/teams/%s", url.PathEscape(teamExternalID)), nil, &response); err != nil {
		return nil, err
	}

	players := make([]Player, 0, len(response.Squad))
	for _, item := range response.Squad {
		players = append(players, Player{
			ExternalID:     fmt.Sprintf("%d", item.ID),
			TeamExternalID: teamExternalID,
			Name:           item.Name,
			Position:       item.Position,
		})
	}
	return players, nil
}

func (c *FootballDataClient) ListCompetitionMatches(ctx context.Context, competitionCode string, dateFrom, dateTo string) ([]Match, error) {
	query := url.Values{}
	if dateFrom != "" {
		query.Set("dateFrom", dateFrom)
	}
	if dateTo != "" {
		query.Set("dateTo", dateTo)
	}

	var response struct {
		Matches []struct {
			ID       int    `json:"id"`
			Status   string `json:"status"`
			UTCDate  string `json:"utcDate"`
			Matchday *int   `json:"matchday"`
			Season   struct {
				StartDate string `json:"startDate"`
			} `json:"season"`
			Competition struct {
				Name string `json:"name"`
			} `json:"competition"`
			HomeTeam struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				ShortName *string `json:"shortName"`
				Crest     *string `json:"crest"`
			} `json:"homeTeam"`
			AwayTeam struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				ShortName *string `json:"shortName"`
				Crest     *string `json:"crest"`
			} `json:"awayTeam"`
			Score struct {
				FullTime struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"fullTime"`
			} `json:"score"`
			Venue *string `json:"venue"`
		} `json:"matches"`
	}
	if err := c.getJSON(ctx, fmt.Sprintf("/competitions/%s/matches", url.PathEscape(competitionCode)), query, &response); err != nil {
		return nil, err
	}

	matches := make([]Match, 0, len(response.Matches))
	for _, item := range response.Matches {
		kickoffAt, err := time.Parse(time.RFC3339, item.UTCDate)
		if err != nil {
			return nil, fmt.Errorf("parse match kickoff time: %w", err)
		}
		season := ""
		if len(item.Season.StartDate) >= 4 {
			season = item.Season.StartDate[:4]
		}
		var round *string
		if item.Matchday != nil {
			value := fmt.Sprintf("Matchday %d", *item.Matchday)
			round = &value
		}
		matches = append(matches, Match{
			ExternalID:  fmt.Sprintf("%d", item.ID),
			Competition: item.Competition.Name,
			Season:      season,
			Round:       round,
			Status:      item.Status,
			KickoffAt:   kickoffAt,
			HomeTeam: MatchTeam{
				ExternalID: fmt.Sprintf("%d", item.HomeTeam.ID),
				Name:       item.HomeTeam.Name,
				ShortName:  item.HomeTeam.ShortName,
				CrestURL:   item.HomeTeam.Crest,
			},
			AwayTeam: MatchTeam{
				ExternalID: fmt.Sprintf("%d", item.AwayTeam.ID),
				Name:       item.AwayTeam.Name,
				ShortName:  item.AwayTeam.ShortName,
				CrestURL:   item.AwayTeam.Crest,
			},
			HomeScore: item.Score.FullTime.Home,
			AwayScore: item.Score.FullTime.Away,
			Venue:     item.Venue,
		})
	}
	return matches, nil
}

func (c *FootballDataClient) GetMatchRoster(ctx context.Context, externalMatchID string) ([]MatchRosterPlayer, error) {
	// football-data.org free tier lineup support is limited; keep this explicit for now.
	return nil, fmt.Errorf("match roster sync is not implemented for football-data provider yet")
}

func (c *FootballDataClient) getJSON(ctx context.Context, path string, query url.Values, target any) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("provider request failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}
	return nil
}

type Limiter struct {
	interval time.Duration
	mu       sync.Mutex
	last     time.Time
}

func NewLimiter(limit int, per time.Duration) *Limiter {
	if limit <= 0 {
		limit = 1
	}
	return &Limiter{interval: per / time.Duration(limit)}
}

func (l *Limiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if l.last.IsZero() {
		l.last = now
		return nil
	}

	next := l.last.Add(l.interval)
	if !next.After(now) {
		l.last = now
		return nil
	}

	timer := time.NewTimer(next.Sub(now))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		l.last = time.Now()
		return nil
	}
}
