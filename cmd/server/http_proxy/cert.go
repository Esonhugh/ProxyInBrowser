package http_proxy

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

type certGenerator struct {
	caCert *x509.Certificate
	caKey  crypto.PrivateKey
}

func newCertGenerator(publicKeyFile, privateKeyFile string) (*certGenerator, error) {
	tlsCert, err := tls.LoadX509KeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		return nil, err
	}
	caCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &certGenerator{
		caCert: caCert,
		caKey:  tlsCert.PrivateKey,
	}, nil
}

func (cg *certGenerator) Get(hostname string) (*tls.Config, error) {
	serial, err := getRandomSerialNumber()
	if err != nil {
		return nil, err
	}

	hostCert := &x509.Certificate{
		SerialNumber: serial,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		DNSNames:              []string{hostname},
		NotBefore:             time.Now().Add(-time.Second * 300),
		NotAfter:              time.Now().Add(time.Hour * 24 * 30),
		Subject: pkix.Name{
			CommonName: hostname,
		},
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	certDER, err := x509.CreateCertificate(rand.Reader, hostCert, cg.caCert, key.Public(), cg.caKey)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{
					certDER,
				},
				PrivateKey: key,
			},
		},
		NextProtos: []string{
			"http/1.1",
			"h2",
		},
	}
	return tlsConfig, nil
}

func createCA(caCertFile, caKeyFile string) error {
	serial, err := getRandomSerialNumber()
	if err != nil {
		return err
	}

	caCert := &x509.Certificate{
		SerialNumber:          serial,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature |
			x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365 * 10),
		Subject: pkix.Name{
			CommonName: "Proxy CA",
		},
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	caCertDER, err := x509.CreateCertificate(rand.Reader, caCert, caCert, key.Public(), key)
	if err != nil {
		return err
	}
	caCertPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertDER,
	}
	err = os.WriteFile(caCertFile, pem.EncodeToMemory(caCertPEM), 0640)
	if err != nil {
		return err
	}

	caKeyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	caKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: caKeyDER,
	}
	err = os.WriteFile(caKeyFile, pem.EncodeToMemory(caKeyPEM), 0600)
	if err != nil {
		return err
	}
	return nil
}

func getRandomSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}
