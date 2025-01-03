package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status struct for standardized responses
type Status struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrorCode int         `json:"error_code"`
}

// NewOkStatus creates a response with status "Ok"
func NewOkStatus() *Status {
	return &Status{Status: "Ok"}
}

// NewErrorStatus creates a response with status "Error"
func NewErrorStatus() *Status {
	return &Status{Status: "Error"}
}

// SetData sets the data for the response
func (s *Status) SetData(data interface{}) *Status {
	s.Data = data
	return s
}

// SetMessage sets the message for the response
func (s *Status) SetMessage(message string) *Status {
	s.Message = message
	return s
}

// SendResponse sends the response in a consistent format
func SendResponse(ctx *gin.Context, httpStatusCode int, status *Status) {
	ctx.JSON(httpStatusCode, status)
}

// SendErrorResponse sends an error response with a message and optional data
func SendErrorResponse(ctx *gin.Context, message string, data interface{}) {
	status := NewErrorStatus().
		SetMessage(message).
		SetData(data)
	SendResponse(ctx, http.StatusBadRequest, status)
}

// SendSuccessResponse sends a success response with a message and data
func SendSuccessResponse(ctx *gin.Context, message string, data interface{}) {
	status := NewOkStatus().
		SetMessage(message).
		SetData(data)
	SendResponse(ctx, http.StatusOK, status)
}
