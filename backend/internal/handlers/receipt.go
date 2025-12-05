package handlers

import (
	"database/sql"
	"io"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/pkg/service"
	"github.com/google/uuid"
)

const (
	MaxUploadSize = 5 << 20
)

var allowedType map[string]bool = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

type ReceiptHandler struct {
	receiptRepo     store.ReceiptRepo
	transactionRepo store.TransactionRepo
	cld             service.CloudinaryService
	kol             service.KolosalService
}

type ReceiptHandlerConfig struct {
	ReceiptRepo     store.ReceiptRepo
	TransactionRepo store.TransactionRepo
	Cld             service.CloudinaryService
	Kol             service.KolosalService
}

func NewReceiptHandler(cfg ReceiptHandlerConfig) ReceiptHandler {
	return ReceiptHandler{
		receiptRepo:     cfg.ReceiptRepo,
		transactionRepo: cfg.TransactionRepo,
		cld:             cfg.Cld,
		kol:             cfg.Kol,
	}
}

// CreateReceipt godoc
// @Summary      Create a new receipt
// @Description  Create a new receipt with items for the authenticated user
// @Tags         Receipts
// @Accept       x-www-form-urlencoded
// @Param        image   formData  image  true  "receipt to scan"
// @Produce      json
// @Security     BearerAuth
// @Success      201      {object}  utils.Response{data=models.ReceiptResponse}  "Receipt created successfully"
// @Failure      400      {object}  utils.Response{message=string}  "Invalid request data"
// @Failure      401      {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500      {object}  utils.Response{message=string}  "Internal server error"
// @Router       /receipts [post]
func (h *ReceiptHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		if err == io.ErrUnexpectedEOF {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Ukuran data terlalu besar, Max 5 MB",
			})
			return
		}

		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Gagal membaca data",
		})
		return
	}

	media, header, err := r.FormFile("image")
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Gagal menerima gambar",
		})
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format tidak sesuai",
		})
		return
	}

	ctx := r.Context()

	resp, err := h.kol.OCRForm(media, header.Filename)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Terjadi kesalahan pada OCR",
		})
		return
	}

	secureUrl, publidID, err := h.cld.UploadMedia(ctx, "receipts", media)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Terjadi kesalahan dalam menyimpan gambar",
		})
		return
	}

	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	receiptID, _ := uuid.NewV7()
	receipt := models.Receipt{
		ID:         receiptID.String(),
		UserID:     userID,
		StoreName:  resp.Issuer.Name,
		TotalItems: uint32(len(resp.InvoiceItems)),
		TotalPrice: resp.Total,
		ImageURL:   secureUrl,
		PublicID:   publidID,
	}

	if err := h.receiptRepo.Create(ctx, &receipt); err != nil {
		h.cld.DeleteMedia(ctx, publidID)
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat receipt",
		})
		return
	}

	// Create receipt items
	items := []models.ReceiptItem{}
	for _, itemPayload := range resp.InvoiceItems {
		itemID, _ := uuid.NewV7()
		item := models.ReceiptItem{
			ID:        itemID.String(),
			ReceiptID: receiptID.String(),
			Name:      itemPayload.Description,
			Price:     itemPayload.Total,
		}
		items = append(items, item)
	}

	if err := h.receiptRepo.CreateItems(ctx, items); err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat item receipt",
		})
		return
	}

	// Create Transaction
	tscID, _ := uuid.NewV7()
	date, err := time.Parse("2006-01-02", resp.InvoiceDate)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat transaksi",
		})
		return
	}

	tsc := models.Transaction{
		ID:              tscID.String(),
		ReceiptID:       &receipt.ID,
		UserID:          userID,
		Type:            "expense",
		Source:          "receipt",
		Amount:          resp.Total,
		OrderID:         nil,
		TransactionDate: date,
	}

	if err := h.transactionRepo.AddTransaction(ctx, &tsc); err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat transaksi",
		})
		return
	}

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil membuat receipt",
		Data: models.ReceiptResponse{
			Receipt: receipt,
			Items:   items,
		},
	})
}

// GetReceipts godoc
// @Summary      Get paginated receipts
// @Description  Get a paginated list of receipts for the authenticated user
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int  false  "Page number (default: 1)"
// @Param        per_page  query     int  false  "Items per page (default: 10)"
// @Success      200       {object}  utils.ResponsePaginate{data=models.ReceiptListResponse}  "Receipts retrieved successfully"
// @Failure      401       {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      500       {object}  utils.Response{message=string}  "Internal server error"
// @Router       /receipts [get]
func (h *ReceiptHandler) GetReceipts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	// Parse pagination parameters
	page := 1
	perPage := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}

	receipts, totalCount, err := h.receiptRepo.GetReceiptsPaginate(ctx, userID, uint(page), uint(perPage))
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data receipt",
		})
		return
	}

	totalPages := uint(math.Ceil(float64(totalCount) / float64(perPage)))

	utils.ResponseJson(w, http.StatusOK, utils.ResponsePaginate{
		Message: "Berhasil mengambil data receipt",
		Data: models.ReceiptListResponse{
			Receipts: receipts,
		},
		Meta: utils.Meta{
			Page:        uint(page),
			TotalPage:   totalPages,
			TotalData:   totalCount,
			DataperPage: uint(perPage),
		},
	})

}

// GetReceiptByID godoc
// @Summary      Get receipt by ID
// @Description  Get a specific receipt with its items by ID for the authenticated user
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Receipt ID"
// @Success      200  {object}  utils.Response{data=models.ReceiptResponse}  "Receipt retrieved successfully"
// @Failure      401  {object}  utils.Response{message=string}  "Unauthorized"
// @Failure      404  {object}  utils.Response{message=string}  "Receipt not found"
// @Failure      500  {object}  utils.Response{message=string}  "Internal server error"
// @Router       /receipts/{id} [get]
func (h *ReceiptHandler) GetReceiptByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	receiptID := r.PathValue("id")
	if receiptID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Receipt ID tidak valid",
		})
		return
	}

	receipt, err := h.receiptRepo.GetReceiptByID(ctx, receiptID, userID)
	if err == sql.ErrNoRows {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Receipt tidak ditemukan",
		})
		return
	}
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil data receipt",
		})
		return
	}

	items, err := h.receiptRepo.GetReceiptItemsByReceiptID(ctx, receiptID)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil item receipt",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data receipt",
		Data: models.ReceiptResponse{
			Receipt: *receipt,
			Items:   items,
		},
	})
}

func (h *ReceiptHandler) GetItemsByRecieptID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receiptID := r.PathValue("id")
	if receiptID == "" {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Receipt ID tidak valid",
		})
		return
	}

	items, err := h.receiptRepo.GetReceiptItemsByReceiptID(ctx, receiptID)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengambil item receipt",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengambil data receipt",
		Data: models.ItemsResponse{
			Items: items,
		},
	})
}
