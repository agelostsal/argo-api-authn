package utils

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type APIErrorsTestSuite struct {
	suite.Suite
}

func (suite *APIErrorsTestSuite) TestAllErrors() {

	testMsg := "errMsg"
	testPlc := "errPlace"

	errBadRequest := &APIError{"Poorly formatted JSON. errMsg", 400, "BAD REQUEST"}
	errUnauthorized := &APIError{"errMsg", 401, "UNAUTHORIZED"}
	errNotFound := &APIError{"errMsg was not found", 404, "NOT FOUND"}
	errConflict := &APIError{"string object with errMsg: errMsg already exists", 409, "CONFLICT"}
	errMissingRequired := &APIError{"errMsg", 422, "UNPROCESSABLE ENTITY"}
	errInvalidField := &APIError{"Field: errMsg contains invalid data. reason", 422, "UNPROCESSABLE ENTITY"}
	errUnsupportedContent := &APIError{"errPlace: errMsg is not yet supported", 422, "UNPROCESSABLE ENTITY"}
	errDatabase := &APIError{"Database Error: errMsg", 500, "INTERNAL SERVER ERROR"}
	errGenericMissing := "string object contains an empty value for field: errMsg"
	errGenericInternal := &APIError{"Internal Error: errMsg", 500, "INTERNAL SERVER ERROR"}

	suite.Equal(errBadRequest, APIErrBadRequest(testMsg))
	suite.Equal(errUnauthorized, APIErrUnauthorized(testMsg))
	suite.Equal(errNotFound, APIErrNotFound(testMsg))
	suite.Equal(errConflict, APIErrConflict(testMsg, testMsg, testMsg))
	suite.Equal(errMissingRequired, APIErrEmptyRequiredField(testMsg))
	suite.Equal(errInvalidField, APIErrInvalidFieldContent(testMsg, "reason"))
	suite.Equal(errUnsupportedContent, APIErrUnsupportedContent(testPlc, testMsg))
	suite.Equal(errDatabase, APIErrDatabase(testMsg))
	suite.Equal(errGenericMissing, GenericEmptyRequiredField(testMsg, "errMsg").Error())
	suite.Equal(errGenericInternal, APIGenericInternalError(testMsg))
}

func TestAPIErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(APIErrorsTestSuite))
}
