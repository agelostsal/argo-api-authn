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
	expMap1 := map[string]interface{}{"type": "api-key", "service": "s1", "host": "host1", "port": 9000.0, "path": "test_path_1", "access_key": "key1"}
	apk1, err1 := FindApiKeyMethod("s1", "host1", mockstore)

	// not found case
	apk2, err2 := FindApiKeyMethod("no service", "no host", mockstore)

	// path not declared case
	apk3, err3 := FindApiKeyMethod("s2", "host3", mockstore)

	// port not declared case
	apk4, err4 := FindApiKeyMethod("s2", "host4", mockstore)

	// access key not declared case
	apk5, err5 := FindApiKeyMethod("s1", "host2", mockstore)

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

func TestApiKeyAuthTestSuite(t *testing.T) {
	suite.Run(t, new(TestApiKeyAuthSuite))
}
