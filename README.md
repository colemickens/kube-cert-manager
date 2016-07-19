# Kubernetes Certificate Manager

Status: Almost working prototype

This is not an official Google Project.

`kube-cert-manager` is currently a prototype with the following features:

* Manage Lets Encrypt certificates based on a ThirdParty `certificate` resource.
* Will only ever support the dns-01 challenge for Google Cloud DNS.
* Saves Lets Encrypt certificates as Kubernetes secrets.

This repository will also include a end-to-end tutorial on how to dynamically load TLS certificates.

## Usage

Add the Certificate ThirdParty resource

```
kubectl create -f kubernetes/extensions/certificate.yaml 
```

Create a `certificate` object:

```
apiVersion: "stable.hightower.com/v1"
kind: "Certificate"
metadata:
  name: "hightowerlabs-dot-com"
spec:
  domain: "hightowerlabs.com"
  email: "kelsey.hightower@gmail.com"
  project: "hightowerlabs"
  serviceAccount: "hightowerlabs"
```

Create A Kubernetes secret for the hightowerlabs Google Cloud service account:

```
kubectl create secret generic hightowerlabs \
  --from-file=/Users/khightower/Desktop/service-account.json
```

> The secret key must be named `service-account.json`

```
kubectl describe secret hightowerlabs
```
```
Name:        hightowerlabs
Namespace:   default
Labels:      <none>
Annotations: <none>

Type:        Opaque

Data
====
service-account.json:   3915 bytes
```

### The Results

```
kubectl get secrets hightowerlabs.com
```
```
NAME                TYPE                DATA      AGE
hightowerlabs.com   kubernetes.io/tls   2         10m
```

```
kubectl describe secrets hightowerlabs.com
```
```
Name:        hightowerlabs.com
Namespace:   default
Labels:      <none>
Annotations: <none>

Type:        kubernetes.io/tls

Data
====
tls.crt:     1761 bytes
tls.key:     1679 bytes
```
