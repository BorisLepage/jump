package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"src/facture/internal/receipts"
	"src/helper_api"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/user/{id}", GetUser)
	r.Post("/invoices", CreateInvoice)
	r.Put("/invoices/{id}/_paid", MarkInvoiceAsPaid)

	fmt.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		helper_api.SendErrorResponse(w, "bad_request", err.Error(), http.StatusBadRequest)
		return
	}

	user, err := receipts.GetUser(userIDInt)
	if errors.Is(err, receipts.UserNotFound) {
		helper_api.SendErrorResponse(w, "not_found", err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		helper_api.SendErrorResponse(w, "internal_server_error", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		helper_api.SendErrorResponse(w, "internal_server_error", err.Error(), http.StatusInternalServerError)
		return
	}
}

type CreateInvoicePayload receipts.CreateInvoiceRequest

func (p CreateInvoicePayload) Validate() error {
	if p.UserID <= 0 {
		return errors.New("user_id must be greater than 0")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if p.Label == "" {
		return errors.New("label is required")
	}
	return nil
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoicePayload

	err := helper_api.ReadAndValidate(r, &req)
	if err != nil {
		helper_api.SendErrorResponse(w, "bad_request", err.Error(), http.StatusBadRequest)
		return
	}

	invoice, err := receipts.CreateInvoice(receipts.CreateInvoiceRequest(req))
	if errors.Is(err, receipts.UserNotFound) {
		helper_api.SendErrorResponse(w, "not_found", "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		helper_api.SendErrorResponse(w, "internal_server_error", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func MarkInvoiceAsPaid(w http.ResponseWriter, r *http.Request) {
	invoiceIDStr := chi.URLParam(r, "id")

	invoiceID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		helper_api.SendErrorResponse(w, "bad_request", "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	invoice, err := receipts.MarkInvoiceAsPaid(invoiceID)
	if errors.Is(err, receipts.InvoiceNotFound) {
		helper_api.SendErrorResponse(w, "not_found", "Invoice not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, receipts.InvoiceAlreadyPaid) {
		helper_api.SendErrorResponse(w, "conflict", "Invoice is already paid", http.StatusConflict)
		return
	}
	if err != nil {
		helper_api.SendErrorResponse(w, "internal_server_error", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}
