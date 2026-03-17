package test

import (
	"testing"

	"github.com/nskforward/trading"
	"github.com/nskforward/trading/finam"
)

func TestCore(t *testing.T) {
	broker := finam.NewBroker()
	core := trading.NewCore(broker)
	err := core.Run()
	if err != nil {
		t.Fatal(err)
	}
}
