package stores

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	Mockstore *Mockstore
}

// SetUpTestSuite assigns the mock store to be used in the querying tests
// It should be used on each test case so CRUD operations don't need to be reverted
func (suite *StoreTestSuite) SetUpStoreTestSuite() {

	mockstore := &Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()
	suite.Mockstore = mockstore
}

// TestSetUp tests if the mockstore setup has been completed successfully
func (suite *StoreTestSuite) TestSetUp() {

	suite.SetUpStoreTestSuite()

	mockstore := &Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	var qServices []QServiceType
	var qBindings []QBinding
	var qAuthms []QAuthMethod

	// Populate qServices
	service1 := QServiceType{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}
	service2 := QServiceType{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "headers", UUID: "uuid2", CreatedOn: "2018-05-05T18:04:05Z", Type: "web-api"}
	serviceSame1 := QServiceType{Name: "same_name"}
	serviceSame2 := QServiceType{Name: "same_name"}
	qServices = append(qServices, service1, service2, serviceSame1, serviceSame2)

	// Populate Bindings
	binding1 := QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding2 := QBinding{Name: "b2", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid2", AuthIdentifier: "test_dn_2", UniqueKey: "unique_key_2", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding3 := QBinding{Name: "b3", ServiceUUID: "uuid1", Host: "host2", UUID: "b_uuid3", AuthIdentifier: "test_dn_3", UniqueKey: "unique_key_3", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding4 := QBinding{Name: "b4", ServiceUUID: "uuid2", Host: "host3", UUID: "b_uuid4", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}

	qBindings = append(qBindings, binding1, binding2, binding3, binding4)

	// Populate AuthMethods
	amb1 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	am1 := &QApiKeyAuthMethod{QBasicAuthMethod: amb1, AccessKey: "access_key"}
	amb2 := QBasicAuthMethod{ServiceUUID: "uuid2", Host: "host3", Port: 9000, Type: "headers", UUID: "am_uuid_2", CreatedOn: ""}
	am2 := &QHeadersAuthMethod{QBasicAuthMethod: amb2, Headers: map[string]string{"x-api-key": "key-1", "Accept": "application/json"}}
	qAuthms = append(qAuthms, am1, am2)

	suite.Equal(mockstore.Session, true)
	suite.Equal(mockstore.Database, "test_db")
	suite.Equal(mockstore.Server, "localhost")
	suite.Equal(mockstore.ServiceTypes, qServices)
	suite.Equal(mockstore.Bindings, qBindings)
	suite.Equal(mockstore.AuthMethods, qAuthms)
}

func (suite *StoreTestSuite) TestClose() {

	suite.SetUpStoreTestSuite()

	suite.Mockstore.Close()
	suite.Equal(false, suite.Mockstore.Session)
}

func (suite *StoreTestSuite) TestClone() {

	suite.SetUpStoreTestSuite()

	tempStore := suite.Mockstore.Clone()

	suite.Equal(suite.Mockstore, tempStore)

}

func (suite *StoreTestSuite) TestQueryServiceTypes() {

	suite.SetUpStoreTestSuite()

	// normal case outcome - 1 service
	expQServices1 := []QServiceType{{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}}
	qServices1, err1 := suite.Mockstore.QueryServiceTypes("s1")
	expQServices2 := []QServiceType{{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "headers", UUID: "uuid2", CreatedOn: "2018-05-05T18:04:05Z", Type: "web-api"}}
	qServices2, err2 := suite.Mockstore.QueryServiceTypes("s2")

	// normal case outcome - all services
	expQServicesAll := []QServiceType{
		{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"},
		{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "headers", UUID: "uuid2", CreatedOn: "2018-05-05T18:04:05Z", Type: "web-api"},
		{Name: "same_name"},
		{Name: "same_name"},
	}
	qServicesAll, errAll := suite.Mockstore.QueryServiceTypes("")

	// was not found
	var expQService3 []QServiceType
	qServices3, err3 := suite.Mockstore.QueryServiceTypes("wrong_name")

	// tests the normal case - 1 service type
	suite.Equal(expQServices1, qServices1)
	suite.Nil(err1)
	suite.Equal(expQServices2, qServices2)
	suite.Nil(err2)

	// tests the normal case - all service types
	suite.Equal(expQServicesAll, qServicesAll)
	suite.Nil(errAll)

	// tests the not found case
	suite.Equal(expQService3, qServices3)
	suite.Nil(err3)
}

func (suite *StoreTestSuite) TestQueryServiceTypesByUUID() {

	suite.SetUpStoreTestSuite()

	// normal case outcome
	expQServices1 := []QServiceType{{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}}
	qServices1, err1 := suite.Mockstore.QueryServiceTypesByUUID("uuid1")

	// was not found
	var expQServices2 []QServiceType
	qServices2, err2 := suite.Mockstore.QueryServiceTypesByUUID("wrong_uuid")

	suite.Equal(expQServices1, qServices1)
	suite.Equal(expQServices2, qServices2)

	suite.Nil(err1)
	suite.Nil(err2)

}

func (suite *StoreTestSuite) TestQueryApiKeyAuthMethods() {

	suite.SetUpStoreTestSuite()

	// normal case
	var expApiAms []QApiKeyAuthMethod
	amb1 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	qapi := QApiKeyAuthMethod{amb1, "access_key"}
	expApiAms = append(expApiAms, qapi)
	apiAms, err1 := suite.Mockstore.QueryApiKeyAuthMethods("uuid1", "host1")

	// not found - empty list
	apiAms2, err2 := suite.Mockstore.QueryApiKeyAuthMethods("unknown", "unknown")

	// query all
	apiAms3, err3 := suite.Mockstore.QueryApiKeyAuthMethods("", "")

	suite.Equal(expApiAms, apiAms)
	suite.Equal(0, len(apiAms2))
	suite.Equal(expApiAms, apiAms3)

	suite.Nil(err1)
	suite.Nil(err2)
	suite.Nil(err3)
}

func (suite *StoreTestSuite) TestQueryHeadersMethods() {

	suite.SetUpStoreTestSuite()

	// normal case
	var expApiAms []QHeadersAuthMethod
	amb2 := QBasicAuthMethod{ServiceUUID: "uuid2", Host: "host3", Port: 9000, Type: "headers", UUID: "am_uuid_2", CreatedOn: ""}
	am2 := QHeadersAuthMethod{QBasicAuthMethod: amb2, Headers: map[string]string{"x-api-key": "key-1", "Accept": "application/json"}}
	expApiAms = append(expApiAms, am2)
	apiAms, err1 := suite.Mockstore.QueryHeadersAuthMethods("uuid2", "host3")

	// not found - empty list
	apiAms2, err2 := suite.Mockstore.QueryHeadersAuthMethods("unknown", "unknown")

	// query all
	apiAms3, err3 := suite.Mockstore.QueryHeadersAuthMethods("", "")

	suite.Equal(expApiAms, apiAms)
	suite.Equal(0, len(apiAms2))
	suite.Equal(expApiAms, apiAms3)

	suite.Nil(err1)
	suite.Nil(err2)
	suite.Nil(err3)
}

func (suite *StoreTestSuite) TestQueryBindingsByDN() {

	suite.SetUpStoreTestSuite()

	// normal case
	expBinding1 := []QBinding{{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}}
	qBinding1, err1 := suite.Mockstore.QueryBindingsByAuthID("test_dn_1", "uuid1", "host1", "x509")

	// not found case
	var expBinding2 []QBinding
	qBinding2, err2 := suite.Mockstore.QueryBindingsByAuthID("wrong_dn", "wrong_uuid", "wrong_host", "x509")

	// tests the normal case
	suite.Equal(expBinding1, qBinding1)
	suite.Nil(err1)

	//tests the not found case
	suite.Equal(expBinding2, qBinding2)
	suite.Nil(err2)
}

func (suite *StoreTestSuite) TestQueryBindingsByUUID() {

	suite.SetUpStoreTestSuite()

	// normal case
	expBinding1 := []QBinding{{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}}
	qBinding1, err1 := suite.Mockstore.QueryBindingsByUUID("b_uuid1")

	// not found case
	var expBinding2 []QBinding
	qBinding2, err2 := suite.Mockstore.QueryBindingsByUUID("wrong_uuid")

	// tests the normal case
	suite.Equal(expBinding1, qBinding1)
	suite.Nil(err1)

	//tests the not found case
	suite.Equal(expBinding2, qBinding2)
	suite.Nil(err2)
}

func (suite *StoreTestSuite) TestQueryBindings() {

	suite.SetUpStoreTestSuite()

	// normal case - with parameters
	expBindings1 := []QBinding{
		{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b2", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid2", AuthIdentifier: "test_dn_2", UniqueKey: "unique_key_2", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
	}
	qBindings1, err1 := suite.Mockstore.QueryBindings("uuid1", "host1")

	// normal case - without parameters
	expBindings2 := []QBinding{
		{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b2", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid2", AuthIdentifier: "test_dn_2", UniqueKey: "unique_key_2", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b3", ServiceUUID: "uuid1", Host: "host2", UUID: "b_uuid3", AuthIdentifier: "test_dn_3", UniqueKey: "unique_key_3", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b4", ServiceUUID: "uuid2", Host: "host3", UUID: "b_uuid4", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
	}
	qBindings2, err2 := suite.Mockstore.QueryBindings("", "")

	// ot result case - with parameters
	var expBindings3 []QBinding
	qBindings3, err3 := suite.Mockstore.QueryBindings("wrong_service", "wrong_host")

	// tests the normal case - with parameters
	suite.Equal(expBindings1, qBindings1)
	suite.Nil(err1)

	// tests the normal case - without parameters
	suite.Equal(expBindings2, qBindings2)
	suite.Nil(err2)

	// tests the no result case - with parameters
	suite.Equal(expBindings3, qBindings3)
	suite.Nil(err3)
}

func (suite *StoreTestSuite) TestInsertAuthMethod() {

	suite.SetUpStoreTestSuite()

	// insert an QApiKeyAuthMethod and then query the datastore to see if it was inserted
	amb1 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host2", Port: 9000, UUID: "am_uuid_1", CreatedOn: ""}
	expApiAms := []QApiKeyAuthMethod{{amb1, "access_key"}}

	amIns := &QApiKeyAuthMethod{amb1, "access_key"}
	errIns := suite.Mockstore.InsertAuthMethod(amIns)

	apiAms, _ := suite.Mockstore.QueryApiKeyAuthMethods("uuid1", "host2")

	suite.Equal(expApiAms, apiAms)

	suite.Nil(errIns)

}

func (suite *StoreTestSuite) TestInsertServiceType() {

	suite.SetUpStoreTestSuite()

	_, err1 := suite.Mockstore.InsertServiceType("sIns", []string{"host1", "host2", "host3"}, []string{"x509", "oidc"}, "api-key", "uuid_ins", "2018-05-05T18:04:05Z", "ams")

	expQServices1 := []QServiceType{{Name: "sIns", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid_ins", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}}
	qServices1, err1 := suite.Mockstore.QueryServiceTypes("sIns")

	suite.Equal(expQServices1[0], qServices1[0])
	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestInsertBinding() {

	suite.SetUpStoreTestSuite()

	var expBinding1 QBinding
	_, err1 := suite.Mockstore.InsertBinding("bIns", "uuid1", "host1", "b_uuid", "test_dn_ins", "unique_key_ins", "x509")
	// check if the new binding can be found
	expBindings, _ := suite.Mockstore.QueryBindingsByAuthID("test_dn_ins", "uuid1", "host1", "x509")
	expBinding1 = expBindings[0]

	suite.Equal("bIns", expBinding1.Name)
	suite.Equal("uuid1", expBinding1.ServiceUUID)
	suite.Equal("host1", expBinding1.Host)
	suite.Equal("b_uuid", expBinding1.UUID)
	suite.Equal("test_dn_ins", expBinding1.AuthIdentifier)
	suite.Equal("unique_key_ins", expBinding1.UniqueKey)
	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestUpdateBinding() {

	suite.SetUpStoreTestSuite()

	original := QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	updated := QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_upd", UniqueKey: "unique_key_upd", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}

	_, err1 := suite.Mockstore.UpdateBinding(original, updated)

	expBindings, _ := suite.Mockstore.QueryBindingsByAuthID("test_dn_upd", "uuid1", "host1", "x509")
	expBinding1 := expBindings[0]

	suite.Equal(expBinding1, updated)
	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestUpdateServiceType() {

	suite.SetUpStoreTestSuite()

	original := QServiceType{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}

	updated := QServiceType{Name: "s_updated", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}

	_, err1 := suite.Mockstore.UpdateServiceType(original, updated)

	expSVTs, _ := suite.Mockstore.QueryServiceTypesByUUID("uuid1")
	expSVT1 := expSVTs[0]

	suite.Equal(expSVT1, updated)
	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestUpdateAuthMethod() {

	suite.SetUpStoreTestSuite()

	amb1 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	original := &QApiKeyAuthMethod{AccessKey: "access_key"}
	original.QBasicAuthMethod = amb1

	amb2 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	updated := &QApiKeyAuthMethod{AccessKey: "access_key_2"}
	updated.QBasicAuthMethod = amb2

	uqam1, err1 := suite.Mockstore.UpdateAuthMethod(original, updated)

	// query the datastore to see if the update was successful
	apiAms, _ := suite.Mockstore.QueryApiKeyAuthMethods("uuid1", "host1")
	ambExp := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	expApiAms := []QApiKeyAuthMethod{{ambExp, "access_key_2"}}

	suite.Equal(uqam1, updated)
	suite.Equal(expApiAms, apiAms)

	suite.Nil(err1)

}

func (suite *StoreTestSuite) TestDeleteServiceTypeByUUID() {

	suite.SetUpStoreTestSuite()

	err1 := suite.Mockstore.DeleteServiceTypeByUUID("uuid1")

	// query to check if the service type with uuid1 still exists
	st, _ := suite.Mockstore.QueryServiceTypesByUUID("uuid1")

	suite.Nil(err1)
	suite.Nil(st)
}

func (suite *StoreTestSuite) TestDeleteBindingByUUID() {

	suite.SetUpStoreTestSuite()

	err1 := suite.Mockstore.DeleteBindingByServiceUUID("uuid1")

	// since 3 of the 4 bindings belong to the service with uuid1

	suite.Equal(1, len(suite.Mockstore.Bindings))
	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestDeleteBinding() {

	suite.SetUpStoreTestSuite()

	qBinding := QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}

	err1 := suite.Mockstore.DeleteBinding(qBinding)

	// check the slice containing the bindings to see if the qBinding was removed
	expBindings := []QBinding{
		{Name: "b2", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid2", AuthIdentifier: "test_dn_2", UniqueKey: "unique_key_2", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b3", ServiceUUID: "uuid1", Host: "host2", UUID: "b_uuid3", AuthIdentifier: "test_dn_3", UniqueKey: "unique_key_3", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
		{Name: "b4", ServiceUUID: "uuid2", Host: "host3", UUID: "b_uuid4", AuthIdentifier: "test_dn_1", UniqueKey: "unique_key_1", AuthType: "x509", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""},
	}
	qBindings, _ := suite.Mockstore.QueryBindings("", "")

	suite.Equal(expBindings, qBindings)

	suite.Nil(err1)
}

func (suite *StoreTestSuite) TestDeleteAuthMethodByUUID() {

	suite.SetUpStoreTestSuite()

	err1 := suite.Mockstore.DeleteAuthMethodByServiceUUID("uuid1")

	suite.Nil(err1)
	suite.Equal(1, len(suite.Mockstore.AuthMethods))
}

func (suite *StoreTestSuite) TestDeleteAuthMethod() {

	suite.SetUpStoreTestSuite()
	var expAMS []QAuthMethod

	// add a temporary auth method
	amb1 := QBasicAuthMethod{ServiceUUID: "ins_uuid", Host: "ins_host", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	am1 := QApiKeyAuthMethod{AccessKey: "access_key"}
	am1.QBasicAuthMethod = amb1
	suite.Mockstore.AuthMethods = append(suite.Mockstore.AuthMethods, &am1)

	err1 := suite.Mockstore.DeleteAuthMethod(&am1)

	amb := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Type: "api-key", UUID: "am_uuid_1", CreatedOn: ""}
	expApiAms := &QApiKeyAuthMethod{amb, "access_key"}
	amb2 := QBasicAuthMethod{ServiceUUID: "uuid2", Host: "host3", Port: 9000, Type: "headers", UUID: "am_uuid_2", CreatedOn: ""}
	expHeaderam := &QHeadersAuthMethod{QBasicAuthMethod: amb2, Headers: map[string]string{"x-api-key": "key-1", "Accept": "application/json"}}
	expAMS = append(expAMS, expApiAms, expHeaderam)

	suite.Equal(expAMS, suite.Mockstore.AuthMethods)
	suite.Nil(err1)

}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
