package policy

import (
	"time"

	"final-whistle/backend/internal/config"
)

type MatchWindow string

const (
	MatchWindowInactive  MatchWindow = "inactive"
	MatchWindowFarMatch  MatchWindow = "far_match"
	MatchWindowPreMatch  MatchWindow = "pre_match"
	MatchWindowLive      MatchWindow = "live_window"
	MatchWindowPostMatch MatchWindow = "post_match"
)

type MatchSchedulePolicy struct {
	Lookback            time.Duration
	Lookahead           time.Duration
	FarMatchWindow      time.Duration
	PreMatchWindow      time.Duration
	LiveWindow          time.Duration
	PostMatchWindow     time.Duration
	FarMatchEvery       time.Duration
	PreMatchEvery       time.Duration
	LiveEvery           time.Duration
	PostMatchEvery      time.Duration
	RosterBeforeKickoff time.Duration
	RosterAfterKickoff  time.Duration
	RosterEvery         time.Duration
}

func New(cfg config.SyncConfig) MatchSchedulePolicy {
	return MatchSchedulePolicy{
		Lookback:            time.Duration(cfg.MatchLookbackHours) * time.Hour,
		Lookahead:           time.Duration(cfg.MatchLookaheadDays) * 24 * time.Hour,
		FarMatchWindow:      time.Duration(cfg.WindowFarMatchDays) * 24 * time.Hour,
		PreMatchWindow:      time.Duration(cfg.WindowPreMatchMinutes) * time.Minute,
		LiveWindow:          time.Duration(cfg.WindowLiveAfterKickoffMinutes) * time.Minute,
		PostMatchWindow:     time.Duration(cfg.WindowPostMatchMinutes) * time.Minute,
		FarMatchEvery:       time.Duration(cfg.ScheduleFarMatchEveryMinutes) * time.Minute,
		PreMatchEvery:       time.Duration(cfg.SchedulePreMatchEveryMinutes) * time.Minute,
		LiveEvery:           time.Duration(cfg.ScheduleLiveEveryMinutes) * time.Minute,
		PostMatchEvery:      time.Duration(cfg.SchedulePostMatchEveryMinutes) * time.Minute,
		RosterBeforeKickoff: time.Duration(cfg.RosterWindowBeforeKickoffMinutes) * time.Minute,
		RosterAfterKickoff:  time.Duration(cfg.RosterWindowAfterKickoffMinutes) * time.Minute,
		RosterEvery:         time.Duration(cfg.RosterScheduleEveryMinutes) * time.Minute,
	}
}

func (p MatchSchedulePolicy) ClassifyMatchWindow(now, kickoff time.Time) MatchWindow {
	untilKickoff := kickoff.Sub(now)
	if untilKickoff > p.FarMatchWindow {
		return MatchWindowInactive
	}
	if untilKickoff > p.PreMatchWindow {
		return MatchWindowFarMatch
	}
	if now.Before(kickoff) {
		return MatchWindowPreMatch
	}
	if now.Before(kickoff.Add(p.LiveWindow)) {
		return MatchWindowLive
	}
	if now.Before(kickoff.Add(p.LiveWindow).Add(p.PostMatchWindow)) {
		return MatchWindowPostMatch
	}
	return MatchWindowInactive
}

func (p MatchSchedulePolicy) IntervalForWindow(window MatchWindow) time.Duration {
	switch window {
	case MatchWindowFarMatch:
		return p.FarMatchEvery
	case MatchWindowPreMatch:
		return p.PreMatchEvery
	case MatchWindowLive:
		return p.LiveEvery
	case MatchWindowPostMatch:
		return p.PostMatchEvery
	default:
		return 0
	}
}

func (p MatchSchedulePolicy) InRosterWindow(now, kickoff time.Time) bool {
	return !now.Before(kickoff.Add(-p.RosterBeforeKickoff)) && !now.After(kickoff.Add(p.RosterAfterKickoff))
}
