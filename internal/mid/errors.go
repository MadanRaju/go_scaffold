package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"inventory-optimisation-server/internal/platform/web"

	"github.com/pkg/errors"
)

// ErrorHandler for catching and responding to errors.
func ErrorHandler(next web.Handler) web.Handler {

	// Create the handler that will be attached in the middleware chain.
	h := func(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

		v := ctx.Value(web.KeyValues).(*web.Values)

		// In the event of a panic, we want to capture it here so we can send an
		// error down the stack.
		defer func() {
			if r := recover(); r != nil {

				// Indicate this request had an error.
				v.Error = true

				// Log the panic.
				log.Printf("ERROR : Panic Caught : %s\n", r)

				// Respond with the error.
				web.RespondError(ctx, log, w, errors.New("unhandled"), http.StatusInternalServerError)

				// Print out the stack.
				log.Printf("ERROR : Stacktrace\n%s\n", debug.Stack())
			}
		}()

		if err := next(ctx, log, w, r, params); err != nil {

			// Indicate this request had an error.
			v.Error = true

			// What is the root error.
			err = errors.Cause(err)

			if err != web.ErrNotFound {

				// Log the error.
				log.Printf("ERROR : %v\n", err)
			}

			// Respond with the error.
			web.Error(ctx, log, w, err)

			// The error has been handled so we can stop propagating it.
			return nil
		}

		return nil
	}

	return h
}
