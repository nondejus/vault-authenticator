# vaultlibs


[![Status](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQ1hMb3dNeE4zYjFzUzZYMi9QekU1SlN0bksyVVY4QnE3WjRhVGw5MnB3T2U2cTZvQ2hYMlRqb3pNWXJoLytQR1N6WCtDY01pVUJHdVV4MkpuQnVKaE5RPSIsIml2UGFyYW1ldGVyU3BlYyI6InA4YVpZdU4ybDNRUmEvbE8iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)](https://us-east-2.console.aws.amazon.com/codesuite/codebuild/projects/Vaultlibs/history?region=us-east-2)


Useful golang functions for interacting with Vault.

Vault is a great tool, but programming against it sometimes requires one to go more deeply than one wants to in order to navigate these deep waters.

This library abstracts some of the work and provided some high level bindings so that the author of a tool that _uses_ Vault doesn't need to be an expert in Vault.

The crown jewel is the `authenticator` object which has has one main method: Auth().  This method tries to authenticate to Vault in a number of ways and returns an authenticated Vault client for the first one that succeeds.

## Configuration

To configure `authenticator`, create the object via it's constructor:

    auth = authenticator.NewAuthenticator()
    
    
Then set the address of the Vault server:

	auth.SetAddress("https://vault.corp.scribd.com")
	
	
Set a private CA if you're using one:

	auth.SetCACertificate(`-----BEGIN CERTIFICATE-----
	...
    -----END CERTIFICATE-----
    `)


Set Auth methods.  These will be tried in order:

	auth.SetAuthMethods([]string{
		"iam",
		"k8s",
		"tls",
		"ldap",
	})
	
If your usernames don't necessarily map to posix users on the system:

	auth.SetUsernameFunc(somelib.GetUsername)
	

Finally, if using TLS Auth, set the locations of the client certs:

	auth.SetTlsClientCrtPath("/path/to/cert.crt")
	auth.SetTlsClientKeyPath("/path/to/key.key")
	
	
After that, simply run:

    client, err := auth.Auth()
    if err != nil {
      log.Fatalf("Auth failed: %s", err)
    }
    
    path := "/secret/foo
    
    secret, err := authenticator.GetSecret(client, path)
    if err != nil {
      log.Fatalf("Failed getting secret from %s: %s", path, err)
    }
    
    ... do something with secret ...
