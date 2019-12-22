package handlers

import (
	"context"
	"log"
	"net/http"

	"inventory-optimisation-server/internal/platform/db"
	"inventory-optimisation-server/internal/platform/web"
)

// Health provides support for orchestration health checks.
type Health struct {
	MasterDB *db.DB
}

// Check validates the service is healthy and ready to accept requests.
func (h *Health) Check(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := h.MasterDB.Copy()
	defer dbConn.Close()

	if err := dbConn.StatusCheck(ctx); err != nil {
		return err
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	web.Respond(ctx, log, w, status, http.StatusOK)
	return nil
}
