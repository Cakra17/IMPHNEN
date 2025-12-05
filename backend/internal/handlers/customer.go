package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/internal/validation"
)

type CustomerHandler struct {
	customerRepo store.CustomerRepo
}

type CustomerHandlerConfig struct {
	CustomerRepo store.CustomerRepo
}

func NewCustomerHandler(cfg CustomerHandlerConfig) CustomerHandler {
	return CustomerHandler{
		customerRepo: cfg.CustomerRepo,
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer with the provided details
// @Tags Customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=models.Customer}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload models.CreateCustomerRequest
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format data tidak valid",
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		if errs, ok := err.(validation.ValidationErrors); ok {
			for _, e := range errs {
				utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
					Message: fmt.Sprintf("%s %s", e.Field, e.Message),
				})
				return
			}
		}
	}

	customer := &models.Customer{
		ID:      payload.ID,
		Name:    payload.Name,
		Address: payload.Address,
		Phone:   payload.Phone,
	}

	err := h.customerRepo.CreateCustomer(ctx, customer)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal menambahkan customer",
		})
		return
	}

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil menambahkan customer",
		Data:    customer,
	})
}

// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Get a specific customer by their ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.Customer}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerIDStr := r.PathValue("id")

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "ID customer tidak valid",
		})
		return
	}

	customer, err := h.customerRepo.GetCustomerByID(ctx, customerID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Customer tidak ditemukan",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mendapatkan customer",
		Data:    customer,
	})
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update an existing customer's details
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body models.Customer true "Customer details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.Customer}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerIDStr := r.PathValue("id")

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "ID customer tidak valid",
		})
		return
	}

	// Check if customer exists
	_, err = h.customerRepo.GetCustomerByID(ctx, customerID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Customer tidak ditemukan",
		})
		return
	}

	var customer models.Customer
	if err := utils.ParseJson(r, &customer); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format data tidak valid",
		})
		return
	}

	customer.ID = customerID

	err = h.customerRepo.UpdateCustomer(ctx, customer)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengupdate customer",
		})
		return
	}

	// Fetch updated customer
	updatedCustomer, _ := h.customerRepo.GetCustomerByID(ctx, customerID)

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengupdate customer",
		Data:    updatedCustomer,
	})
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer by their ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerIDStr := r.PathValue("id")

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "ID customer tidak valid",
		})
		return
	}

	// Check if customer exists
	_, err = h.customerRepo.GetCustomerByID(ctx, customerID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Customer tidak ditemukan",
		})
		return
	}

	err = h.customerRepo.DeleteCustomer(ctx, customerID)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal menghapus customer",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil menghapus customer",
	})
}
