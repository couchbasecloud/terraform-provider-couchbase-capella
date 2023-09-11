package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) PostProject(ctx context.Context, project CreateProjectRequest, organizationId string) (*CreateProjectResponse, error) {
	rb, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v4/organizations/%s/projects", c.HostURL, organizationId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newProject := CreateProjectResponse{}
	err = json.Unmarshal(resp.Body, &newProject)
	if err != nil {
		return nil, err
	}

	return &newProject, nil
}

func (c *Client) GetProject(ctx context.Context, organizationId, projectId string) (*GetProjectResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	project := GetProjectResponse{}
	err = json.Unmarshal(resp.Body, &project)
	if err != nil {
		return nil, err
	}

	project.Etag = resp.HTTPResponse.Header.Get("ETag")

	return &project, nil
}

func (c *Client) GetProjects(ctx context.Context, organizationId string) (*GetProjectsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v4/organizations/%s/projects", c.HostURL, organizationId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	projects := GetProjectsResponse{}
	err = json.Unmarshal(resp.Body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func (c *Client) UpdateProject(ctx context.Context, project PutProjectRequest, organizationId, projectId string, params *PutProjectParams) error {
	rb, err := json.Marshal(project)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	if params != nil && params.IfMatch != nil {
		req.Header.Set("If-Match", *params.IfMatch)
	}

	_, err = c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteProject(ctx context.Context, organizationId, projectId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v4/organizations/%s/projects/%s", c.HostURL, organizationId, projectId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
