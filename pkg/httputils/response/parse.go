package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

var ErrDestinationNotPointer error = fmt.Errorf("destination MUST be a pointer value")

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
