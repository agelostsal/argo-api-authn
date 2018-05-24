package handlers

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"net/http"
	"bytes"
	"github.com/ARGOeu/argo-api-authn/stores"
	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"net/http/httptest"
	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/bindings"
)

type BindingHandlersSuite struct {
	suite.Suite
}

// TestBindingCreate tests the normal case of a binding creation
func (suite *BindingHandlersSuite) TestBindingCreate(){

	postJSON := `{
	"name": "new_binding",
	"service_uuid": "uuid1",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key": "uni_key"
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(201, w.Code)
	createdBind := bindings.Binding{}
	_ = json.Unmarshal([]byte(w.Body.String()), &createdBind)

	suite.Equal("uuid1", createdBind.ServiceUUID)
	suite.Equal("host1", createdBind.Host)
	suite.Equal("new_binding", createdBind.Name)
	suite.Equal("test_dn", createdBind.DN)
	suite.Equal("", createdBind.OIDCToken)
	suite.Equal("uni_key", createdBind.UniqueKey)

}

// TestBindingCreateMissingFieldName tests the case where the binding doesn't contain the  name field
func (suite *BindingHandlersSuite) TestBindingCreateMissingFieldName(){

	postJSON := `{
	"service_uuid": "uuid1",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key": "uni_key"
}`

	expRespJSON := `{
 "error": {
  "message": "bindings.Binding object contains an empty value for field: Name",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateInvalidJSON tests the case where the request body is not a vlaid json
func (suite *BindingHandlersSuite) TestBindingCreateInvalidJSON(){

	postJSON := `{
	"service_uuid": "uuid1",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key": "uni_key"
`

	expRespJSON := `{
 "error": {
  "message": "Poorly formatted JSON. unexpected EOF",
  "code": 400,
  "status": "BAD REQUEST"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(400, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateMissingFieldServiceUUID tests the case where the binding doesn't contain the service_uuid field
func (suite *BindingHandlersSuite) TestBindingCreateMissingFieldServiceUUID(){

	postJSON := `{
    "name": "b1",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key": "uni_key"
}`

	expRespJSON := `{
 "error": {
  "message": "bindings.Binding object contains an empty value for field: ServiceUUID",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateMissingFieldHost tests the case where the binding doesn't contain the host field
func (suite *BindingHandlersSuite) TestBindingCreateMissingFieldHost(){

	postJSON := `{
    "name": "b1",
	"service_uuid": "uuid1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key": "uni_key"
}`

	expRespJSON := `{
 "error": {
  "message": "bindings.Binding object contains an empty value for field: Host",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateMissingFieldDNAndOIDC tests the case where the binding doesn't contain both dn and oidc fields
func (suite *BindingHandlersSuite) TestBindingCreateMissingFieldDNAndOIDC(){

	postJSON := `{
    "name": "b1",
	"service_uuid": "uuid1",
    "host": "host1",
    "unique_key": "uni_key"
}`

	expRespJSON := `{
 "error": {
  "message": "Both DN and OIDC Token fields are empty",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateMissingFieldUniqueKey tests the case where the binding doesn't contain the service_uuid field
func (suite *BindingHandlersSuite) TestBindingCreateMissingFieldUniqueKey(){

	postJSON := `{
    "name": "b1",
    "service_uuid": "uuid1",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": ""
}`

	expRespJSON := `{
 "error": {
  "message": "bindings.Binding object contains an empty value for field: UniqueKey",
  "code": 422,
  "status": "UNPROCESSABLE ENTITY"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(422, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateUnknownService tests the case where the service uuid is not known
func (suite *BindingHandlersSuite) TestBindingCreateUnknownService(){

	postJSON := `{
    "name": "b1",
    "service_uuid": "unknown",
    "host": "host1",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key":"key1"
}`

	expRespJSON := `{
 "error": {
  "message": "ServiceType was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateUnknownHost tests the case where the host is not known to the specified service
func (suite *BindingHandlersSuite) TestBindingCreateUnknownHost(){

	postJSON := `{
    "name": "b1",
    "service_uuid": "uuid1",
    "host": "unknown",
    "dn": "test_dn",
    "oidc_token": "",
    "unique_key":"key1"
}`

	expRespJSON := `{
 "error": {
  "message": "Host was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestBindingCreateDNAlreadyExists tests the case where the given dn is already used by another binding
func (suite *BindingHandlersSuite) TestBindingCreateDNAlreadyExists(){

	postJSON := `{
    "name": "b1",
    "service_uuid": "uuid1",
    "host": "host1",
    "dn": "test_dn_1",
    "oidc_token": "",
    "unique_key":"key1"
}`

	expRespJSON := `{
 "error": {
  "message": "bindings.Binding object with dn: test_dn_1 already exists",
  "code": 409,
  "status": "CONFLICT"
 }
}`

	req, err := http.NewRequest("POST", "http://localhost:8080/bindings", bytes.NewBuffer([]byte(postJSON)))
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/bindings", WrapConfig(BindingCreate, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(409, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

func TestBindingHandlersSuite(t *testing.T){
	suite.Run(t, new(BindingHandlersSuite))
}