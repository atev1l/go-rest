package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func (s *PostgressStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccounts(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewPostgressStore() (*PostgressStore, error) {
	connStr := "user=postgres password=gym dbname=postgres sslmode=verify-full sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgressStore{db: db}, nil

}

func (s *PostgressStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgressStore) CreateAccountTable() error {
	query := `create table if not exists account (
    	id serial primary key,
    	first_name varchar(50),
        last_name varchar(50),
        number serial,
        balance serial,
        created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgressStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name, last_name, number, balance, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	fmt.Println("Account created successfully")
	return nil
}

func (s *PostgressStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgressStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number=$1", number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}
	return nil, fmt.Errorf("account %d not found", number)
}

func (s *PostgressStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`

	_, err := s.db.Query(query, id)

	return err
}
func (s *PostgressStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id=$1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccounts(rows *sql.Rows) (*Account, error) {
	account := new(Account)

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
