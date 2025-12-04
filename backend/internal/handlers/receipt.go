package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	cld             *cloudinary.Cloudinary
}

type ReceiptHandlerConfig struct {
	ReceiptRepo     store.ReceiptRepo
	TransactionRepo store.TransactionRepo
	Cld             *cloudinary.Cloudinary
}

func NewReceiptHandler(cfg ReceiptHandlerConfig) ReceiptHandler {
	return ReceiptHandler{
		receiptRepo:     cfg.ReceiptRepo,
		transactionRepo: cfg.TransactionRepo,
		cld:             cfg.Cld,
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

	resp, err := ocr(media, header.Filename)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Terjadi kesalahan pada OCR",
		})
		return
	}

	ctx := r.Context()
	secureUrl, err := h.uploadImage(ctx, media)
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
	}

	if err := h.receiptRepo.Create(ctx, receipt); err != nil {
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

func (h *ReceiptHandler) uploadImage(ctx context.Context, image multipart.File) (string, error) {
	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			log.Printf("[ERROR] Failed to seek file: %v", err)
			return "", fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	resp, err := h.cld.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder:         "imphnen",
		UniqueFilename: api.Bool(true),
		UseFilename:    api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		log.Printf("[ERROR] Failed to upload photo to claudinary: %v", err)
		return "", fmt.Errorf("Failed to save photo: %s", err.Error())
	}

	return resp.SecureURL, nil
}

func ocr(image multipart.File, filename string) (models.OCR, error) {
	apiKey := os.Getenv("KOLOSAL_API_KEY")
	if apiKey == "" {
		return models.OCR{}, fmt.Errorf("KOLOSAL_API_KEY not set")
	}

	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return models.OCR{}, fmt.Errorf("failed to seek file: %w", err)
		}
	}

	buffer := make([]byte, 512)
	n, err := image.Read(buffer)
	if err != nil && err != io.EOF {
		return models.OCR{}, fmt.Errorf("failed to read file for type detection: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])

	if seeker, ok := image.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return models.OCR{}, fmt.Errorf("failed to seek file: %w", err)
		}
	}

	if !strings.HasPrefix(contentType, "image/") {
		return models.OCR{}, fmt.Errorf("file is not an image, detected type: %s", contentType)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, image); err != nil {
		return models.OCR{}, fmt.Errorf("failed to copy image data: %w", err)
	}

	if err := writer.WriteField("invoice", "true"); err != nil {
		return models.OCR{}, fmt.Errorf("failed to write invoice field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return models.OCR{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.kolosal.ai/ocr/form", &body)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.OCR{}, fmt.Errorf("OCR API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.OCR{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.OCR{}, fmt.Errorf("OCR API returned status %d: %s",
			resp.StatusCode, string(respBody))
	}

	var payload models.OCR
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return models.OCR{}, fmt.Errorf("failed to decode OCR response: %w, body: %s",
			err, string(respBody))
	}

	return payload, nil
}
