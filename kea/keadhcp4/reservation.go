package keadhcp4

import (
	"context"

	"github.com/venthe/terraform-provider-kea/kea"
	"github.com/venthe/terraform-provider-kea/kea/keaquery"
)

func (c *Client) GetReservations(ctx context.Context, subnetID uint32) ([]Reservation, error) {
	result, err := kea.ExecWithResponse[struct {
		Hosts []Reservation `json:"hosts"`
	}](ctx, c.transport, "reservation-get-all", struct {
		SubnetID uint32 `json:"subnet-id"`
	}{SubnetID: subnetID})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return []Reservation{}, nil
	}

	return result.Hosts, nil
}

func (c *Client) GetReservation(ctx context.Context, query keaquery.ReservationQuery) (*Reservation, error) {
	return kea.ExecWithResponse[Reservation](ctx, c.transport, "reservation-get", query)
}

func (c *Client) AddReservation(ctx context.Context, reservation Reservation) error {
	return kea.Exec(ctx, c.transport, "reservation-add", struct {
		Reservation Reservation `json:"reservation"`
	}{
		Reservation: reservation,
	})
}

func (c *Client) UpdateReservation(ctx context.Context, reservation Reservation) error {
	return kea.Exec(ctx, c.transport, "reservation-update", struct {
		Reservation Reservation `json:"reservation"`
	}{
		Reservation: reservation,
	})
}

func (c *Client) DeleteReservation(ctx context.Context, query keaquery.ReservationQuery) error {
	return kea.Exec(ctx, c.transport, "reservation-del", query)
}
