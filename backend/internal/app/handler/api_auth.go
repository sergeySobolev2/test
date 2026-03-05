package handler

import (
	"errors"
	"net/http"
	"strings"

	"partitionlab/internal/app/ds"
	"partitionlab/internal/app/middleware"
	"partitionlab/internal/app/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ApiRegisterUser регистрация нового пользователя
// @Summary Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{login=string,password=string} true "Данные для регистрации"
// @Success 200 {object} object{user=ds.User}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /api/users/register [post]
func (h *Handler) ApiRegisterUser(ctx *gin.Context) {
	type requestBody struct {
		Login    string `json:"login" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var body requestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	if _, err := h.Repository.GetUserByLogin(body.Login); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Создаем пользователя
	isModerator := strings.HasPrefix(strings.ToLower(body.Login), "moder")
	user := &ds.User{
		Login:       body.Login,
		Password:    string(hashedPassword),
		IsModerator: isModerator,
	}

	if err := h.Repository.CreateUser(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(ctx, gin.H{"user": user}, 1, gin.H{})
}

// ApiLogin вход пользователя
// @Summary Вход пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{login=string,password=string} true "Данные для входа"
// @Success 200 {object} object{user=ds.User,token=string,session_id=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /api/users/login [post]
func (h *Handler) ApiLogin(ctx *gin.Context) {
	type requestBody struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body requestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Находим пользователя
	user, err := h.Repository.GetUserByLogin(body.Login)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Генерируем JWT токен
	token, err := h.JWTService.Generate(user.ID, user.Login, user.IsModerator)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Создаем сессию в Redis
	sessionID := uuid.New().String()
	sessionData := auth.SessionData{
		UserID:      user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}
	if err := h.SessionService.Create(ctx.Request.Context(), sessionID, sessionData); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Устанавливаем cookie с session_id
	ctx.SetCookie("session_id", sessionID, 86400, "/", "", false, true)

	jsonResponse(ctx, gin.H{
		"user":       user,
		"token":      token,
		"session_id": sessionID,
	}, 1, gin.H{})
}

// ApiLogout выход пользователя
// @Summary Выход пользователя
// @Tags auth
// @Security BearerAuth
// @Security CookieAuth
// @Produce json
// @Success 200 {object} object{message=string}
// @Router /api/users/logout [post]
func (h *Handler) ApiLogout(ctx *gin.Context) {
	// Удаляем сессию из Redis
	if sessionID, err := ctx.Cookie("session_id"); err == nil && sessionID != "" {
		_ = h.SessionService.Delete(ctx.Request.Context(), sessionID)
	}

	// Удаляем cookie
	ctx.SetCookie("session_id", "", -1, "/", "", false, true)

	jsonResponse(ctx, gin.H{"message": "logged out"}, 1, gin.H{})
}

// ApiGetProfile получение профиля текущего пользователя
// @Summary Получение профиля текущего пользователя
// @Tags auth
// @Security BearerAuth
// @Security CookieAuth
// @Produce json
// @Success 200 {object} object{user=ds.User}
// @Failure 401 {object} object{error=string}
// @Router /api/users/profile [get]
func (h *Handler) ApiGetProfile(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(ctx, gin.H{"user": user}, 1, gin.H{})
}

// ApiUpdateProfile обновление профиля текущего пользователя
// @Summary Обновление профиля
// @Tags auth
// @Security BearerAuth
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request body object{login=string} true "Новые данные профиля"
// @Success 200 {object} object{user=ds.User}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /api/users/profile [put]
func (h *Handler) ApiUpdateProfile(ctx *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	type requestBody struct {
		Login *string `json:"login,omitempty"`
	}

	var body requestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	fields := map[string]interface{}{}
	if body.Login != nil {
		fields["login"] = *body.Login
	}

	if len(fields) > 0 {
		if err := h.Repository.UpdateUser(userID, fields); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	user, _ := h.Repository.GetUserByID(userID)
	jsonResponse(ctx, gin.H{"user": user}, 1, gin.H{})
}
