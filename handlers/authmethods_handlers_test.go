package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/authmethods"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/gorilla/mux"
	LOGGER "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthMethodsHandlersTestSuite struct {
	suite.Suite
}

// TestAuthMethodCreate tests the default case of creating an auth method of type api-key and service type of ams
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreate() {

	var expAm = &authmethods.ApiKeyAuthMethod{}

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(201, w.Code)

	// unmarshal the response
	json.Unmarshal([]byte(w.Body.String()), expAm)
	suite.Equal("uuid1", expAm.ServiceUUID)
	suite.Equal("host2", expAm.Host)
	suite.Equal(9000, expAm.Port)
	suite.Equal("token", expAm.RetrievalField)
	suite.Equal("/v1/users:byUUID/{{identifier}}?key={{access_key}}", expAm.Path)
	suite.Equal("api-key", expAm.Type)
	suite.NotEqual("", expAm.UUID)
	suite.NotEqual("", expAm.CreatedOn)
}

// TestAuthMethodCreateOverrideDefaults tests the default case of creating an auth method of type api-key and service type of ams while overriding path and auth retrieval field
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaults() {

	var expAm = &authmethods.ApiKeyAuthMethod{}

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "path": "/some/other/{{identifier}}?key={{access_key}}",
 "retrieval_field": "some_token"
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(201, w.Code)

	// unmarshal the response
	json.Unmarshal([]byte(w.Body.String()), expAm)
	suite.Equal("uuid1", expAm.ServiceUUID)
	suite.Equal("host2", expAm.Host)
	suite.Equal(9000, expAm.Port)
	suite.Equal("some_token", expAm.RetrievalField)
	suite.Equal("/some/other/{{identifier}}?key={{access_key}}", expAm.Path)
	suite.NotEqual("", expAm.UUID)
	suite.NotEqual("", expAm.CreatedOn)
}

// TestAuthMethodCreateOverrideDefaultsServiceUUID tests the default case of creating an auth method of type api-key and service type of ams while overriding the service uuid WHICH shall not work
// the service uuid is assigned through the specified service on the request
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsServiceUUID() {

	var expAm = &authmethods.ApiKeyAuthMethod{}

	reqBody := `{
 "access_key": "key1",
 "host": "host3",
 "port": 9000,
 "service_uuid": "uuid2",
 "path": "/some/other/{{identifier}}?key={{access_key}}",
 "retrieval_field": "some_token"
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(201, w.Code)

	// unmarshal the response
	json.Unmarshal([]byte(w.Body.String()), expAm)
	suite.Equal("uuid1", expAm.ServiceUUID) //uuid2 is being ignored
	suite.Equal("host3", expAm.Host)
	suite.Equal(9000, expAm.Port)
	suite.Equal("some_token", expAm.RetrievalField)
	suite.Equal("/some/other/{{identifier}}?key={{access_key}}", expAm.Path)
	suite.Equal("api-key", expAm.Type)
	suite.NotEqual("", expAm.UUID)
	suite.NotEqual("", expAm.CreatedOn)
}

// TestAuthMethodCreateOverrideDefaultsInvalidPathAccessKey tests the default case of creating an auth method of type api-key and service type of ams while overriding path into an invalid one
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsInvalidPathAccessKey() {

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "path": "/some/other/{{identifier}}?key=",
 "retrieval_field": "some_token"
}`

	expRespJSON := `{
 "error": {
  "message": "Field: path contains invalid data. Missing {{access_key}} interpolation",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateOverrideDefaultsInvalidPathAccessKey tests the default case of creating an auth method of type api-key and service type of ams while overriding path into an invalid one
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsInvalidPathIdentifier() {

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "path": "/some/other/?key={{access_key}}",
 "retrieval_field": "some_token"
}`

	expRespJSON := `{
 "error": {
  "message": "Field: path contains invalid data. Missing {{identifier}} interpolation",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateOverrideDefaultsInvalidEmptyPath tests the default case of creating an auth method of type api-key and service type of ams while overriding path into an invalid one
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsInvalidEmptyPath() {

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "path": "",
 "retrieval_field": "some_token"
}`

	expRespJSON := `{
 "error": {
  "message": "auth method object contains empty fields. empty value for field: path",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateOverrideDefaultsInvalidEmptyRetrievalField tests the default case of creating an auth method of type api-key and service type of ams while overriding path into an invalid one
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsInvalidEmptyRetrievalField() {

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "retrieval_field": ""
}`

	expRespJSON := `{
 "error": {
  "message": "auth method object contains empty fields. empty value for field: retrieval_field",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateInvalidURL tests the default case of creating an auth method of type api-key and service type of ams where host + path don't resemble a valid url
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateOverrideDefaultsInvalidURL() {

	reqBody := `{
 "access_key": "key1",
 "host": "host2",
 "port": 9000,
 "path": ":::fff/"
}`

	expRespJSON := `{
 "error": {
  "message": "The url to access resources in invalid. URL: https://host2:9000:::fff/",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateAlreadyExists tests the case where there is an already existing auth method for the given service type and host
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateAlreadyExists() {

	reqBody := `{
 "access_key": "key1",
 "host": "host1",
 "port": 9000
}`

	expRespJSON := `{
 "error": {
  "message": "Auth method object with host: host1 already exists",
  "code": 409,
  "status": "CONFLICT"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(409, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateUnknownServiceType tests the case where the provided service type is unknown
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateUnknownServiceType() {

	reqBody := `{
 "access_key": "key1",
 "host": "host1",
 "port": 9000
}`

	expRespJSON := `{
 "error": {
  "message": "Service-type was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/unknown/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateUnknownHost tests the case where the provided host isn't declared for the given service type
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateUnknownHost() {

	reqBody := `{
 "access_key": "key1",
 "host": "unknown",
 "port": 9000
}`

	expRespJSON := `{
 "error": {
  "message": "Host was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateInvalidJSON tests the case where the request body is malformed
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateInvalidJSON() {

	reqBody := `{
 "access_key": "key1",
 "host": "host1",
 "port": 9000
` // missing closing bracket

	expRespJSON := `{
 "error": {
  "message": "Poorly formatted JSON. unexpected EOF",
  "code": 400,
  "status": "BAD REQUEST"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(400, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateEmptyAccessKey tests the case where the request body contains an empty data for field access_key
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateEmptyAccessKey() {

	reqBody := `{
 "host": "host2",
 "port": 9000
}`

	expRespJSON := `{
 "error": {
  "message": "auth method object contains empty fields. empty value for field: access_key",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateEmptyHost tests the case where the request body contains an empty data for field host
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateEmptyHost() {

	reqBody := `{
 "access_key": "key",
 "port": 9000
}`

	expRespJSON := `{
 "error": {
  "message": "auth method object contains empty fields. empty value for field: host",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodCreateEmptyHost tests the case where the request body contains an empty data for field port
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodCreateEmptyPort() {

	reqBody := `{
 "access_key": "key",
 "host": "host1,"
}`

	expRespJSON := `{
 "error": {
  "message": "auth method object contains empty fields. empty value for field: port",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/service-types/s1/authm", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/authm", WrapConfig(AuthMethodCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListOne tests the normal case of finding the auth method associated with the given service type and host
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListOne() {

	expRespJSON := `{
 "service_uuid": "uuid1",
 "port": 9000,
 "host": "host1",
 "retrieval_field": "",
 "path": "test_path_1",
 "type": "api-key",
 "uuid": "am_uuid_1",
 "created_on": "",
 "access_key": "access_key"
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/service-types/s1/hosts/host1/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}/authm", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListOneUnknownServiceType tests the case where the provided service type is unknown
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListOneUnknownServiceType() {

	expRespJSON := `{
 "error": {
  "message": "Service-type was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/service-types/unknown/hosts/host1/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}/authm", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListOneUnknownHost tests the case where the provided host is not associated with the given service-type
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListOneUnknownHost() {

	expRespJSON := `{
 "error": {
  "message": "Host was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/service-types/s1/hosts/unknown/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}/authm", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListOneNotFound tests the case where there is no registered auth method under the given service-type and host
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListOneNotFound() {

	expRespJSON := `{
 "error": {
  "message": "Auth method was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/service-types/s1/hosts/host3/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}/authm", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListAll tests the normal case and returns all auth methods in the service type
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListAll() {

	expRespJSON := `{
 "auth_methods": [
  {
   "service_uuid": "uuid1",
   "port": 9000,
   "host": "host1",
   "retrieval_field": "",
   "path": "test_path_1",
   "type": "api-key",
   "uuid": "am_uuid_1",
   "created_on": "",
   "access_key": "access_key"
  }
 ]
}`
	req, err := http.NewRequest("GET", "http://localhost:8080/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/authm", WrapConfig(AuthMethodListAll, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthMethodListAllEmptyList tests the normal case where there are no auth methods in the service yet
func (suite *AuthMethodsHandlersTestSuite) TestAuthMethodListAllEmptyList() {

	expRespJSON := `{
 "auth_methods": []
}`
	req, err := http.NewRequest("GET", "http://localhost:8080/authm", nil)
	if err != nil {
		LOGGER.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// empty the mockstore
	mockstore.AuthMethods = []stores.QAuthMethod{}

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/authm", WrapConfig(AuthMethodListAll, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

func TestAuthMethodsHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMethodsHandlersTestSuite))
}
