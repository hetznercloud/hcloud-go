package transaction

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/actionutil"
)

type ServerOffTx struct {
	client *hcloud.Client
	server *hcloud.Server
}

func NewServerOffTx(client *hcloud.Client, server *hcloud.Server) *ServerOffTx {
	return &ServerOffTx{client: client, server: server}
}

func (t *ServerOffTx) WithShutdown(ctx context.Context, nextFunc func() error) error {
	return t.do(ctx, t.client.Server.Shutdown, nextFunc)
}

func (t *ServerOffTx) WithPowerOff(ctx context.Context, nextFunc func() error) error {
	return t.do(ctx, t.client.Server.Poweroff, nextFunc)
}

func (t *ServerOffTx) do(
	ctx context.Context,
	shutdownFunc func(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error),
	nextFunc func() error,
) error {

	if err := t.ensureOff(ctx, shutdownFunc); err != nil {
		return err
	}

	if err := nextFunc(); err != nil {
		return err
	}

	if err := t.ensureOn(ctx); err != nil {
		return err
	}

	return nil
}

func (t *ServerOffTx) ensureOff(
	ctx context.Context,
	shutdownFunc func(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error),
) error {
	switch t.server.Status {
	case hcloud.ServerStatusOff:
		// Nothing to do

	case hcloud.ServerStatusStopping:
		actions, err := actionutil.RunningForResource(t.client, ctx, hcloud.ActionResourceTypeServer, t.server.ID)
		if err != nil {
			return err
		}

		if err := t.client.Action.WaitFor(ctx, actions...); err != nil {
			return err
		}

	default:
		shutdown, _, err := shutdownFunc(ctx, t.server)
		if err != nil {
			return err
		}

		if err := t.client.Action.WaitFor(ctx, shutdown); err != nil {
			return err
		}
	}

	return nil
}

func (t *ServerOffTx) ensureOn(
	ctx context.Context,
) error {
	switch t.server.Status {
	case hcloud.ServerStatusRunning:
		// Nothing to do

	case hcloud.ServerStatusStarting:
		actions, err := actionutil.RunningForResource(t.client, ctx, hcloud.ActionResourceTypeServer, t.server.ID)
		if err != nil {
			return err
		}

		if err := t.client.Action.WaitFor(ctx, actions...); err != nil {
			return err
		}

	default:
		powerOn, _, err := t.client.Server.Poweron(ctx, t.server)
		if err != nil {
			return err
		}

		if err := t.client.Action.WaitFor(ctx, powerOn); err != nil {
			return err
		}
	}

	return nil
}
