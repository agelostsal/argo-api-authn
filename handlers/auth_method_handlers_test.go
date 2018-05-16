package handlers

import (
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthMethodHandlersTestSuite struct {
	suite.Suite
}

// TestAuthMethodListOne tests the normal case and returns the information of the auth method under the given service and host
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOne() {

	expRespJSON := `{
 "access_key": "key1",
 "host": "host1",
 "path": "test_path_1",
 "port": 9000,
 "service": "s1",
 "type": "api-key"
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s1/hosts/host1/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

//TestAuthMethodListOneUndeclaredAccessKey tests the case where the auth method doesn't contain the required access key
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOneUndeclaredAccessKey() {

	expRespJSON := `{
 "error": {
  "message": "Database Error: Access Key was not found in the ApiKeyAuth object",
  "code": 500,
  "status": "INTERNAL SERVER ERROR"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s1/hosts/host2/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(500, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

//TestAuthMethodListOneUndeclaredPath tests the case where the auth method doesn't contain the required path
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOneUndeclaredPath() {

	expRespJSON := `{
 "error": {
  "message": "Database Error: Path was not found in the ApiKeyAuth object",
  "code": 500,
  "status": "INTERNAL SERVER ERROR"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s2/hosts/host3/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(500, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

//TestAuthMethodListOneUndeclaredPort tests the case where the auth method doesn't contain the required port
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOneUndeclaredPort() {

	expRespJSON := `{
 "error": {
  "message": "Database Error: Port was not found in the ApiKeyAuth object",
  "code": 500,
  "status": "INTERNAL SERVER ERROR"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s2/hosts/host4/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(500, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestAuthMethodListOneUnknownService tests the case where the given service doesn't exist
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOneUnknownService() {

	expRespJSON := `{
 "error": {
  "message": "Service was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/unknown_service/hosts/host4/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestAuthMethodListOneUnknownHost tests the case where the given host is associated with the given service
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListOneUnknownHost() {

	expRespJSON := `{
 "error": {
  "message": "Host was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	req, err := http.NewRequest("GET", "http://localhost:8080/services/s1/hosts/host_unknown/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/services/{service}/hosts/{host}/authM", WrapConfig(AuthMethodListOne, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())

}

// TestAuthMethodListAll tests the normal case and returns all auth methods in the service
func (suite *AuthMethodHandlersTestSuite) TestAuthMethodListAll() {

	expRespJSON := `{
 "auth_methods": [
  {
   "access_key": "key1",
   "host": "host1",
   "path": "test_path_1",
   "port": 9000,
   "service": "s1",
   "type": "api-key"
  },
  {
   "host": "host2",
   "path": "test_path_1",
   "port": 9000,
   "service": "s1",
   "type": "api-key"
  },
  {
   "access_key": "key1",
   "host": "host3",
   "port": 9000,
   "service": "s2",
   "type": "api-key"
  },
  {
   "access_key": "key1",
   "host": "host4",
   "path": "test_path_1",
   "service": "s2",
   "type": "api-key"
  }
 ]
}`
	req, err := http.NewRequest("GET", "http://localhost:8080/authM", nil)
	if err != nil {
		log.Error(err.Error())
	}

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/authM", WrapConfig(AuthMethodListAll, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

func TestAuthMethodHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMethodHandlersTestSuite))
}
