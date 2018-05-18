package handlers

import (
	"github.com/ARGOeu/argo-api-authn/auth-methods"
	"github.com/ARGOeu/argo-api-authn/services"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthMethodListOne(w http.ResponseWriter, r *http.Request) {

	var err error
	var authM map[string]interface{}
	var service services.Service

	//context references
	store := context.Get(r, "stores").(stores.Store)

	// url vars
	vars := mux.Vars(r)

	// check if the service exists
	if service, err = services.FindServiceByName(vars["service"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// check if the host exists
	if !service.HasHost(vars["host"]) {
		err = utils.APIErrNotFound("Host")
		utils.RespondError(w, err)
		return
	}

	// depending on the service's declared auth method, grab the respective finder
	authMFinder := auth_methods.AuthMethodFinders[service.AuthMethod]

	if authM, err = authMFinder(vars["service"], vars["host"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// if everything went ok, return the auth method
	utils.RespondOk(w, 200, authM)
}

func AuthMethodListAll(w http.ResponseWriter, r *http.Request) {

	//context references
	store := context.Get(r, "stores").(stores.Store)

	var authMs []map[string]interface{}
	var err error

	if authMs, err = auth_methods.FindAllAuthMethods(store); err != nil {
		utils.RespondError(w, err)
		return
	}

	logrus.Info(authMs)

	aml := auth_methods.AuthMethodsList{authMs}

	utils.RespondOk(w, 200, aml)
}
