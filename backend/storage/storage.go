package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/afafpratama/ideaimajitest/types"
)

// ---------------------- SETUP CONFIGURATION ---------------------- //

// Define usable storage function
type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(int, *types.Account) error
	GetAccounts(string, int, int) ([]*types.Account, int, int, error)
	GetAccountByID(int) (*types.Account, error)
	GetAccountByUsername(string) (*types.Account, error)

	CreateCustomer(*types.Customer) error
	DeleteCustomer(int) error
	UpdateCustomer(int, *types.Customer) error
	GetCustomers(string, int, int) ([]*types.Customer, int, int, error)
	GetCustomerByID(int) (*types.Customer, error)

	CreateOrder(*types.Order) error
	DeleteOrder(int) error
	UpdateOrder(int, *types.Order) error
	GetOrders(string, int, int) ([]*types.Order, int, int, error)
	GetOrderByID(int) (*types.Order, error)
}

// Initiate storage structure with database
type PostgresStore struct {
	db *sql.DB
}

// Setup connection to local database with designated password and without ssl
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=ideaimajitest sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

// Define what to do in initiation process
func (s *PostgresStore) Init() error {
	// return s.CreateAccountTable(), s.CreateCustomerTable(), s.CreateOrderTable()

	if err := s.CreateAccountTable(); err != nil {
		return err
	}
	if err := s.CreateCustomerTable(); err != nil {
		return err
	}
	if err := s.CreateOrderTable(); err != nil {
		return err
	}

	return nil
}

// ---------------------- ACCOUNT ---------------------- //

// Create sys_account table with sql query
func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS sys_account (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		phone VARCHAR(50),
		username VARCHAR(100),
		password TEXT,
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

// Get all data in sys_account table with sql query and designated params
func (s *PostgresStore) GetAccounts(search string, page, limit int) ([]*types.Account, int, int, error) {
	var count int

	offset := (page - 1) * limit

	if search != "" {
		search = "%" + search + "%"

		err := s.db.QueryRow(`SELECT COUNT(*) FROM sys_account WHERE name iLIKE $1`, search).Scan(&count)
		if err != nil {
			return nil, 0, 0, err
		}

		rows, err := s.db.Query(`SELECT * FROM sys_account WHERE name iLIKE $1 LIMIT $2 OFFSET $3`, search, limit, offset)
		if err != nil {
			return nil, 0, 0, err
		}

		totalPages := count / limit
		if count%limit != 0 {
			totalPages++
		}

		fmt.Println("totalPages search", totalPages)

		accounts := []*types.Account{}
		for rows.Next() {
			account, err := scanIntoAccount(rows)
			if err != nil {
				return nil, 0, 0, err
			}
			accounts = append(accounts, account)
		}

		return accounts, count, totalPages, nil
	}

	err := s.db.QueryRow(`SELECT COUNT(*) FROM sys_account`).Scan(&count)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := s.db.Query(`SELECT * FROM sys_account LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := count / limit
	if count%limit != 0 {
		totalPages++
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, 0, 0, err
		}
		accounts = append(accounts, account)
	}

	return accounts, count, totalPages, nil
}

// Get data in sys_account table by ID with sql query
func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM sys_account WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account with id [%d] not found", id)
}

// Get data in sys_account table by username with sql query
func (s *PostgresStore) GetAccountByUsername(username string) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM sys_account WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account with username [%s] not found", username)
}

// Create new data in sys_account table with sql query
func (s *PostgresStore) CreateAccount(req *types.Account) error {
	query := `INSERT INTO 
		sys_account (name, phone, username, password, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		req.Name,
		req.Phone,
		req.Username,
		req.Password,
		req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update data in sys_account table with sql query
func (s *PostgresStore) UpdateAccount(id int, req *types.Account) error {
	query := `UPDATE sys_account
		SET name = $2, phone = $3, password = $4
		WHERE id = $1`

	_, err := s.db.Query(
		query,
		id,
		req.Name,
		req.Phone,
		req.Password)

	if err != nil {
		return err
	}

	return nil
}

// Delete data in sys_account table by ID with sql query
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM sys_account WHERE id = $1", id)
	return err
}

// Map account columns to requiring request
func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	req := new(types.Account)
	err := rows.Scan(
		&req.ID,
		&req.Name,
		&req.Phone,
		&req.Username,
		&req.Password,
		&req.CreatedAt)

	return req, err
}

// ---------------------- CUSTOMER ---------------------- //

// Create dt_customer table with sql query
func (s *PostgresStore) CreateCustomerTable() error {
	query := `CREATE TABLE IF NOT EXISTS dt_customer (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		phone VARCHAR(50),
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

// Get all data in dt_customer table with sql query and designated params
func (s *PostgresStore) GetCustomers(search string, page, limit int) ([]*types.Customer, int, int, error) {
	var count int

	offset := (page - 1) * limit

	if search != "" {
		search = "%" + search + "%"

		err := s.db.QueryRow(`SELECT COUNT(*) FROM dt_customer WHERE name iLIKE $1`, search).Scan(&count)
		if err != nil {
			return nil, 0, 0, err
		}

		rows, err := s.db.Query(`SELECT * FROM dt_customer WHERE name iLIKE $1 LIMIT $2 OFFSET $3`, search, limit, offset)
		if err != nil {
			return nil, 0, 0, err
		}

		totalPages := count / limit
		if count%limit != 0 {
			totalPages++
		}

		fmt.Println("totalPages search", totalPages)

		customers := []*types.Customer{}
		for rows.Next() {
			customer, err := scanIntoCustomer(rows)
			if err != nil {
				return nil, 0, 0, err
			}
			customers = append(customers, customer)
		}

		return customers, count, totalPages, nil
	}

	err := s.db.QueryRow(`SELECT COUNT(*) FROM dt_customer`).Scan(&count)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := s.db.Query(`SELECT * FROM dt_customer LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := count / limit
	if count%limit != 0 {
		totalPages++
	}

	customers := []*types.Customer{}
	for rows.Next() {
		customer, err := scanIntoCustomer(rows)
		if err != nil {
			return nil, 0, 0, err
		}
		customers = append(customers, customer)
	}

	return customers, count, totalPages, nil
}

// Get data in dt_customer table by ID with sql query
func (s *PostgresStore) GetCustomerByID(id int) (*types.Customer, error) {
	rows, err := s.db.Query(`SELECT * FROM dt_customer WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoCustomer(rows)
	}

	return nil, fmt.Errorf("Customer with id [%d] not found", id)
}

// Create new data in dt_customer table with sql query
func (s *PostgresStore) CreateCustomer(req *types.Customer) error {
	query := `INSERT INTO 
		dt_customer (name, phone, created_at) 
		VALUES ($1, $2, $3)`

	_, err := s.db.Query(
		query,
		req.Name,
		req.Phone,
		req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update data in dt_customer table with sql query
func (s *PostgresStore) UpdateCustomer(id int, req *types.Customer) error {
	query := `UPDATE dt_customer
		SET name = $2, phone = $3
		WHERE id = $1`

	_, err := s.db.Query(
		query,
		id,
		req.Name,
		req.Phone)

	if err != nil {
		return err
	}

	return nil
}

// Delete data in dt_customer table by ID with sql query
func (s *PostgresStore) DeleteCustomer(id int) error {
	_, err := s.db.Query("DELETE FROM dt_customer WHERE id = $1", id)
	return err
}

// Map customer columns to requiring request
func scanIntoCustomer(rows *sql.Rows) (*types.Customer, error) {
	req := new(types.Customer)
	err := rows.Scan(
		&req.ID,
		&req.Name,
		&req.Phone,
		&req.CreatedAt)

	return req, err
}

// ---------------------- ORDER ---------------------- //

// Create dt_order table with sql query
func (s *PostgresStore) CreateOrderTable() error {
	query := `CREATE TABLE IF NOT EXISTS dt_order (
		id SERIAL PRIMARY KEY,
		customer_id INT,
		service TEXT,
		amount INT,
		unit VARCHAR(100),
		price INT,
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

// Get all data in dt_order table with sql query and designated params
func (s *PostgresStore) GetOrders(search string, page, limit int) ([]*types.Order, int, int, error) {
	var count int

	offset := (page - 1) * limit

	if search != "" {
		search = "%" + search + "%"

		err := s.db.QueryRow(`SELECT COUNT(dt_order.id) FROM dt_order JOIN dt_customer ON dt_order.customer_id = dt_customer.id WHERE service iLIKE $1`, search).Scan(&count)
		if err != nil {
			return nil, 0, 0, err
		}

		rows, err := s.db.Query(`SELECT 
				o.id AS id,
				o.customer_id AS customer_id,
				c.name AS name,
				c.phone AS phone,
				o.service AS service,
				o.amount AS amount,
				o.unit AS unit,
				o.price AS price,
				o.created_at AS created_at
			FROM 
				dt_order o
			JOIN
				dt_customer c
			ON
				o.customer_id = c.id WHERE service iLIKE $1 LIMIT $2 OFFSET $3`, search, limit, offset)
		if err != nil {
			return nil, 0, 0, err
		}

		totalPages := count / limit
		if count%limit != 0 {
			totalPages++
		}

		orders := []*types.Order{}
		for rows.Next() {
			order, err := scanIntoOrder(rows)
			if err != nil {
				return nil, 0, 0, err
			}
			orders = append(orders, order)
		}

		return orders, count, totalPages, nil
	}

	err := s.db.QueryRow(`SELECT COUNT(dt_order.id) FROM dt_order JOIN dt_customer ON dt_order.customer_id = dt_customer.id`).Scan(&count)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := s.db.Query(`SELECT 
				o.id AS id,
				o.customer_id AS customer_id,
				c.name AS name,
				c.phone AS phone,
				o.service AS service,
				o.amount AS amount,
				o.unit AS unit,
				o.price AS price,
				o.created_at AS created_at
			FROM 
				dt_order o
			JOIN
				dt_customer c
			ON
				o.customer_id = c.id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := count / limit
	if count%limit != 0 {
		totalPages++
	}

	orders := []*types.Order{}
	for rows.Next() {
		order, err := scanIntoOrder(rows)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, order)
	}

	return orders, count, totalPages, nil
}

// Get data in dt_order table by ID with sql query
func (s *PostgresStore) GetOrderByID(id int) (*types.Order, error) {
	rows, err := s.db.Query(`SELECT 
				o.id AS id,
				o.customer_id AS customer_id,
				c.name AS name,
				c.phone AS phone,
				o.service AS service,
				o.amount AS amount,
				o.unit AS unit,
				o.price AS price,
				o.created_at AS created_at
			FROM 
				dt_order o
			JOIN
				dt_customer c
			ON
				o.customer_id = c.id WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, fmt.Errorf("Order with id [%d] not found", id)
}

// Create new data in dt_order table with sql query
func (s *PostgresStore) CreateOrder(req *types.Order) error {
	query := `INSERT INTO 
		dt_order (customer_id, service, amount, unit, price, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		req.Customer,
		req.Service,
		req.Amount,
		req.Unit,
		req.Price,
		req.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update data in dt_order table with sql query
func (s *PostgresStore) UpdateOrder(id int, req *types.Order) error {
	query := `UPDATE dt_order
		SET customer_id = $2, service = $3, amount = $4, unit = $5, price = $6
		WHERE id = $1`

	_, err := s.db.Query(
		query,
		id,
		req.Customer,
		req.Service,
		req.Amount,
		req.Unit,
		req.Price)

	if err != nil {
		return err
	}

	return nil
}

// Delete data in dt_order table by ID with sql query
func (s *PostgresStore) DeleteOrder(id int) error {
	_, err := s.db.Query("DELETE FROM dt_order WHERE id = $1", id)
	return err
}

// Map order columns to requiring request
func scanIntoOrder(rows *sql.Rows) (*types.Order, error) {
	req := new(types.Order)
	err := rows.Scan(
		&req.ID,
		&req.Customer,
		&req.Name,
		&req.Phone,
		&req.Service,
		&req.Amount,
		&req.Unit,
		&req.Price,
		&req.CreatedAt)

	return req, err
}
