package main

import (
	"log"
	"time"
)

func main() {
	for {
		certificates, err := getCertificates()
		if err != nil {
			log.Println("Get certificates failed:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, c := range certificates {
			serviceAccount, err := getServiceAccountFromSecret(c.Spec.ServiceAccount)
			if err != nil {
				log.Println("Get service account failed:", err)
				continue
			}

			config := &CertificateConfig{
				Email:          c.Spec.Email,
				Domain:         c.Spec.Domain,
				Project:        c.Spec.Project,
				ServiceAccount: serviceAccount,
			}

			cert, key, err := RequestCertificate(config)
			if err != nil {
				log.Println("Request LetsEncrypt certificate failed:", err)
				continue
			}

			err = createSecret(c.Spec.Domain, cert, key)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
