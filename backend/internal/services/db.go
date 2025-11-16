package services

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB(path string) (*sqlx.DB, error) {
	// If the path looks like a Postgres URL, use pgx driver
	if strings.HasPrefix(path, "postgres://") || strings.HasPrefix(path, "postgresql://") {
		db, err := sqlx.Connect("pgx", path)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	// Default to sqlite file
	dsn := fmt.Sprintf("file:%s?_foreign_keys=1", path)
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDB(db *sqlx.DB) error {
	driver := db.DriverName()
	var schema string
	if strings.Contains(driver, "pg") || strings.Contains(driver, "pgx") {
		schema = `
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  symbol TEXT,
  side TEXT,
  quantity INTEGER,
  remaining INTEGER,
  price NUMERIC,
  status TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fills (
  id SERIAL PRIMARY KEY,
  order_id INTEGER REFERENCES orders(id),
  user_id INTEGER REFERENCES users(id),
  symbol TEXT,
  side TEXT,
  quantity INTEGER,
  price NUMERIC,
  filled_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_fills_user ON fills(user_id);
`
	} else {
		schema = `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  symbol TEXT,
  side TEXT,
  quantity INTEGER,
  remaining INTEGER,
  price REAL,
  status TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS fills (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  order_id INTEGER,
  user_id INTEGER,
  symbol TEXT,
  side TEXT,
  quantity INTEGER,
  price REAL,
  filled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(order_id) REFERENCES orders(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_fills_user ON fills(user_id);
`
	}

	_, err := db.Exec(schema)
	return err
}
