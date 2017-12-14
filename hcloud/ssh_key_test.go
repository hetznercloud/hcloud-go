package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestSSHKeyClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/ssh_keys/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"ssh_key": {
				"id": 1,
				"name": "My key",
				"fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
				"public_key": "ssh-rsa AAAjjk76kgf...Xt"
			}
		}`)
	})

	ctx := context.Background()
	sshKey, _, err := env.Client.SSHKey.Get(ctx, 1)
	if err != nil {
		t.Fatalf("SSHKey.Get failed: %s", err)
	}
	if sshKey == nil {
		t.Fatal("no SSH key")
	}
	if sshKey.ID != 1 {
		t.Errorf("unexpected SSH key ID: %v", sshKey.ID)
	}
}

func TestSSHKeyClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/ssh_keys", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		fmt.Fprint(w, `{
			"ssh_keys": [
				{
					"id": 1,
					"name": "My key",
					"fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
					"public_key": "ssh-rsa AAAjjk76kgf...Xt"
				},
				{
					"id": 2,
					"name": "Another key",
					"fingerprint": "c7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
					"public_key": "ssh-rsa AAAjjk76kgf...XX"
				}
			]
		}`)
	})

	opts := SSHKeyListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	sshKeys, _, err := env.Client.SSHKey.List(ctx, opts)
	if err != nil {
		t.Fatalf("SSHKey.List failed: %s", err)
	}
	if len(sshKeys) != 2 {
		t.Fatal("unexpected number of SSH keys")
	}
	if sshKeys[0].ID != 1 || sshKeys[1].ID != 2 {
		t.Fatalf("unexpected SSH key IDs: %d, %d", sshKeys[0].ID, sshKeys[1].ID)
	}
}

func TestSSHKeyAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/ssh_keys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			SSHKeys []schema.SSHKey `json:"ssh_keys"`
			Meta    schema.Meta     `json:"meta"`
		}{
			SSHKeys: []schema.SSHKey{
				{
					ID:          1,
					Name:        "My key",
					Fingerprint: "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
					PublicKey:   "ssh-rsa AAAjjk76kgf...Xt",
				},
				{
					ID:          2,
					Name:        "Another key",
					Fingerprint: "c7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
					PublicKey:   "ssh-rsa AAAjjk76kgf...XX",
				},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      2,
					TotalEntries: 2,
				},
			},
		})
	})

	ctx := context.Background()
	sshKeys, err := env.Client.SSHKey.All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sshKeys) != 2 {
		t.Fatalf("expected 2 SSH keys; got %d", len(sshKeys))
	}
	if sshKeys[0].ID != 1 || sshKeys[1].ID != 2 {
		t.Error("got wrong SSH keys")
	}
}

func TestSSHKeyClientDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/ssh_keys/1", func(w http.ResponseWriter, r *http.Request) {})

	ctx := context.Background()
	_, err := env.Client.SSHKey.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("SSHKey.Delete failed: %s", err)
	}
}

func TestSSHKeyClientCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/ssh_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"ssh_key": {
				"id": 1,
				"name": "My key",
				"fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
				"public_key": "ssh-rsa AAAjjk76kgf...Xt"
			}
		}`)
	})

	ctx := context.Background()
	opts := SSHKeyCreateOpts{
		Name:      "My key",
		PublicKey: "ssh-rsa AAAjjk76kgf...Xt",
	}
	sshKey, _, err := env.Client.SSHKey.Create(ctx, opts)
	if err != nil {
		t.Fatalf("SSHKey.Get failed: %s", err)
	}
	if sshKey.ID != 1 {
		t.Errorf("unexpected SSH key ID: %v", sshKey.ID)
	}
}
