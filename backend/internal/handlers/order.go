package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
)

// OrderHandler handles merchant-specific order operations
// Merchants can only see and manage orders that belong to them
type OrderHandler struct {
	orderRepo store.OrderRepo
}

type OrderHandlerConfig struct {
	OrderRepo store.OrderRepo
}

func NewOrderHandler(cfg OrderHandlerConfig) OrderHandler {
	return OrderHandler{
		orderRepo: cfg.OrderRepo,
	}
}

// GetOrders retrieves orders for the authenticated user (merchant) with filtering
// @Summary Get all orders for merchant
// @Description Retrieve orders with optional filtering by customer and status
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param customer_id query string false "Filter by customer ID"
// @Param status query string false "Filter by status (pending, confirmed, cancelled)"
// @Success 200 {object} utils.ResponsePaginate{data=[]models.Order}
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
			Message: "Unauthorized",
		})
		return
	}
	userID, _ := claims["user_id"].(string)

	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")
	customerIDStr := r.URL.Query().Get("customer_id")
	statusStr := r.URL.Query().Get("status")

	page := uint(1)
	perPage := uint(10)

	if pageStr != "" {
		if p, err := strconv.ParseUint(pageStr, 10, 32); err == nil && p > 0 {
			page = uint(p)
		}
	}

	if perPageStr != "" {
		if pp, err := strconv.ParseUint(perPageStr, 10, 32); err == nil && pp > 0 {
			perPage = uint(pp)
		}
	}

	filter := models.OrderFilter{
		UserID:  userID,
		Page:    page,
		PerPage: perPage,
	}

	if customerIDStr != "" {
		filter.CustomerID = &customerIDStr
	}

	if statusStr != "" {
		status := models.OrderStatus(statusStr)
		if status == models.OrderStatusPending ||
			status == models.OrderStatusConfirmed ||
			status == models.OrderStatusCancelled {
			filter.Status = &status
		}
	}

	orders, total, err := h.orderRepo.GetOrders(ctx, filter)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan pesanan",
		})
		return
	}

	totalPage := uint(math.Ceil(float64(total / perPage)))

	utils.ResponseJson(w, http.StatusOK, utils.ResponsePaginate{
		Message: "Berhasil mendapatkan pesanan",
		Data:    orders,
		Meta: utils.Meta{
			Page:        page,
			TotalData:   total,
			DataperPage: perPage,
			TotalPage:   totalPage,
		},
	})
}

// GetOrderByID retrieves a specific order by ID
// @Summary Get order by ID
// @Description Retrieve a specific order with its items and customer details
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order}
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
			Message: "Unauthorized",
		})
		return
	}
	userID, _ := claims["user_id"].(string)

	orderID := r.PathValue("id")
	if orderID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Order ID diperlukan",
		})
		return
	}

	order, err := h.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.ResponseJson(w, http.StatusNotFound, utils.Response{
				Message: "Pesanan tidak ditemukan",
			})
			return
		}
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan pesanan",
		})
		return
	}

	if order.UserID != userID {
		utils.ResponseJson(w, http.StatusForbidden, utils.Response{
			Message: "Tidak memiliki akses ke pesanan ini",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mendapatkan pesanan",
		Data:    order,
	})
}

// GetOrdersByCustomer retrieves all orders for a specific customer
// @Summary Get orders by customer
// @Description Retrieve all orders for a specific customer
// @Tags orders
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=models.OrderListResponse}
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /orders/customer/{customer_id} [get]
func (h *OrderHandler) GetOrdersByCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
			Message: "Unauthorized",
		})
		return
	}
	userID, _ := claims["user_id"].(string)

	customerID := r.PathValue("customer_id")
	if customerID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Customer ID diperlukan",
		})
		return
	}

	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := uint(1)
	perPage := uint(10)

	if pageStr != "" {
		if p, err := strconv.ParseUint(pageStr, 10, 32); err == nil && p > 0 {
			page = uint(p)
		}
	}

	if perPageStr != "" {
		if pp, err := strconv.ParseUint(perPageStr, 10, 32); err == nil && pp > 0 {
			perPage = uint(pp)
		}
	}

	orders, total, err := h.orderRepo.GetOrdersByCustomer(ctx, userID, customerID, page, perPage)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan pesanan",
		})
		return
	}

	totalPage := uint(math.Ceil(float64(total / perPage)))

	utils.ResponseJson(w, http.StatusOK, utils.ResponsePaginate{
		Message: "Berhasil mendapatkan pesanan untuk customer",
		Data:    orders,
		Meta: utils.Meta{
			Page:        page,
			TotalData:   total,
			DataperPage: perPage,
			TotalPage:   totalPage,
		},
	})
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Update order status (pending, confirmed, cancelled) with automatic stock restoration
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Security BearerAuth
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
			Message: "Unauthorized",
		})
		return
	}
	userID, _ := claims["user_id"].(string)

	orderID := r.PathValue("id")
	if orderID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Order ID diperlukan",
		})
		return
	}

	order, err := h.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.ResponseJson(w, http.StatusNotFound, utils.Response{
				Message: "Pesanan tidak ditemukan",
			})
			return
		}
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan pesanan",
		})
		return
	}

	if order.UserID != userID {
		utils.ResponseJson(w, http.StatusForbidden, utils.Response{
			Message: "Tidak memiliki akses ke pesanan ini",
		})
		return
	}

	var payload models.UpdateOrderStatusRequest
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Data yang dikirim tidak sesuai",
		})
		return
	}

	if payload.Status != models.OrderStatusPending &&
		payload.Status != models.OrderStatusConfirmed &&
		payload.Status != models.OrderStatusCancelled {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Status tidak valid",
		})
		return
	}

	err = h.orderRepo.UpdateOrderStatus(ctx, orderID, payload.Status)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: err.Error(),
		})
		return
	}

	updatedOrder, err := h.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan data pesanan terbaru",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengupdate status pesanan",
		Data:    updatedOrder,
	})
}
