package okta

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

// GroupService deals with Okta groups.
type GroupService service

type Group struct {
	ID   string
	Name string
}

type group struct {
	ID                    string    `json:"id"`
	Created               time.Time `json:"created"`
	LastUpdated           time.Time `json:"lastUpdated"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
	ObjectClass           []string  `json:"objectClass"`
	Type                  string    `json:"type"`
	Profile               struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"profile"`
	Links struct {
		Logo []struct {
			Name string `json:"name"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"logo"`
		Users struct {
			Href string `json:"href"`
		} `json:"users"`
		Apps struct {
			Href string `json:"href"`
		} `json:"apps"`
	} `json:"_links"`
}

type getGroupsQuery struct {
	Limit int `url:"limit,omitempty"`
}

// GetGroups returns all the Okta groups.
func (s *GroupService) GetGroups(ctx context.Context) ([]*Group, error) {
	u := "/api/v1/groups"

	uu, err := addOptions(u, &getGroupsQuery{Limit: 200})
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

	var rawGroups []*group
	_, err = s.client.Do(ctx, req, &rawGroups)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	for _, g := range rawGroups {
		// skipping anything that is not an OKTA_GROUP (not sure when this can happen).
		if g.Type != "OKTA_GROUP" {
			continue
		}

		groups = append(groups, &Group{
			ID:   g.ID,
			Name: g.Profile.Name,
		})
	}

	return groups, nil
}

// GetGroupMembership returns all users from a group.
func (s *GroupService) GetGroupMembership(ctx context.Context, groupID string) ([]*User, error) {
	u := fmt.Sprintf("/api/v1/groups/%v/users", groupID)

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

func (s *GroupService) GetUserGroups(ctx context.Context, userID string) ([]*Group, error) {
	u := fmt.Sprintf("/api/v1/users/%v/groups", userID)

	uu, err := addOptions(u, &getGroupsQuery{Limit: 200})
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

	var rawGroups []*group
	_, err = s.client.Do(ctx, req, &rawGroups)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	for _, g := range rawGroups {
		// skipping anything that is not an OKTA_GROUP (not sure when this can happen).
		if g.Type != "OKTA_GROUP" {
			continue
		}

		groups = append(groups, &Group{
			ID:   g.ID,
			Name: g.Profile.Name,
		})
	}

	return groups, nil
}
