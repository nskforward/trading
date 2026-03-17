package types

import "time"

type Schedule struct {
	Slots []Session
}

func (s Schedule) Current() (Session, bool) {
	now := time.Now().Unix()

	for _, slot := range s.Slots {
		if slot.Contains(now) {
			return slot, true
		}
	}

	return Session{}, false
}
