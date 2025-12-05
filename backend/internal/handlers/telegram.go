package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/internal/validation"
	"github.com/google/uuid"
)

// TelegramHandler handles telegram bot operations for customers/clients
// Allows customers to: list products, create orders, list their orders, and cancel/delete orders
type TelegramHandler struct {
	orderRepo   store.OrderRepo
	productRepo store.ProductRepo
}

type TelegramHandlerConfig struct {
	OrderRepo   store.OrderRepo
	ProductRepo store.ProductRepo
}

func NewTelegramHandler(cfg TelegramHandlerConfig) TelegramHandler {
	return TelegramHandler{
		orderRepo:   cfg.OrderRepo,
		productRepo: cfg.ProductRepo,
	}
}

// ListProductsByMerchant lists all products from a specific merchant
// @Summary List products by merchant (for Telegram bot)
// @Description Get all products available from a specific merchant
// @Tags telegram
// @Accept json
// @Produce json
// @Param merchant_id path string true "Merchant User ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.ResponsePaginate{data=[]models.Product,meta=utils.Meta{}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /telegram/merchants/{merchant_id}/products [get]
func (h *TelegramHandler) ListProductsByMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	merchantID := r.PathValue("merchant_id")
	if merchantID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Merchant ID diperlukan",
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

	products, total, err := h.productRepo.GetUserProductsPaginated(ctx, merchantID, page, perPage)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan produk",
		})
		return
	}

	totalPage := uint(math.Ceil(float64(total) / float64(perPage)))

	utils.ResponseJson(w, http.StatusOK, utils.ResponsePaginate{
		Message: "Berhasil mendapatkan produk",
		Data:    products,
		Meta: utils.Meta{
			Page:        page,
			TotalData:   total,
			DataperPage: perPage,
			TotalPage:   totalPage,
		},
	})
}

// CreateOrderForCustomer creates a new order from telegram bot (customer side)
// @Summary Create order for customer (Telegram bot)
// @Description Create a new order for a customer via telegram bot
// @Tags telegram
// @Accept json
// @Produce json
// @Param order body models.CreateTelegramOrderRequest true "Order details"
// @Success 201 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /telegram/orders [post]
func (h *TelegramHandler) CreateOrderForCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload models.CreateTelegramOrderRequest
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Data yang dikirim tidak sesuai",
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

	if len(payload.Items) > 0 {
		firstProduct, err := h.productRepo.GetProductByID(ctx, payload.Items[0].ProductID)
		if err != nil {
			utils.ResponseJson(w, http.StatusNotFound, utils.Response{
				Message: "Produk tidak ditemukan",
			})
			return
		}
		payload.MerchantID = firstProduct.UserID
	}

	customer, err := h.orderRepo.GetCustomerByID(ctx, payload.CustomerID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Customer tidak ditemukan",
		})
		return
	}

	orderID, _ := uuid.NewV7()
	now := time.Now()
	order := &models.Order{
		ID:         orderID.String(),
		UserID:     payload.MerchantID,
		CustomerID: strconv.Itoa(payload.CustomerID),
		Status:     models.OrderStatusPending,
		OrderDate:  now,
		CreatedAt:  now,
	}

	orderItems := make([]models.OrderItem, len(payload.Items))
	for i, item := range payload.Items {
		itemID, _ := uuid.NewV7()
		orderItems[i] = models.OrderItem{
			ID:        itemID.String(),
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CreatedAt: now,
		}
	}

	err = h.orderRepo.CreateOrder(ctx, order, orderItems)
	if err != nil {
		if err.Error() == "customer not found" {
			utils.ResponseJson(w, http.StatusNotFound, utils.Response{
				Message: "Customer tidak ditemukan",
			})
			return
		}
		if err.Error() == "one or more products not found or don't belong to user" {
			utils.ResponseJson(w, http.StatusNotFound, utils.Response{
				Message: "Satu atau lebih produk tidak ditemukan atau bukan milik merchant",
			})
			return
		}
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: err.Error(),
		})
		return
	}

	order.Customer = customer

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil membuat pesanan",
		Data:    order,
	})
}

// ListCustomerOrders lists all orders for a specific customer
// @Summary List customer orders (Telegram bot)
// @Description Get all orders for a specific customer
// @Tags telegram
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.ResponsePaginate{data=[]models.Order,meta=utils.Meta{}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /telegram/customers/{customer_id}/orders [get]
func (h *TelegramHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	filter := models.OrderFilter{
		CustomerID: &customerID,
		Page:       page,
		PerPage:    perPage,
	}

	orders, total, err := h.orderRepo.GetOrdersByCustomerOnly(ctx, filter)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan pesanan",
		})
		return
	}

	totalPage := uint(math.Ceil(float64(total) / float64(perPage)))

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

// CancelCustomerOrder cancels an order (customer side)
// @Summary Cancel customer order (Telegram bot)
// @Description Cancel an order and restore stock
// @Tags telegram
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /telegram/orders/{order_id}/cancel [patch]
func (h *TelegramHandler) CancelCustomerOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := r.PathValue("order_id")
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

	if order.Status == models.OrderStatusCancelled {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Pesanan sudah dibatalkan",
		})
		return
	}

	err = h.orderRepo.UpdateOrderStatus(ctx, orderID, models.OrderStatusCancelled)
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
		Message: "Berhasil membatalkan pesanan",
		Data:    updatedOrder,
	})
}

// DeleteCustomerOrder deletes an order (customer side)
// @Summary Delete customer order (Telegram bot)
// @Description Delete a pending or cancelled order
// @Tags telegram
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /telegram/orders/{order_id} [delete]
func (h *TelegramHandler) DeleteCustomerOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := r.PathValue("order_id")
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

	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusCancelled {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Hanya pesanan pending atau cancelled yang dapat dihapus",
		})
		return
	}

	err = h.orderRepo.DeleteOrder(ctx, orderID)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: err.Error(),
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil menghapus pesanan",
	})
}
