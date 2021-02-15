package rvdb_clickhouse

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

/*
Usage example:
    chHosts := strings.Split(stringCHHosts, ",")
    chDB, err = rvdb_clickhouse.Connect(chHosts, authCHUser, authCHPass)
    if err != nil {
        log.Fatalf("Can't connect to ClickHouse: %v", err)
    }
    go rvdb_clickhouse.RunPinger(chDB, func(err error) {
        if err != nil {
            log.Printf("ClickHouse ping: error %v", err)
        } else {
            log.Printf("ClickHouse ping: OK")
        }
    })

Inserter example:

func (rb RecordBatchMyData) FlushToDBStmt(stmt *sql.Stmt) (rvdb_clickhouse.Result, error) {
    res := rvdb_clickhouse.Result{}

    for _, v := range rb.Data {
        // ch driver does not support sql.Result fields
        _, err := stmt.Exec(
            rb.UserID,
            rb.SessionID,
            v.TimestampMillis,
            v.MyData,
        )

        if err != nil {
            // silence row errors, if any, just count
            res.RowsError++
            continue
        }

        res.RowsProcessed++
    }

    return res, nil
}
*/

type Result struct {
	RowsProcessed int64
	RowsError     int64
}

type StorageMeta struct {
	Schema     string
	TableWrite string
	TableRead  string
	Fields     []string
}

type Inserter interface {
	FlushToDBStmt(stmt *sql.Stmt) (Result, error)
}

func InsertRecords(db *sql.DB, meta StorageMeta, recs []Inserter, silenceRecFlushErrors bool) (Result, error) {
	tx, err := db.Begin()
	if err != nil {
		// fatal error, return
		return Result{}, err
	}
	defer func() { _ = tx.Rollback() }()

	placeholders := make([]string, len(meta.Fields))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES (%s)",
		meta.Schema,
		meta.TableWrite,
		strings.Join(meta.Fields, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		// fatal error, return
		return Result{}, err
	}
	defer stmt.Close()

	res := Result{}
	for _, rec := range recs {
		recRes, err := rec.FlushToDBStmt(stmt)
		res.RowsProcessed += recRes.RowsProcessed
		res.RowsError += recRes.RowsError

		if err != nil && !silenceRecFlushErrors {
			return res, err
		}
	}

	if err := tx.Commit(); err != nil {
		return res, err
	}
	return res, nil
}

func Connect(chHosts []string, authCHUser, authCHPass string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse",
		fmt.Sprintf(
			// TODO: dsn-config, like conf = ParseDSN(string); conf.Timeout = 5; sql.Open(conf.String())
			"tcp://%s?timeout=5&read_timeout=300&write_timeout=300&pool_size=200&altHosts=%s&username=%s&password=%s",
			chHosts[0], strings.Join(chHosts[1:], ","), authCHUser, authCHPass,
		),
	)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunPinger(db *sql.DB, callback func(error)) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := db.Ping()
		go callback(err)
	}
}
