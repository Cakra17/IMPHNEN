package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/internal/validation"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionStore *store.TransactionRepo
}

type TransactionHandlerConfig struct {
	TransactionStore *store.TransactionRepo
}

func NewTransactionHandler(cfg TransactionHandlerConfig) TransactionHandler {
	return TransactionHandler{
		transactionStore: cfg.TransactionStore,
	}
}

// CreateTransaction godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateTransactionPayload  true  "Transaction details"
// @Success      201      {object}  utils.Response{data=models.TransactionResponse}  "Transaction created successfully"
// @Failure      400      {object}  utils.Response{message=string}  "Invalid request data"
// @Failure      401      {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500      {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateTransactionPayload
	ctx := r.Context()

	if err := utils.ParseJson(r, &payload); err != nil {
		log.Println(err)
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

	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	tscID, _ := uuid.NewV7()
	date, err := time.Parse("2006-01-02", payload.TransactionDate)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat transaksi",
		})
		return
	}

	transaction := models.Transaction{
		ID:              tscID.String(),
		UserID:          userID,
		Type:            payload.Type,
		Source:          payload.Source,
		Amount:          payload.Amount,
		TransactionDate: date,
		ReceiptID:       payload.ReceiptID,
		OrderID:         payload.OrderID,
	}

	if err := h.transactionStore.AddTransaction(ctx, &transaction); err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat transaksi",
		})
		return
	}

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil membuat transaksi",
		Data: models.TransactionResponse{
			Transaction: transaction,
		},
	})
}

// GetTransactionsByDate godoc
// @Summary      Get transactions by date
// @Description  Get all transactions for a specific date for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        date   query     string  true  "Date in YYYY-MM-DD format"
// @Success      200    {object}  utils.Response{data=models.TransactionListResponse}  "Transactions retrieved successfully"
// @Failure      400    {object}  utils.Response{message=string}  "Invalid date format"
// @Failure      401    {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500    {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/date [get]
func (h *TransactionHandler) GetTransactionsByDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter date diperlukan",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format date tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}

	transactions, err := h.transactionStore.GetTransactionsByDate(ctx, userID, date)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data transaksi",
		})
		return
	}

	transactionList := make([]models.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = t
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data transaksi",
		Data: models.TransactionListResponse{
			Transactions: transactionList,
		},
	})
}

// GetTransactionsByRange godoc
// @Summary      Get transactions by date range
// @Description  Get all transactions within a date range for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date  query     string  true  "Start date in YYYY-MM-DD format"
// @Param        end_date    query     string  true  "End date in YYYY-MM-DD format"
// @Success      200         {object}  utils.Response{data=models.TransactionListResponse}  "Transactions retrieved successfully"
// @Failure      400         {object}  utils.Response{message=string}  "Invalid date format"
// @Failure      401         {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500         {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/range [get]
func (h *TransactionHandler) GetTransactionsByRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter start_date dan end_date diperlukan",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format start_date tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format end_date tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}

	// Add time to end date to include the entire day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	transactions, err := h.transactionStore.GetTransactionsByRange(ctx, userID, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data transaksi",
		})
		return
	}

	transactionList := make([]models.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = t
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data transaksi",
		Data: models.TransactionListResponse{
			Transactions: transactionList,
		},
	})
}

// GetTransactionsByDays godoc
// @Summary      Get transactions for last N days
// @Description  Get all transactions for the last N days for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  true  "Number of days (e.g., 7, 30, 90)"
// @Success      200   {object}  utils.Response{data=models.TransactionListResponse}  "Transactions retrieved successfully"
// @Failure      400   {object}  utils.Response{message=string}  "Invalid days parameter"
// @Failure      401   {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500   {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/days [get]
func (h *TransactionHandler) GetTransactionsByDays(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	daysStr := r.URL.Query().Get("days")
	if daysStr == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter days diperlukan",
		})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter days harus berupa angka positif",
		})
		return
	}

	transactions, err := h.transactionStore.GetTransactionsByDays(ctx, userID, days)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data transaksi",
		})
		return
	}

	transactionList := make([]models.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = t
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data transaksi",
		Data: models.TransactionListResponse{
			Transactions: transactionList,
		},
	})
}

// GetTransactionStats godoc
// @Summary      Get transaction statistics
// @Description  Get transaction statistics for a date range for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date  query     string  false  "Start date in YYYY-MM-DD format (default: 30 days ago)"
// @Param        end_date    query     string  false  "End date in YYYY-MM-DD format (default: today)"
// @Success      200         {object}  utils.Response{data=models.TransactionStats}  "Statistics retrieved successfully"
// @Failure      400         {object}  utils.Response{message=string}  "Invalid date format"
// @Failure      401         {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500         {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/stats [get]
func (h *TransactionHandler) GetTransactionStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		start, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format start_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		startDate = start
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Default 30 days ago
	}

	if endDateStr != "" {
		end, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format end_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		endDate = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now() // Default today
	}

	stats, err := h.transactionStore.GetTransactionStats(ctx, userID, startDate, endDate)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil statistik transaksi",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil statistik transaksi",
		Data:    stats,
	})
}

// GetTransactionStatsByDays godoc
// @Summary      Get transaction statistics for last N days
// @Description  Get transaction statistics for the last N days for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Number of days (default: 30)"
// @Success      200   {object}  utils.Response{data=models.TransactionStats}  "Statistics retrieved successfully"
// @Failure      400   {object}  utils.Response{message=string}  "Invalid days parameter"
// @Failure      401   {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500   {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/stats/days [get]
func (h *TransactionHandler) GetTransactionStatsByDays(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	daysStr := r.URL.Query().Get("days")
	days := 30 // Default 30 days

	if daysStr != "" {
		d, err := strconv.Atoi(daysStr)
		if err != nil || d <= 0 {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Parameter days harus berupa angka positif",
			})
			return
		}
		days = d
	}

	stats, err := h.transactionStore.GetTransactionStatsByDays(ctx, userID, days)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil statistik transaksi",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil statistik transaksi",
		Data:    stats,
	})
}

// GetTransactionsByType godoc
// @Summary      Get transactions by type
// @Description  Get transactions by type for a date range for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        type       query     string  true  "Transaction type (income/expense)"
// @Param        start_date query     string  false  "Start date in YYYY-MM-DD format (default: 30 days ago)"
// @Param        end_date   query     string  false  "End date in YYYY-MM-DD format (default: today)"
// @Success      200        {object}  utils.Response{data=models.TransactionListResponse}  "Transactions retrieved successfully"
// @Failure      400        {object}  utils.Response{message=string}  "Invalid parameters"
// @Failure      401        {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500        {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/type [get]
func (h *TransactionHandler) GetTransactionsByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	transactionType := r.URL.Query().Get("type")
	if transactionType == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter type diperlukan (income/expense)",
		})
		return
	}

	if transactionType != "income" && transactionType != "expense" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter type harus 'income' atau 'expense'",
		})
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		start, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format start_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		startDate = start
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Default 30 days ago
	}

	if endDateStr != "" {
		end, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format end_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		endDate = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now() // Default today
	}

	fmt.Println(startDate, endDate, transactionType)

	transactions, err := h.transactionStore.GetTransactionsByType(ctx, userID, transactionType, startDate, endDate)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data transaksi",
		})
		return
	}

	transactionList := make([]models.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = t
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data transaksi",
		Data: models.TransactionListResponse{
			Transactions: transactionList,
		},
	})
}

// GetTransactionsBySource godoc
// @Summary      Get transactions by source
// @Description  Get transactions by source for a date range for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        source     query     string  true  "Transaction source (receipt/whatsapp/manual)"
// @Param        start_date query     string  false  "Start date in YYYY-MM-DD format (default: 30 days ago)"
// @Param        end_date   query     string  false  "End date in YYYY-MM-DD format (default: today)"
// @Success      200        {object}  utils.Response{data=models.TransactionListResponse}  "Transactions retrieved successfully"
// @Failure      400        {object}  utils.Response{message=string}  "Invalid parameters"
// @Failure      401        {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500        {object}  utils.Response{message=string}  "Internal server error"
// @Router       /transactions/source [get]
func (h *TransactionHandler) GetTransactionsBySource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	source := r.URL.Query().Get("source")
	if source == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter source diperlukan (receipt/bot/manual)",
		})
		return
	}

	if source != "receipt" && source != "bot" && source != "manual" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Parameter source harus 'receipt', 'bot', atau 'manual'",
		})
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		start, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format start_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		startDate = start
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Default 30 days ago
	}

	if endDateStr != "" {
		end, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format end_date tidak valid, gunakan YYYY-MM-DD",
			})
			return
		}
		endDate = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now() // Default today
	}

	transactions, err := h.transactionStore.GetTransactionsBySource(ctx, userID, source, startDate, endDate)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data transaksi",
		})
		return
	}

	transactionList := make([]models.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = t
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data transaksi",
		Data: models.TransactionListResponse{
			Transactions: transactionList,
		},
	})
}
