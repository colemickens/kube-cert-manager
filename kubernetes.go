package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	apiHost              = "http://127.0.0.1:8001"
	certificatesEndpoint = "/apis/stable.hightower.com/v1/namespaces/default/certificates"
	secretsEndpoint      = "/api/v1/namespaces/default/secrets"
)

type Certificate struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       CertificateSpec   `json:"spec"`
}

type CertificateSpec struct {
	Domain         string `json:"domain"`
	Email          string `json:"email"`
	Project        string `json:"project"`
	ServiceAccount string `json:"serviceAccount"`
}

type CertificateList struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Items      []Certificate     `json:"items"`
}

type Secret struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	Metadata   map[string]string `json:"metadata"`
	Data       map[string]string `json:"data"`
	Type       string            `json:"type"`
}

func getCertificates() ([]Certificate, error) {
	resp, err := http.Get(apiHost + certificatesEndpoint)
	if err != nil {
		return nil, err
	}

	var certList CertificateList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&certList)
	if err != nil {
		return nil, err
	}
	return certList.Items, nil
}

func getServiceAccountFromSecret(name string) ([]byte, error) {
	resp, err := http.Get(apiHost + secretsEndpoint + "/" + name)
	if err != nil {
		return nil, err
	}
	var secret Secret
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&secret)
	if err != nil {
		return nil, err
	}

	data := secret.Data["service-account.json"]
	serviceAccount, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func checkSecret(name string) (bool, error) {
	resp, err := http.Get(apiHost + secretsEndpoint + "/" + name)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}

func createSecret(domain string, cert, key []byte) error {
	metadata := make(map[string]string)
	metadata["name"] = domain

	data := make(map[string]string)
	data["tls.crt"] = base64.StdEncoding.EncodeToString(cert)
	data["tls.key"] = base64.StdEncoding.EncodeToString(key)

	secret := &Secret{
		ApiVersion: "v1",
		Data:       data,
		Kind:       "Secret",
		Metadata:   metadata,
		Type:       "kubernetes.io/tls",
	}

	b := make([]byte, 0)
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(secret)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiHost+secretsEndpoint, "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return errors.New("Secrets: Unexpected HTTP status code" + resp.Status)
	}
	return nil
}
