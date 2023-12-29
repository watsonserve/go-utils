package goutils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

func addCert(certPool *x509.CertPool, cert string) error {
	certBytes, err := os.ReadFile(cert)
	if nil == err {
		if !certPool.AppendCertsFromPEM(certBytes) {
			err = errors.New("failed to parse root certificate")
		}
	}
	return err
}

func LoadCryptoObj(crt, key, ca string) ([]tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if nil != err {
		return nil, nil, err
	}

	caCertPool := x509.NewCertPool()
	if nil != addCert(caCertPool, ca) {
		return nil, nil, err
	}

	return []tls.Certificate{cert}, caCertPool, nil
}

func GenTlsSrvConfig(crt, key, ca string, forClient bool) (*tls.Config, error) {
	certs, caCertPool, err := LoadCryptoObj(crt, key, ca)
	if nil != err {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: certs,
	}

	if forClient {
		cfg.RootCAs = caCertPool
	} else {
		cfg.ClientCAs = caCertPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert

	}

	return cfg, nil
}
