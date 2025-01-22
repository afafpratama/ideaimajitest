package types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ---------------------- ACCOUNT ---------------------- //

// Map login response after successfully submitted, getting JWT token
type LoginResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

// Map required login request to be submitted
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Map designated columns to create sys_account table
type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// Map required columns to create / register account data
type CreateAccountRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Map required columns to update account data
type UpdateAccountRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Password validation function compare hashedEncryptedPassword with submitted password by byte
func (a *Account) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) == nil
}

// Define new account variables with encrypted password and auto-generated created time by timestamp
func NewAccount(Name string, Phone string, Username string, Password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		Name:      Name,
		Phone:     Phone,
		Username:  Username,
		Password:  string(encpw),
		CreatedAt: time.Now().Local().UTC(),
	}, nil
}

// Define updated account variables with encrypted password
func UpdateAccount(ID int, Name string, Phone string, Password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:       ID,
		Name:     Name,
		Phone:    Phone,
		Password: string(encpw),
	}, nil
}

// ---------------------- CUSTOMER ---------------------- //

// Map designated columns to create dt_customer table
type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

// Map required columns to create customer data
type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Map required columns to update customer data
type UpdateCustomerRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Define new customer variables with encrypted password and auto-generated created time by timestamp
func NewCustomer(Name string, Phone string) (*Customer, error) {
	return &Customer{
		Name:      Name,
		Phone:     Phone,
		CreatedAt: time.Now().Local().UTC(),
	}, nil
}

// Define updated customer variables with encrypted password
func UpdateCustomer(ID int, Name string, Phone string) (*Customer, error) {
	return &Customer{
		ID:    ID,
		Name:  Name,
		Phone: Phone,
	}, nil
}

// ---------------------- ORDER ---------------------- //

// Map designated columns to create dt_order table
type Order struct {
	ID        int       `json:"id"`
	Customer  int       `json:"customer_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Service   string    `json:"service"`
	Amount    int       `json:"amount"`
	Unit      string    `json:"unit"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// Map required columns to create order data
type CreateOrderRequest struct {
	Customer int    `json:"customer_id"`
	Service  string `json:"service"`
	Amount   int    `json:"amount"`
	Unit     string `json:"unit"`
	Price    int    `json:"price"`
}

// Map required columns to update order data
type UpdateOrderRequest struct {
	ID       int    `json:"id"`
	Customer int    `json:"customer_id"`
	Service  string `json:"service"`
	Amount   int    `json:"amount"`
	Unit     string `json:"unit"`
	Price    int    `json:"price"`
}

// Define new order variables with encrypted password and auto-generated created time by timestamp
func NewOrder(Customer int, Service string, Amount int, Unit string, Price int) (*Order, error) {
	return &Order{
		Customer:  Customer,
		Service:   Service,
		Amount:    Amount,
		Unit:      Unit,
		Price:     Price,
		CreatedAt: time.Now().Local().UTC(),
	}, nil
}

// Define updated order variables with encrypted password
func UpdateOrder(ID int, Customer int, Service string, Amount int, Unit string, Price int) (*Order, error) {
	return &Order{
		ID:       ID,
		Customer: Customer,
		Service:  Service,
		Amount:   Amount,
		Unit:     Unit,
		Price:    Price,
	}, nil
}

// ---------------------- RESPONSE ---------------------- //

type Response struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Count      int         `json:"count"`
	TotalPages int         `json:"total"`
}
