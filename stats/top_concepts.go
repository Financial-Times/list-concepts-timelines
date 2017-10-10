package stats

import "time"

const DateLayout = "2006-01-02"

type TopConcepts struct {
	start time.Time
	stop  time.Time
}

func NewTopConcepts(start time.Time, stop time.Time) *TopConcepts {
	return &TopConcepts{start, stop}
}

func (tc *TopConcepts) Start() string {
	return tc.start.Format(DateLayout)
}

func (tc *TopConcepts) Stop() string {
	return tc.stop.Format(DateLayout)
}
