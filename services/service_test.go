package services

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

	var emptyService Service

	// test the normal case
	s1 := Service{"sCr", []string{"host1", "host2"}, []string{"x509", "oidc"}, "api-key", "token", ""}
	_, err := CreateService(s1, mockstore, *cfg)
	res1, _ := mockstore.QueryServices("sCr")

	// test the case where the name already exists
	s2 := Service{"s1", []string{"host1", "host2"}, []string{"x509", "oidc"}, "api-key", "token", ""}
	res2, err2 := CreateService(s2, mockstore, *cfg)

	// test the case of unsupported auth type
	s3 := Service{"sCr", []string{"host1", "host2"}, []string{"unsup_type", "oidc"}, "api-key", "token", ""}
	res3, err3 := CreateService(s3, mockstore, *cfg)

	// test the case of empty auth type list
	s4 := Service{"sCr", []string{"host1", "host2"}, []string{}, "api-key", "token", ""}
	res4, err4 := CreateService(s4, mockstore, *cfg)

	// test the case of unsupported auth method
	s5 := Service{"sCr", []string{"host1", "host2"}, []string{"x509", "oidc"}, "unsup_method", "token", ""}
	res5, err5 := CreateService(s5, mockstore, *cfg)

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
	suite.Equal("services.Service object with name: s1 already exists", err2.Error())
	suite.Equal("Authentication Type: unsup_type is not yet supported", err3.Error())
	suite.Equal("Authentication Type: empty is not yet supported", err4.Error())
	suite.Equal("Authentication Method: unsup_method is not yet supported", err5.Error())

}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
