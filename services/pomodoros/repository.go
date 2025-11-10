package pomodoros

import (
	"backend/types"
	"database/sql"
	"fmt"
	"strings"
)

type PomodoroRepoImpl struct {
	db *sql.DB
}

func NewPomodoroRepoImpl(db *sql.DB) *PomodoroRepoImpl {
	return &PomodoroRepoImpl{db: db}
}

func (p *PomodoroRepoImpl) AddPomodoro(payload types.AddingPomodoroPayload) (*types.Pomodoro, error) {
	res, err := p.db.Exec(
		"insert into pomodoros(user_id,type,completed,session_duration,start_time,end_time) values (?,?,?,?,?,?)",
		payload.UserId,
		payload.Type,
		payload.Completed,
		payload.SessionDuration,
		payload.StartTime,
		payload.EndTime,
	)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return nil, fmt.Errorf("invalid user_id: user don't exist")
		}
		return nil, err
	}
	id, _ := res.LastInsertId()
	return p.getPomodoroById(id)
}

func (p *PomodoroRepoImpl) getPomodoroById(id int64) (*types.Pomodoro, error) {
	var pomodoro types.Pomodoro
	res := p.db.QueryRow("Select * from pomodoros where id = ?", id)

	if err := res.Err(); err != nil {
		return nil, err
	}
	err := scanRowIntoPomodoro(res, &pomodoro)
	if err != nil {
		return nil, err
	}
	return &pomodoro, nil
}
func scanRowIntoPomodoro(row *sql.Row, pomodoro *types.Pomodoro) error {
	err := row.Scan(
		&pomodoro.Id,
		&pomodoro.UserId,
		&pomodoro.Type,
		&pomodoro.Completed,
		&pomodoro.SessionDuration,
		&pomodoro.StartTime,
		&pomodoro.EndTime,
		&pomodoro.CreatedAt,
	)
	return err
}
