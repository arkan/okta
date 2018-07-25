package okta

import (
	"context"
	"fmt"
)

// UserService handles users operations.
type UserService service

// User represents a Okta user.
type User struct {
	ID              string            `json:"id"`
	Status          string            `json:"status"`
	LastLogin       string            `json:"last_login"`
	Created         string            `json:"created"`
	LastUpdated     string            `json:"last_updated"`
	PasswordChanged string            `json:"password_changed"`
	Profile         map[string]string `json:"profile"`
}

type getUsersQuery struct {
	Limit int `url:"limit,omitempty"`
}

// GetUsers returns all the users.
func (s *UserService) GetUsers(ctx context.Context) ([]*User, error) {
	u := "/api/v1/users"

	uu, err := addOptions(u, &getUsersQuery{Limit: 200})
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", uu, nil)
	if err != nil {
		return nil, err
	}

	if err := s.client.AddAuthorization(ctx, req); err != nil {
		return nil, err
	}

	var users []*User
	_, err = s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser returns a user.
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	u := fmt.Sprintf("/api/v1/users/%v", id)

	var user User

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	if err := s.client.AddAuthorization(ctx, req); err != nil {
		return nil, err
	}

	_, err = s.client.Do(ctx, req, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateCustomAttributes returns a user.
func (s *UserService) UpdateCustomAttributes(ctx context.Context, id string, attributes map[string]string) error {
	u := fmt.Sprintf("/api/v1/users/%v", id)

	post := struct {
		Profile map[string]string `json:"profile"`
	}{
		Profile: attributes,
	}
	req, err := s.client.NewRequest("POST", u, post)
	if err != nil {
		return err
	}

	if err := s.client.AddAuthorization(ctx, req); err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
