package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"pivium/internal/config"
)

// ProxmoxClient is a client for the Proxmox API.

type ProxmoxClient struct {
	httpClient *http.Client
	apiURL     string
	node       string
}

// NewProxmoxClient creates a new Proxmox API client.
func NewProxmoxClient(apiURL string) (*ProxmoxClient, error) {
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	node, err := GetNodeName()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Proxmox uses self-signed certs
	}
	client := &ProxmoxClient{
		httpClient: &http.Client{Transport: tr},
		apiURL:     apiURL,
		node:       node,
	}

	// Add the authorization header to all requests.
	// We need to create a custom RoundTripper to do this.
	client.httpClient.Transport = &authTransport{
		Transport: tr,
		token:     apiToken,
	}

	return client, nil
}

// authTransport is a custom http.RoundTripper that adds the Proxmox API token to each request.
type authTransport struct {
	Transport http.RoundTripper
	token     string
}

// RoundTrip adds the Authorization header to the request.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "PVEAPIToken="+t.token)
	return t.Transport.RoundTrip(req)
}

// GetNodeName gets the hostname of the current machine.
func GetNodeName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// GetVMs gets a list of all VMs on the node.
func (c *ProxmoxClient) GetVMs() ([]config.ProxmoxResource, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api2/json/nodes/%s/qemu", c.apiURL, c.node))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []config.ProxmoxResource `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetLXC gets a list of all LXC containers on the node.
func (c *ProxmoxClient) GetLXC() ([]config.ProxmoxResource, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api2/json/nodes/%s/lxc", c.apiURL, c.node))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []config.ProxmoxResource `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetResources gets a list of all resources on the node.
func (c *ProxmoxClient) GetResources() ([]config.ProxmoxResource, error) {
	vms, err := c.GetVMs()
	if err != nil {
		return nil, err
	}

	lxcs, err := c.GetLXC()
	if err != nil {
		return nil, err
	}

	return append(vms, lxcs...), nil
}
