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

	if "" == ca {
		return []tls.Certificate{cert}, nil, nil
	}

	caCertPool := x509.NewCertPool()
	if nil != addCert(caCertPool, ca) {
		return nil, nil, err
	}

	return []tls.Certificate{cert}, caCertPool, nil
}

type TlsFlag int

const (
	TLSFLAG_IGNORE = 0 // 000
	TLSFLAG_CLIENT = 1 // 001
	TLSFLAG_SERVER = 2 // 010
	TLSFLAG_VERIFY = 3 // 011
)

func GenTlsConfig(flag TlsFlag, crt, key, ca string) (*tls.Config, error) {
	if TLSFLAG_IGNORE == flag {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}

	certs, caCertPool, err := LoadCryptoObj(crt, key, ca)
	if nil != err {
		return nil, err
	}

	cfg := &tls.Config{Certificates: certs}
	if TLSFLAG_VERIFY == flag {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if "" == ca {
		return cfg, nil
	}

	switch flag {
	case TLSFLAG_CLIENT:
		cfg.RootCAs = caCertPool
	case TLSFLAG_SERVER:
		cfg.ClientCAs = caCertPool
	}

	return cfg, nil
}
