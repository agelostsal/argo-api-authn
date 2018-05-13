package auth

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// load_CAs reads the root certificates from a directory within the filesystem, and creates the trusted root CA chain
func Load_CAs(dir string) (roots *x509.CertPool) {
	log.Info("Building the root CA chain...")
	pattern := "*.pem"
	roots = x509.NewCertPool()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		if ok, _ := filepath.Match(pattern, info.Name()); ok {
			bytes, _ := ioutil.ReadFile(filepath.Join(dir, info.Name()))
			if ok = roots.AppendCertsFromPEM(bytes); !ok {
				return errors.New("Something went wrong while parsing certificate: " + filepath.Join(dir, info.Name()))
			}
		}
		// if info.IsDir() {
		// 	log.Infof("Skipping a dir without errors: %+v \n", info.Name())
		// }
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dir, err)
	} else {
		log.Info("API", "\t", "All certificates parsed successfully.")
	}

	return

}
