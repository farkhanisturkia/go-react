package structs

type PaginatedResponse[T any] struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    []T         `json:"data"`
    Page    int         `json:"page"`
    Size    int         `json:"size"`
    Total   int64       `json:"total"`
}