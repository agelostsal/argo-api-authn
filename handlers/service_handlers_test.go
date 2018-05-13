package handlers

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"

	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/services"
	"github.com/ARGOeu/argo-api-authn/stores"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http/httptest"
)

type ServiceHandlersSuite struct {
	suite.Suite
}

// TestServiceCreate tests the normal case of a service creation
func (suite *ServiceHandlersSuite) TestServiceCreate() {

	postJSON := `{
	"name": "service1",
	"hosts": ["127.0.0.1"],
	"auth_types": ["x509", "oidc"],
	"auth_method": "api-key",
	"retrieval_field": "token"
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(201, w.Code)
	createdSer := services.Service{}
	_ = json.Unmarshal([]byte(w.Body.String()), &createdSer)

	suite.Equal("service1", createdSer.Name)
	suite.Equal([]string{"127.0.0.1"}, createdSer.Hosts)
	suite.Equal([]string{"x509", "oidc"}, createdSer.AuthTypes)
	suite.Equal("api-key", createdSer.AuthMethod)
	suite.Equal("token", createdSer.RetrievalField)

}

// TestServiceCreateInvalidName tests the case where the project's name already exists
func (suite *ServiceHandlersSuite) TestServiceCreateInvalidName() {

	postJSON := `{
	"name": "s1",
	"hosts": ["127.0.0.1"],
	"auth_types": ["x509", "oidc"],
	"auth_method": "api-key",
	"retrieval_field": "token"
}`

	expRespJSON := `{
 "error": {
  "message": "services.Service object with name: s1 already exists",
  "status_code": 409,
  "status": "CONFLICT"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(409, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestServiceCreateInvalidAuthTypes tests the case where the project's auth types are not yet supported
func (suite *ServiceHandlersSuite) TestServiceCreateInvalidAuthTypes() {

	postJSON := `{
	"name": "s1",
	"hosts": ["127.0.0.1"],
	"auth_types": ["unsup_type", "oidc"],
	"auth_method": "api-key",
	"retrieval_field": "token"
}`

	expRespJSON := `{
 "error": {
  "message": "Authentication Type: unsup_type is not yet supported",
  "status_code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestServiceCreateInvalidAuthMethod tests the case where the project's auth method are not yet supported
func (suite *ServiceHandlersSuite) TestServiceCreateInvalidAuthMethod() {

	postJSON := `{
	"name": "s1",
	"hosts": ["127.0.0.1"],
	"auth_types": ["x509", "oidc"],
	"auth_method": "unsup_method",
	"retrieval_field": "token"
}`

	expRespJSON := `{
 "error": {
  "message": "Authentication Method: unsup_method is not yet supported",
  "status_code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestServiceCreateEmptyAuthTypes tests the case where the project's auth types are empty
func (suite *ServiceHandlersSuite) TestServiceCreateEmptyAuthTypes() {

	postJSON := `{
	"name": "s1",
	"hosts": ["127.0.0.1"],
	"auth_types": [],
	"auth_method": "api-key",
	"retrieval_field": "token"
}`

	expRespJSON := `{
 "error": {
  "message": "Authentication Type: empty is not yet supported",
  "status_code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestServiceCreateInvalidJSON tests the case of the request containing an invalid json body
func (suite *ServiceHandlersSuite) TestServiceCreateInvalidJSON() {

	postJSON := `{
	"name": "service1",
	"hosts": ["127.0.0.1"],
	"auth_types": ["x509", "oidc"],
	"auth_method": "api-key",
	"retrieval_field": "token",
}`

	expResJSON := `{
 "error": {
  "message": "Poorly formatted JSON. invalid character '}' looking for beginning of object key string",
  "status_code": 400,
  "status": "BAD REQUEST"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(400, w.Code)
	suite.Equal(expResJSON, w.Body.String())
}

// TestServiceCreateMissingField tests the case of the request containing an incomplete json body
func (suite *ServiceHandlersSuite) TestServiceCreateMissingField() {

	postJSON := `{
	"hosts": ["127.0.0.1"],
	"auth_types": ["x509", "oidc"],
	"auth_method": "api-key",
	"retrieval_field": "token"
}`

	expResJSON := `{
 "error": {
  "message": "services.Service object contains an empty value for field: Name",
  "status_code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/services", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expResJSON, w.Body.String())
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceHandlersSuite))
}
