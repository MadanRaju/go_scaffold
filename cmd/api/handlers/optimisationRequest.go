package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"

	"inventory-optimisation-server/internal/constants"
	"inventory-optimisation-server/internal/optimisationRequest"
	"inventory-optimisation-server/internal/platform/db"
	"inventory-optimisation-server/internal/platform/web"
)

// OptimisationRequest Handler
type OptimisationRequest struct {
	MasterDB *db.DB
}

// Create will validate the request and then add it to queue
func (o *OptimisationRequest) Create(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	dbConn := o.MasterDB.Copy()
	defer dbConn.Close()

	v := ctx.Value(web.KeyValues).(*web.Values)

	err := r.ParseMultipartForm(100000)
	if err != nil {
		return errors.Wrap(err, "")
	}

	// conf := aws.Config{Region: aws.String("us-west-2")}
	// sess := session.New(&conf)
	// svc := s3manager.NewUploader(sess)
	// bucket := ""

	fileTypes := []string{constants.PRODUCT_DATA_FILE, constants.FACTORY_DATA_FILE}
	requestInput := []optimisationRequest.NewRequestInput{}
	name := r.Form["name"][0]
	for _, fileType := range fileTypes {
		for _, file := range r.MultipartForm.File[fileType] {
			errFile, valid := optimisationRequest.Validate(file, fileType)
			if valid == false {
				web.Respond(ctx, log, w, errFile, http.StatusUnprocessableEntity)
			} else {
				// Upload to S3 Code
				requestInput = append(requestInput, optimisationRequest.NewRequestInput{
					Type:     fileType,
					Location: "s3-path",
				})
			}
		}
	}

	newRequest := optimisationRequest.NewRequest{
		Name:  name,
		Input: requestInput,
	}

	request, err := optimisationRequest.Create(ctx, dbConn, &newRequest, v.Now)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Request: %+v", &request)
	}

	web.Respond(ctx, log, w, request, http.StatusCreated)
	return nil
}

// Validate will validate the excel input
func (o *OptimisationRequest) Validate(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	fileType := r.URL.Query().Get("type")
	if len(r.MultipartForm.File[fileType]) == 0 {
		return web.ErrValidation
	}
	result, valid := optimisationRequest.Validate(r.MultipartForm.File[fileType][0], fileType)
	if valid == false {
		web.Respond(ctx, log, w, result, http.StatusUnprocessableEntity)
	}
	web.Respond(ctx, log, w, nil, http.StatusNoContent)
	return nil
}

// List will list all the requests received
func (o *OptimisationRequest) List(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := o.MasterDB.Copy()
	defer dbConn.Close()

	requests, err := optimisationRequest.List(ctx, dbConn)
	if err = translate(err); err != nil {
		return errors.Wrap(err, "")
	}

	web.Respond(ctx, log, w, requests, http.StatusOK)
	return nil
}

// Retrieve returns the specified request from the system.
func (o *OptimisationRequest) Retrieve(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn := o.MasterDB.Copy()
	defer dbConn.Close()

	request, err := optimisationRequest.Retrieve(ctx, dbConn, params["id"])
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, log, w, request, http.StatusOK)
	return nil
}
