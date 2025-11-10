package stats

import (
	"backend/config"
	"backend/types"
	"database/sql"
	"fmt"
	"time"
)

type StatsRepoImpl struct {
	db *sql.DB
}

func NewStatsRepoImpl(db *sql.DB) *StatsRepoImpl {
	return &StatsRepoImpl{db: db}
}

func (s *StatsRepoImpl) GetUserStats(id int) (*types.ExtendedStats, error) {
	var stats types.Stats
	var extendedStats types.ExtendedStats
	err := s.GetUserStatsRow(id, &stats)
	if err != nil {
		return nil, err
	}
	// copy the stats fields into extendedStats
	extendedStats.UserID = id
	extendedStats.LongestStreak = stats.LongestStreak
	extendedStats.CurrentStreak = stats.CurrentStreak
	extendedStats.CreatedAt = stats.CreatedAt
	extendedStats.LastUpdated = stats.LastUpdated
	// here do all the aggregation needed to get
	extendedStats.XPMultiplier = config.BaseMultiplier + config.GrowthFactor*float64(stats.CurrentStreak-1)
	// TotalPomodoros
	if err := s.getUserTotalPomodoros(id, &extendedStats.TotalPomodoros); err != nil {
		return nil, err
	}
	// TotalFocusMinutes
	if err := s.getUserTotalFocusMinutes(id, &extendedStats.TotalFocusMinutes); err != nil {
		return nil, err
	}
	// BestDay
	if err := s.getUserBestDay(id, &extendedStats.BestDay); err != nil {
		return nil, err
	}
	return &extendedStats, nil
}
func (s *StatsRepoImpl) GetUserStatsRow(id int, stats *types.Stats) error {
	row := s.db.QueryRow(
		"SELECT user_id, longest_streak, current_streak, last_updated, created_at FROM stats WHERE user_id = ?",
		id,
	)
	if err := scanRowIntoStats(row, stats); err != nil {
		return fmt.Errorf("failed to scan stats for user %d: %w", id, err)
	}
	return nil
}
// Best Day = the day where the user made the longest time
func (s *StatsRepoImpl) getUserBestDay(id int, bestDayMinutes *int) error {
	row := s.db.QueryRow(
		`Select Max(total) as m from (
			Select SUM(session_duration) as total from pomodoros 
			where user_id = ? AND type = 'pomodoro'
			group by DATE(start_time)
		) sub`,
		id,
	)
	return row.Scan(bestDayMinutes)

}
func (s *StatsRepoImpl) getUserTotalFocusMinutes(id int, total *int) error {
	row := s.db.QueryRow(
		"Select SUM(session_duration) from pomodoros where user_id = ? AND type = 'pomodoro'",
		id,
	)
	return row.Scan(total)
}
func (s *StatsRepoImpl) getUserTotalPomodoros(id int, count *int) error {
	row := s.db.QueryRow(
		"Select COUNT(*) from pomodoros where user_id = ? AND type = 'pomodoro'",
		id,
	)
	return row.Scan(count)
}
func scanRowIntoStats(row *sql.Row, stats *types.Stats) error {
	return row.Scan(
		&stats.UserID,
		&stats.LongestStreak,
		&stats.CurrentStreak,
		&stats.LastUpdated,
		&stats.CreatedAt,
	)
}
func (s *StatsRepoImpl) AddUserStats(payload *types.Stats) (*types.Stats, error) {
	// insert into stats table
	_, err := s.db.Exec("Insert into stats values (?, ?, ?, ?, ?)",
		payload.UserID,
		payload.LongestStreak,
		payload.CurrentStreak,
		payload.LastUpdated,
		payload.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
func (s *StatsRepoImpl) UpdateUserStats(payload *types.Stats) (*types.Stats, error) {
	_, err := s.db.Exec("Update stats set longest_streak = ?, current_streak = ?, last_updated = ? where user_id = ?",
		payload.LongestStreak,
		payload.CurrentStreak,
		payload.LastUpdated,
		payload.UserID,
	)
	if err != nil {
		return nil, err
	}
	err = s.GetUserStatsRow(payload.UserID, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
func (s *StatsRepoImpl) GetUserHeatmap(p *types.HeatMapPayload) ([]types.HeatMapEntry, error) {
	var list []types.HeatMapEntry
	rows, err := s.db.Query(`
		SELECT date, count 
		FROM heatmap
		WHERE user_id = ? AND date BETWEEN ? AND ?
		ORDER BY date ASC
		`, p.UserID, p.StartDate, p.EndDate,
	)
	if err != nil {
		return nil, err
	}
	heatmap := make(map[time.Time]int)
	for rows.Next() {
		var date time.Time
		var count int
		rows.Scan(&date, &count)
		heatmap[date] = count
	}
	// now fill the empty dates
	for d := p.StartDate; !d.After(p.EndDate); d = d.AddDate(0, 0, 1) {
		list = append(list, types.HeatMapEntry{
			Date:  d,
			Count: heatmap[d],
		})
	}
	return list, nil
}
func (s *StatsRepoImpl) UpsertUserHeatmapEntry(p *types.HeatMap) error {
	_, err := s.db.Exec(`
	INSERT INTO heatmap (user_id, date, count)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE count = VALUES(count);`,
		p.UserID, p.Date, p.Count,
	)
	return err
}
