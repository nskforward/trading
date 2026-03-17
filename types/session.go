package types

type Session struct {
	Type  SessionType
	Start int64
	End   int64
}

type SessionType string

const (
	SessionClosed     SessionType = "closed"
	SessionPremarket  SessionType = "premarket"
	SessionMain       SessionType = "main"
	SessionPostmarket SessionType = "postmarket"
)

func (s Session) Contains(timestamp int64) bool {
	return timestamp >= s.Start && timestamp <= s.End
}
