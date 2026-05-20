package response

import apperr "github.com/gis-mdm/server-backend-go/internal/shared/errors"

// FromError maps errors to AppError for HTTP responses.
func FromError(err error) *apperr.AppError {
	if e, ok := err.(*apperr.AppError); ok {
		return e
	}
	return apperr.Internal("unexpected error", err)
}
