package apierrors_test

import (
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/apierrors"
	"github.com/cultureamp/ca-go/ref"
	"github.com/stretchr/testify/assert"
)

func ExampleResponse_AddError() {
	er := apierrors.Response{}.AddError("title")

	fmt.Println(string(er.MustSerialize()))
	// Output:
	// {"errors":[{"title":"title"}]}
}

func ExampleResponse_Add() {
	er := apierrors.Response{}.Add(apierrors.ErrorInstance{
		Title:  "Unexpected error",
		Detail: ref.String("Further error details (without including internal information)"),
		ID:     ref.String("Request ID or similar"),
	})

	fmt.Println(string(er.MustSerialize()))

	// Output:
	// {"errors":[{"title":"Unexpected error","detail":"Further error details (without including internal information)","id":"Request ID or similar"}]}
}

func ExampleNewUnexpectedError() {
	er := apierrors.NewUnexpectedError(ref.String("Request ID or similar"))

	fmt.Println(string(er.MustSerialize()))
	// Output:
	// {"errors":[{"title":"Unexpected error","id":"Request ID or similar"}]}
}

func TestResponse_Add(t *testing.T) {
	er := apierrors.Response{}.AddError("500")
	assert.NotNil(t, er)
	assert.Len(t, er.Errors, 1)
	assert.Equal(t, apierrors.ErrorInstance{Title: "500"}, er.Errors[0])
}
