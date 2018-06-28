// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"awesomeProject/restapi/operations"
	"awesomeProject/restapi/operations/users"
	"sync"
	"sync/atomic"
	"awesomeProject/models"
	"github.com/go-openapi/swag"
)

//go:generate swagger generate server --target .. --name  --spec ../swagger.yml

var items = make(map[int64]*models.User)
var lastID int64

var itemsLock = &sync.Mutex{}

func newItemID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

func addItem(item *models.User) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	newID := newItemID()
	item.ID = newID
	items[newID] = item

	return nil
}

//func updateItem(id int64, item *models.Item) error {
//	if item == nil {
//		return errors.New(500, "item must be present")
//	}
//
//	itemsLock.Lock()
//	defer itemsLock.Unlock()
//
//	_, exists := items[id]
//	if !exists {
//		return errors.NotFound("not found: item %d", id)
//	}
//
//	item.ID = id
//	items[id] = item
//	return nil
//}

//func deleteItem(id int64) error {
//	itemsLock.Lock()
//	defer itemsLock.Unlock()
//
//	_, exists := items[id]
//	if !exists {
//		return errors.NotFound("not found: item %d", id)
//	}
//
//	delete(items, id)
//	return nil
//}

func allItems() (result []*models.User) {
	result = make([]*models.User, 0)
	for id, item := range items {
		if id >= 0 {
			result = append(result, item)
		}

	}
	return
}


func configureFlags(api *operations.AwesomeProjectAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AwesomeProjectAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder {
		if err := addItem(params.Body); err != nil {
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return users.NewAddOneCreated().WithPayload(params.Body)
	})
	api.UsersDestroyOneHandler = users.DestroyOneHandlerFunc(func(params users.DestroyOneParams) middleware.Responder {
		return middleware.NotImplemented("operation users.DestroyOne has not yet been implemented")
	})
	api.UsersFindUsersHandler = users.FindUsersHandlerFunc(func(params users.FindUsersParams) middleware.Responder {
		return users.NewFindUsersOK().WithPayload(allItems())
	})
	api.UsersUpdateOneHandler = users.UpdateOneHandlerFunc(func(params users.UpdateOneParams) middleware.Responder {
		return middleware.NotImplemented("operation users.UpdateOne has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
