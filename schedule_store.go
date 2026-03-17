package trading

import (
	"fmt"

	"github.com/nskforward/trading/types"
)

type ScheduleStore struct {
	broker    types.Broker
	schedules map[string]*types.Schedule
}

func NewScheduleStore(broker types.Broker) *ScheduleStore {
	return &ScheduleStore{
		schedules: make(map[string]*types.Schedule),
		broker:    broker,
	}
}

func (store *ScheduleStore) CurrentSession(symbol string) (types.Session, error) {
	schedule, ok := store.schedules[symbol]
	if ok {
		session, ok := schedule.Current()
		if ok {
			return session, nil
		}
	}

	schedule, err := store.broker.GetSchedule(symbol)
	if err != nil {
		return types.Session{}, err
	}

	session, ok := schedule.Current()
	if !ok {
		return types.Session{}, fmt.Errorf("broker returned schedule that does not contain the session with interval of current timestamp")
	}

	store.schedules[symbol] = schedule

	return session, nil
}
