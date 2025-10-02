package receipts

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	UserNotFound       = errors.New("user not found")
	InvoiceNotFound    = errors.New("invoice not found")
	InvoiceAlreadyPaid = errors.New("invoice already paid")
)

type User struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Balance   int64  `db:"balance" json:"balance"`
}

type Invoice struct {
	ID     int    `db:"id" json:"id"`
	UserID int    `db:"user_id" json:"user_id"`
	Status string `db:"status" json:"status"`
	Label  string `db:"label" json:"label"`
	Amount int64  `db:"amount" json:"amount"`
}

func GetUser(id int) (User, error) {
	var user User

	db, err := getDBConnection()
	if err != nil {
		return user, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, first_name, last_name, balance FROM users WHERE id = $1", id)

	err = row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Balance)
	if err == sql.ErrNoRows {
		return user, UserNotFound
	}
	if err != nil {
		return user, fmt.Errorf("scanning user row: %w", err)
	}

	return user, nil
}

type CreateInvoiceRequest struct {
	UserID int    `json:"user_id"`
	Amount int64  `json:"amount"`
	Label  string `json:"label"`
}

func CreateInvoice(req CreateInvoiceRequest) (Invoice, error) {
	var invoice Invoice

	db, err := getDBConnection()
	if err != nil {
		return invoice, err
	}
	defer db.Close()

	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserID).Scan(&userExists)
	if err != nil {
		return invoice, err
	}
	if !userExists {
		return invoice, UserNotFound
	}

	query := `INSERT INTO invoices (user_id, label, amount) VALUES ($1, $2, $3) RETURNING id, user_id, status, label, amount`
	err = db.QueryRow(query, req.UserID, req.Label, req.Amount).Scan(
		&invoice.ID, &invoice.UserID, &invoice.Status, &invoice.Label, &invoice.Amount)
	if err != nil {
		return invoice, fmt.Errorf("creating invoice: %w", err)
	}

	return invoice, nil
}

func MarkInvoiceAsPaid(invoiceID int) (Invoice, error) {
	var invoice Invoice

	db, err := getDBConnection()
	if err != nil {
		return invoice, err
	}
	defer db.Close()

	var currentStatus string
	err = db.QueryRow("SELECT status FROM invoices WHERE id = $1", invoiceID).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return invoice, InvoiceNotFound
	}
	if err != nil {
		return invoice, fmt.Errorf("checking invoice status: %w", err)
	}

	if currentStatus == "paid" {
		return invoice, InvoiceAlreadyPaid
	}

	query := `UPDATE invoices SET status = 'paid' WHERE id = $1 RETURNING id, user_id, status, label, amount`
	err = db.QueryRow(query, invoiceID).Scan(
		&invoice.ID, &invoice.UserID, &invoice.Status, &invoice.Label, &invoice.Amount)
	if err != nil {
		return invoice, fmt.Errorf("updating invoice status: %w", err)
	}

	return invoice, nil
}

func getDBConnection() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "receipts-db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "jump"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "receipts"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", psqlInfo)
}
