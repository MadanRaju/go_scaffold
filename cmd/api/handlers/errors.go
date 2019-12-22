package handlers

import (
	"inventory-optimisation-server/internal/platform/web"
	"inventory-optimisation-server/internal/user"

	"github.com/pkg/errors"
)

// translate looks for certain error types and transforms
// them into web errors. We are losing the trace when this
// error is converted. But we don't log traces for these.
func translate(err error) error {
	switch errors.Cause(err) {
	case user.ErrNotFound:
		return web.ErrNotFound
	case user.ErrInvalidID:
		return web.ErrInvalidID
	case user.ErrAuthenticationFailure:
		return web.ErrUnauthorized
	case user.ErrForbidden:
		return web.ErrForbidden
	}
	return err
}
