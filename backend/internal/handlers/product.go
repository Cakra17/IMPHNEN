package handlers

import (
	"io"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/pkg/service"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productRepo store.ProductRepo
	cld         service.CloudinaryService
}

type ProductHandlerConfig struct {
	ProductRepo store.ProductRepo
	Cld         service.CloudinaryService
}

func NewProductHandler(cfg ProductHandlerConfig) ProductHandler {
	return ProductHandler{
		productRepo: cfg.ProductRepo,
		cld:         cfg.Cld,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with image upload
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Product image"
// @Param name formData string true "Product name"
// @Param price formData string true "Product price"
// @Param stock formData string true "Product stock"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=models.Product}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	priceVal := r.FormValue("price")
	stockVal := r.FormValue("stock")
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	price, err := strconv.ParseFloat(priceVal, 64)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format harga tidak sesuai",
		})
		return
	}

	stock, err := strconv.ParseInt(stockVal, 10, 64)
	if err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Format stock tidak sesuai",
		})
		return
	}

	secureUrl, publicID, err := h.cld.UploadMedia(ctx, "products", media)
	if err != nil {
		log.Printf("%s", err.Error())
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengupload gambar",
		})
		return
	}

	id, _ := uuid.NewV7()
	product := &models.Product{
		ID:       id.String(),
		UserID:   userID,
		Name:     name,
		Price:    price,
		Stock:    int(stock),
		ImageURL: secureUrl,
		PublicID: publicID,
	}

	err = h.productRepo.AddProduct(ctx, product)
	if err != nil {
		h.cld.DeleteMedia(ctx, publicID)
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Gagal menambahkan product",
		})
		return
	}

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil menambahkan product",
		Data:    product,
	})
}

// GetProducts godoc
// @Summary Get user products
// @Description Get paginated list of products for the authenticated user
// @Tags Product
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} utils.ResponsePaginate{data=models.ProductListResponse}
// @Failure 500 {object} utils.Response
// @Router /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)
	userID, _ := claims["user_id"].(string)

	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := uint(1)
	perPage := uint(10)

	if pageStr != "" {
		if p, err := strconv.ParseUint(pageStr, 10, 32); err == nil {
			page = uint(p)
		}
	}

	if perPageStr != "" {
		if pp, err := strconv.ParseUint(perPageStr, 10, 32); err == nil {
			perPage = uint(pp)
		}
	}

	products, totalCount, err := h.productRepo.GetUserProductsPaginated(ctx, userID, page, perPage)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mendapatkan produk",
		})
		return
	}

	totalPages := uint(math.Ceil(float64(totalCount) / float64(perPage)))

	utils.ResponseJson(w, http.StatusOK, utils.ResponsePaginate{
		Message: "Berhasil mengambil data receipt",
		Data: models.ProductListResponse{
			Products: products,
		},
		Meta: utils.Meta{
			Page:        page,
			TotalPage:   totalPages,
			TotalData:   totalCount,
			DataperPage: perPage,
		},
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.Product}
// @Failure 404 {object} utils.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.PathValue("id")

	product, err := h.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Produk tidak ditemukan",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mendapatkan produk",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product with optional image upload
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param image formData file false "Product image"
// @Param name formData string false "Product name"
// @Param price formData string false "Product price"
// @Param stock formData string false "Product stock"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.Product}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.PathValue("id")

	existingProduct, err := h.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Produk tidak ditemukan",
		})
		return
	}
	oldUrl := existingProduct.ImageURL

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

	name := r.FormValue("name")
	priceVal := r.FormValue("price")
	stockVal := r.FormValue("stock")

	if name != "" {
		existingProduct.Name = name
	}

	if priceVal != "" {
		price, err := strconv.ParseFloat(priceVal, 64)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format harga tidak sesuai",
			})
			return
		}
		existingProduct.Price = price
	}

	if stockVal != "" {
		stock, err := strconv.ParseInt(stockVal, 10, 64)
		if err != nil {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format stock tidak sesuai",
			})
			return
		}
		existingProduct.Stock = int(stock)
	}

	media, header, err := r.FormFile("image")
	if err == nil {
		defer media.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedType[ext] {
			utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
				Message: "Format tidak sesuai",
			})
			return
		}

		secureUrl, publicID, err := h.cld.UploadMedia(ctx, "products", media)
		if err != nil {
			log.Printf("%s", err.Error())
			utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
				Message: "Gagal mengupload gambar",
			})
			return
		}
		if err := h.cld.DeleteMedia(ctx, existingProduct.PublicID); err != nil {
			log.Printf("%s", err.Error())
			utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
				Message: "Gagal mengdate gambar",
			})
			return
		}

		existingProduct.ImageURL = secureUrl
		existingProduct.PublicID = publicID
	}

	err = h.productRepo.UpdateProduct(ctx, existingProduct)
	if err != nil {
		if oldUrl != existingProduct.ImageURL {
			h.cld.DeleteMedia(ctx, existingProduct.PublicID)
		}
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengupdate produk",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil mengupdate produk",
		Data:    existingProduct,
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.PathValue("id")

	product, err := h.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Produk tidak ditemukan",
		})
		return
	}

	err = h.productRepo.DeleteProduct(ctx, productID)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal menghapus produk",
		})
		return
	}

	if err := h.cld.DeleteMedia(ctx, product.PublicID); err != nil {
		log.Printf("%s", err.Error())
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal mengdate gambar",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Berhasil menghapus produk",
	})
}
