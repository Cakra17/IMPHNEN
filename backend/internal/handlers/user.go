package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/models"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/Cakra17/imphnen/internal/validation"
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

// Register godoc
// @Summary      Register a new user account
// @Description  Create a new user account with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterPayload  true  "Registration details"
// @Success      201      {object}  utils.Response{message=string}  "Account created successfully"
// @Failure      400      {object}  utils.Response{message=string}  "Invalid request data"
// @Failure      409      {object}  utils.Response{message=string}  "Email already registered"
// @Failure      500      {object}  utils.Response{message=string}  "Internal server error"
// @Router       /auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterPayload
	ctx := r.Context()

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
		FirstName:    payload.FistName,
		LastName:     payload.LastName,
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

// Login godoc
// @Summary      Login to user account
// @Description  Authenticate user with email and password, returns access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.AuthPayload  true  "Login credentials"
// @Success      200      {object}  utils.Response{data=models.UserResponse}  "Login successful with access token"
// @Failure      400      {object}  utils.Response{message=string}  "Invalid request data or wrong password"
// @Failure      404      {object}  utils.Response{message=string}  "User not found"
// @Failure      500      {object}  utils.Response{message=string}  "Failed to generate token"
// @Router       /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.AuthPayload
	ctx := r.Context()

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
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
		},
	})
}

// Session godoc
// @Summary      Get current user session
// @Description  Retrieve the authenticated user's information from their JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response{data=models.SessionResponse}  "Session valid with user data"
// @Failure      401  {object}  utils.Response{message=string}  "Unauthorized - invalid or missing token"
// @Failure      404  {object}  utils.Response{message=string}  "User not found"
// @Router       /users/me [get]
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
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
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
