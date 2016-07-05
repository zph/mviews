package main

import (
	"fmt"
	"github.com/gchaincl/dotsql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
)

type Look interface {
	RefreshData(*sqlx.DB)
	SelectPrimary() string
	SelectSecondary() string
	RebuildStagingQuery() string
	ReplacePrimaryQuery() string
	Statements() *dotsql.DotSql
}

func looker(s string) string { return fmt.Sprintf("/looker/%s.json", s) }
func named(s string) string  { return fmt.Sprintf("/named/%s.json", s) }

func getQuery(q Look, name string) string {
	s, err := q.Statements().Raw(name)
	if err != nil {
		log.Fatalf("Failed reading %v with error %v", s, err)
	}
	return s
}

func setupDb() {
	var err error
	dbx, err = sqlx.Open("postgres", connection())
	if err != nil {
		log.Fatalf("DB couldn't connect %v with err %v\n", dbx, err)
	}
}

func connection() string {
	conn := os.Getenv("REDSHIFT_CONNECTION")
	u := os.Getenv("REDSHIFT_USERNAME")
	p := os.Getenv("REDSHIFT_PASSWORD")
	return fmt.Sprintf(conn, u, p)
}

func setRefresher(c map[string]refresher) {
	go func() {
		for {
			select {
			case n := <-refreshChan:
				c[n].refetchData()
				log.Printf("Refreshed cache data for %s", n)
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
}
