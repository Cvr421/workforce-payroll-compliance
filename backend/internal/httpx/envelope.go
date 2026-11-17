package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorResp  `json:"error,omitempty"`
	Meta    Meta        `json:"meta"`
}

type ErrorResp struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type Meta struct {
	RequestID string `json:"request_id"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: true,
		Data:    data,
		Meta:    Meta{RequestID: requestID},
	})
}

func WriteError(w http.ResponseWriter, status int, code, message string, details map[string]string, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: false,
		Error: &ErrorResp{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: Meta{RequestID: requestID},
	})
}

func WritePaginated(w http.ResponseWriter, status int, data interface{}, page, pageSize int, total int64, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: true,
		Data:    data,
		Meta: Meta{
			RequestID: requestID,
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
		},
	})
}