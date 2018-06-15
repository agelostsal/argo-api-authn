package handlers

import (
	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	"github.com/gorilla/context"

	"github.com/ARGOeu/argo-api-authn/servicetypes"
	"github.com/gorilla/mux"
	"net/http"
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
		err := utils.APIErrEmptyRequiredField("service-type", err.Error())
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
	var servList servicetypes.ServiceTypesList

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

// ServiceTypeUpdate updates a service type
func ServiceTypeUpdate(w http.ResponseWriter, r *http.Request) {

	var err error
	var originalSt servicetypes.ServiceType
	var updatedSt servicetypes.ServiceType
	var tempST servicetypes.TempServiceType

	//context references
	store := context.Get(r, "stores").(stores.Store)
	cfg := context.Get(r, "config").(config.Config)


	// url vars
	vars := mux.Vars(r)

	if originalSt, err = servicetypes.FindServiceTypeByName(vars["service-type"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// first, fill the temporary binding with the fields of the original binding
	if err := utils.CopyFields(originalSt, &tempST); err != nil {
		err = utils.APIGenericInternalError(err.Error())
		utils.RespondError(w, err)
	}

	// check the validity of the JSON and updated the provided fields
	if err = json.NewDecoder(r.Body).Decode(&tempST); err != nil {
		err := utils.APIErrBadRequest(err.Error())
		utils.RespondError(w, err)
		return
	}

	if updatedSt, err = servicetypes.UpdateServiceType(originalSt, tempST, store, cfg); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondOk(w, 200, updatedSt)

}
