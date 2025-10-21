package hcloud_test

import (
	"context"
	"log"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func ExampleStorageBoxClient_Create() {
	ctx := context.Background()
	client := hcloud.NewClient(hcloud.WithToken("token"))

	opts := hcloud.StorageBoxCreateOpts{
		Name:           "my-storage-box",
		StorageBoxType: &hcloud.StorageBoxType{Name: "bx11"},
		Location:       &hcloud.Location{Name: "fsn1"},
		Password:       "my-secure-password",
		SSHKeys: []*hcloud.SSHKey{
			{
				PublicKey: "ssh-rsa AAAAB3NzaC1yc2E...", // Your full SSH public key
			},
		},
	}

	result, _, err := client.StorageBox.Create(ctx, opts)
	if err != nil {
		log.Fatalf("error creating Storage Box: %s\n", err)
	}

	if err := client.Action.WaitFor(ctx, result.Action); err != nil {
		log.Fatalf("error waiting for Storage Box creation: %s\n", err)
	}
}

func ExampleStorageBoxClient_Create_fromAPI() {
	ctx := context.Background()
	client := hcloud.NewClient(hcloud.WithToken("token"))

	sshKey, _, err := client.SSHKey.Get(ctx, "my-key")
	if err != nil {
		log.Fatalf("error fetching SSH Key: %s\n", err)
	}

	opts := hcloud.StorageBoxCreateOpts{
		Name:           "my-storage-box",
		StorageBoxType: &hcloud.StorageBoxType{Name: "bx11"},
		Location:       &hcloud.Location{Name: "fsn1"},
		Password:       "my-secure-password",
		SSHKeys:        []*hcloud.SSHKey{sshKey}, // Your existing SSH key fetched from the API
	}

	result, _, err := client.StorageBox.Create(ctx, opts)
	if err != nil {
		log.Fatalf("error creating Storage Box: %s\n", err)
	}

	if err := client.Action.WaitFor(ctx, result.Action); err != nil {
		log.Fatalf("error waiting for Storage Box creation: %s\n", err)
	}
}
