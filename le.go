package main

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/kelseyhightower/kube-cert-manager/provider/dns/googlecloud"

	"github.com/xenolf/lego/acme"
)

type CertificateConfig struct {
	Email          string
	Domain         string
	Project        string
	ServiceAccount []byte
}

var acmeURL = "https://acme-staging.api.letsencrypt.org/directory"

func RequestCertificate(config *CertificateConfig) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	user := User{
		Email: config.Email,
		key:   privateKey,
	}
	client, err := acme.NewClient(acmeURL, &user, acme.RSA2048)
	if err != nil {
		return nil, nil, err
	}

	provider, err := googlecloud.NewDNSProvider(config.Project, config.ServiceAccount)
	if err != nil {
		return nil, nil, err
	}

	client.SetChallengeProvider(acme.DNS01, provider)
	client.ExcludeChallenges([]acme.Challenge{
		acme.HTTP01,
		acme.TLSSNI01,
	})

	user.Registration, err = client.Register()
	if err != nil {
		return nil, nil, err
	}

	err = client.AgreeToTOS()
	if err != nil {
		return nil, nil, err
	}

	cr, failures := client.ObtainCertificate([]string{config.Domain}, false, nil)
	if len(failures) > 0 {
		err := failures[config.Domain]
		return nil, nil, err
	}

	return cr.Certificate, cr.PrivateKey, nil
}
