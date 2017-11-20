package stats

//import (
//	log "github.com/sirupsen/logrus"
//	"sort"
//	"time"
//)
//
//const DateLayout = "2006-01-02"
//
//type TopConcepts struct {
//	Start    time.Time
//	Stop     time.Time
//	ListUUID string
//	Concepts map[string]*Concept
//}
//
//type Concept struct {
//	PrefLabel string
//	TimeLine  map[time.Time]int
//	Total     int
//}
//
//func GetTopConcepts(start time.Time, stop time.Time, listUUID string) *TopConcepts {
//
//	tc, found := cache.get(start, stop, listUUID)
//	if found {
//		return tc
//	}
//	tc = &TopConcepts{
//		Start:    start,
//		Stop:     stop,
//		ListUUID: listUUID,
//	}
//
//	tc.init()
//	go cache.put(start, stop, listUUID, tc)
//	return tc
//}
//
//func (tc *TopConcepts) init() {
//	log.Info("Getting data...")
//	rows, err := db.Query(
//		"SELECT annotations_cdc.id,annotations_cdc.pref_label,annotations_cdc.type,l.entry_timestamp,l.exit_timestamp FROM "+
//			"(SELECT content_id,entry_timestamp,exit_timestamp FROM ftarticlesdb.content_list_cdc WHERE list_id=$1 AND entry_timestamp>=$2 AND entry_timestamp<=$3) as l "+
//			"JOIN ftarticlesdb.articles_v2_cdc ON l.content_id = articles_v2_cdc.article_uuid "+
//			"JOIN ftarticlesdb.annotations_cdc ON annotations_cdc.article_uuid = articles_v2_cdc.article_uuid;", tc.ListUUID, tc.Start, tc.Stop)
//	if err != nil {
//		log.WithError(err).Error("Error in getting data")
//	}
//	defer rows.Close()
//
//	tc.Concepts = make(map[string]*Concept)
//
//	for rows.Next() {
//		var conceptUUID string
//		var conceptPrefLabel string
//		var conceptType string
//		var entryTime time.Time
//		var exitTime time.Time
//		if err := rows.Scan(&conceptUUID, &conceptPrefLabel, &conceptType, &entryTime, &exitTime); err != nil {
//			log.WithError(err).Error()
//		}
//		log.WithField("conceptUUID", conceptUUID).
//			WithField("conceptPrefLabel", conceptPrefLabel).
//			WithField("entry", entryTime).
//			WithField("exit", exitTime).
//			Info()
//
//		tc.addConceptEntry(conceptUUID, conceptPrefLabel, conceptType, entryTime, exitTime)
//	}
//	if err := rows.Err(); err != nil {
//		log.WithError(err).Error()
//	}
//}
//
//func (tc *TopConcepts) addConceptEntry(conceptUUID string, conceptPrefLabel string, conceptType string, entryTime time.Time, exitTime time.Time) {
//	if exitTime == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
//		exitTime = tc.Stop
//	}
//	concept, found := tc.Concepts[conceptUUID]
//	if !found {
//		concept = &Concept{
//			PrefLabel: conceptPrefLabel,
//			TimeLine:  make(map[time.Time]int),
//		}
//		tc.Concepts[conceptUUID] = concept
//	}
//	timeStamp := tc.Start
//
//	for timeStamp.Before(tc.Stop) {
//
//		_, found := concept.TimeLine[timeStamp]
//		if !found {
//			concept.TimeLine[timeStamp] = 0
//		}
//
//		if timeStamp.After(entryTime) && timeStamp.Before(exitTime) {
//			concept.TimeLine[timeStamp] = concept.TimeLine[timeStamp] + 1
//		}
//		timeStamp = timeStamp.Add(1 * time.Minute)
//	}
//
//	concept.Total++
//}
//
//func (tc *TopConcepts) GetTimeLine() *DataTable {
//	dt := newDataTable()
//	dt.addCol("Concept", "Concept", "", "string")
//	dt.addCol("Appearances", "Appearances", "", "number")
//
//	for _, c := range tc.top(10) {
//		dt.addRow([]Value{{V: c.PrefLabel}, {V: c.Total}})
//	}
//	return dt
//}
//
//func (tc *TopConcepts) top(size int) []*Concept {
//	r := []*Concept{}
//	for _, c := range tc.Concepts {
//		r = append(r, c)
//	}
//	sort.Sort(ByTotal(r))
//	return r[:size]
//}
//
//type ByTotal []*Concept
//
//func (a ByTotal) Len() int           { return len(a) }
//func (a ByTotal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ByTotal) Less(i, j int) bool { return a[i].Total > a[j].Total }
//
//type DataTable struct {
//	Cols []Column `json:"cols"`
//	Rows []Row    `json:"rows"`
//}
//
//type Column struct {
//	ID      string `json:"id"`
//	Label   string `json:"label"`
//	Pattern string `json:"pattern"`
//	Type    string `json:"type"`
//}
//
//type Row struct {
//	C []Value `json:"c"`
//}
//
//type Value struct {
//	V interface{} `json:"v"`
//	F interface{} `json:"f"`
//}
//
//func newDataTable() *DataTable {
//	return &DataTable{[]Column{}, []Row{}}
//}
//
//func (dt *DataTable) addCol(id, label, pattern, t string) {
//	dt.Cols = append(dt.Cols, Column{id, label, pattern, t})
//}
//
//func (dt *DataTable) addRow(values []Value) {
//	dt.Rows = append(dt.Rows, Row{values})
//}
