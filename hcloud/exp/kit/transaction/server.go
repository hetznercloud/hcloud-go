package transaction

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func WithServerShutdown(ctx context.Context, client *hcloud.Client, server *hcloud.Server, nextFunc func() error) error {
	return withServerShutdown(ctx, client, server, client.Server.Shutdown, nextFunc)
}

func WithServerPowerOff(ctx context.Context, client *hcloud.Client, server *hcloud.Server, nextFunc func() error) error {
	return withServerShutdown(ctx, client, server, client.Server.Poweroff, nextFunc)
}

func withServerShutdown(
	ctx context.Context,
	client *hcloud.Client,
	server *hcloud.Server,
	shutdownFunc func(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error),
	nextFunc func() error,
) error {
	shutdown, _, err := shutdownFunc(ctx, server)
	if err != nil {
		return err
	}

	if err := client.Action.WaitFor(ctx, shutdown); err != nil {
		return err
	}

	if err := nextFunc(); err != nil {
		return err
	}

	powerOn, _, err := client.Server.Poweron(ctx, server)
	if err != nil {
		return err
	}

	if err := client.Action.WaitFor(ctx, powerOn); err != nil {
		return err
	}

	return nil
}
