package keadhcp4

import "github.com/venthe/terraform-provider-kea/kea"

type Client struct {
	transport kea.Transport
}

func NewClient(transport kea.Transport) *Client {
	return &Client{transport}
}
