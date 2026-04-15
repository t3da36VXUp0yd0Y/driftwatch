package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds the relevant runtime state of a container.
type ContainerInfo struct {
	Name  string
	Image string
	State string
}

// Client wraps the Docker SDK client.
type Client struct {
	cli *client.Client
}

// NewClient creates a new Docker client using environment variables.
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker: failed to create client: %w", err)
	}
	return &Client{cli: cli}, nil
}

// Close releases resources held by the client.
func (c *Client) Close() error {
	return c.cli.Close()
}

// GetRunningContainers returns info for all running containers whose names
// match the provided service names. The map key is the service name.
func (c *Client) GetRunningContainers(ctx context.Context, serviceNames []string) (map[string]ContainerInfo, error) {
	f := filters.NewArgs()
	f.Add("status", "running")

	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{Filters: f})
	if err != nil {
		return nil, fmt.Errorf("docker: failed to list containers: %w", err)
	}

	wanted := make(map[string]struct{}, len(serviceNames))
	for _, s := range serviceNames {
		wanted[s] = struct{}{}
	}

	result := make(map[string]ContainerInfo)
	for _, ctr := range containers {
		for _, name := range ctr.Names {
			// Docker prefixes names with "/"
			clean := trimLeadingSlash(name)
			if _, ok := wanted[clean]; ok {
				result[clean] = ContainerInfo{
					Name:  clean,
					Image: ctr.Image,
					State: ctr.State,
				}
			}
		}
	}
	return result, nil
}

func trimLeadingSlash(s string) string {
	if len(s) > 0 && s[0] == '/' {
		return s[1:]
	}
	return s
}
