package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultCredentialsPath = "Library/Application Support/Granola/supabase.json"

type supabaseFile struct {
	WorkOSTokens string `json:"workos_tokens"`
}

type workOSTokens struct {
	AccessToken string `json:"access_token"`
}

// LoadCredentials loads the Granola access token.
// Resolution order: credentialsPath arg > GRANOLA_CREDENTIALS env > default path
func LoadCredentials(credentialsPath string) (string, error) {
	path := resolveCredentialsPath(credentialsPath)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("credentials file not found at %s (is Granola installed and logged in?)", path)
		}
		return "", fmt.Errorf("failed to read credentials file: %w", err)
	}

	var sf supabaseFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return "", fmt.Errorf("failed to parse credentials file: %w", err)
	}

	if sf.WorkOSTokens == "" {
		return "", fmt.Errorf("no workos_tokens found in credentials file")
	}

	var tokens workOSTokens
	if err := json.Unmarshal([]byte(sf.WorkOSTokens), &tokens); err != nil {
		return "", fmt.Errorf("failed to parse workos_tokens: %w", err)
	}

	if tokens.AccessToken == "" {
		return "", fmt.Errorf("no access_token found in workos_tokens")
	}

	return tokens.AccessToken, nil
}

func resolveCredentialsPath(override string) string {
	if override != "" {
		return override
	}

	if envPath := os.Getenv("GRANOLA_CREDENTIALS"); envPath != "" {
		return envPath
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return defaultCredentialsPath
	}

	return filepath.Join(home, defaultCredentialsPath)
}
