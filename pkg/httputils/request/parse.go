package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var ErrDestinationNotPointer error = errors.New("destination MUST be a pointer value")

// parse request http's request body to the dest object
// r = pointer http.Request struct
// dest = destination object to which the request body will be copied
// note: dest should be pointer i.e &x else, the function will return an error
// note: for parse body use tag "json" to specify which member correspond to what field on the struct
func ParseRequestBody(r *http.Request, dest any) error {
	defer r.Body.Close()

	dstVal := reflect.ValueOf(dest)
	if dstVal.Kind() != reflect.Pointer {
		return ErrDestinationNotPointer
	}
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return nil
}

// parse request http's path parameters to the dest object
// ex: api/v1/users/{user_id}
// the function will find for fields that have "path" tag in it, find the corresponding path param, and copy its value to the struct field
// ex: UserId int64 `path:"user_id"` <-- this will find path param named "user_id" and copy the value to the UserId member variable
// note: dest should be pointer i.e &x else, the function will return an error
func ParsePathParam(r *http.Request, dest any) error {
	rVal := reflect.ValueOf(dest)
	if rVal.Kind() != reflect.Pointer {
		return ErrDestinationNotPointer
	}
	rVal = rVal.Elem()
	rType := rVal.Type()
	for i := 0; i < rType.NumField(); i++ {
		pathName, ok := rType.Field(i).Tag.Lookup("path")
		if !ok {
			continue
		}

		field := rVal.Field(i)
		if !field.CanSet() {
			continue
		}

		value := chi.URLParam(r, pathName)
		if value == "" {
			return fmt.Errorf("path param '%v' not found", pathName)
		}
		switch field.Kind() {
		case reflect.Int64:
			iValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to convert value from path '%v', value=%v", pathName, value)
			}
			field.SetInt(iValue)
		case reflect.String:
			field.SetString(value)
		default:
			// add more if needed
			return fmt.Errorf("value of kind %v is not implemeted yet", field.Kind())
		}
	}
	return nil
}

// parse request http's query parameters to the dest object
// the function will find for fields that have "schema" tag in it, find the corresponding query parameter, and copy its value to the struct field
// ex: Filter string `schema:"filter"` <-- this will find query named "filter" and copy the value to the Filter member variable
// note: dest should be pointer i.e &x else, the function will return an error
func ParseQueryParam(r *http.Request, dest any) error {
	rVal := reflect.ValueOf(dest)
	if rVal.Kind() != reflect.Pointer {
		return ErrDestinationNotPointer
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dest, r.URL.Query()); err != nil {
		return err
	}

	return nil
}
