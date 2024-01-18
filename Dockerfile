FROM alpine
ENTRYPOINT ["/usr/bin/vault-hoa"]
COPY vault-hoa /usr/bin/vault-hoa