package stores

import (
	"github.com/ARGOeu/argo-api-authn/utils"
)

type Mockstore struct {
	Session     bool
	Server      string
	Database    string
	Services    []QService
	Bindings    []QBinding
	AuthMethods []map[string]interface{}
}

// SetUp is used to initialize the mock store
func (mock *Mockstore) SetUp() {

	mock.Session = true

	// Populate services
	service1 := QService{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", RetrievalField: "token", CreatedOn: "2018-05-05T18:04:05Z"}
	service2 := QService{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "api-key", RetrievalField: "user_token", CreatedOn: "2018-05-05T18:04:05Z"}
	serviceSame1 := QService{Name: "same_name"}
	serviceSame2 := QService{Name: "same_name"}
	mock.Services = append(mock.Services, service1, service2, serviceSame1, serviceSame2)

	// Populate Bindings
	binding1 := QBinding{Name: "b1", Service: "s1", Host: "host1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "unique_key_1", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding2 := QBinding{Name: "b2", Service: "s1", Host: "host1", DN: "test_dn_2", OIDCToken: "", UniqueKey: "unique_key_2", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding3 := QBinding{Name: "b3", Service: "s2", Host: "host2", DN: "test_dn_3", OIDCToken: "", UniqueKey: "unique_key_3", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	mock.Bindings = append(mock.Bindings, binding1, binding2, binding3)

	// Populate AuthMethods
	mock.AuthMethods = []map[string]interface{}{{"service": "s1", "host": "host1", "port": 9000.0, "path": "test_path_1", "access_key": "key1", "type": "api-key"},
		{"host": "host2", "port": 9000.0, "path": "test_path_1", "type": "api-key", "service": "s1"},
		{"access_key": "key1", "type": "api-key", "service": "s2", "host": "host3", "port": 9000.0},
		{"path": "test_path_1", "access_key": "key1", "type": "api-key", "service": "s2", "host": "host4"}}
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

func (mock *Mockstore) QueryAuthMethods(service string, host string, typeName string) ([]map[string]interface{}, error) {

	var qAuthMs []map[string]interface{}
	var authM map[string]interface{}

	if service == "" && host == "" && typeName == "" {
		return mock.AuthMethods, nil
	}

	for _, authM = range mock.AuthMethods {
		if authM["service"] == service && authM["host"] == host && authM["type"] == typeName {
			qAuthMs = append(qAuthMs, authM)
			break
		}
	}

	return qAuthMs, nil
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

func (mock *Mockstore) InsertService(name string, hosts []string, authTypes []string, authMethod string, retrievalField string, createdOn string) (QService, error) {

	qService := QService{Name: name, Hosts: hosts, AuthTypes: authTypes, AuthMethod: authMethod, RetrievalField: retrievalField, CreatedOn: createdOn}

	mock.Services = append(mock.Services, qService)

	return qService, nil
}

func (mock *Mockstore) InsertBinding(name string, service string, host string, dn string, oidcToken string, uniqueKey string) (QBinding, error) {

	qBinding := QBinding{Name: name, Service: service, Host: host, DN: dn, OIDCToken: oidcToken, UniqueKey: uniqueKey, CreatedOn: utils.ZuluTimeNow()}

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
