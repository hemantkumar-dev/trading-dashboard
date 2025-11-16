package services

import (
	"database/sql"
	"errors"

	"trading-dashboard/backend/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrNoUser = errors.New("user not found")

// EnsureUser returns user id, creates user if not exists (simple demo)
func EnsureUser(db *sqlx.DB, username string) (int64, error) {
	var id int64
	err := db.Get(&id, "SELECT id FROM users WHERE username = ?", username)
	if err == nil {
		return id, nil
	}
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO users(username) VALUES(?)", username)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	return 0, err
}

func PlaceOrder(db *sqlx.DB, username string, o models.Order) (int64, error) {
	uid, err := EnsureUser(db, username)
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(`INSERT INTO orders(user_id,symbol,side,quantity,remaining,price,status)
		VALUES(?,?,?,?,?,?,'open')`, uid, o.Symbol, o.Side, o.Quantity, o.Quantity, o.Price)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

type DBOrder struct {
	ID        int64   `db:"id" json:"id"`
	UserID    int64   `db:"user_id" json:"-"`
	Symbol    string  `db:"symbol" json:"symbol"`
	Side      string  `db:"side" json:"side"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Remaining int     `db:"remaining" json:"remaining"`
	Price     float64 `db:"price" json:"price"`
	Status    string  `db:"status" json:"status"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}

func GetAllOrders(db *sqlx.DB) ([]DBOrder, error) {
	var out []DBOrder
	err := db.Select(&out, "SELECT * FROM orders ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	return out, nil
}

// helper to record fill
func recordFill(db *sqlx.DB, orderID int64, userID int64, symbol, side string, qty int, price float64) error {
	_, err := db.Exec(`INSERT INTO fills(order_id,user_id,symbol,side,quantity,price) VALUES(?,?,?,?,?,?)`,
		orderID, userID, symbol, side, qty, price)
	return err
}
