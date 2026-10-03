package keadhcp4

import (
	"context"

	"github.com/venthe/terraform-provider-kea/kea"
)

func (c *Client) GetSubnets(ctx context.Context) ([]Subnet, error) {
	result, err := kea.ExecWithResponse[struct {
		Subnets []Subnet `json:"subnets"`
	}](ctx, c.transport, "subnet4-list", nil)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return []Subnet{}, nil
	}

	return result.Subnets, nil
}
