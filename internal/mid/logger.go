package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"inventory-optimisation-server/internal/platform/web"
)

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	h := func(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

		err := next(ctx, log, w, r, params)

		v := ctx.Value(web.KeyValues).(*web.Values)

		log.Printf("(%d) : %s %s -> %s (%s)",
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now),
		)

		// For consistency return the error we received.
		return err
	}

	return h
}
