package services

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB(path string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=1", path)
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDB(db *sqlx.DB) error {
	schema := `
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
	_, err := db.Exec(schema)
	return err
}
