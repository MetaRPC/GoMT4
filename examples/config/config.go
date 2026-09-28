package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Login         int    `json:"Login"`
	Password      string `json:"Password"`
	Server        string `json:"Server"`
	DefaultSymbol string `json:"DefaultSymbol"`
	ApiKey        string `json:"ApiKey"`
}

type demoAccountResponse struct {
	ResultCode int    `json:"resultCode"`
	Login      any    `json:"login"`
	Password   string `json:"password"`
	Server     string `json:"server"`
	Error      string `json:"error"`
}

// OpenDemoAccount provisions a new live demo account on the specified server
func OpenDemoAccount(server, apiKey string) (int, string, string, error) {
	if server == "" || server == "RoboForex-Demo" {
		server = "MetaQuotes-Demo"
	}
	if apiKey == "" {
		apiKey = "TRIAL"
	}
	u := fmt.Sprintf("https://mt4.mrpc.pro/DemoAccount/Open?server=%s", url.QueryEscape(server))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, "", "", err
	}
	req.Header.Set("APIKey", apiKey)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", "", err
	}
	var res demoAccountResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, "", "", fmt.Errorf("failed to parse demo account response: %w", err)
	}
	var login int
	switch v := res.Login.(type) {
	case float64:
		login = int(v)
	case string:
		parsed, _ := strconv.Atoi(v)
		login = parsed
	}
	srv := res.Server
	if srv == "" {
		srv = server
	}
	return login, res.Password, srv, nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.ApiKey == "" || cfg.ApiKey == "TRIAL" {
		if envKey := os.Getenv("MRPC_API_KEY"); envKey != "" {
			cfg.ApiKey = envKey
		} else if cfg.ApiKey == "" {
			cfg.ApiKey = "TRIAL"
		}
	}
	if cfg.Login == 0 {
		fmt.Printf("  Auto-provisioning live demo account on %s...\n", cfg.Server)
		login, password, srv, err := OpenDemoAccount(cfg.Server, cfg.ApiKey)
		if err == nil && login != 0 {
			cfg.Login = login
			cfg.Password = password
			cfg.Server = srv
			fmt.Printf("✓ Live Demo Account Provisioned: #%d (Server: %s)\n", cfg.Login, cfg.Server)
		}
	}
	return &cfg, nil
}
