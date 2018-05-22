package servicetypes

import (
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (suite *ServiceTestSuite) TestCreateService() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	var emptyService ServiceType

	// test the normal case
	s1 := ServiceType{"sCr", []string{"host1", "host2"}, []string{"x509", "oidc"}, "api-key", "token", ""}
	_, err := CreateServiceType(s1, mockstore, *cfg)
	res1, _ := mockstore.QueryServiceTypes("sCr")

	// test the case where the name already exists
	s2 := ServiceType{"s1", []string{"host1", "host2"}, []string{"x509", "oidc"}, "api-key", "token", ""}
	res2, err2 := CreateServiceType(s2, mockstore, *cfg)

	// test the case of unsupported auth type
	s3 := ServiceType{"sCr", []string{"host1", "host2"}, []string{"unsup_type", "oidc"}, "api-key", "token", ""}
	res3, err3 := CreateServiceType(s3, mockstore, *cfg)

	// test the case of empty auth type list
	s4 := ServiceType{"sCr", []string{"host1", "host2"}, []string{}, "api-key", "token", ""}
	res4, err4 := CreateServiceType(s4, mockstore, *cfg)

	// test the case of unsupported auth method
	s5 := ServiceType{"sCr", []string{"host1", "host2"}, []string{"x509", "oidc"}, "unsup_method", "token", ""}
	res5, err5 := CreateServiceType(s5, mockstore, *cfg)

	suite.Equal(s1.Name, res1[0].Name)
	suite.Equal(s1.Hosts, res1[0].Hosts)
	suite.Equal(s1.AuthTypes, res1[0].AuthTypes)
	suite.Equal(s1.AuthMethod, res1[0].AuthMethod)
	suite.Equal(s1.RetrievalField, res1[0].RetrievalField)
	suite.Equal(emptyService, res2)
	suite.Equal(emptyService, res3)
	suite.Equal(emptyService, res4)
	suite.Equal(emptyService, res5)

	suite.Nil(err)
	suite.Equal("servicetypes.ServiceType object with name: s1 already exists", err2.Error())
	suite.Equal("Authentication Type: unsup_type is not yet supported", err3.Error())
	suite.Equal("Authentication Type: empty is not yet supported", err4.Error())
	suite.Equal("Authentication Method: unsup_method is not yet supported", err5.Error())

}

func (suite *ServiceTestSuite) TestFindServiceByName() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// normal case
	expS1 := ServiceType{"s1", []string{"host1", "host2", "host3"}, []string{"x509", "oidc"}, "api-key", "token", "2018-05-05T18:04:05Z"}
	ser1, err1 := FindServiceTypeByName("s1", mockstore)

	// not found case
	var expS2 ServiceType
	ser2, err2 := FindServiceTypeByName("not_found", mockstore)

	// same name
	var expS3 ServiceType
	ser3, err3 := FindServiceTypeByName("same_name", mockstore)

	suite.Equal(expS1, ser1)
	suite.Equal(expS2, ser2)
	suite.Equal(expS3, ser3)

	suite.Nil(err1)
	suite.Equal("ServiceType was not found", err2.Error())
	suite.Equal("Database Error: Multiple service-types with the same name: same_name", err3.Error())
}

func (suite *ServiceTestSuite) TestFindAllServices() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	// normal case outcome - all services
	expQServicesAll := []ServiceType{
		{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", RetrievalField: "token", CreatedOn: "2018-05-05T18:04:05Z"},
		{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "api-key", RetrievalField: "user_token", CreatedOn: "2018-05-05T18:04:05Z"},
		{Name: "same_name"},
		{Name: "same_name"},
	}
	expServList := ServiceList{expQServicesAll}
	serAll1, err1 := FindAllServiceTypes(mockstore)

	// normal case outcome - empty list
	var empServ ServiceList
	mockstore.Services = []stores.QServiceType{}
	serAll2, err2 := FindAllServiceTypes(mockstore)

	suite.Equal(expServList, serAll1)
	suite.Equal(empServ, serAll2)

	suite.Nil(err1)
	suite.Nil(err2)
}

func (suite *ServiceTestSuite) TestServiceHasHost() {

	mockstore := &stores.Mockstore{Server: "localhost", Database: "test_db"}
	mockstore.SetUp()

	cfg := &config.Config{}
	_ = cfg.ConfigSetUp("../config/configuration-test-files/test-conf.json")

	ser := ServiceType{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", RetrievalField: "token", CreatedOn: "2018-05-05T18:04:05Z"}

	suite.Equal(true, ser.HasHost("host1"))
	suite.Equal(false, ser.HasHost("host_unknown"))

}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
