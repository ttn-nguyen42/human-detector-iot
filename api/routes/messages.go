package routes

import "net/http"

var (
	/*
	 * Default response template for 500 Internal Server Error
	 */
	MessageInternalServerError = MessageResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	}

	/*
	 * Default response template for 400 Not Found
	 */
	MessageNotFound = MessageResponse{
		Message: http.StatusText(http.StatusNotFound),
	}
)

/*
 * Use in endpoints for returning an error message in a structured form.
 *
 * For example: {"message": "Internal server error" }
 */
type MessageResponse struct {
	Message string `json:"message" binding:"required"`
}

/*
 * Use in endpoints for returning an ID in a structured form.
 * 
 * For example: {"id": "noice" }
 */
type IdResponse struct {
	Id string `json:"id" binding:"required"`
}
