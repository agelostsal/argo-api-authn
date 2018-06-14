package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ARGOeu/argo-api-authn/utils"
	log "github.com/Sirupsen/logrus"
)

type Config struct {
	ServicePort            int      `json:"service_port" required:"true"`
	MongoHost              string   `json:"mongo_host" required:"true"`
	MongoDB                string   `json:"mongo_db" required:"true"`
	CertificateAuthorities string   `json:"certificate_authorities" required:"true"`
	Certificate            string   `json:"certificate" required:"true"`
	CertificateKey         string   `json:"certificate_key" required:"true"`
	ServiceToken           string   `json:"service_token" required:"true"`
	SupportedAuthTypes     []string `json:"supported_auth_types" required:"true"`
	SupportedAuthMethods   []string `json:"supported_auth_methods" required:"true"`
	VerifySSL              bool     `json:"verify_ssl" `
}

// ConfigSetUp unmarshals a json file specified by the input parameter into the config object
func (cfg *Config) ConfigSetUp(path string) error {

	var data []byte
	var err error

	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.New("Something went wrong while marshaling the json data. Error: " + err.Error())
	}

	log.Info(fmt.Sprintf("%+v", cfg))

	if err = utils.ValidateRequired(*cfg); err != nil {
		return err
	}
	return nil
}
