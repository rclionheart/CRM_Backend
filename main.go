package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

type CustomerInput struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type CustomerDeletion struct {
	ID int `json:"id"`
}

var customers = []Customer{
	{ID: 1, Name: "Bob", Role: "Member", Email: "bobsyouruncle@gmail.com", Phone: "0678304019", Contacted: true},
	{ID: 2, Name: "Scotty", Role: "Premium Member", Email: "scottsterling@gmail.com", Phone: "0683240193", Contacted: true},
	{ID: 3, Name: "Michael", Role: "New Member", Email: "michaeljakson@gmail.com", Phone: "0745320845", Contacted: false},
	{ID: 4, Name: "Andria", Role: "Member", Email: "andria@gmail.com", Phone: "0836891203", Contacted: false},
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var newID int = len(customers) + 1
	var input CustomerInput

	if err := readJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	if err := input.validateInput(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var newCustomer = []Customer{
		{ID: newID, Name: input.Name, Role: input.Role, Email: input.Email, Phone: input.Phone, Contacted: false},
	}

	customers = append(customers, newCustomer...)
	writeJSON(w, http.StatusCreated, customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	rawID := r.URL.Query().Get("id")

	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be positive integer")
		return
	}

	for _, customer := range customers {
		if customer.ID == id {
			writeJSON(w, http.StatusOK, customer)
			return
		}
	}

	writeError(w, http.StatusNotFound, "Customer not found")
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var update Customer

	if err := readJSON(w, r, &update); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
	}

	if err := update.validateUpdate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
	}

	for i := range customers {
		if customers[i].ID == update.ID {
			customers[i].Name = update.Name
			customers[i].Role = update.Role
			customers[i].Email = update.Email
			customers[i].Phone = update.Phone
			customers[i].Contacted = update.Contacted
			writeJSON(w, http.StatusOK, customers[i])
			return
		}
	}

	writeError(w, http.StatusNotFound, "Customer not found")
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	for i, customer := range customers {
		if customer.ID == id {
			customers = append(customers[:i], customers[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeError(w, http.StatusNotFound, "Customer not found")
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("body must contain exactly one JSON value")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		log.Printf("Write Response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func (c *CustomerInput) validateInput() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is missing")
	}
	if strings.TrimSpace(c.Role) == "" {
		return errors.New("role is missing")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("email is missing")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("phone is missing")
	}
	return nil
}

func (c *Customer) validateUpdate() error {
	if c.ID <= 0 {
		return errors.New("id is not positive/valid")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is missing")
	}
	if strings.TrimSpace(c.Role) == "" {
		return errors.New("role is missing")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("email is missing")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("phone is missing")
	}
	if !c.Contacted {
		return errors.New("contacted is missing")
	}
	return nil
}

func newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /customers", CreateCustomer)
	mux.HandleFunc("GET /customers", GetCustomer)
	mux.HandleFunc("PUT /customers", UpdateCustomer)
	mux.HandleFunc("DELETE /customer", DeleteCustomer)

	return mux
}

func main() {
	port := "8080"
	addr := ":" + port

	server := &http.Server{
		Addr:              addr,
		Handler:           newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("CRM API listening on http://localhost%s", addr)
	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
