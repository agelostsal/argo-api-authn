package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/ARGOeu/argo-api-authn/auth-methods"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CertificateHandlerSuite struct {
	suite.Suite
}

// MockServiceTypeEndpoint mocks the behavior of a service type endpoint and returns a response containing the requested resource
func MockServiceTypeEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{\"token\": \"some-value\"}"))
}

func AuthViaCertSetUp(reqPath string) (*http.Request, *stores.Mockstore, *config.Config, error) {

	var err error
	var crt *x509.Certificate
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	// create a new certificate from the string literal
	cert := `-----BEGIN CERTIFICATE-----
MIIC8jCCAdoCCQCdC824csOlXTANBgkqhkiG9w0BAQsFADA7MQswCQYDVQQGEwJU
QzELMAkGA1UECAwCVE4xDTALBgNVBAcMBENJVFkxEDAOBgNVBAoMB0NPTVBBTlkw
HhcNMTgwNjA0MDkwMzA5WhcNMTkwNjA0MDkwMzA5WjA7MQswCQYDVQQGEwJUQzEL
MAkGA1UECAwCVE4xDTALBgNVBAcMBENJVFkxEDAOBgNVBAoMB0NPTVBBTlkwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDYahRnf7gxkjz81VX9JjQ/PiB7
UpGInckBCl5mah/Q3ucr3OvaLLdO4pfmAMTanLxgcTsP7k/yyvZF17IMhQ5wpNzn
zfLlAswQT6sqkNyJUx1MXI5mjAiDDMpUh0c9CnGaZa/LrnTXmQhsv8uzLDPUYb37
iHHjP7isYiG+7YE2CpRwHazj0SYba4HAYw8Z1L8Z6kI1gfdIiqI5DFMBdmQlac3P
YCtYytQd3swCsxf57/M9X+Ct7DVcuPKmR5vv4ONL2YBwtvULX8DA8aApdxIrHCZE
NRaaCUmMSo4JXHbY5CVP6AXdo3Iz+3v485qtF8lk+XU/fQGNTdgs1hmsTNMRAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAAyvW6yVbfCLMxuQ1Nt61OmKA96fTtdTpLuI
nq9C0XVoFqgrEeobxdH4QbwxduRHWpHEuFskJCBnMbX0d8v63KEN/6I0Ub4niaeP
nvykp3uoKrRwIZo4OxJFuuuLuUw3aAwkKeqsy5HZsKqi9QscHExRKbcIdlkgxzRW
IEPEGk7acMlT20ECjc4zbdza8PKQeBeEVLINJVMRGPZIlo/6z6BxSnINQiWBk1WZ
lXTmJMXGB7/0ECDz2JT1Mbs/q2ijlZywz0xsp+Zdsp1I01wvwqw5M12PHf3buM1w
SoPmZKiBeb+2OQ2n7+FI8ftkqxWw6zjh651brAoy/0zqLTRPh+c=
-----END CERTIFICATE-----
`
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		panic("failed to parse certificate PEM")
	}

	if crt, err = x509.ParseCertificate(block.Bytes); err != nil {
		log.Error(err.Error())
		return req, mockstore, cfg, err
	}

	// create a new request and add the created certificate
	if req, err = http.NewRequest("GET", reqPath, nil); err != nil {
		return req, mockstore, cfg, err
	}

	req.TLS = &tls.ConnectionState{}
	req.TLS.PeerCertificates = append(req.TLS.PeerCertificates, crt)

	// set up the mockstore
	mockstore = &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()
	// append a service type to be used only in auth via cert tests
	qSt := stores.QServiceType{Name: "s_auth_cert", Hosts: []string{"h1_auth_cert"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "mock-api-key", UUID: "uuid_auth_cert", RetrievalField: "token", CreatedOn: "2018-05-05T18:04:05Z"}
	qSt2 := stores.QServiceType{Name: "s_auth_cert_incorrect", Hosts: []string{"h1_auth_cert"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "mock-api-key", UUID: "uuid_auth_cert_incorrect", RetrievalField: "incorrect_field", CreatedOn: "2018-05-05T18:04:05Z"}
	mockstore.ServiceTypes = append(mockstore.ServiceTypes, qSt, qSt2)
	// append a binding to be used only in auth via cert tests
	qB := stores.QBinding{Name: "b_auth_cert", ServiceUUID: "uuid_auth_cert", Host: "h1_auth_cert", DN: "O=COMPANY,L=CITY,ST=TN,C=TC", OIDCToken: "", UniqueKey: "unique_key_1", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	qB2 := stores.QBinding{Name: "b_auth_cert_incorrect", ServiceUUID: "uuid_auth_cert_incorrect", Host: "h1_auth_cert", DN: "O=COMPANY,L=CITY,ST=TN,C=TC", OIDCToken: "", UniqueKey: "unique_key_1", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	mockstore.Bindings = append(mockstore.Bindings, qB, qB2)

	// set up cfg
	cfg = &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	// add a mock auth method handler
	auth_methods.AuthMethodHandlers["mock-api-key"] =
		func(data map[string]interface{}, store stores.Store, config *config.Config) (*http.Response, error) {

			var req2 *http.Request
			var err error

			// mock the request that will take place against the given service type
			if req2, err = http.NewRequest("GET", "http://localhost:8080/some_endpoint", nil); err != nil {
				log.Error(err.Error())
			}
			router := mux.NewRouter().StrictSlash(true)
			w := httptest.NewRecorder()
			router.HandleFunc("/some_endpoint", MockServiceTypeEndpoint)
			router.ServeHTTP(w, req2)
			return w.Result(), err
		}

	return req, mockstore, cfg, err
}

// TestAuthViaCert tests the normal case
func (suite *CertificateHandlerSuite) TestAuthViaCert() {

	var err error
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	expRespJSON := `{
 "token": "some-value"
}`

	if req, mockstore, cfg, err = AuthViaCertSetUp("http://localhost:8080/service-types/s_auth_cert/hosts/h1_auth_cert:authX509"); err != nil {
		log.Error(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}:authX509", WrapConfig(AuthViaCert, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(200, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthViaCertIncorrectRetrievalField tests the case where the response from the service type didn't contain the specified retrieval field
func (suite *CertificateHandlerSuite) TestAuthViaCertIncorrectRetrievalField() {

	var err error
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	expRespJSON := `{
 "error": {
  "message": "Internal Error: The specified retrieval field: incorrect_field was not found in the response body of the service type",
  "code": 500,
  "status": "INTERNAL SERVER ERROR"
 }
}`

	if req, mockstore, cfg, err = AuthViaCertSetUp("http://localhost:8080/service-types/s_auth_cert_incorrect/hosts/h1_auth_cert:authX509"); err != nil {
		log.Error(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}:authX509", WrapConfig(AuthViaCert, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(500, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthViaCertUnknownServiceType tests the case where the provided service type is unknown
func (suite *CertificateHandlerSuite) TestAuthViaCertUnknownServiceType() {

	var err error
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	expRespJSON := `{
 "error": {
  "message": "Service-type was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	if req, mockstore, cfg, err = AuthViaCertSetUp("http://localhost:8080/service-types/unknown/hosts/h1_auth_cert:authX509"); err != nil {
		log.Error(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}:authX509", WrapConfig(AuthViaCert, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthViaCertUnknownHost tests the case where the provided host is unknown
func (suite *CertificateHandlerSuite) TestAuthViaCertUnknownHost() {

	var err error
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	expRespJSON := `{
 "error": {
  "message": "Host was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	if req, mockstore, cfg, err = AuthViaCertSetUp("http://localhost:8080/service-types/s_auth_cert/hosts/unknown:authX509"); err != nil {
		log.Error(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}:authX509", WrapConfig(AuthViaCert, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

// TestAuthViaCertUnknownDN tests the case where the provided certificate's dn is not assigned to any binding
func (suite *CertificateHandlerSuite) TestAuthViaCertUnknownDN() {

	var err error
	var mockstore *stores.Mockstore
	var cfg *config.Config
	var req *http.Request

	expRespJSON := `{
 "error": {
  "message": "Binding was not found",
  "code": 404,
  "status": "NOT FOUND"
 }
}`

	if req, mockstore, cfg, err = AuthViaCertSetUp("http://localhost:8080/service-types/s_auth_cert/hosts/h1_auth_cert:authX509"); err != nil {
		log.Error(err.Error())
	}

	// empty the mockstore, so no dn will match
	mockstore.Bindings = []stores.QBinding{}

	router := mux.NewRouter().StrictSlash(true)
	w := httptest.NewRecorder()
	router.HandleFunc("/service-types/{service-type}/hosts/{host}:authX509", WrapConfig(AuthViaCert, mockstore, cfg))
	router.ServeHTTP(w, req)
	suite.Equal(404, w.Code)
	suite.Equal(expRespJSON, w.Body.String())
}

func TestAuthViaCert(t *testing.T) {
	suite.Run(t, new(CertificateHandlerSuite))
}
