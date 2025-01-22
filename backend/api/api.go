package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"

	"github.com/afafpratama/ideaimajitest/storage"
	"github.com/afafpratama/ideaimajitest/types"
)

// ---------------------- SETUP CONFIGURATION ---------------------- //

// Initiate API server structure
type APIServer struct {
	listenAddr string
	store      storage.Storage
}

// Create new connection for API server
func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// API initiation and router setup
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/register", makeHTTPHandlerFunc(s.handleRegister))

	router.HandleFunc("/account", withJWTAuth(makeHTTPHandlerFunc(s.handleAccount), s.store))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGetAccountByID), s.store))

	router.HandleFunc("/customer", withJWTAuth(makeHTTPHandlerFunc(s.handleCustomer), s.store))
	router.HandleFunc("/customer/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGetCustomerByID), s.store))

	router.HandleFunc("/order", withJWTAuth(makeHTTPHandlerFunc(s.handleOrder), s.store))
	router.HandleFunc("/order/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGetOrderByID), s.store))

	// Handle CORS in a middleware function
	router.Use(corsMiddleware)

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-JWT-TOKEN, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---------------------- ACCOUNT ---------------------- //

// Handling login method validation, request verification, JWT creation, and response
// @Tags Account
// @Summary API Login
// @Router /login [post]
// @Param request body types.LoginRequest true "Payload body [RAW]"
// @Accept json
// @Produces json
// @Success 200 {array} types.LoginRequest
// @Failure 400
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		fmt.Errorf("Method not allowed %s", r.Method)
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByUsername(req.Username)
	if err != nil {
		return err
	}

	if !acc.ValidatePassword(req.Password) {
		fmt.Errorf("Not authenticated")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Username: acc.Username,
		ID:       acc.ID,
		Name:     acc.Name,
		Phone:    acc.Phone,
		Token:    token,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// Handling account registration required params, transaction, and response
// @Tags Account
// @Summary API Register
// @Router /register [post]
// @Param request body types.CreateAccountRequest true "Payload body [RAW]"
// @Accept json
// @Produces json
// @Success 200 {array} types.CreateAccountRequest
// @Failure 400
func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := types.NewAccount(req.Name, req.Phone, req.Username, req.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// Handling account request method to specific request and response
// @Tags Account
// @Summary API Get Account
// @Router /account [get]
// @Param search query string false "Search query"
// @Param limit query int false "Limit results"
// @Param page query int false "Page number"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Account
// @Failure 400
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling account get request params, transaction, and response
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	var search string
	if r.FormValue("search") != "" {
		search = r.FormValue("search")
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	if err != nil {
		return err
	}

	accounts, count, totalPages, err := s.store.GetAccounts(search, page, limit)
	if err != nil {
		return err
	}

	responseData := types.Response{
		Page:       page,
		Limit:      limit,
		Count:      count,
		TotalPages: totalPages,
		Data:       accounts,
	}

	// return WriteJSON(w, http.StatusOK, accounts)
	return WriteJSON(w, http.StatusOK, responseData)
}

// Handling account get request by ID params, method, transaction, and response
// @Tags Account
// @Summary API Get Account by ID
// @Router /account/{id} [post]
// @Param id path int true "Account ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Account
// @Failure 400
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling account create required params, transaction, and response
// @Tags Account
// @Summary API Create Account
// @Router /account [post]
// @Param request body types.CreateAccountRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.CreateAccountRequest
// @Failure 400
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := types.NewAccount(req.Name, req.Phone, req.Username, req.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// Handling account update required params, transaction, and response
// @Tags Account
// @Summary API Update Account
// @Router /account/{id} [put]
// @Param id path int true "Account ID"
// @Param request body types.UpdateAccountRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.UpdateAccountRequest
// @Failure 400
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	req := new(types.UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := types.UpdateAccount(req.ID, req.Name, req.Phone, req.Password)
	if err != nil {
		return err
	}
	if err := s.store.UpdateAccount(id, account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"updated account": id})
}

// Handling account delete request by ID params, transaction, and response
// @Tags Account
// @Summary API Delete Account
// @Router /account/{id} [delete]
// @Param id path int true "Account ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Account
// @Failure 400
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted account": id})
}

// ---------------------- CUSTOMER ---------------------- //

// Handling customer request method to specific request and response
// @Tags Customer
// @Summary API Get Customer
// @Router /customer [get]
// @Param search query string false "Search query"
// @Param limit query int false "Limit results"
// @Param page query int false "Page number"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Customer
// @Failure 400
func (s *APIServer) handleCustomer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetCustomer(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateCustomer(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling customer get request params, transaction, and response
func (s *APIServer) handleGetCustomer(w http.ResponseWriter, r *http.Request) error {
	var search string
	if r.FormValue("search") != "" {
		search = r.FormValue("search")
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	if err != nil {
		return err
	}

	customers, count, totalPages, err := s.store.GetCustomers(search, page, limit)
	if err != nil {
		return err
	}

	responseData := types.Response{
		Page:       page,
		Limit:      limit,
		Count:      count,
		TotalPages: totalPages,
		Data:       customers,
	}

	// return WriteJSON(w, http.StatusOK, customers)
	return WriteJSON(w, http.StatusOK, responseData)
}

// Handling customer get request by ID params, method, transaction, and response
// @Tags Customer
// @Summary API Get Customer by ID
// @Router /customer/{id} [post]
// @Param id path int true "Customer ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Customer
// @Failure 400
func (s *APIServer) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		customer, err := s.store.GetCustomerByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, customer)
	}

	if r.Method == "PUT" {
		return s.handleUpdateCustomer(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteCustomer(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling customer create required params, transaction, and response
// @Tags Customer
// @Summary API Create Customer
// @Router /customer [post]
// @Param request body types.CreateCustomerRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.CreateCustomerRequest
// @Failure 400
func (s *APIServer) handleCreateCustomer(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	customer, err := types.NewCustomer(req.Name, req.Phone)
	if err != nil {
		return err
	}
	if err := s.store.CreateCustomer(customer); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, customer)
}

// Handling customer update required params, transaction, and response
// @Tags Customer
// @Summary API Update Customer
// @Router /customer/{id} [put]
// @Param id path int true "Customer ID"
// @Param request body types.UpdateCustomerRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.UpdateCustomerRequest
// @Failure 400
func (s *APIServer) handleUpdateCustomer(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	req := new(types.UpdateCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	customer, err := types.UpdateCustomer(req.ID, req.Name, req.Phone)
	if err != nil {
		return err
	}
	if err := s.store.UpdateCustomer(id, customer); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"updated customer": id})
}

// Handling customer delete request by ID params, transaction, and response
// @Tags Customer
// @Summary API Delete Customer
// @Router /customer/{id} [delete]
// @Param id path int true "Customer ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Customer
// @Failure 400
func (s *APIServer) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteCustomer(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted customer": id})
}

// ---------------------- ORDER ---------------------- //

// Handling order request method to specific request and response
// @Tags Order
// @Summary API Get Order
// @Router /order [get]
// @Param search query string false "Search query"
// @Param limit query int false "Limit results"
// @Param page query int false "Page number"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Order
// @Failure 400
func (s *APIServer) handleOrder(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetOrder(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateOrder(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling order get request params, transaction, and response
func (s *APIServer) handleGetOrder(w http.ResponseWriter, r *http.Request) error {
	var search string
	if r.FormValue("search") != "" {
		search = r.FormValue("search")
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	if err != nil {
		return err
	}

	orders, count, totalPages, err := s.store.GetOrders(search, page, limit)
	if err != nil {
		return err
	}

	responseData := types.Response{
		Page:       page,
		Limit:      limit,
		Count:      count,
		TotalPages: totalPages,
		Data:       orders,
	}

	// return WriteJSON(w, http.StatusOK, orders)
	return WriteJSON(w, http.StatusOK, responseData)
}

// Handling order get request by ID params, method, transaction, and response
// @Tags Order
// @Summary API Get Order by ID
// @Router /order/{id} [post]
// @Param id path int true "Order ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Order
// @Failure 400
func (s *APIServer) handleGetOrderByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		order, err := s.store.GetOrderByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, order)
	}

	if r.Method == "PUT" {
		return s.handleUpdateOrder(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteOrder(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

// Handling order create required params, transaction, and response
// @Tags Order
// @Summary API Create Order
// @Router /order [post]
// @Param request body types.CreateOrderRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.CreateOrderRequest
// @Failure 400
func (s *APIServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateOrderRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	order, err := types.NewOrder(req.Customer, req.Service, req.Amount, req.Unit, req.Price)
	if err != nil {
		return err
	}
	if err := s.store.CreateOrder(order); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, order)
}

// Handling order update required params, transaction, and response
// @Tags Order
// @Summary API Update Order
// @Router /order/{id} [put]
// @Param id path int true "Order ID"
// @Param request body types.UpdateOrderRequest true "Payload body [RAW]"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.UpdateOrderRequest
// @Failure 400
func (s *APIServer) handleUpdateOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	req := new(types.UpdateOrderRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	order, err := types.UpdateOrder(req.ID, req.Customer, req.Service, req.Amount, req.Unit, req.Price)
	if err != nil {
		return err
	}
	if err := s.store.UpdateOrder(id, order); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"updated order": id})
}

// Handling order delete request by ID params, transaction, and response
// @Tags Order
// @Summary API Delete Order
// @Router /order/{id} [delete]
// @Param id path int true "Order ID"
// @Param X-JWT-TOKEN header string true "JWT token"
// @Accept json
// @Produces json
// @Success 200 {array} types.Order
// @Failure 400
func (s *APIServer) handleDeleteOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteOrder(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted order": id})
}

// ---------------------- GLOBAL FUNCTIONS ---------------------- //

// Send back data transaction with JSON format
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// Create JWT based on account username, 15000 in unix expiration time, and JWT_SECRET in .env
func createJWT(account *types.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":       15000,
		"accountUsername": account.Username,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Handle denied request with default designated JSON response
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission denied"})
}

// Wrap secured APIs router with JWT authentication requirement
func withJWTAuth(handlerFunc http.HandlerFunc, s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("X-JWT-TOKEN")

		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		// Need to rework, validate only if change individual password or something else

		// userID, err := getID(r)
		// if err != nil {
		// 	permissionDenied(w)
		// 	return
		// }

		// account, err := s.GetAccountByID(userID)
		// if err != nil {
		// 	permissionDenied(w)
		// 	return
		// }

		// claims := token.Claims.(jwt.MapClaims)
		// if account.Username != claims["accountUsername"] {
		// 	permissionDenied(w)
		// 	return
		// }

		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

// Validate JWT token sent from request
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

// Initiate API function with response and request
type APIFunc func(http.ResponseWriter, *http.Request) error

// Initiate API error with JSON structure
type APIError struct {
	Error string `json:"error"`
}

// Handle HTTP API response and request
func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

// Convert ID sent by request to integer
func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid ID given %s", idStr)
	}

	return id, nil
}

// Convert Name sent by request to string
func getName(r *http.Request) (string, error) {
	name := mux.Vars(r)["name"]

	return name, nil
}
