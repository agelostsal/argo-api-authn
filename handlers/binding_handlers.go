package handlers

import (
	"encoding/json"
	"github.com/ARGOeu/argo-api-authn/bindings"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	"github.com/gorilla/context"
	"net/http"
)

func BindingCreate(w http.ResponseWriter, r *http.Request) {

	var err error

	//context references
	store := context.Get(r, "stores").(stores.Store)

	var binding bindings.Binding

	// check the validity of the JSON
	if err = json.NewDecoder(r.Body).Decode(&binding); err != nil {
		err := utils.APIErrBadRequest(err.Error())
		utils.RespondError(w, err)
		return
	}

	if binding, err = bindings.CreateBinding(binding, store); err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.RespondOk(w, 201, binding)
}
