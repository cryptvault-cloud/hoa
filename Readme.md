# CryptVault Higher Order Application (HOA)

Easy way to load all secrets into the application via environment variables

Check out [CryptVault](https://cryptvault.cloud).

```
NAME:
   vault-hoa - vault-client

USAGE:
   vault-hoa [global options] command [command options] {application commend to executed }

COMMANDS:
   keys     Show all Vault keys related to given identity
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --serverUrl value                     vault-server endpoint (default: "https://api.cryptvault.cloud/query") [$VAULT_SERVERURL]
   --logLevel value                      Loglevel debug, info, warn, error (default: "info") [$VAULT_LOGLEVEL]
   --identity_key value                  Identity Private Key [$VAULT_IDENTITY_KEY]
   --vault_id value, --vault value       Vaultid to connect [$VAULT_ID]
   --key_value_map value, --remap value  path to Key=value map file (default: "vault.env") [$VAULT_KEY_VALUE_MAP]
   --executer value, --sh value          Shell to execute child application (default: "/bin/sh") [$VAULT_EXECUTER]
   --help, -h                            show help
```

# Configuration 

It is possible to map the keys of the environment variables to new names via configuration.

vault.env
```
VAULT_VALUE_KEY_A=NEW_KEY_NAME
VAULT_VALUE_KEY_B=OTHER_NEW_KEY_NAME
```

to find out current environment keys execute `vault-hoa keys`

# How to install

### With go

```
go install github.com/cryptvault-cloud/vault-hoa@latest
```
# With Docker
```sh
docker pull ghcr.io/cryptvault-cloud/vault-hoa:latest

docker run -e VAULT_ID=[your-vault-id] -e VAULT_IDENTITY_KEY=[identity private key]  ghcr.io/cryptvault-cloud/vault-hoa:latest env  
```

## Use multistage build

you can simply use mutistage docker to transfer the vault hoa binary from the existing docker image to your target image

```dockerfile
# hoa container
FROM ghcr.io/cryptvault-cloud/vault-hoa:latest as hoa

# your target dockerfile
FROM alpine:latest
WORKDIR /my_dir
COPY ./your_app_execution ./app
COPY --from=hoa /usr/bin/vault-hoa ./vault-hoa
CMD ["./vault-hoa", "./app"]

```

or by downloading the binary 

(For example, for target images that have a different CPU structure )

```dockerfile
FROM alpine:latest as hoa
RUN wget -O /usr/bin/vault-hoa https://github.com/cryptvault-cloud/vault-hoa/releases/download/v0.0.10/vault-hoa_0.0.10_linux_arm64
RUN chmod 755 /usr/bin/vault-hoa
# your target dockerfile
FROM alpine:latest
WORKDIR /your_app_execution
COPY ./hoa_test ./app
COPY --from=hoa /usr/bin/vault-hoa ./vault-hoa
CMD ["./vault-hoa", "./app"]


```

# Over download
Go to the release and download the binary you need for your OS


# Getting start

Follow the documentation at [CryptVault.cloud](https://cryptvault.cloud/guides/create_your_cryptvault/overview)

