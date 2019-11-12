package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tomnomnom/linkheader"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed. Please check username/password.")
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

type authenticationResponse struct {
	ExpiresAt    string `json:"expiresAt"`
	Status       string `json:"status"`
	RelayState   string `json:"relayState"`
	SessionToken string `json:"sessionToken"`
	Embedded     struct {
		User struct {
			ID              string    `json:"id"`
			PasswordChanged time.Time `json:"passwordChanged"`
			Profile         struct {
				Login     string `json:"login"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    string `json:"locale"`
				TimeZone  string `json:"timeZone"`
			} `json:"profile"`
		} `json:"user"`
	} `json:"_embedded"`
}

// Authenticate the user with username and password.
// relayState can be used to add additional information.
func (s *UserService) Authenticate(ctx context.Context, username, password, relayState string) (*User, error) {
	u := "/api/v1/authn"

	post := struct {
		Username   string                 `json:"username"`
		Password   string                 `json:"password"`
		RelayState string                 `json:"relayState"`
		Options    map[string]interface{} `json:"options"`
	}{
		username,
		password,
		relayState,
		map[string]interface{}{
			"warnBeforePasswordExpired": false,
			"multiOptionalFactorEnroll": false,
		},
	}
	req, err := s.client.NewRequest("POST", u, post)
	if err != nil {
		return nil, err
	}

	var resp authenticationResponse
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != "SUCCESS" {
		return nil, ErrAuthenticationFailed
	}

	return s.GetUser(ctx, resp.Embedded.User.ID)
}

// GetUsersOptions allows to specify query options.
type GetUsersOptions struct {
	// PerPage is the max number of result in a page,
	// max value defined by Okta is 200.
	// If set to less than 1 or more than 200, it will be reset to 200.
	PerPage int

	// Pages is the max number of pages to retrieve.
	// If 0, it will page until there's no more results to retrieve.
	Pages int
}

// GetUsers returns all the users.
func (s *UserService) GetUsers(ctx context.Context, options *GetUsersOptions) ([]*User, error) {
	u := "/api/v1/users"

	perPage := options.PerPage
	if perPage < 1 || perPage > 200 {
		perPage = 200
	}

	uu, err := addOptions(u, &getUsersQuery{Limit: perPage})
	if err != nil {
		return nil, err
	}

	var users []*User
	var index = 0
LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context done")
		default:
			req, err := s.client.NewRequest("GET", uu, nil)
			if err != nil {
				return nil, err
			}

			if err := s.client.AddAuthorization(ctx, req); err != nil {
				return nil, err
			}

			var usersBatch []*User
			resp, err := s.client.Do(ctx, req, &usersBatch)
			if err != nil {
				return nil, err
			}
			users = append(users, usersBatch...)

			// Okta returns the next URL to page in the Link header,
			// see https://developer.okta.com/docs/reference/api-overview/#link-header.
			links := linkheader.Parse(strings.Join(resp.Header["Link"], ","))
			var next string
			for _, link := range links {
				if link.Rel == "next" {
					next = link.URL
				}
			}

			index++
			// Breaking if next url is empty, meaning there's no next results,
			// or if we've paged enough based on the provided options.
			if next == "" || index == options.Pages {
				break LOOP
			}

			uu = next
		}
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
