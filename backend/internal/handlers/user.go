package handlers

import (
	"net/http"
	"time"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo      store.UserRepo
	jwtSecret     string
	tokenDuration time.Duration
}

type UserHandlerConfig struct {
	UserRepo      store.UserRepo
	JwtSecret     string
	TokenDuration time.Duration
}

func NewUserHandler(cfg UserHandlerConfig) UserHandler {
	return UserHandler{
		userRepo:      cfg.UserRepo,
		jwtSecret:     cfg.JwtSecret,
		tokenDuration: cfg.TokenDuration,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload models.AuthPayload
	ctx := r.Context()

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Data yang dikirim tidak sesuai",
		})
		return
	}

	userExist, _ := h.userRepo.GetUserbyEmail(ctx, payload.Email)
	if userExist != nil {
		utils.ResponseJson(w, http.StatusConflict, utils.Response{
			Message: "Email sudah terdaftar",
		})
		return
	}

	id, _ := uuid.NewV7()
	hashedPassword, _ := hashPassword(payload.Password)

	user := models.User{
		ID:           id.String(),
		Email:        payload.Email,
		PasswordHash: hashedPassword,
	}

	err := h.userRepo.Create(ctx, user)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: err.Error(),
		})
		return
	}

	utils.ResponseJson(w, http.StatusCreated, utils.Response{
		Message: "Berhasil membuat akun",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.AuthPayload
	ctx := r.Context()

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Data yang dikirim tidak sesuai",
		})
		return
	}

	user, _ := h.userRepo.GetUserbyEmail(ctx, payload.Email)
	if user == nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Email salah atau User belum terdaftar",
		})
		return
	}

	if !comparePassword(payload.Password, user.PasswordHash) {
		utils.ResponseJson(w, http.StatusBadRequest, utils.Response{
			Message: "Password salah",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.tokenDuration, h.jwtSecret)
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.Response{
			Message: "Gagal membuat token",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Login Berhasil",
		Data: models.UserResponse{
			Token: models.Token{
				AccessToken: token,
			},
			User: models.User{
				ID:    user.ID,
				Email: user.Email,
			},
		},
	})
}

func (h *UserHandler) Session(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, _ := middleware.GetClaims(ctx)

	email, _ := claims["email"].(string)

	user, _ := h.userRepo.GetUserbyEmail(ctx, email)
	if user == nil {
		utils.ResponseJson(w, http.StatusNotFound, utils.Response{
			Message: "Email salah atau User belum terdaftar",
		})
		return
	}

	utils.ResponseJson(w, http.StatusOK, utils.Response{
		Message: "Session Valid",
		Data: models.SessionResponse{
			User: models.User{
				ID:    user.ID,
				Email: user.Email,
			},
		},
	})
}

func hashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
