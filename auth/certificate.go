package auth

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	LOGGER "github.com/sirupsen/logrus"
)

var ExtraAttributeNames = map[string]string{
	"0.9.2342.19200300.100.1.25": "DC",
}

// load_CAs reads the root certificates from a directory within the filesystem, and creates the trusted root CA chain
func LoadCAs(dir string) (roots *x509.CertPool) {
	LOGGER.Info("Building the root CA chain...")
	pattern := "*.pem"
	roots = x509.NewCertPool()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LOGGER.Errorf("Prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		if ok, _ := filepath.Match(pattern, info.Name()); ok {
			bytes, _ := ioutil.ReadFile(filepath.Join(dir, info.Name()))
			if ok = roots.AppendCertsFromPEM(bytes); !ok {
				return errors.New("Something went wrong while parsing certificate: " + filepath.Join(dir, info.Name()))
			}
		}
		// if info.IsDir() {
		// 	LOGGER.Infof("Skipping a dir without errors: %+v \n", info.Name())
		// }
		return nil
	})

	if err != nil {
		LOGGER.Errorf("error walking the path %q: %v\n", dir, err)
	} else {
		LOGGER.Info("API", "\t", "All certificates parsed successfully.")
	}

	return

}

// ExtractEnhancedRDNSequenceToString extracts a certificate's RDNs to a string using what's provided in the standard library
// and then adding extra attribute names that we have defined
func ExtractEnhancedRDNSequenceToString(cert *x509.Certificate) string {

	var ers string

	ers = cert.Subject.ToRDNSequence().String()

	// we loop the extra attributes in reverse order since certificates from goc db have the RDNs reversed

	for i := len(cert.Subject.Names); i > 0; i-- {
		atv := cert.Subject.Names[i-1]
		if value, ok := ExtraAttributeNames[atv.Type.String()]; ok {
			ers += "," + value + "=" + atv.Value.(string)
		}

	}

	return ers

}
