package stores

import (
	"github.com/ARGOeu/argo-api-authn/utils"
	"reflect"
)

type Mockstore struct {
	Session      bool
	Server       string
	Database     string
	ServiceTypes []QServiceType
	Bindings     []QBinding
	AuthMethods  []QAuthMethod
}

// SetUp is used to initialize the mock store
func (mock *Mockstore) SetUp() {

	mock.Session = true

	// Populate services
	service1 := QServiceType{Name: "s1", Hosts: []string{"host1", "host2", "host3"}, AuthTypes: []string{"x509", "oidc"}, AuthMethod: "api-key", UUID: "uuid1", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}
	service2 := QServiceType{Name: "s2", Hosts: []string{"host3", "host4"}, AuthTypes: []string{"x509"}, AuthMethod: "api-key", UUID: "uuid2", CreatedOn: "2018-05-05T18:04:05Z", Type: "ams"}
	serviceSame1 := QServiceType{Name: "same_name"}
	serviceSame2 := QServiceType{Name: "same_name"}
	mock.ServiceTypes = append(mock.ServiceTypes, service1, service2, serviceSame1, serviceSame2)

	// Populate Bindings
	binding1 := QBinding{Name: "b1", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid1", DN: "test_dn_1", OIDCToken: "", UniqueKey: "unique_key_1", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding2 := QBinding{Name: "b2", ServiceUUID: "uuid1", Host: "host1", UUID: "b_uuid2", DN: "test_dn_2", OIDCToken: "", UniqueKey: "unique_key_2", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	binding3 := QBinding{Name: "b3", ServiceUUID: "uuid1", Host: "host2", UUID: "b_uuid3", DN: "test_dn_3", OIDCToken: "", UniqueKey: "unique_key_3", CreatedOn: "2018-05-05T15:04:05Z", LastAuth: ""}
	mock.Bindings = append(mock.Bindings, binding1, binding2, binding3)

	// Populate AuthMethods
	amb1 := QBasicAuthMethod{ServiceUUID: "uuid1", Host: "host1", Port: 9000, Path: "/path/{{identifier}}?key={{access_key}}", Type: "api-key", UUID: "am_uuid_1", CreatedOn: "", RetrievalField: "token"}
	am1 := &QApiKeyAuthMethod{AccessKey: "access_key"}
	am1.QBasicAuthMethod = amb1
	mock.AuthMethods = append(mock.AuthMethods, am1)
}

func (mock *Mockstore) Close() {
	mock.Session = false
}

func (mock *Mockstore) QueryServiceTypes(name string) ([]QServiceType, error) {

	var qServices []QServiceType

	if name != "" {
		for _, service := range mock.ServiceTypes {
			if service.Name == name {
				qServices = append(qServices, service)
			}
		}
	} else {
		qServices = mock.ServiceTypes
	}

	return qServices, nil
}

func (mock *Mockstore) QueryServiceTypesByUUID(uuid string) ([]QServiceType, error) {

	var qServices []QServiceType

	for _, service := range mock.ServiceTypes {
		if service.UUID == uuid {
			qServices = append(qServices, service)
		}
	}

	return qServices, nil
}

func (mock *Mockstore) QueryApiKeyAuthMethods(serviceUUID string, host string) ([]QApiKeyAuthMethod, error) {

	var qAuthms []QApiKeyAuthMethod
	var err error
	var ok bool
	var qAuthm *QApiKeyAuthMethod

	if serviceUUID == "" && host == "" {
		for _, am := range mock.AuthMethods {
			if qAuthm, ok = am.(*QApiKeyAuthMethod); ok {
				qAuthms = append(qAuthms, *qAuthm)
			}
		}
		return qAuthms, nil
	}

	for _, am := range mock.AuthMethods {
		if qAuthm, ok = am.(*QApiKeyAuthMethod); ok {
			if qAuthm.ServiceUUID == serviceUUID && qAuthm.Host == host {
				qAuthms = append(qAuthms, *qAuthm)
			}
		}
	}

	return qAuthms, err

}

func (mock *Mockstore) QueryBindingsByDN(dn string, serviceUUID string, host string) ([]QBinding, error) {

	var qBindings []QBinding

	for _, qBinding := range mock.Bindings {
		if qBinding.DN == dn && qBinding.Host == host && qBinding.ServiceUUID == serviceUUID {
			qBindings = append(qBindings, qBinding)
		}
	}

	return qBindings, nil
}

func (mock *Mockstore) QueryBindingsByUUID(uuid string) ([]QBinding, error) {

	var qBindings []QBinding

	for _, qBinding := range mock.Bindings {
		if qBinding.UUID == uuid {
			qBindings = append(qBindings, qBinding)
		}
	}

	return qBindings, nil
}

func (mock *Mockstore) QueryBindings(serviceUUID string, host string) ([]QBinding, error) {

	var qBindings []QBinding

	if serviceUUID == "" && host == "" {
		qBindings = mock.Bindings
		return qBindings, nil
	}

	for _, qBinding := range mock.Bindings {
		if qBinding.ServiceUUID == serviceUUID && qBinding.Host == host {
			qBindings = append(qBindings, qBinding)
		}
	}

	return qBindings, nil
}

func (mock *Mockstore) InsertAuthMethod(am QAuthMethod) error {

	mock.AuthMethods = append(mock.AuthMethods, am)

	return nil
}

func (mock *Mockstore) InsertServiceType(name string, hosts []string, authTypes []string, authMethod string, uuid string, createdOn string, sType string) (QServiceType, error) {

	qService := QServiceType{Name: name, Hosts: hosts, AuthTypes: authTypes, AuthMethod: authMethod, UUID: uuid, CreatedOn: createdOn, Type: sType}

	mock.ServiceTypes = append(mock.ServiceTypes, qService)

	return qService, nil
}

func (mock *Mockstore) InsertBinding(name string, serviceUUID string, host string, uuid string, dn string, oidcToken string, uniqueKey string) (QBinding, error) {

	qBinding := QBinding{Name: name, ServiceUUID: serviceUUID, Host: host, DN: dn, UUID: uuid, OIDCToken: oidcToken, UniqueKey: uniqueKey, CreatedOn: utils.ZuluTimeNow()}

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

func (mock *Mockstore) UpdateServiceType(original QServiceType, updated QServiceType) (QServiceType, error) {

	// find the service type in the list and replace it
	for idx, sv := range mock.ServiceTypes {
		if reflect.DeepEqual(original, sv) { // requires DeepEqual because structs with []string as fields can't be compared
			mock.ServiceTypes[idx] = updated
			break
		}
	}

	return updated, nil
}

func (mock *Mockstore) UpdateAuthMethod(original QAuthMethod, updated QAuthMethod) (QAuthMethod, error) {

	// find the auth method in the list and replace it
	for idx, sv := range mock.AuthMethods {
		if reflect.DeepEqual(original, sv) {
			mock.AuthMethods[idx] = updated
			break
		}
	}

	return updated, nil
}

func (mock *Mockstore) DeleteServiceTypeByUUID(uuid string) error {

	for idx, st := range mock.ServiceTypes {
		if st.UUID == uuid {
			mock.ServiceTypes = append(mock.ServiceTypes[:idx], mock.ServiceTypes[idx+1:]...)
			break
		}
	}

	return nil
}

// DeleteBinding removes the given qBinding from the slice of bindings
func (mock *Mockstore) DeleteBinding(qBinding QBinding) error {

	// find the  binding in the list and replace it
	for idx, qb := range mock.Bindings {
		if qb == qBinding {
			mock.Bindings = append(mock.Bindings[:idx], mock.Bindings[idx+1:]...)
			break
		}
	}

	return nil
}

func (mock *Mockstore) DeleteBindingByServiceUUID(serviceUUID string) error {

	var remainingBindings []QBinding

	// extract all the bindings that don't match the service's uuid
	for _, qb := range mock.Bindings {
		if qb.ServiceUUID != serviceUUID {
			remainingBindings = append(remainingBindings, qb)
		}
	}

	mock.Bindings = remainingBindings
	return nil
}

func (mock *Mockstore) DeleteAuthMethod(am QAuthMethod) error {

	// loop through the slice of auth methods
	// and delete

	for idx, ami := range mock.AuthMethods {
		if reflect.DeepEqual(ami, am) {
			mock.AuthMethods = append(mock.AuthMethods[:idx], mock.AuthMethods[idx+1:]...)
			break
		}
	}

	return nil

}

func (mock *Mockstore) DeleteAuthMethodByServiceUUID(serviceUUID string) error {

	var remainingQAM []QAuthMethod

	for _, qam := range mock.AuthMethods {
		i, _ := utils.GetFieldValueByName(qam, "ServiceUUID")
		if i.(string) != serviceUUID {
			remainingQAM = append(remainingQAM, qam)
		}
	}
	mock.AuthMethods = remainingQAM
	return nil
}