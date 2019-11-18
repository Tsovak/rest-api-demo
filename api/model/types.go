package model

// Account represents an account entity
type Account struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Currency string `db:"currency"`
	Balance  int64  `db:"balance"`
}

// Payment define payment entity
type Payment struct {
	ID            int64  `db:"id"`
	Amount        int64  `db:"amount"`
	ToAccountID   string `db:"to_account_id" json:"to_id" `
	FromAccountID string `db:"from_account_id" json:"from_id"`
	Direction     string `db:"direction"`
}
