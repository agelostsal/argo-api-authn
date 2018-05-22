package auth_methods

import (
	"github.com/ARGOeu/argo-api-authn/stores"
	"github.com/ARGOeu/argo-api-authn/utils"
	log "github.com/Sirupsen/logrus"
	"strings"
)

func FindApiKeyAuthMethod(service string, host string, store stores.Store) (map[string]interface{}, error) {

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
		log.Error("FindApiKeyAuthMethod", "\t", "Path was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Path was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	if _, ok = apiKeyAuthMaps[0]["port"]; ok == false {
		log.Error("FindApiKeyAuthMethod", "\t", "Port was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Port was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	if _, ok = apiKeyAuthMaps[0]["access_key"]; ok == false {
		log.Error("FindApiKeyAuthMethod", "\t", "Access key was not found in the apiKeyAuthMap")
		err = utils.APIErrDatabase("Access Key was not found in the ApiKeyAuth object")
		return apiKeyAuthMap, err
	}

	return apiKeyAuthMaps[0], err
}

func CreateApiKeyAuthMethod(authM map[string]interface{}, store stores.Store) (map[string]interface{}, error) {

	var err error
	var ok bool

	// extra required variables
	if _, ok = authM["access_key"]; ok == false {
		log.Error("FindApiKeyAuthMethod", "\t", "Access key was not found in the apiKeyAuthMap")
		err = utils.APIErrEmptyRequiredField("access_key was not found in the request body")
		return authM, err
	}

	// check if the path contains at least the two interpolations, {{identifier}}, {{access_key}}
	if !strings.Contains(authM["path"].(string), "{{identifier}}") {
		err = utils.APIErrInvalidFieldContent("path", "Doesn't contain {{identifier}}")
		return authM, err
	}

	if !strings.Contains(authM["path"].(string), "{{access_key}}") {
		err = utils.APIErrInvalidFieldContent("path", "Doesn't contain {{access_key}}")
		return authM, err
	}

	if err = store.InsertAuthMethod(authM); err != nil {
		return authM, err
	}

	return authM, err

}
