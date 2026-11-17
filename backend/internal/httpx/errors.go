package httpx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrValidation      = errors.New("validation error")
	ErrRateLimited     = errors.New("rate limited")
	ErrInternal        = errors.New("internal error")
)

type APIError struct {
	Status  int
	Code    string
	Message string
	Details map[string]string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, code, message string, details map[string]string) *APIError {
	return &APIError{Status: status, Code: code, Message: message, Details: details}
}

func MapError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return mapPGError(pgErr)
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	case errors.Is(err, ErrUnauthenticated):
		return NewAPIError(http.StatusUnauthorized, "UNAUTHENTICATED", err.Error(), nil)
	case errors.Is(err, ErrForbidden):
		return NewAPIError(http.StatusForbidden, "FORBIDDEN", err.Error(), nil)
	case errors.Is(err, ErrNotFound):
		return NewAPIError(http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	case errors.Is(err, ErrConflict):
		return NewAPIError(http.StatusConflict, "CONFLICT", err.Error(), nil)
	case errors.Is(err, ErrValidation):
		return NewAPIError(http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error(), nil)
	case errors.Is(err, ErrRateLimited):
		return NewAPIError(http.StatusTooManyRequests, "RATE_LIMITED", err.Error(), nil)
	default:
		return NewAPIError(http.StatusInternalServerError, "INTERNAL", "internal server error", nil)
	}
}

func mapPGError(pgErr *pgconn.PgError) *APIError {
	switch pgErr.Code {
	case pgerrcode.UniqueViolation:
		return NewAPIError(http.StatusConflict, "DUPLICATE", "resource already exists", nil)
	case pgerrcode.ForeignKeyViolation:
		return NewAPIError(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "referenced resource not found", map[string]string{pgErr.ColumnName: "invalid reference"})
	case pgerrcode.CheckViolation:
		return NewAPIError(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "constraint violation", map[string]string{pgErr.ConstraintName: "check failed"})
	case pgerrcode.ExclusionViolation:
		return NewAPIError(http.StatusConflict, "CONFLICT", "overlapping range", nil)
	case "23000": // custom trigger error for immutability
		return NewAPIError(http.StatusConflict, "CONFLICT", pgErr.Message, nil)
	default:
		return NewAPIError(http.StatusInternalServerError, "INTERNAL", "database error", nil)
	}
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status == http.StatusNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.NoDataFound
	}
	return false
}

func ValidationError(details map[string]string) *APIError {
	return NewAPIError(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validation failed", details)
}

func ConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, "CONFLICT", message, nil)
}

func InvalidTransitionError(from, to string) *APIError {
	return NewAPIError(http.StatusConflict, "INVALID_TRANSITION", fmt.Sprintf("cannot transition from %s to %s", from, to), nil)
}

func MakerCheckerViolationError() *APIError {
	return NewAPIError(http.StatusConflict, "MAKER_CHECKER_VIOLATION", "approver must differ from creator and submitter", nil)
}

func MonthLockedError() *APIError {
	return NewAPIError(http.StatusConflict, "MONTH_LOCKED", "attendance month is locked", nil)
}

func EvidenceRequiredError() *APIError {
	return NewAPIError(http.StatusConflict, "EVIDENCE_REQUIRED", "at least one accepted evidence is required", nil)
}

func BlockingIssuesError() *APIError {
	return NewAPIError(http.StatusConflict, "BLOCKING_ISSUES", "payroll run has blocking issues", nil)
}