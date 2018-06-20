package auth

import (
	"github.com/stretchr/testify/suite"
	"crypto/x509"
	"encoding/pem"
	log "github.com/Sirupsen/logrus"
	"testing"
	"crypto/x509/pkix"
	"encoding/asn1"
)

type CertificateTestSuite struct {
	suite.Suite
}

func (suite *CertificateTestSuite) TestExtractEnhancedRDNSequenceToString() {

	// create a new certificate from the string literal - (doesn't contain extra attribute names)
	commonCert := `-----BEGIN CERTIFICATE-----
MIIC8jCCAdoCCQCdC824csOlXTANBgkqhkiG9w0BAQsFADA7MQswCQYDVQQGEwJU
QzELMAkGA1UECAwCVE4xDTALBgNVBAcMBENJVFkxEDAOBgNVBAoMB0NPTVBBTlkw
HhcNMTgwNjA0MDkwMzA5WhcNMTkwNjA0MDkwMzA5WjA7MQswCQYDVQQGEwJUQzEL
MAkGA1UECAwCVE4xDTALBgNVBAcMBENJVFkxEDAOBgNVBAoMB0NPTVBBTlkwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDYahRnf7gxkjz81VX9JjQ/PiB7
UpGInckBCl5mah/Q3ucr3OvaLLdO4pfmAMTanLxgcTsP7k/yyvZF17IMhQ5wpNzn
zfLlAswQT6sqkNyJUx1MXI5mjAiDDMpUh0c9CnGaZa/LrnTXmQhsv8uzLDPUYb37
iHHjP7isYiG+7YE2CpRwHazj0SYba4HAYw8Z1L8Z6kI1gfdIiqI5DFMBdmQlac3P
YCtYytQd3swCsxf57/M9X+Ct7DVcuPKmR5vv4ONL2YBwtvULX8DA8aApdxIrHCZE
NRaaCUmMSo4JXHbY5CVP6AXdo3Iz+3v485qtF8lk+XU/fQGNTdgs1hmsTNMRAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAAyvW6yVbfCLMxuQ1Nt61OmKA96fTtdTpLuI
nq9C0XVoFqgrEeobxdH4QbwxduRHWpHEuFskJCBnMbX0d8v63KEN/6I0Ub4niaeP
nvykp3uoKrRwIZo4OxJFuuuLuUw3aAwkKeqsy5HZsKqi9QscHExRKbcIdlkgxzRW
IEPEGk7acMlT20ECjc4zbdza8PKQeBeEVLINJVMRGPZIlo/6z6BxSnINQiWBk1WZ
lXTmJMXGB7/0ECDz2JT1Mbs/q2ijlZywz0xsp+Zdsp1I01wvwqw5M12PHf3buM1w
SoPmZKiBeb+2OQ2n7+FI8ftkqxWw6zjh651brAoy/0zqLTRPh+c=
-----END CERTIFICATE-----
`

	// tests the case where the certificate doesn't contain extra attributes names
	var crt *x509.Certificate
	var err error

	block, _ := pem.Decode([]byte(commonCert))
	if block == nil {
		panic("failed to parse certificate PEM")
	}

	if crt, err = x509.ParseCertificate(block.Bytes); err != nil {
		log.Error(err.Error())
	}

	ers := ExtractEnhancedRDNSequenceToString(crt)

	// add some extra attribute names to the certificate
	obj := asn1.ObjectIdentifier{0,9,2342,19200300,100,1,25}
	extraAttributeValue1 := pkix.AttributeTypeAndValue{Type:obj, Value:"v1"}
	extraAttributeValue2 := pkix.AttributeTypeAndValue{Type:obj, Value:"v2"}
	enhancedCert := crt
	enhancedCert.Subject.Names = append(enhancedCert.Subject.Names, extraAttributeValue1, extraAttributeValue2)
	ers2 := ExtractEnhancedRDNSequenceToString(enhancedCert)


	suite.Equal("O=COMPANY,L=CITY,ST=TN,C=TC", ers)
	suite.Equal("O=COMPANY,L=CITY,ST=TN,C=TC,DC=v2,DC=v1", ers2)
}

func TestCertificateTestSuite(t *testing.T) {
	suite.Run(t, new(CertificateTestSuite))
}