package handlers

import (
	"context"
	"log"
	"net/http"

	"inventory-optimisation-server/internal/platform/auth"
	"inventory-optimisation-server/internal/platform/db"
	"inventory-optimisation-server/internal/platform/web"
	"inventory-optimisation-server/internal/user"

	"github.com/pkg/errors"
)

// User represents the User API method handler set.
type User struct {
	MasterDB       *db.DB
	TokenGenerator user.TokenGenerator

	// ADD OTHER STATE LIKE THE LOGGER AND CONFIG HERE.
}

// List returns all the existing users in the system.
func (u *User) List(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	usrs, err := user.List(ctx, dbConn)
	if err = translate(err); err != nil {
		return errors.Wrap(err, "")
	}

	web.Respond(ctx, log, w, usrs, http.StatusOK)
	return nil
}

// Retrieve returns the specified user from the system.
func (u *User) Retrieve(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.ErrUnauthorized
	}

	usr, err := user.Retrieve(ctx, claims, dbConn, params["id"])
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, log, w, usr, http.StatusOK)
	return nil
}

// Create inserts a new user into the system.
func (u *User) Create(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	v := ctx.Value(web.KeyValues).(*web.Values)

	var newU user.NewUser
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Wrap(err, "")
	}

	usr, err := user.Create(ctx, dbConn, &newU, v.Now)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "User: %+v", &usr)
	}

	web.Respond(ctx, log, w, usr, http.StatusCreated)
	return nil
}

// Update updates the specified user in the system.
func (u *User) Update(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	v := ctx.Value(web.KeyValues).(*web.Values)

	var upd user.UpdateUser
	if err := web.Unmarshal(r.Body, &upd); err != nil {
		return errors.Wrap(err, "")
	}

	err := user.Update(ctx, dbConn, params["id"], &upd, v.Now)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Id: %s  User: %+v", params["id"], &upd)
	}

	web.Respond(ctx, log, w, nil, http.StatusNoContent)
	return nil
}

// Delete removes the specified user from the system.
func (u *User) Delete(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	err := user.Delete(ctx, dbConn, params["id"])
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, log, w, nil, http.StatusNoContent)
	return nil
}

// Token handles a request to authenticate a user. It expects a request using
// Basic Auth with a user's email and password. It responds with a JWT.
func (u *User) Token(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	v := ctx.Value(web.KeyValues).(*web.Values)

	email, pass, ok := r.BasicAuth()
	if !ok {
		return web.ErrUnauthorized
	}

	tkn, err := user.Authenticate(ctx, dbConn, u.TokenGenerator, v.Now, email, pass)
	if err = translate(err); err != nil {
		return errors.Wrap(err, "authenticating")
	}

	web.Respond(ctx, log, w, tkn, http.StatusOK)
	return nil
}
