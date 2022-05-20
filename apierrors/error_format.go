package apierrors

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/cultureamp/ca-go/ref"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

type errorIDGenerator func() string

// NewApiErrorFormatter returns a Goa error formatter
// that ensures the shape of errors returned by
// the service conform to the JSON API error spec.
//
// The additional consideration here is to avoid
// details of unexpected errors from appearing in
// the service output.
//
// If a service returns a conforming error, it's
// expected that care has been taken to only output
// known, generic information and not the output of
// Error().
//
// It is assumed that logging unhandled errors
// is dealt with by a separate component.
func NewApiErrorFormatter() func(error) goahttp.Statuser {
	defaultErrorGenerator := goa.NewErrorID
	return func(err error) goahttp.Statuser {
		return formatApiError(defaultErrorGenerator, err)
	}
}

func formatApiError(idGenerator errorIDGenerator, err error) goahttp.Statuser {
	//
	// Using this function as a formatter bypasses
	// both the generic error handling logic and the
	// generated serialisation for the types specified
	// in the design.
	//
	// The idea is that it will convert Goa ServiceErrors
	// (as these will contain validation information),
	// map from generated types that appear to conform to JSON API,
	// and output a generic error for the rest.
	//

	// for a goa error, map it to an API error
	if serr, ok := err.(*goa.ServiceError); ok { //nolint:errorlint // we're deliberately only wanting to test for top-level ServiceError
		return mapServiceErrorToApiErrorResponse(serr)
	}

	// if it looks like an API error from generated code, map it
	if isGeneratedErrorType(err) {
		return mapGeneratedErrorToApiErrorResponse(err)
	}

	// fall back to a generic
	return Response{
		Errors: []ErrorInstance{
			{
				Title: "Internal server error",
				ID:    ref.String(idGenerator()),
			},
		},
		Status: http.StatusInternalServerError,
	}
}

// map a goa service eror to the JSON API
func mapServiceErrorToApiErrorResponse(serviceError *goa.ServiceError) Response {
	statusCode := serviceErrorStatusCode(serviceError)

	detail := serviceError.Message
	if serviceError.Fault {
		// be sure to mask internal errors
		detail = "Internal error"
	}

	return Response{
		Errors: []ErrorInstance{
			{
				Title: serviceError.Name,
				ID:    &serviceError.ID,
				// FIXME support field pointer?
				Detail: &detail,
			},
		},
		Status: statusCode,
	}
}

func serviceErrorStatusCode(serviceError *goa.ServiceError) int {
	// infer the status code in the same way that the Goa core normally does
	if serviceError.Fault {
		return http.StatusInternalServerError
	}
	if serviceError.Timeout {
		if serviceError.Temporary {
			return http.StatusGatewayTimeout
		}
		return http.StatusRequestTimeout
	}
	if serviceError.Temporary {
		return http.StatusServiceUnavailable
	}
	return http.StatusBadRequest
}

func isGeneratedErrorType(err error) bool {
	// does it look like a Goa-generated JSON API error?
	// not attempting perfection, having "gen" in the package name
	// and an array field called "errors" is enough

	et := reflect.TypeOf(err)
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}

	// is it a generated struct?
	if !strings.Contains(et.PkgPath(), "/gen/") {
		return false
	}

	// does it have an errors member?
	f, found := et.FieldByName("Errors")
	if !found {
		return false
	}

	// is the errors member an array?
	switch f.Type.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	}

	return false
}

func mapGeneratedErrorToApiErrorResponse(err error) Response {
	ev := reflect.ValueOf(err)
	if ev.Kind() == reflect.Ptr {
		// assume non-null
		ev = ev.Elem()
	}

	ev = ev.FieldByName("Errors")

	// iterate the array
	l := ev.Len()
	errorInstances := []ErrorInstance{}
	for i := 0; i < l; i++ {
		ei := ev.Index(i)
		if ei.Kind() == reflect.Ptr {
			if ei.IsNil() {
				continue
			}
			ei = ei.Elem()
		}

		errorInstances = append(errorInstances, ErrorInstance{
			Title:  strFieldValue(ei, "Title"),
			Detail: strPtrFieldValue(ei, "Detail"),
			ID:     strPtrFieldValue(ei, "ID"),
			Code:   strPtrFieldValue(ei, "Code"),
		})
	}

	return Response{
		Errors: errorInstances,
	}
}

func strFieldValue(structValue reflect.Value, fieldName string) string {
	strValue := strPtrFieldValue(structValue, fieldName)

	if strValue == nil {
		return ""
	}

	return *strValue
}

func strPtrFieldValue(structValue reflect.Value, fieldName string) *string {
	fieldValue := structValue.FieldByName(fieldName)

	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil
		}

		fieldValue = fieldValue.Elem()
	}

	// This will print bogus values if there happen to be structs, but the assumption
	// here is that we're mapping types that look basically like the JSON API error.
	strValue := fieldValue.String()
	return &strValue
}
