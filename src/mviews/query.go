package main

import (
	"fmt"
	"github.com/gchaincl/dotsql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

///////
/// Boilerplate
///////

type Query struct {
	name        string
	lookerRoute string
	namedRoute  string
	// Defaults
	refreshInterval int
	refreshedAt     int
	activeQuery     string
	statements      *dotsql.DotSql
	C               chan string
}

func newQuery(name string, lookId string, refreshInterval int) Query {
	sql := fmt.Sprintf("sql/%s.sql", name)
	dot, err := dotsql.LoadFromFile(sql)
	if err != nil {
		log.Fatalf("Dotsql: unable to find or load %v", sql)
	}

	q := Query{
		lookerRoute:     looker(lookId),
		name:            name,
		namedRoute:      named(name),
		refreshInterval: refreshInterval,
		statements:      dot,
		C:               refreshChan,
	}
	q.activeQuery = q.SelectPrimary()
	q.StartTicker()
	return q
}

type query interface {
	StartTicker()
}

type refresher interface {
	refetchData() error
	Reset()
	LookerRoute() string
	NamedRoute() string
	ActiveQuery() string
	Name() string
	Handler(http.ResponseWriter, *http.Request)
}

func (q Query) Name() string        { return q.name }
func (q Query) ActiveQuery() string { return q.activeQuery }
func (q Query) NamedRoute() string  { return q.namedRoute }
func (q Query) LookerRoute() string { return q.lookerRoute }

func (q Query) StartTicker() {
	i := time.Duration(q.refreshInterval) * time.Second
	t := time.NewTicker(i)

	go func() {
		for {
			select {
			case <-t.C:
				q.RefreshData(dbx)
				q.C <- q.name
			}
		}
	}()
}

func (q Query) Statements() *dotsql.DotSql {
	return q.statements
}

func (q Query) RefreshData(dbx *sqlx.DB) {
	log.Printf("Starting to drop %s table and then rebuild", q.name)

	// Ensure we're reading from primary table
	q.activeQuery = q.SelectPrimary()
	log.Printf("ActiveQuery now: %s\n", q.activeQuery)
	log.Printf("Rebuilding staging\n")
	_, err := dbx.Query(q.RebuildStagingQuery())
	if err != nil {
		log.Fatalf("Rebuilding %s table failed because of %v", q.name, err)
	}
	// Set reads to read from staging table
	q.activeQuery = q.SelectSecondary()
	log.Printf("Switching to secondary %s\n", q.activeQuery)
	_, err = dbx.Query(q.ReplacePrimaryQuery())
	if err != nil {
		log.Printf("Rebuilding table failed because of %v", err)
	}
	// Return to reading from primary table
	q.activeQuery = q.SelectPrimary()
	log.Printf("Switching to primary %s\n", q.activeQuery)
	log.Printf("Regenerated %s table", q.name)
}

func (q Query) SelectPrimary() string {
	return getQuery(q, "select-primary")
}

func (q Query) SelectSecondary() string {
	return getQuery(q, "select-secondary")
}

func (q Query) RebuildStagingQuery() string {
	return getQuery(q, "rebuild-staging-query")
}

func (q Query) ReplacePrimaryQuery() string {
	return getQuery(q, "replace-primary-query")
}
