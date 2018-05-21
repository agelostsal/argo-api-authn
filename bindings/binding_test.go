package bindings

import (
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BindingTestSuite struct {
	suite.Suite
}

func (suite *BindingTestSuite) TestCreateBinding() {

		mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
		mockstore.SetUp()

		// test the normal case
		b1 := Binding{Name: "bins", ServiceUUID: "uuid1", Host: "host1", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err1 := CreateBinding(b1, mockstore)
		res1, _ := mockstore.QueryBindingsByDN("dn_ins", "uuid1", "host1")

		// tests the case of missing field name
		b2 := Binding{ServiceUUID: "uuid1", Host: "host1", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err2 := CreateBinding(b2, mockstore)

		// tests the case with missing field serviceUUID
		b3 := Binding{Name: "bins", Host: "host1", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err3 := CreateBinding(b3, mockstore)

		// tests the case with missing field host
		b4 := Binding{Name: "bins", ServiceUUID: "uuid1", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err4 := CreateBinding(b4, mockstore)

		// tests the case with missing field uniquekey
		b5 := Binding{Name: "bins", ServiceUUID: "uuid1", Host: "host1", DN: "dn_ins", OIDCToken: ""}
		_, err5 := CreateBinding(b5, mockstore)

		// tests the case with missing field dn and oidctoken
		b6 := Binding{Name: "bins", ServiceUUID: "uuid1", Host: "host1", UniqueKey: "key"}
		_, err6 := CreateBinding(b6, mockstore)

		// tests the case with unknown service uuid
		b7 := Binding{Name: "bins", ServiceUUID: "unknown", Host: "host1", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err7 := CreateBinding(b7, mockstore)

		// tests the case with unknown host
		b8 := Binding{Name: "bins", ServiceUUID: "uuid1", Host: "unknown", DN: "dn_ins", OIDCToken: "", UniqueKey: "key"}
		_, err8 := CreateBinding(b8, mockstore)

		// tests the case where a binding with the given dn already exists
		b9 := Binding{Name: "bins", ServiceUUID: "uuid1", Host: "host1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "key"}
		_, err9 := CreateBinding(b9, mockstore)

		suite.Equal(b1.Name, res1[0].Name)
		suite.Equal(b1.ServiceUUID, res1[0].ServiceUUID)
		suite.Equal(b1.Host, res1[0].Host)
		suite.Equal(b1.OIDCToken, res1[0].OIDCToken)
		suite.Equal(b1.UniqueKey, res1[0].UniqueKey)

		suite.Nil(err1)
		suite.Equal("bindings.Binding object contains an empty value for field: Name", err2.Error())
		suite.Equal("bindings.Binding object contains an empty value for field: ServiceUUID", err3.Error())
		suite.Equal("bindings.Binding object contains an empty value for field: Host", err4.Error())
		suite.Equal("bindings.Binding object contains an empty value for field: UniqueKey", err5.Error())
		suite.Equal("Both DN and OIDC Token fields are empty", err6.Error())
		suite.Equal("ServiceType was not found", err7.Error())
		suite.Equal("Host was not found", err8.Error())
		suite.Equal("bindings.Binding object with dn: test_dn_1 already exists", err9.Error())

}

func (suite *BindingTestSuite) TestFindBindingByDN() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// tests the normal case
	expBinding1 := Binding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "unique_key_1"}
	b1, err1 := FindBindingByDN("test_dn_1","uuid1", "host1", mockstore)

	suite.Equal(expBinding1.Name, b1.Name)
	suite.Equal(expBinding1.ServiceUUID, b1.ServiceUUID)
	suite.Equal(expBinding1.Host, b1.Host)
	suite.Equal(expBinding1.DN, b1.DN)
	suite.Equal(expBinding1.UniqueKey, b1.UniqueKey)

	// tests the not found case
	_, err2 := FindBindingByDN("unknown", "unknown", "unknown", mockstore)

	// tests the case of more than 2 bindigs with the same dn exist under the same host and service
	// append one more binding , same to an existing
	mockstore.Bindings = append(mockstore.Bindings, stores.QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "unique_key_1"})
	_, err3 := FindBindingByDN("test_dn_1","uuid1", "host1", mockstore)


	suite.Nil(err1)
	suite.Equal("Binding was not found", err2.Error())
	suite.Equal("Database Error: More than 1 bindings found under the service type: uuid1 and host: host1 using the same DN: test_dn_1", err3.Error())

}

func TestBindingTestSuite(t *testing.T) {
	suite.Run(t, new(BindingTestSuite))
}
