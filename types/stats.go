package types

import "time"

type StatsRepo interface {
	AddUserStats(*Stats) (*Stats, error)
	UpdateUserStats(*Stats) (*Stats, error)
	GetUserStats(int) (*ExtendedStats, error)
	// Heatmap operations
	GetUserHeatmap(*HeatMapPayload) ([]HeatMapEntry, error)
	UpsertUserHeatmapEntry(*HeatMap) error
}

type Stats struct {
	UserID        int       `json:"user_id"`
	LongestStreak int       `json:"longest_streak"`
	CurrentStreak int       `json:"current_streak"`
	LastUpdated   time.Time `json:"last_updated"`
	CreatedAt     time.Time `json:"created_at"`
}
type HeatMap struct {
	UserID int       `json:"user_id"`
	Count  int       `json:"count"`
	Date   time.Time `json:"date"`
}
type HeatMapEntry struct {
	Count int       `json:"count"`
	Date  time.Time `json:"date"`
}
type HeatMapPayload struct {
	UserID    int       `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type HeatMapResponse struct {
	UserID int            `json:"user_id"`
	Data   []HeatMapEntry `json:"data"`
}

type StatsPayload struct {
	UserID        int `json:"user_id"`
	LongestStreak int `json:"longest_streak"`
	CurrentStreak int `json:"current_streak"`
}

type ExtendedStats struct {
	UserID            int       `json:"user_id"`
	LongestStreak     int       `json:"longest_streak"`
	CurrentStreak     int       `json:"current_streak"`
	XPMultiplier      float64   `json:"xp_multiplier"`
	BestDay           int       `json:"best_day"`
	TotalPomodoros    int       `json:"total_pomodoros"`
	TotalFocusMinutes int       `json:"total_focus_minutes"`
	LastUpdated       time.Time `json:"last_updated"`
	CreatedAt         time.Time `json:"created_at"`
}
