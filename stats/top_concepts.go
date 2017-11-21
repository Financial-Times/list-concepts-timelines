package stats

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const DateLayout = "2006-01-02"

var mutex = &sync.Mutex{}

type TopConcepts struct {
	Start    time.Time
	Stop     time.Time
	ListUUID string
	Concepts map[string]map[string][]*Concept
}

type Concept struct {
	PrefLabel string
	TimeLine  map[time.Time]int
	Total     int
}

func GetTopConcepts(start time.Time, stop time.Time, listUUID string) *TopConcepts {
	mutex.Lock()
	defer mutex.Unlock()

	tc, found := cache.get(start, stop, listUUID)
	if found {
		return tc
	}
	tc = &TopConcepts{
		Start:    start,
		Stop:     stop,
		ListUUID: listUUID,
	}

	tc.init()
	go cache.put(start, stop, listUUID, tc)
	return tc
}

func (tc *TopConcepts) init() {
	log.Info("Getting data...")
	rows, err := db.Query(
		"SELECT count(articles_v2_cdc.article_uuid) as n, annotations_cdc.id, annotations_cdc.pref_label, annotations_cdc.type, articles_v2_cdc.articletype FROM "+
			"(SELECT DISTINCT content_id "+
			"FROM ftarticlesdb.content_list_cdc "+
			"WHERE list_id=$1 AND entry_timestamp>=$2 AND entry_timestamp<=$3) as l "+
			"JOIN ftarticlesdb.articles_v2_cdc ON l.content_id = articles_v2_cdc.article_uuid "+
			"JOIN ftarticlesdb.annotations_cdc ON annotations_cdc.article_uuid = articles_v2_cdc.article_uuid "+
			"GROUP BY annotations_cdc.id, annotations_cdc.pref_label, annotations_cdc.type, articles_v2_cdc.articletype "+
			"ORDER BY n DESC",
		tc.ListUUID, tc.Start, tc.Stop)
	if err != nil {
		log.WithError(err).Error("Error in getting data")
	}
	defer rows.Close()

	tc.Concepts = make(map[string]map[string][]*Concept)

	for rows.Next() {
		var n int
		var conceptUUID string
		var conceptPrefLabel string
		var conceptType string
		var contentType string
		if err := rows.Scan(&n, &conceptUUID, &conceptPrefLabel, &conceptType, &contentType); err != nil {
			log.WithError(err).Error()
		}
		log.WithField("n", n).
			WithField("conceptUUID", conceptUUID).
			WithField("conceptPrefLabel", conceptPrefLabel).
			WithField("conceptType", conceptType).
			WithField("contentType", contentType).
			Info()

		tc.addConceptEntry(conceptUUID, conceptPrefLabel, conceptType, n, contentType)
	}
	if err := rows.Err(); err != nil {
		log.WithError(err).Error()
	}
}

func (tc *TopConcepts) addConceptEntry(conceptUUID string, conceptPrefLabel string, conceptType string, n int, contentType string) {
	//global
	c := &Concept{
		PrefLabel: conceptPrefLabel,
		Total:     n,
	}

	if conceptType == "ORGANISATION" || conceptType == "PERSON" || conceptType == "LOCATION" || conceptType == "TOPIC" {
		conceptMap, found := tc.Concepts["concepts"]
		if !found {
			conceptMap = make(map[string][]*Concept)
		}
		conceptList, found := conceptMap["All"]
		if !found {
			conceptList = []*Concept{}
		}
		conceptList = append(conceptList, c)
		conceptMap["All"] = conceptList

		conceptList, found = conceptMap[contentType]
		if !found {
			conceptList = []*Concept{}
		}
		conceptList = append(conceptList, c)
		conceptMap[contentType] = conceptList

		tc.Concepts["concepts"] = conceptMap
	}

	conceptMap, found := tc.Concepts[conceptType]
	if !found {
		conceptMap = make(map[string][]*Concept)
	}
	conceptList, found := conceptMap["All"]
	if !found {
		conceptList = []*Concept{}
	}
	conceptList = append(conceptList, c)
	conceptMap["All"] = conceptList

	conceptList, found = conceptMap[contentType]
	if !found {
		conceptList = []*Concept{}
	}
	conceptList = append(conceptList, c)
	conceptMap[contentType] = conceptList

	tc.Concepts[conceptType] = conceptMap
}

func (tc *TopConcepts) GetDataTable(conceptType string, limit int, contentType string) *DataTable {
	dt := newDataTable()
	dt.addCol("Concept", "Concept", "", "string")
	dt.addCol("Appearances", "Appearances", "", "number")

	cList, found := tc.Concepts[conceptType][contentType]

	if found {
		if limit != -1 {
			for i := 0; i < limit && i < len(cList); i++ {
				dt.addRow([]Value{{V: cList[i].PrefLabel}, {V: cList[i].Total}})
			}
		} else {
			for _, c := range tc.Concepts[conceptType][contentType] {
				dt.addRow([]Value{{V: c.PrefLabel}, {V: c.Total}})
			}
		}

	}

	return dt
}

type DataTable struct {
	Cols []Column `json:"cols"`
	Rows []Row    `json:"rows"`
}

type Column struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type Row struct {
	C []Value `json:"c"`
}

type Value struct {
	V interface{} `json:"v"`
	F interface{} `json:"f"`
}

func newDataTable() *DataTable {
	return &DataTable{[]Column{}, []Row{}}
}

func (dt *DataTable) addCol(id, label, pattern, t string) {
	dt.Cols = append(dt.Cols, Column{id, label, pattern, t})
}

func (dt *DataTable) addRow(values []Value) {
	dt.Rows = append(dt.Rows, Row{values})
}
