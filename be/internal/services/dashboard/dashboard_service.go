// Package dashboard computes the generic admin overview KPIs.
package dashboard

import (
	"context"
	"math"
	"time"

	"be/common/utils"
	admindto "be/internal/dto/admin"
	"be/internal/repository/interfaces"
)

// LoginAction is the audit action counted as a sign-in.
const LoginAction = "auth.login"

type Service struct {
	stats interfaces.StatsRepository
	audit interfaces.AuditLogRepository
	now   func() time.Time
}

func NewService(stats interfaces.StatsRepository, audit interfaces.AuditLogRepository) *Service {
	return &Service{stats: stats, audit: audit, now: time.Now}
}

// Overview compares the last 7 days with the 7 days before:
//   - totalUsers: users now vs users 7 days ago
//   - activeSessions: live sessions now; change = sign-ins this week vs last week
//   - newUsers: sign-ups this week vs last week
//   - roles: count (no trend)
//
// weeklyActivity is sign-ins per UTC day for the last 7 days.
func (s *Service) Overview(ctx context.Context) (*admindto.DashboardOverview, error) {
	now := s.now().UTC()
	weekAgo := now.AddDate(0, 0, -7)
	twoWeeksAgo := now.AddDate(0, 0, -14)

	totalNow, err := s.stats.CountUsers(ctx, now.Add(time.Second))
	if err != nil {
		return nil, err
	}
	totalWeekAgo, err := s.stats.CountUsers(ctx, weekAgo)
	if err != nil {
		return nil, err
	}
	newThisWeek, err := s.stats.CountUsersCreated(ctx, weekAgo, now.Add(time.Second))
	if err != nil {
		return nil, err
	}
	newLastWeek, err := s.stats.CountUsersCreated(ctx, twoWeeksAgo, weekAgo)
	if err != nil {
		return nil, err
	}
	sessions, err := s.stats.CountActiveSessions(ctx, now)
	if err != nil {
		return nil, err
	}
	loginsThisWeek, err := s.audit.CountByAction(ctx, LoginAction, weekAgo, now.Add(time.Second))
	if err != nil {
		return nil, err
	}
	loginsLastWeek, err := s.audit.CountByAction(ctx, LoginAction, twoWeeksAgo, weekAgo)
	if err != nil {
		return nil, err
	}
	roles, err := s.stats.CountRoles(ctx)
	if err != nil {
		return nil, err
	}

	firstDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -6)
	perDay, err := s.audit.CountByActionPerDay(ctx, LoginAction, firstDay)
	if err != nil {
		return nil, err
	}
	activity := make([]admindto.DashboardPoint, 0, 7)
	for i := 0; i < 7; i++ {
		day := firstDay.AddDate(0, 0, i).Format("2006-01-02")
		activity = append(activity, admindto.DashboardPoint{Date: day, Value: perDay[day]})
	}

	return &admindto.DashboardOverview{
		KPIs: []admindto.DashboardKPI{
			{Key: "totalUsers", Value: totalNow, ChangePercent: change(totalNow, totalWeekAgo)},
			{Key: "activeSessions", Value: sessions, ChangePercent: change(loginsThisWeek, loginsLastWeek)},
			{Key: "newUsers", Value: newThisWeek, ChangePercent: change(newThisWeek, newLastWeek)},
			{Key: "roles", Value: roles, ChangePercent: 0},
		},
		WeeklyActivity: activity,
		UpdatedAt:      utils.FormatRFC3339(now),
	}, nil
}

// change is the percent change from previous to current, one decimal. From 0
// it reports 100 for any growth (there is no meaningful ratio).
func change(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return math.Round(float64(current-previous)/float64(previous)*1000) / 10
}

// SetClock replaces the time source (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }
