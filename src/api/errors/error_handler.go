package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

var (
	badRequestErrorType = reflect.TypeOf(BadRequestError(""))
	forbiddenErrorType  = reflect.TypeOf(ForbiddenError(""))
	notFoundErrorType   = reflect.TypeOf(NotFoundError(""))

	mappings = map[reflect.Type]func(e error) ApiError{
		badRequestErrorType: func(e error) ApiError { return badRequestApiError(e.Error()) },
		forbiddenErrorType:  func(e error) ApiError { return forbiddenApiError(e.Error()) },
		notFoundErrorType:   func(e error) ApiError { return notFoundApiError(e.Error()) },

		// json errors => bad request error
		reflect.TypeOf(&json.SyntaxError{}):           func(e error) ApiError { return badRequestApiError(e.Error()) },
		reflect.TypeOf(&json.MarshalerError{}):        func(e error) ApiError { return badRequestApiError(e.Error()) },
		reflect.TypeOf(&json.UnmarshalTypeError{}):    func(e error) ApiError { return badRequestApiError(e.Error()) },
		reflect.TypeOf(&json.InvalidUnmarshalError{}): func(e error) ApiError { return badRequestApiError(e.Error()) },
		reflect.TypeOf(&json.UnsupportedTypeError{}):  func(e error) ApiError { return badRequestApiError(e.Error()) },
		reflect.TypeOf(&json.UnsupportedValueError{}): func(e error) ApiError { return badRequestApiError(e.Error()) },
	}
)

func Handle(w http.ResponseWriter, e error) {
	if handleFunc, ok := mappings[reflect.TypeOf(e)]; ok {
		sendApiError(w, handleFunc(e))
		return
	}

	sendApiError(w, internalApiError(e))
}

func sendApiError(w http.ResponseWriter, e ApiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Sprintf("Error writing HTTP response - %s", err.Error()))
	}
}
