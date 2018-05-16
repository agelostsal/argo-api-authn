package auth_methods

import (
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	log "github.com/Sirupsen/logrus"
)

func FindApiKeyMethod(service string, host string, store stores.Store) (map[string]interface{}, error) {

	var err error
	var apiKeyAuthMap map[string]interface{}
	var apiKeyAuthMaps []map[string]interface{}
	var ok bool

	if apiKeyAuthMaps, err = store.QueryAuthMethods(service, host, "api-key"); err != nil {
		return apiKeyAuthMap, err
	}

	if len(apiKeyAuthMaps) == 0 {
		err = utils.APIErrNotFound("Auth method")
		return apiKeyAuthMap, err
	}

	// required variables
	if _, ok = apiKeyAuthMaps[0]["path"]; ok == false {
		log.Error("FindApiKeyMethod", "\t", "Path was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Path was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	if _, ok = apiKeyAuthMaps[0]["port"]; ok == false {
		log.Error("FindApiKeyMethod", "\t", "Port was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Port was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	if _, ok = apiKeyAuthMaps[0]["access_key"]; ok == false {
		log.Error("FindApiKeyMethod", "\t", "Access key was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Access Key was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	return apiKeyAuthMaps[0], err
}

