package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-react/backend/database"
	"go-react/backend/helpers"
	"go-react/backend/models"
	"go-react/backend/pkg/redis"
	"go-react/backend/structs"
)

const (
	userListCacheKey = "users:list"
	userCachePrefix  = "user:id:"
	cacheTTL         = 5 * time.Minute
	userTTL          = 10 * time.Minute
)

// FindUsers - List semua user (dengan cache)
func FindUsers(c *gin.Context) {
	var users []models.User

	// Coba ambil dari Redis dulu
	val, err := redis.Client.Get(redis.Ctx, userListCacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &users) == nil {
			c.JSON(http.StatusOK, structs.SuccessResponse{
				Success: true,
				Message: "Lists Data Users (from cache)",
				Data:    users,
			})
			return
		}
	}

	// Cache miss → ambil dari DB
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan ke Redis
	if data, err := json.Marshal(users); err == nil {
		redis.Client.Set(redis.Ctx, userListCacheKey, data, cacheTTL)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Users",
		Data:    users,
	})
}

// CreateUser - Buat user + hapus cache list
func CreateUser(c *gin.Context) {
	var req structs.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Invalidate cache list
	redis.Client.Del(redis.Ctx, userListCacheKey)

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}


// FindUserById - Detail user (dengan cache)
func FindUserById(c *gin.Context) {
	id := c.Param("id")
	cacheKey := userCachePrefix + id

	var user models.User

	// Coba dari cache
	val, err := redis.Client.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &user) == nil {
			c.JSON(http.StatusOK, structs.SuccessResponse{
				Success: true,
				Message: "User Found (from cache)",
				Data: structs.UserResponse{
					Id:        user.Id,
					Name:      user.Name,
					Username:  user.Username,
					Email:     user.Email,
					CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				},
			})
			return
		}
	}

	// Cache miss
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan ke cache
	if data, err := json.Marshal(user); err == nil {
		redis.Client.Set(redis.Ctx, cacheKey, data, userTTL)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// UpdateUser - Update + invalidate cache related
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	if req.Password != "" {
		user.Password = helpers.HashPassword(req.Password)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Invalidate cache
	cacheKey := userCachePrefix + id
	redis.Client.Del(redis.Ctx, cacheKey, userListCacheKey)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User updated successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteUser - Hapus + invalidate cache
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Invalidate cache
	cacheKey := userCachePrefix + id
	redis.Client.Del(redis.Ctx, cacheKey, userListCacheKey)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}