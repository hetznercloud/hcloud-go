package metadata

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	Client *Client
}

func (env *testEnv) Teardown() {
	env.Server.Close()
	env.Server = nil
	env.Mux = nil
	env.Client = nil
}

func newTestEnv() testEnv {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(
		WithEndpoint(server.URL),
		// This makes sure that our instrumentation does not cause any panics/errors
		WithInstrumentation(prometheus.DefaultRegisterer),
	)
	return testEnv{
		Server: server,
		Mux:    mux,
		Client: client,
	}
}

func TestClient_Base(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	env.Mux.HandleFunc("/sanitized", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sanitized  \n"))
	})
	env.Mux.HandleFunc("/no-content", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	env.Mux.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
	})

	if body, err := env.Client.get("/ok"); assert.NoError(t, err) {
		assert.Equal(t, "ok", body)
	}
	if body, err := env.Client.get("/sanitized"); assert.NoError(t, err) {
		assert.Equal(t, "sanitized", body)
	}
	if body, err := env.Client.get("/no-content"); assert.NoError(t, err) {
		assert.Equal(t, "", body)
	}
	if body, err := env.Client.get("/not-found"); assert.EqualError(t, err, "response status was 404") {
		assert.Equal(t, "not found", body)
	}
}

func TestClient_IsHcloudServer(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/hostname", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("my-server"))
	})

	isHcloudServer := env.Client.IsHcloudServer()

	assert.True(t, isHcloudServer)
}

func TestClient_IsHcloudServer_EmptyReturn(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/hostname", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	isHcloudServer := env.Client.IsHcloudServer()

	assert.False(t, isHcloudServer)
}

func TestClient_IsHcloudServer_HttpError(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/hostname", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	isHcloudServer := env.Client.IsHcloudServer()

	assert.False(t, isHcloudServer)
}

func TestClient_Hostname(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/hostname", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("my-server"))
	})

	hostname, err := env.Client.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if hostname != "my-server" {
		t.Fatalf("Unexpected hostname %s", hostname)
	}
}

func TestClient_InstanceID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env := newTestEnv()
		env.Mux.HandleFunc("/instance-id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("123456"))
		})

		instanceID, err := env.Client.InstanceID()
		require.NoError(t, err)
		require.Equal(t, int64(123456), instanceID)
	})

	t.Run("success sanitized", func(t *testing.T) {
		env := newTestEnv()
		env.Mux.HandleFunc("/instance-id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("123456\n"))
		})

		instanceID, err := env.Client.InstanceID()
		require.NoError(t, err)
		require.Equal(t, int64(123456), instanceID)
	})
}

func TestClient_PublicIPv4(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env := newTestEnv()
		env.Mux.HandleFunc("/public-ipv4", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("127.0.0.1"))
		})

		publicIPv4, err := env.Client.PublicIPv4()
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1", publicIPv4.String())
	})

	t.Run("success sanitized", func(t *testing.T) {
		env := newTestEnv()
		env.Mux.HandleFunc("/public-ipv4", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("127.0.0.1\n"))
		})

		publicIPv4, err := env.Client.PublicIPv4()
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1", publicIPv4.String())
	})
}

func TestClient_Region(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("eu-central"))
	})

	region, err := env.Client.Region()
	if err != nil {
		t.Fatal(err)
	}
	if region != "eu-central" {
		t.Fatalf("Unexpected region %s", region)
	}
}

func TestClient_AvailabilityZone(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/availability-zone", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fsn1-dc14"))
	})

	availabilityZone, err := env.Client.AvailabilityZone()
	if err != nil {
		t.Fatal(err)
	}
	if availabilityZone != "fsn1-dc14" {
		t.Fatalf("Unexpected availabilityZone %s", availabilityZone)
	}
}

func TestClient_PrivateNetworks(t *testing.T) {
	env := newTestEnv()
	env.Mux.HandleFunc("/private-networks", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`- ip: 10.0.0.2
  alias_ips: [10.0.0.3, 10.0.0.4]
  interface_num: 1
  mac_address: 86:00:00:2a:7d:e0
  network_id: 1234
  network_name: nw-test1
  network: 10.0.0.0/8
  subnet: 10.0.0.0/24
  gateway: 10.0.0.1
- ip: 192.168.0.2
  alias_ips: []
  interface_num: 2
  mac_address: 86:00:00:2a:7d:e1
  network_id: 4321
  network_name: nw-test2
  network: 192.168.0.0/16
  subnet: 192.168.0.0/24
  gateway: 192.168.0.1`))
	})

	privateNetworks, err := env.Client.PrivateNetworks()
	if err != nil {
		t.Fatal(err)
	}
	expectedNetworks := `- ip: 10.0.0.2
  alias_ips: [10.0.0.3, 10.0.0.4]
  interface_num: 1
  mac_address: 86:00:00:2a:7d:e0
  network_id: 1234
  network_name: nw-test1
  network: 10.0.0.0/8
  subnet: 10.0.0.0/24
  gateway: 10.0.0.1
- ip: 192.168.0.2
  alias_ips: []
  interface_num: 2
  mac_address: 86:00:00:2a:7d:e1
  network_id: 4321
  network_name: nw-test2
  network: 192.168.0.0/16
  subnet: 192.168.0.0/24
  gateway: 192.168.0.1`
	assert.Equal(t, expectedNetworks, privateNetworks)
}
