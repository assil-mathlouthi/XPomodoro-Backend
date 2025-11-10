package types

import (
	"time"
)

type PomodoroRepo interface {
	AddPomodoro(AddingPomodoroPayload) (*Pomodoro, error)
}

type Pomodoro struct {
	Id              int       `json:"id"`
	UserId          int       `json:"user_id"`
	Type            string    `json:"type"`
	Completed       bool      `json:"completed"`
	SessionDuration int       `json:"session_duration"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	CreatedAt       time.Time `json:"created_at"`
}

type AddingPomodoroPayload struct {
	UserId          int       `json:"user_id"`
	Type            string    `json:"type"`
	Completed       bool      `json:"completed"`
	SessionDuration int       `json:"session_duration"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
}
