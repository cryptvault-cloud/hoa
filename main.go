package main

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	client "github.com/cryptvault-cloud/api"
	"github.com/cryptvault-cloud/helper"
	"github.com/cryptvault-cloud/vault-hoa/logger"
	"github.com/urfave/cli/v2"
)

const (
	CliLogLevel    = "logLevel"
	CliIdentityKey = "identity_key"
	CliVaultId     = "vault_id"

	CliKeyValueMap = "key_value_map"
	CliServerUrl   = "serverUrl"
	CliExecuter    = "executer"

	App = "VAULT"
)

func getFlagEnvByFlagName(flagName string) string {
	fName := strings.ToUpper(flagName)
	fName = strings.ReplaceAll(fName, App+"_", "")
	return fmt.Sprintf("%s_%s", App, fName)
}

func main() {
	runner := Runner{}
	app := &cli.App{
		Usage:     "vault-client",
		ArgsUsage: "{application commend to executed }",
		Action:    runner.Inject,
		Before:    runner.Before,
		Flags: []cli.Flag{

			&cli.StringFlag{
				Name:    CliServerUrl,
				EnvVars: []string{getFlagEnvByFlagName(CliServerUrl)},
				Value:   "https://api.cryptvault.cloud/query",
				Usage:   "vault-server endpoint",
			},
			&cli.StringFlag{
				Name:    CliLogLevel,
				EnvVars: []string{getFlagEnvByFlagName(CliLogLevel)},
				Value:   "info",
				Usage:   "Loglevel debug, info, warn, error",
			},
			&cli.StringFlag{
				Name:     CliIdentityKey,
				EnvVars:  []string{getFlagEnvByFlagName(CliIdentityKey)},
				Usage:    "Identity Private Key",
				Required: true,
			},
			&cli.StringFlag{
				Name:     CliVaultId,
				Aliases:  []string{"vault"},
				EnvVars:  []string{getFlagEnvByFlagName(CliVaultId)},
				Usage:    "Vaultid to connect",
				Required: true,
			},
			&cli.PathFlag{
				Name:    CliKeyValueMap,
				Aliases: []string{"remap"},
				EnvVars: []string{getFlagEnvByFlagName(CliKeyValueMap)},
				Usage:   "path to Key=value map file",
				Value:   "vault.env",
			},
			&cli.StringFlag{
				Name:    CliExecuter,
				Aliases: []string{"sh"},
				EnvVars: []string{getFlagEnvByFlagName(CliExecuter)},
				Value:   "/bin/sh",
				Usage:   "Shell to execute child application",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "keys",
				Usage:  "Show all Vault keys related to given identity",
				Action: runner.ShowKeys,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println("Application start failed", "error", err)
	}
}

type Runner struct {
	client      client.Api
	pClient     client.ProtectedApi
	publickey   helper.Base64PublicPem
	privateKey  *ecdsa.PrivateKey
	identityID  string
	keyValueMap map[string]string
}

func (r *Runner) Before(c *cli.Context) error {
	logs, err := logger.Initialize(c.String(CliLogLevel))
	if err != nil {
		return err
	}
	readFile, err := os.Open(c.Path(CliKeyValueMap))
	if err != nil {
		logs.Info("No Key.value map file found. Vault name will be used as key")
		r.keyValueMap = nil
	} else {

		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)
		keyValueMap := make(map[string]string)
		for fileScanner.Scan() {
			line := strings.TrimSpace(fileScanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			keyValue := strings.Split(line, "=")
			if len(keyValue) != 2 {
				continue
			}
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			if key == "" || value == "" {
				continue
			}
			key = strings.ReplaceAll(key, ".", "_")

			keyValueMap[strings.ToUpper(key)] = strings.ToUpper(value)
		}
		r.keyValueMap = keyValueMap
		err := readFile.Close()
		if err != nil {
			return err
		}
	}

	r.client = client.NewApi(c.String(CliServerUrl), http.DefaultClient)

	key, err := helper.GetPrivateKeyFromB64String(c.String(CliIdentityKey))
	if err != nil {
		return err
	}

	b64PubKey, err := helper.NewBase64PublicPem(&key.PublicKey)
	if err != nil {
		return err
	}
	r.pClient = *r.client.GetProtectedApi(key, c.String(CliVaultId))
	r.publickey = b64PubKey
	r.privateKey = key
	id, err := r.publickey.GetIdentityId(c.String(CliVaultId))
	if err != nil {
		return err
	}
	r.identityID = id
	return nil
}

func (r *Runner) Inject(c *cli.Context) error {
	executionCommand := c.Args().Get(0)
	m, err := r.getAlleValues()
	if err != nil {
		return err
	}
	envs := make([]string, 0)
	for k, v := range m {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	envs = append(envs, os.Environ()...)
	cmd := exec.Command(c.String(CliExecuter), "-c", executionCommand)
	cmd.Env = envs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}

func (r *Runner) ShowKeys(c *cli.Context) error {
	m, err := r.getAlleValues()
	if err != nil {
		return err
	}
	for k := range m {
		fmt.Printf("%s=\"****\"\r\n", k)
	}

	return nil
}

func (r *Runner) getAlleValues() (map[string]string, error) {
	resp, err := r.pClient.GetAllRelatedValues(r.identityID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, v := range resp {
		iv, err := r.pClient.GetIdentityValueById(v)
		if err != nil {
			return nil, err
		}
		keyName := strings.ToUpper(fmt.Sprintf("%s.%s", App, iv.Name))
		keyName = strings.ReplaceAll(keyName, ".", "_")
		if v, ok := r.keyValueMap[keyName]; ok {
			keyName = v
		}
		result[keyName] = iv.Value
	}
	return result, nil
}
