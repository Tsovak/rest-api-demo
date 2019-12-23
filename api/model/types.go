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
	ToAccountID   string `db:"to_account_id" `
	FromAccountID string `db:"from_account_id"`
	Direction     string `db:"direction"`
}

// AccountRequest define a request entity for create a new account
type AccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

// AccountRequest define a request entity for create a new payment
type PaymentRequest struct {
	Amount        int64  `json:"amount"`
	ToAccountID   string `json:"to_id" `
	FromAccountID string `json:"from_id"`
}

// AccountResponse define a response entity after account creation
type AccountResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

// PaymentResponse define a response entity after payment creation
type PaymentResponse struct {
	ID            int64  `json:"id"`
	Amount        int64  `json:"amount"`
	ToAccountID   string `json:"to_account" `
	FromAccountID string `json:"from_account"`
	Direction     string `json:"direction"`
}
