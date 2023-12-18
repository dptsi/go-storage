package its

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getOidcWellKnownConfig(providerUrl string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/.well-known/openid-configuration", providerUrl)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get oidc well-known config: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read oidc well-known config: %w", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode oidc well-known config: %w", err)
	}

	return data, nil
}

func getOidcTokenEndpoint(providerUrl string) (string, error) {
	data, err := getOidcWellKnownConfig(providerUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get oidc token endpoint: %w", err)
	}

	tokenEndpoint, ok := data["token_endpoint"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get oidc token endpoint: %w", err)
	}

	return tokenEndpoint, nil
}
