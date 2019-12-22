package optimisationRequest

import (
	"context"
	"fmt"
	"inventory-optimisation-server/internal/platform/db"
	"mime/multipart"
	"time"

	"inventory-optimisation-server/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrNotFound abstracts the mgo not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

const requestsCollection = "requests"

// Create inserts a new optimisation request into the database.
func Create(ctx context.Context, dbConn *db.DB, newRequest *NewRequest, now time.Time) (*Request, error) {
	now = now.Truncate(time.Millisecond)

	var requestInput []RequestInput
	for _, input := range newRequest.Input {
		requestInput = append(requestInput, RequestInput{
			Type:     input.Type,
			Location: input.Location,
		})
	}

	request := Request{
		ID:          bson.NewObjectId(),
		Name:        newRequest.Name,
		Input:       requestInput,
		DateCreated: now,
	}

	f := func(collection *mgo.Collection) error {
		return collection.Insert(&request)
	}
	if err := dbConn.Execute(ctx, requestsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.requests.insert(%s)", db.Query(&request)))
	}

	return &request, nil
}

// UploadToS3 will upload excel file to S3
func UploadToS3(svc *s3manager.Uploader, bucket string, file *multipart.FileHeader) (*s3manager.UploadOutput, error) {
	actualFile, err := file.Open()
	defer actualFile.Close()
	if err != nil {
		return nil, err
	}

	return svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Filename),
		Body:   actualFile,
	})
}

// Validate will check if excel file is valid or not
func Validate(file *multipart.FileHeader, fileType string) (*multipart.FileHeader, bool) {
	if fileType == constants.FACTORY_DATA_FILE {
		return file, false
	}
	return nil, true
}

// List will list all the requests served
func List(ctx context.Context, dbConn *db.DB) ([]Request, error) {

	r := []Request{}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&r)
	}
	if err := dbConn.Execute(ctx, requestsCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.requests.find()")
	}

	return r, nil
}

// Retrieve gets the specified request from the database.
func Retrieve(ctx context.Context, dbConn *db.DB, id string) (*Request, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, ErrInvalidID
	}

	q := bson.M{"_id": bson.ObjectIdHex(id)}

	var r *Request
	f := func(collection *mgo.Collection) error {
		return collection.Find(q).One(&r)
	}
	if err := dbConn.Execute(ctx, requestsCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.requests.find(%s)", db.Query(q)))
	}

	return r, nil
}
