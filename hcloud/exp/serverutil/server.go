package serverutil

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/actionutil"
)

func EnsureShutdown(ctx context.Context, client *hcloud.Client, server *hcloud.Server) error {
	return ensureOff(ctx, client, client.Server.Shutdown, server)
}

func EnsurePowerOff(ctx context.Context, client *hcloud.Client, server *hcloud.Server) error {
	return ensureOff(ctx, client, client.Server.Poweroff, server)
}

func ensureOff(
	ctx context.Context,
	client *hcloud.Client,
	clientShutdownFunc func(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error),
	server *hcloud.Server,
) error {
	switch server.Status {
	case hcloud.ServerStatusOff:
		return nil // Nothing to do

	case hcloud.ServerStatusStopping:
		actions, err := actionutil.RunningForResource(ctx, client, hcloud.ActionResourceTypeServer, server.ID)
		if err != nil {
			return err
		}
		if err := client.Action.WaitFor(ctx, actions...); err != nil {
			return err
		}

	default:
		shutdown, _, err := clientShutdownFunc(ctx, server)
		if err != nil {
			return err
		}
		if err := client.Action.WaitFor(ctx, shutdown); err != nil {
			return err
		}
	}

	return nil
}

func EnsurePowerOn(ctx context.Context, client *hcloud.Client, server *hcloud.Server) error {
	switch server.Status {
	case hcloud.ServerStatusRunning:
		return nil // Nothing to do

	case hcloud.ServerStatusStarting:
		actions, err := actionutil.RunningForResource(ctx, client, hcloud.ActionResourceTypeServer, server.ID)
		if err != nil {
			return err
		}
		if err := client.Action.WaitFor(ctx, actions...); err != nil {
			return err
		}

	default:
		powerOn, _, err := client.Server.Poweron(ctx, server)
		if err != nil {
			return err
		}
		if err := client.Action.WaitFor(ctx, powerOn); err != nil {
			return err
		}
	}

	return nil
}
