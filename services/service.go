package services

import (
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
)

type Service struct {
	Name           string   `json:"name" required:"true"`
	Hosts          []string `json:"hosts" required:"true"`
	AuthTypes      []string `json:"auth_types" required:"true"`
	AuthMethod     string   `json:"auth_method" required:"true"`
	RetrievalField string   `json:"retrieval_field" required:"true"`
	CreatedOn      string   `json:"created_on"`
}

type ServiceList struct {
	Services []Service `json:"services"`
}

// CreateService creates a new service after validating the service
func CreateService(service Service, store stores.Store, cfg config.Config) (Service, error) {

	var qServices []stores.QService
	var qService stores.QService
	var err error

	// check if the authentication methods are supported
	if err = service.hasValidAuthMethod(cfg); err != nil {
		return Service{}, err
	}

	// check if the authentication type is supported
	if err = service.hasValidAuthTypes(cfg); err != nil {
		return Service{}, err
	}

	// check that the name of the service is unique
	if qServices, err = store.QueryServices(service.Name); err != nil {
		return Service{}, err
	}

	if len(qServices) > 0 {
		err = utils.APIErrConflict(service, "name", service.Name)
		return Service{}, err
	}

	// insert the service
	if qService, err = store.InsertService(service.Name, service.Hosts, service.AuthTypes, service.AuthMethod, service.RetrievalField, utils.ZuluTimeNow()); err != nil {
		return Service{}, err
	}

	// convert the qService to a Service
	if err = utils.CopyFields(qService, &service); err != nil {
		err = utils.APIErrDatabase(err.Error())
		return Service{}, err
	}

	return service, err
}

// FindServiceByName queries the datastore to find a service associated with the provided argument name
func FindServiceByName(name string, store stores.Store) (Service, error) {

	var qServices []stores.QService
	var service Service
	var err error

	if qServices, err = store.QueryServices(name); err != nil {
		return Service{}, err
	}

	if len(qServices) == 0 {
		err = utils.APIErrNotFound("Service")
		return Service{}, err
	}

	if len(qServices) > 1 {
		err = utils.APIErrDatabase("Multiple services with the same name: " + name)
		return Service{}, err
	}

	if err := utils.CopyFields(qServices[0], &service); err != nil {
		err = utils.APIGenericInternalError(err.Error())
		return Service{}, err
	}

	return service, err
}

// FindAllServices returns all the services from the datastore
func FindAllServices(store stores.Store) (ServiceList, error) {

	var qServices []stores.QService
	var services []Service
	var err error

	if qServices, err = store.QueryServices(""); err != nil {
		return ServiceList{}, err
	}

	for _, qs := range qServices {
		_service := &Service{}
		if err := utils.CopyFields(qs, _service); err != nil {
			err = utils.APIGenericInternalError(err.Error())
			return ServiceList{}, err
		}
		services = append(services, *_service)
	}

	return ServiceList{services}, err

}

// hsHost returns whether or not a host is associated with a service
func (s *Service) HasHost(host string) bool {

	flag := false
	for _, h := range s.Hosts {
		if h == host {
			flag = true
			break
		}
	}

	return flag
}

// hasValidAuthTypes checks whether or not the authentication types of a project are supported
func (s *Service) hasValidAuthTypes(cfg config.Config) error {

	var err error
	var flag bool

	if len(s.AuthTypes) == 0 {
		err = utils.APIErrUnsupportedContent("Authentication Type", "empty")
		return err
	}

	for _, am := range s.AuthTypes {
		flag = false
		for _, cam := range cfg.SupportedAuthTypes {
			if am == cam {
				flag = true
				break
			}
		}
		if !flag {
			err = utils.APIErrUnsupportedContent("Authentication Type", am)
			return err
		}
	}
	return err
}

// hasValidAuthMethod checks whether or not the authentication method of a project is supported
func (s *Service) hasValidAuthMethod(cfg config.Config) error {

	var err error
	for _, am := range cfg.SupportedAuthMethods {
		if am == s.AuthMethod {
			return err
		}
	}

	err = utils.APIErrUnsupportedContent("Authentication Method", s.AuthMethod)
	return err
}
