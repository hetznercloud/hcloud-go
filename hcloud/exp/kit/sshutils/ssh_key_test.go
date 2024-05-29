package sshutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {
	privBytes, pubBytes, err := GenerateKeyPair()
	assert.Nil(t, err)

	priv := string(privBytes)
	pub := string(pubBytes)

	if !(strings.HasPrefix(priv, "-----BEGIN OPENSSH PRIVATE KEY-----\n") &&
		strings.HasSuffix(priv, "-----END OPENSSH PRIVATE KEY-----\n")) {
		assert.Fail(t, "private key is invalid", priv)
	}

	if !strings.HasPrefix(pub, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA") {
		assert.Fail(t, "public key is invalid", pub)
	}
}

func TestGeneratePublicKey(t *testing.T) {
	privBytes, pubBytesOrig, err := GenerateKeyPair()
	require.NoError(t, err)

	pubBytes, err := GeneratePublicKey(privBytes)
	require.NoError(t, err)

	pub := string(pubBytes)
	priv := string(privBytes)

	if !(strings.HasPrefix(priv, "-----BEGIN OPENSSH PRIVATE KEY-----\n") &&
		strings.HasSuffix(priv, "-----END OPENSSH PRIVATE KEY-----\n")) {
		assert.Fail(t, "private key is invalid", priv)
	}

	if !strings.HasPrefix(pub, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA") {
		assert.Fail(t, "public key is invalid", pub)
	}

	assert.Equal(t, pubBytesOrig, pubBytes)
}

func TestGetPublicKeyFingerprint(t *testing.T) {
	fingerprint, err := GetPublicKeyFingerprint([]byte(`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIccHCW76xx2rrPAUrjnuT6IjpEF1O+/U4IByVgv99Oi`))
	require.NoError(t, err)
	assert.Equal(t, "77:79:69:b1:4d:c6:b6:45:6a:e9:52:29:04:3e:59:48", fingerprint)
}
