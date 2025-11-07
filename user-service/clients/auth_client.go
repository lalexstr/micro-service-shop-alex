package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthClient struct {
	BaseURL string
}

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
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

// GetUser fetches user by ID from auth-service
func (c *AuthClient) GetUser(userID uint) (*User, error) {
	url := fmt.Sprintf("%s/api/v1/internal/users/%d", c.BaseURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service error: %s", string(body))
	}

	var authResp AuthServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	userBytes, err := json.Marshal(authResp.Data)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return nil, err
	}

	return &user, nil
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

// ListUsers fetches list of users from auth-service
func (c *AuthClient) ListUsers(page, size int, filters map[string]string) ([]User, int64, error) {
	url := fmt.Sprintf("%s/api/v1/users?page=%d&size=%d", c.BaseURL, page, size)
	
	// Add filters to URL
	for key, value := range filters {
		url += fmt.Sprintf("&%s=%s", key, value)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("auth service error: %s", string(body))
	}

	var authResp AuthServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, 0, err
	}

	// Extract items from response
	itemsBytes, err := json.Marshal(authResp.Data)
	if err != nil {
		return nil, 0, err
	}

	var result struct {
		Items []User `json:"items"`
		Total int64  `json:"total"`
	}
	if err := json.Unmarshal(itemsBytes, &result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.Total, nil
}

// UpdateUserRole updates user role in auth-service
func (c *AuthClient) UpdateUserRole(userID uint, role string) error {
	url := fmt.Sprintf("%s/api/v1/internal/users/%d/role", c.BaseURL, userID)
	
	body := map[string]string{"role": role}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	
	resp, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	resp.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	httpResp, err := client.Do(resp)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	
	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return fmt.Errorf("auth service error: %s", string(body))
	}
	
	return nil
}

