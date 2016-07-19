package main

import (
	"log"
	"time"
)

func main() {
	for {
		// this should be a watch
		certificates, err := getCertificates()
		if err != nil {
			log.Println("Get certificates failed:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, c := range certificates {
			exists, err := checkSecret(c.Spec.Domain)
			if err != nil {
				log.Println("Not able to check if secret exists", err)
				continue
			}
			if exists {
				log.Printf("%s secret already exists skipping...", c.Spec.Domain)
				continue
			}

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
