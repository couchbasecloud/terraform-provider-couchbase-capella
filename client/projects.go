package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) PostProject(project CreateProjectRequest, organizationId string) (*CreateProjectResponse, error) {
	rb, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v4/organizations/%s/projects", c.HostURL, organizationId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	newProject := CreateProjectResponse{}
	err = json.Unmarshal(body, &newProject)
	if err != nil {
		return nil, err
	}

	return &newProject, nil
}

func (c *Client) GetProject(organizationId, projectId string) (*GetProjectResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	project := GetProjectResponse{}
	err = json.Unmarshal(body, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) GetProjects(organizationId string) (*GetProjectsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v4/organizations/%s/projects", c.HostURL, organizationId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	projects := GetProjectsResponse{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func (c *Client) UpdateProject(project PutProjectRequest, organizationId, projectId string) error {
	rb, err := json.Marshal(project)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteProject(organizationId, projectId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
