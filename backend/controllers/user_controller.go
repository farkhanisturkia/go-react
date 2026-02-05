package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-react-vue/backend/database"
	"go-react-vue/backend/helpers"
	"go-react-vue/backend/models"
	"go-react-vue/backend/pkg/redis"
	"go-react-vue/backend/structs"
	"go-react-vue/backend/cache"
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
	cacheKey := fmt.Sprintf("%s:size:%d:role:user", cache.UserListCachePrefix, page, size)

	val, err := redis.Client.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse structs.PaginatedResponse[structs.UserListItemResponse]
		if json.Unmarshal([]byte(val), &cachedResponse) == nil {
			cachedResponse.Message = "Lists Data Users (dari cache)"
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	var total int64
	if err := database.DB.Model(&models.User{}).
		Where("role = ?", "user").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	subQuery := database.DB.Table("user_courses").
		Select("COUNT(*)").
		Where("user_courses.participant_id = users.id")

	var users []structs.UserListItemResponse
	err = database.DB.
		Table("users").
		Select(`
			users.id,
			users.name,
			users.username,
			users.email,
			users.role,
			users.created_at,
			users.updated_at,
			(?) AS enrolled_courses_count
		`, subQuery).
		Where("users.role = ?", "user").
		Offset(offset).
		Limit(size).
		Scan(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Getting lists Data Users is failed",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	response := structs.PaginatedResponse[structs.UserListItemResponse] {
		Success: true,
		Message: "Lists Data Users",
		Data:    users,
		Page:    page,
		Size:    size,
		Total:   total,
	}

	if len(users) > 0 {
		if data, err := json.Marshal(response); err == nil {
			redis.Client.Set(redis.Ctx, cacheKey, data, cache.CacheTTL)
			// ← Track key ini di Set
			redis.Client.SAdd(redis.Ctx, cache.UserListKeysSet, cacheKey)
			redis.Client.Expire(redis.Ctx, cache.UserListKeysSet, cache.CacheTTL)
		}
	}

	redis.Client.Set(redis.Ctx, cache.UserTotalCacheKey, strconv.FormatInt(total, 10), cache.CacheTTL)

	c.JSON(http.StatusOK, response)
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
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
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
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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