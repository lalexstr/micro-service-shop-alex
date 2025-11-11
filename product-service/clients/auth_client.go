package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthClient struct {
	BaseURL string
}

type UserRoleResponse struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

type AuthServiceResponse struct {
	Data interface{} `json:"data"`
}

func NewAuthClient() *AuthClient {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &AuthClient{BaseURL: baseURL}
}

// GetUserRole fetches user role from auth-service
func (c *AuthClient) GetUserRole(userID uint) (string, error) {
	url := fmt.Sprintf("%s/api/v1/internal/users/%d/role", c.BaseURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth service error: %s", string(body))
	}

	var authResp AuthServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	roleBytes, err := json.Marshal(authResp.Data)
	if err != nil {
		return "", err
	}

	var roleResp UserRoleResponse
	if err := json.Unmarshal(roleBytes, &roleResp); err != nil {
		return "", err
	}

	return roleResp.Role, nil
}






