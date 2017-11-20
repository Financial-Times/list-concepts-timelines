package stats

import (
	"sync"
	"time"
)

var cache = &topConceptCache{&sync.Map{}}

type topConceptCache struct {
	*sync.Map
}

func (c *topConceptCache) get(start time.Time, stop time.Time, listUUID string) (*TopConcepts, bool) {
	m, found := c.Load(start)
	if !found {
		return nil, found
	}
	m, found = m.(*sync.Map).Load(stop)
	if !found {
		return nil, found
	}
	tc, found := m.(*sync.Map).Load(listUUID)
	if !found {
		return nil, found
	}
	return tc.(*TopConcepts), found
}

func (c *topConceptCache) put(start time.Time, stop time.Time, listUUID string, tc *TopConcepts) {
	m, found := c.Load(start)
	if !found {
		m = &sync.Map{}
		c.Store(start, m)
	}
	m2, found := m.(*sync.Map).Load(stop)
	if !found {
		m2 = &sync.Map{}
		m.(*sync.Map).Store(stop, m2)
	}
	m2.(*sync.Map).Store(listUUID, tc)
}
