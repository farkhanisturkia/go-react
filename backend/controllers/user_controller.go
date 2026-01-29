package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-react/backend/database"
	"go-react/backend/helpers"
	"go-react/backend/models"
	"go-react/backend/pkg/redis"
	"go-react/backend/structs"
	"go-react/backend/cache"
)

const (
	defaultPageSize      = 10
	maxPageSize          = 50
)

func FindUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", strconv.Itoa(defaultPageSize))

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(sizeStr)
	if size < 1 || size > maxPageSize {
		size = defaultPageSize
	}

	offset := (page - 1) * size
	cacheKey := fmt.Sprintf("%s%d:size:%d", cache.UserListCachePrefix, page, size)

	var users []models.User

	val, err := redis.Client.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &users) == nil {
			total := cache.GetUsersTotalCount()
			c.JSON(http.StatusOK, structs.PaginatedResponse{
				Success: true,
				Message: "Lists Data Users (from cache)",
				Data:    users,
				Page:    page,
				Size:    size,
				Total:   total,
			})
			return
		}
	}

	var total int64
	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.
		Offset(offset).
		Limit(size).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if len(users) > 0 {
		if data, err := json.Marshal(users); err == nil {
			redis.Client.Set(redis.Ctx, cacheKey, data, cache.CacheTTL)
			// ← Track key ini di Set
			redis.Client.SAdd(redis.Ctx, cache.UserListKeysSet, cacheKey)
			redis.Client.Expire(redis.Ctx, cache.UserListKeysSet, cache.CacheTTL)
		}
	}

	redis.Client.Set(redis.Ctx, cache.UserTotalCacheKey, strconv.FormatInt(total, 10), cache.CacheTTL)

	c.JSON(http.StatusOK, structs.PaginatedResponse{
		Success: true,
		Message: "Lists Data Users",
		Data:    users,
		Page:    page,
		Size:    size,
		Total:   total,
	})
}

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

	cache.InvalidateUserListCache()
	userCacheKey := cache.UserCachePrefix + strconv.Itoa(int(user.Id))
	redis.Client.Del(redis.Ctx, userCacheKey)

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

func FindUserById(c *gin.Context) {
	id := c.Param("id")
	cacheKey := cache.UserCachePrefix + id

	var user models.User

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

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if data, err := json.Marshal(user); err == nil {
		redis.Client.Set(redis.Ctx, cacheKey, data, cache.UserTTL)
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

	cache.InvalidateUserListCache()
	cacheKey := cache.UserCachePrefix + id
	redis.Client.Del(redis.Ctx, cacheKey)

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

	cache.InvalidateUserListCache()
	cacheKey := cache.UserCachePrefix + id
	redis.Client.Del(redis.Ctx, cacheKey)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}