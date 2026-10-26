package command

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func endpointAuthorizationFor(ctx context.Context, opts *options, credentialsPath, baseURL string, e domain.Endpoint) (string, string, error) {
	if !domain.IsServerlessFamily(e.FamilyID) {
		return endpointAuthorization(ctx, opts, credentialsPath, baseURL)
	}
	path, err := credentialPath(credentialsPath)
	if err != nil {
		return "", "", err
	}
	store, err := agcapi.LoadCredentials(path)
	if err != nil {
		return "", "", err
	}
	credential, ok, err := resolveCredential(opts, store)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", &agcapi.AuthenticationError{Message: "select an API client profile or provide an explicit token"}
	}
	if e.Authentication != "service-account-or-api-client" && credential.Mode != "api-client" {
		return "", "", &agcapi.AuthenticationError{Message: "this endpoint requires an API client profile; Service Account credentials cannot authorize this service"}
	}
	if e.Authentication == "client-api-token" {
		token, err := agcapi.ServerlessClientToken(ctx, nil, baseURL, credential)
		if err != nil {
			return "", credential.ClientID, &agcapi.AuthenticationError{Message: "Client gateway authentication failed; verify the selected API client and region."}
		}
		return token, credential.ClientID, nil
	}
	authBase, err := agcapi.ServerlessTokenBaseURL(baseURL)
	if err != nil {
		return "", "", err
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("credential token redirects refused") }}
	token, err := agcapi.AccessToken(ctx, client, authBase, credential)
	if err != nil {
		return "", credential.ClientID, &agcapi.AuthenticationError{Message: "serverless authorization failed; verify selected profile, role and regional API client credentials"}
	}
	return token.AccessToken, credential.ClientID, nil
}
