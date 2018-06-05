package bindings

import (
	"github.com/ARGOeu/argo-api-authn/servicetypes"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	log "github.com/Sirupsen/logrus"
	uuid2 "github.com/satori/go.uuid"
)

type Binding struct {
	Name        string `json:"name" required:"true"`
	ServiceUUID string `json:"service_uuid" required:"true"`
	Host        string `json:"host" required:"true"`
	UUID        string `json:"uuid,omitempty"`
	DN          string `json:"dn,omitempty"`
	OIDCToken   string `json:"oidc_token,omitempty"`
	UniqueKey   string `json:"unique_key" required:"true"`
	CreatedOn   string `json:"created_on,omitempty"`
	LastAuth    string `json:"last_auth,omitempty"`
}

type BindingList struct {
	Bindings []Binding `json:"bindings"`
}

//CreateBinding creates a new binding after validating its context
func CreateBinding(binding Binding, store stores.Store) (Binding, error) {

	var qBinding stores.QBinding
	var err error

	// validate the binding
	if err = binding.Validate(store); err != nil {
		return binding, err
	}

	// generate uuid
	uuid := uuid2.NewV4().String()

	if qBinding, err = store.InsertBinding(binding.Name, binding.ServiceUUID, binding.Host, uuid, binding.DN, binding.OIDCToken, binding.UniqueKey); err != nil {
		return binding, err
	}

	if err = utils.CopyFields(qBinding, &binding); err != nil {
		err = utils.APIGenericInternalError(err.Error())
		return binding, err
	}

	return binding, err
}

func (binding *Binding) Validate(store stores.Store) error {

	var err error
	var ok bool
	var serviceType servicetypes.ServiceType

	// check if all required field have been provided
	if err = utils.ValidateRequired(*binding); err != nil {
		err := utils.APIErrEmptyRequiredField(err.Error())
		return err
	}

	// check if one of DN or OIDCToken has been provided
	if binding.DN == "" && binding.OIDCToken == "" {
		err = utils.APIErrEmptyRequiredField("Both DN and OIDC Token fields are empty")
		return err
	}

	// check if the ServiceUUID is aligned with an existing service type
	if serviceType, err = servicetypes.FindServiceTypeByUUID(binding.ServiceUUID, store); err != nil {
		return err
	}

	// check if the provided host is associated with the given service type
	if ok = serviceType.HasHost(binding.Host); ok == false {
		err = utils.APIErrNotFound("Host")
		return err
	}

	// check if the given dn doesn't already exist under the given service type and host
	// first check for all the other errors regrading bindings
	if _, err = FindBindingByDN(binding.DN, binding.ServiceUUID, binding.Host, store); err != nil {
		log.Info(err)
		if err.Error() != "Binding was not found" {
			return err
		}
	}

	// if the error is nil, it means the function found and returned a binding
	if err == nil {
		err = utils.APIErrConflict(*binding, "dn", binding.DN)
		return err
	}

	//TODO check if the given OIDC token already exists

	return nil
}

// FindBindingByDn queries the datastore and returns a binding based on the given dn and host
func FindBindingByDN(dn string, serviceUUID string, host string, store stores.Store) (Binding, error) {

	var qBindings []stores.QBinding
	var err error
	var binding Binding

	if qBindings, err = store.QueryBindingsByDN(dn, serviceUUID, host); err != nil {
		return Binding{}, err
	}

	if len(qBindings) > 1 {
		log.Warning("STORE", "\t", "More than 1 Users found under the host: "+host+" using the same DN: "+dn)
		err = utils.APIErrDatabase("More than 1 bindings found under the service type: " + serviceUUID + " and host: " + host + " using the same DN: " + dn)
		return Binding{}, err
	}

	if len(qBindings) == 0 {
		err = utils.APIErrNotFound("Binding")
		return Binding{}, err
	}

	if err = utils.CopyFields(qBindings[0], &binding); err != nil {
		err = utils.APIGenericInternalError(err.Error())
		return binding, err
	}

	// update LastAuth field
	//if _, err = UpdateBinding(binding, utils.ZuluTimeNow(), store); err != nil {
	//	err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
	//	return Binding{}, err
	//}
	return binding, err
}

// FindAllBindings returns all the bindings in the service
func FindAllBindings(store stores.Store) (BindingList, error) {

	var err error
	var qBindings []stores.QBinding
	var bindings = []Binding{}

	if qBindings, err = store.QueryBindings("", ""); err != nil {
		return BindingList{Bindings: []Binding{}}, err
	}

	// convert the QBindings to Bindings
	for _, qb := range qBindings {
		_binding := &Binding{}
		if err := utils.CopyFields(qb, _binding); err != nil {
			err = utils.APIGenericInternalError(err.Error())
			return BindingList{Bindings: []Binding{}}, err
		}
		bindings = append(bindings, *_binding)
	}

	return BindingList{Bindings: bindings}, err

}

//FindBindingsByServiceTypeAndHost returns all the bindings of a specific service type and host
func FindBindingsByServiceTypeAndHost(serviceUUID string, host string, store stores.Store) (BindingList, error) {

	var qBindings []stores.QBinding
	var bindings = []Binding{}
	var err error

	if qBindings, err = store.QueryBindings(serviceUUID, host); err != nil {
		return BindingList{Bindings: []Binding{}}, err
	}

	for _, qb := range qBindings {
		_binding := &Binding{}
		if err := utils.CopyFields(qb, _binding); err != nil {
			err = utils.APIGenericInternalError(err.Error())
			return BindingList{Bindings: []Binding{}}, err
		}
		bindings = append(bindings, *_binding)
	}

	return BindingList{Bindings: bindings}, err
}
