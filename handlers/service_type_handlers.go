package handlers

import (
	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	"github.com/gorilla/context"

	"github.com/gorilla/mux"
	"net/http"
	"github.com/ARGOeu/argo-api-authn/servicetypes"
)

func ServiceTypeCreate(w http.ResponseWriter, r *http.Request) {

	var err error
	var service servicetypes.ServiceType

	//context references
	store := context.Get(r, "stores").(stores.Store)
	cfg := context.Get(r, "config").(config.Config)

	// check the validity of the JSON
	if err = json.NewDecoder(r.Body).Decode(&service); err != nil {
		err := utils.APIErrBadRequest(err.Error())
		utils.RespondError(w, err)
		return
	}

	// check if all required field have been provided
	if err = utils.ValidateRequired(service); err != nil {
		err := utils.APIErrEmptyRequiredField(err.Error())
		utils.RespondError(w, err)
		return
	}

	// create the service
	if service, err = servicetypes.CreateServiceType(service, store, cfg); err != nil {
		utils.RespondError(w, err)
		return
	}

	// if everything went ok, reflect the created object
	utils.RespondOk(w, 201, service)
}

func ServiceTypesListOne(w http.ResponseWriter, r *http.Request) {

	var err error
	var service servicetypes.ServiceType

	//context references
	store := context.Get(r, "stores").(stores.Store)

	// url vars
	vars := mux.Vars(r)

	// find the service
	if service, err = servicetypes.FindServiceTypeByName(vars["service-type"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// if everything went ok, return the service
	utils.RespondOk(w, 200, service)
}

func ServiceTypeListAll(w http.ResponseWriter, r *http.Request) {

	var err error
	var servList servicetypes.ServiceList

	//context references
	store := context.Get(r, "stores").(stores.Store)

	// find the service
	if servList, err = servicetypes.FindAllServiceTypes(store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// if everything went ok, return the service
	utils.RespondOk(w, 200, servList)
}
