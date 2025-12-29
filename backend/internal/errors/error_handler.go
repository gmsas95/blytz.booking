package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Code    ErrorCode   `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		response := ErrorResponse{
			Error:   http.StatusText(appErr.StatusCode),
			Message: appErr.Message,
			Code:    appErr.Code,
			Details: appErr.Details,
		}
		c.JSON(appErr.StatusCode, response)
		return
	}

	response := ErrorResponse{
		Error:   http.StatusText(http.StatusInternalServerError),
		Message: "An unexpected error occurred",
		Code:    ErrCodeInternal,
	}
	c.JSON(http.StatusInternalServerError, response)
}
