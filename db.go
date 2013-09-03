package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type GoxtraceRecord struct {
	XtrId  string
	Record map[string]interface{}
}

func NewGoxtraceRecord(data map[string]interface{}) *GoxtraceRecord {
	record := GoxtraceRecord{Record: data}
	value, valid := data["X-Trace"].(string)
	if !valid {
		return nil
	}
	record.XtrId = value
	return &record
}

func (r *GoxtraceRecord) TaskId() string {
	return r.XtrId[0:16]
}

type GoxDb struct {
	conn *sql.DB
}

func NewGoxDb(filename string) *GoxDb {
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	v := &GoxDb{conn: conn}
	v.init()
	return v
}

func (g *GoxDb) init() {
	sqls := []string{
		"create table if not exists entries (xtrid varchar not null primary key, taskid varchar, body text)",
	}
	for _, sql := range sqls {
		_, err := g.conn.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
}

func (g *GoxDb) Close() {
	g.conn.Close()
}

func (g *GoxDb) Write(record *GoxtraceRecord) {
	tx, err := g.conn.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("insert into entries(xtrid, taskid, body) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(record.XtrId, record.TaskId(), "<wip>")
	tx.Commit()
}
