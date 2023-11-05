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
# How to install

### With go

```
go install github.com/cryptvault-cloud/vault-hoa@latest
```

# Over download
Go to the release and download the binary you need for your OS

# Getting start

Follow the documentation at [CryptVault.cloud](https://cryptvault.cloud/guides/create_your_cryptvault/overview)

