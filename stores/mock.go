package stores

import (
	"github.com/ARGOeu/argo-api-authn/argo-consts"
	"github.com/ARGOeu/argo-api-authn/utils"
	"time"
)

type Mockstore struct {
	Session   bool
	Server    string
	Database  string
	Services  []QService
	Bindings  []QBinding
	AuthTypes []interface{}
}

// SetUp is used to initialize the mock store
func (mock *Mockstore) Setup(server string, database string) {

	mock.Server = server
	mock.Database = database
	mock.Session = true

	// Populate services
	service1 := QService{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509, oidc"}, AuthMethod: "api-key", RetrievalField: "token"}
	service2 := QService{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "api-key", RetrievalField: "user_token"}
	mock.Services = append(mock.Services, service1, service2)

	// Populate Bindings
	binding1 := QBinding{Name: "b1", Service: "s1", Host: "host1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "unique_key_1", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding2 := QBinding{Name: "b2", Service: "s1", Host: "host1", DN: "test_dn_2", OIDCToken: "", UniqueKey: "unique_key_2", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding3 := QBinding{Name: "b3", Service: "s2", Host: "host2", DN: "test_dn_3", OIDCToken: "", UniqueKey: "unique_key_3", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	mock.Bindings = append(mock.Bindings, binding1, binding2, binding3)

	// Populate AuthTypes
	apiKeyAuth1 := QApiKeyAuth{Type: "api-key", Service: "s1", Host: "host1", Path: "test_path_1", Port: 9000, AccessKey: "key1"}
	mock.AuthTypes = append(mock.AuthTypes, apiKeyAuth1)
}

func (mock *Mockstore) Close() {
	mock.Session = false
}

func (mock *Mockstore) QueryServices(name string) ([]QService, error) {

	var qServices []QService

	if name != "" {
		for _, service := range mock.Services {
			if service.Name == name {
				qServices = append(qServices, service)
			}
		}
	} else {
		qServices = mock.Services
	}

	return qServices, nil
}

func (mock *Mockstore) QueryAuthMethod(service string, host string, typeName string) (map[string]interface{}, error) {

	var qAuthM = make(map[string]interface{})

	for _, authM := range mock.AuthTypes {
		qAuthM = utils.StructToMap(authM)
		if qAuthM["Service"] == service && qAuthM["Host"] == host && qAuthM["Type"] == typeName {
			return qAuthM, nil
		}
	}

	return make(map[string]interface{}), nil
}

func (mock *Mockstore) QueryBindingsByDN(dn string, host string) ([]QBinding, error) {

	var qBindings []QBinding

	for _, qBinding := range mock.Bindings {
		if qBinding.DN == dn && qBinding.Host == host {
			qBindings = append(qBindings, qBinding)
		}
	}

	return qBindings, nil
}

func (mock *Mockstore) QueryBindings(service string, host string) ([]QBinding, error) {

	var qBindings []QBinding

	if service == "" && host == "" {
		qBindings = mock.Bindings
		return qBindings, nil
	}

	for _, qBinding := range mock.Bindings {
		if qBinding.Service == service && qBinding.Host == host {
			qBindings = append(qBindings, qBinding)
		}
	}

	return qBindings, nil
}

func (mock *Mockstore) InsertBinding(name string, service string, host string, dn string, oidcToken string, uniqueKey string) (QBinding, error) {

	qBinding := QBinding{Name: name, Service: service, Host: host, DN: dn, OIDCToken: oidcToken, UniqueKey: uniqueKey, CreatedOn: time.Now().Format(argo_consts.ZULU_FORM)}

	mock.Bindings = append(mock.Bindings, qBinding)

	return qBinding, nil

}

func (mock *Mockstore) UpdateBinding(original QBinding, updated QBinding) (QBinding, error) {

	// find the  binding in the list and replace it
	for idx, qb := range mock.Bindings {
		if qb == original {
			mock.Bindings[idx] = updated
			break
		}
	}

	return updated, nil
}
