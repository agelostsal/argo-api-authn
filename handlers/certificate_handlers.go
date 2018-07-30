package handlers

import (
	"net/http"

	"github.com/ARGOeu/argo-api-authn/auth"
	"github.com/ARGOeu/argo-api-authn/authmethods"
	"github.com/ARGOeu/argo-api-authn/bindings"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/servicetypes"
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func AuthViaCert(w http.ResponseWriter, r *http.Request) {

	var err error
	var ok bool
	var dataRes = make(map[string]interface{})
	var binding bindings.Binding
	var serviceType servicetypes.ServiceType
	var authm authmethods.AuthMethod

	//context references
	store := context.Get(r, "stores").(stores.Store)

	// url vars
	vars := mux.Vars(r)
	cfg := context.Get(r, "config").(config.Config)

	if len(r.TLS.PeerCertificates) == 0 {
		err = &utils.APIError{Message: "No certificate provided", Code: 400, Status: "BAD REQUEST"}
		utils.RespondError(w, err)
		return
	}

	// validate the certificate
	if cfg.VerifyCertificate {
		if err = auth.ValidateClientCertificate(r.TLS.PeerCertificates[0], r.RemoteAddr); err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	// Find information regarding the requested serviceType
	if serviceType, err = servicetypes.FindServiceTypeByName(vars["service-type"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// check if the provided host is associated with the given serviceType type
	if ok = serviceType.HasHost(vars["host"]); ok == false {
		err = utils.APIErrNotFound("Host")
		utils.RespondError(w, err)
		return
	}

	// check if the auth method exists
	if authm, err = authmethods.AuthMethodFinder(serviceType.UUID, vars["host"], serviceType.AuthMethod, store); err != nil {
		utils.RespondError(w, err)
		return
	}

	// Find the binding associated with the provided certificate
	if binding, err = bindings.FindBindingByDN(auth.ExtractEnhancedRDNSequenceToString(r.TLS.PeerCertificates[0]), serviceType.UUID, vars["host"], store); err != nil {
		utils.RespondError(w, err)
		return
	}

	extraRequestData := map[string]interface{}{"binding-identifier": binding.UniqueKey}

	if dataRes, err = authm.RetrieveAuthResource(extraRequestData, &cfg); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondOk(w, 200, dataRes)

}
