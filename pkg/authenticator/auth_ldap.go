package authenticator

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"strings"
)

// LDAPLogin logs the user into vault via LDAP and obtains a token.  (Really only intended for user usage)
func LDAPLogin(authenticator *Authenticator, client *api.Client) (err error) {
	if !authenticator.Prompt {
		err = errors.New("Interactive prompts disabled.  LDAP Auth cannot proceed.")
		return err
	}
	apiConfig := api.DefaultConfig()
	err = apiConfig.ReadEnvironment()
	if err != nil {
		err = errors.Wrapf(err, "failed to inject environment into client authenticator")
		return err
	}

	if apiConfig.Address == "https://127.0.0.1:8200" {
		if authenticator.Address != "" {
			apiConfig.Address = authenticator.Address
		}
	}

	client, err = api.NewClient(apiConfig)

	var username string

	if authenticator.Identifier == "" {
		username, err = authenticator.UsernameFunc()
		if err != nil {
			err = errors.Wrap(err, "failed getting username")
			return err
		}

		verboseOutput(authenticator.Verbose, "Username: %s", username)

		if username == "" {
			err = errors.New("No username.  Cannot authenticate")
			return err
		}
	} else {
		username = authenticator.Identifier
	}

	path := fmt.Sprintf("/auth/ldap/login/%s", username)
	data := make(map[string]interface{})

	fmt.Println("")
	fmt.Printf("Enter Your LDAP password\n")

	passwordBytes, err := terminal.ReadPassword(0)
	if err != nil {
		err = errors.Wrapf(err, "failed to read password from terminal")
		return err
	}

	passwordString := string(passwordBytes)
	passwordString = strings.TrimSuffix(passwordString, "\n")

	data["password"] = passwordString

	resp, err := client.Logical().Write(path, data)
	if err != nil {
		err = errors.Wrapf(err, "failed submitting auth data to vault")
		return err
	}

	if resp != nil {
		auth := resp.Auth
		token := auth.ClientToken

		if token != "" {
			client.SetToken(token)
			homeDir, err := homedir.Dir()
			if err != nil {
				err = errors.Wrapf(err, "failed to derive user home dir")
				return err
			}

			tokenFile := fmt.Sprintf("%s/%s", homeDir, DEFAULT_VAULT_TOKEN_FILE)

			// write the token to the filesystem where expected for future use
			err = ioutil.WriteFile(tokenFile, []byte(token), 0644)
			if err != nil {
				err = errors.Wrapf(err, "failed to write token file: %s", tokenFile)
				return err
			}

			return err
		}
	}

	err = errors.New(VAULT_AUTH_FAIL)

	return err
}
