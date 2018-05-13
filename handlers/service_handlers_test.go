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
  "code": 409,
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
  "code": 422,
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
  "code": 422,
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
  "code": 422,
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
  "code": 400,
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
  "code": 422,
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

// TestServiceListOne tests the normal case
func (suite *ServiceHandlersSuite) TestServiceListOne() {

	expResJSON := `{
 "name": "s1",
 "hosts": [
  "host1",
  "host2",
  "host3"
 ],
 "auth_types": [
  "x509",
  "oidc"
 ],
 "auth_method": "api-key",
 "retrieval_field": "token",
 "created_on": "2018-05-05T18:04:05Z"
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s1", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{name}", WrapConfig(ServiceListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expResJSON, w.Body.String())

}

// TestServiceListOneNameCollision tests the case where two or more service exist with the same name
func (suite *ServiceHandlersSuite) TestServiceListOneNameCollision() {

	expResJSON := `{
 "error": {
  "message": "Database Error: Multiple services with the same name: same_name",
  "code": 500,
  "status": "INTERNAL SERVER ERROR"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/same_name", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{name}", WrapConfig(ServiceListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(500, w.Code)
	suite.Equal(expResJSON, w.Body.String())

}

// TestServiceListOneNotFound tests the case where two or more service exist with the same name
func (suite *ServiceHandlersSuite) TestServiceListOneNotFound() {

	expResJSON := `{
 "error": {
  "message": "Service was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/not_found", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{name}", WrapConfig(ServiceListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expResJSON, w.Body.String())

}

// TestServiceListAll lists all services
func (suite *ServiceHandlersSuite) TestServiceListAll() {

	expResJSON := `{
 "services": [
  {
   "name": "s1",
   "hosts": [
    "host1",
    "host2",
    "host3"
   ],
   "auth_types": [
    "x509",
    "oidc"
   ],
   "auth_method": "api-key",
   "retrieval_field": "token",
   "created_on": "2018-05-05T18:04:05Z"
  },
  {
   "name": "s2",
   "hosts": [
    "host3",
    "host4"
   ],
   "auth_types": [
    "x509"
   ],
   "auth_method": "api-key",
   "retrieval_field": "user_token",
   "created_on": "2018-05-05T18:04:05Z"
  },
  {
   "name": "same_name",
   "hosts": null,
   "auth_types": null,
   "auth_method": "",
   "retrieval_field": "",
   "created_on": ""
  },
  {
   "name": "same_name",
   "hosts": null,
   "auth_types": null,
   "auth_method": "",
   "retrieval_field": "",
   "created_on": ""
  }
 ]
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services", nil)

	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceListAll, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expResJSON, w.Body.String())

}

// TestsServiceListAllEmptyList
func (suite *ServiceHandlersSuite) TestServiceListAllEmptyList() {

	expResJSON := `{
 "services": null
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()
	// empty the store
	mockstore.Services = []stores.QService{}

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services", WrapConfig(ServiceListAll, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expResJSON, w.Body.String())

}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceHandlersSuite))
}
