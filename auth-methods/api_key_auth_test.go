package auth_methods

import (
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestApiKeyAuthSuite struct {
	suite.Suite
}

func (suite *TestApiKeyAuthSuite) TestFindApiKeyMethod() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// normal case
	expMap1 := map[string]interface{}{"type": "api-key", "service_uuid": "uuid1", "host": "host1", "port": 9000.0, "path": "test_path_1", "access_key": "key1"}
	apk1, err1 := FindApiKeyAuthMethod("uuid1", "host1", mockstore)

	// not found case
	apk2, err2 := FindApiKeyAuthMethod("no service", "no host", mockstore)

	// path not declared case
	apk3, err3 := FindApiKeyAuthMethod("uuid2", "host3", mockstore)

	// port not declared case
	apk4, err4 := FindApiKeyAuthMethod("uuid2", "host4", mockstore)

	// access key not declared case
	apk5, err5 := FindApiKeyAuthMethod("uuid1", "host2", mockstore)

	suite.Equal(expMap1["port"], apk1["port"])
	suite.Equal(expMap1["type"], apk1["type"])
	suite.Equal(expMap1["host"], apk1["host"])
	suite.Equal(expMap1["path"], apk1["path"])
	suite.Equal(expMap1["access_key"], apk1["access_key"])
	suite.Equal(expMap1["service"], apk1["service"])

	suite.Equal(0, len(apk2))
	suite.Equal(0, len(apk3))
	suite.Equal(0, len(apk4))
	suite.Equal(0, len(apk5))

	suite.Nil(err1)
	suite.Equal("Auth method was not found", err2.Error())
	suite.Equal("Database Error: Path was not found in the ApiKeyAuth object", err3.Error())
	suite.Equal("Database Error: Port was not found in the ApiKeyAuth object", err4.Error())
	suite.Equal("Database Error: Access Key was not found in the ApiKeyAuth object", err5.Error())

}

func (suite *TestApiKeyAuthSuite) TestCreateApiKeyMethod() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// tests the normal case
	am := map[string]interface{}{"type": "api-key", "service_uuid": "s_temp", "host": "h_temp", "port": 9000.0, "path": "/path/{{identifier}}?key={{access_key}}", "access_key": "key1"}
	_, err1 := CreateApiKeyAuthMethod(am, mockstore)
	am1, _ := FindApiKeyAuthMethod("s_temp", "h_temp", mockstore)

	// tests the case where access_key is not included
	AmMissingAccessKey := map[string]interface{}{"type": "api-key", "service_uuid": "s_temp", "host": "h_temp", "port": 9000.0, "path": "/path/{{identifier}}?key={{access_key}}"}
	_, err2 := CreateApiKeyAuthMethod(AmMissingAccessKey, mockstore)

	// tests the case where {{identifier}} is missing from path
	AmInvalidPath1 := map[string]interface{}{"type": "api-key", "service_uuid": "s_temp", "host": "h_temp", "port": 9000.0, "path": "/path/?key={{access_key}}", "access_key": "key1"}
	_, err3 := CreateApiKeyAuthMethod(AmInvalidPath1, mockstore)

	// tests the case where {{access_key}} is missing from path
	AmInvalidPath2 := map[string]interface{}{"type": "api-key", "service_uuid": "s_temp", "host": "h_temp", "port": 9000.0, "path": "/path/{{identifier}}?key=", "access_key": "key1"}
	_, err4 := CreateApiKeyAuthMethod(AmInvalidPath2, mockstore)

	suite.Equal(am, am1)
	suite.Equal("access_key was not found in the request body", err2.Error())
	suite.Equal("Field: path contains invalid data. Doesn't contain {{identifier}}", err3.Error())
	suite.Equal("Field: path contains invalid data. Doesn't contain {{access_key}}", err4.Error())

	suite.Nil(err1)
}

func TestApiKeyAuthTestSuite(t *testing.T) {
	suite.Run(t, new(TestApiKeyAuthSuite))
}
