package labCode

import (
	"fmt"
	"time"
)

const INTERVAL_PERIOD time.Duration = 2 * time.Second

const HOUR_TO_TICK int = 23
const MINUTE_TO_TICK int = 00
const SECOND_TO_TICK int = 03

type Scheduler struct {
	Timer *time.Timer
}

func (t *Scheduler) RunningRoutine() {
	Scheduler := &Scheduler{}
	Scheduler.updateTimer()
	fmt.Println("running routine")
	for {
		<-Scheduler.Timer.C
		fmt.Println(time.Now(), "- just ticked")
		Scheduler.updateTimer()
	}
}

func (t *Scheduler) updateTimer() {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	fmt.Println(nextTick, "- next tick")
	diff := nextTick.Sub(time.Now())
	if t.Timer == nil {
		t.Timer = time.NewTimer(diff)
	} else {
		t.Timer.Reset(diff)
	}
}
