FROM alpine
ADD kube-cert-manager /kube-cert-manager
ENTRYPOINT ["/kube-cert-manager"]
