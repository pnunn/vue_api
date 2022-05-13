package driver

import (
  "database/sql"
  "fmt"
  _ "github.com/jackc/pgconn"
  _ "github.com/jackc/pgx/v4"
  _ "github.com/jackc/pgx/v4/stdlib"
  "time"
)

type DB struct {
  SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectPostgres(dsn string) (*DB, error) {
  d, err := sql.Open("pgx", dsn)
  if err != nil {
    return nil, err
  }

  d.SetMaxOpenConns(maxOpenDbConn)
  d.SetConnMaxIdleTime(maxIdleDbConn)
  d.SetConnMaxLifetime(maxDbLifetime)

  err = testDB(d)
  if err != nil {
    return nil, err
  }

  dbConn.SQL = d
  return dbConn, nil
}

func testDB(d *sql.DB) error {
  err := d.Ping()
  if err != nil {
    fmt.Println("Error!", err)
    return err
  }

  fmt.Println("*** Ping database successfully ***")

  return err
}
