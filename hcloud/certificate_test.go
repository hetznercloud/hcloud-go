package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestCertificateClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.CertificateGetResponse{
			Certificate: schema.Certificate{
				ID: 1,
			},
		})
	})
	ctx := context.Background()

	certficate, _, err := env.Client.Certificate.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if certficate == nil {
		t.Fatal("no certficate")
	}
	if certficate.ID != 1 {
		t.Errorf("unexpected certficate ID: %v", certficate.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		certficate, _, err := env.Client.Certificate.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if certficate == nil {
			t.Fatal("no certficate")
		}
		if certficate.ID != 1 {
			t.Errorf("unexpected certficate ID: %v", certficate.ID)
		}
	})
}

func TestCertificateClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	certficate, _, err := env.Client.Certificate.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if certficate != nil {
		t.Fatal("expected no certficate")
	}
}

func TestCertificateClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mycert" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.CertificateListResponse{
			Certificates: []schema.Certificate{
				{
					ID:   1,
					Name: "mycert",
				},
			},
		})
	})
	ctx := context.Background()

	certficate, _, err := env.Client.Certificate.GetByName(ctx, "mycert")
	if err != nil {
		t.Fatal(err)
	}
	if certficate == nil {
		t.Fatal("no certficate")
	}
	if certficate.ID != 1 {
		t.Errorf("unexpected certficate ID: %v", certficate.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		certficate, _, err := env.Client.Certificate.Get(ctx, "mycert")
		if err != nil {
			t.Fatal(err)
		}
		if certficate == nil {
			t.Fatal("no certficate")
		}
		if certficate.ID != 1 {
			t.Errorf("unexpected certficate ID: %v", certficate.ID)
		}
	})
}

func TestCertificateClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mycert" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.CertificateListResponse{
			Certificates: []schema.Certificate{},
		})
	})

	ctx := context.Background()
	certficate, _, err := env.Client.Certificate.GetByName(ctx, "mycert")
	if err != nil {
		t.Fatal(err)
	}
	if certficate != nil {
		t.Fatal("unexpected certficate")
	}
}

func TestCertificateCreate(t *testing.T) {
	var (
		ctx = context.Background()
	)

	t.Run("missing required field name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := CertificateCreateOpts{}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err == nil || err.Error() != "missing name" {
			t.Fatalf("Certificate.Create should fail with \"missing name\" but failed with %s", err)
		}
	})

	t.Run("missing required field certificate", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := CertificateCreateOpts{
			Name: "my-cert",
		}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err == nil || err.Error() != "missing certificate" {
			t.Fatalf("Certificate.Create should fail with \"missing certificate\" but failed with %s", err)
		}
	})

	t.Run("missing required field certificate", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := CertificateCreateOpts{
			Name:        "my-cert",
			Certificate: "-----BEGIN CERTIFICATE-----\n...",
		}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err == nil || err.Error() != "missing private key" {
			t.Fatalf("Certificate.Create should fail with \"missing private key\" but failed with %s", err)
		}
	})

	t.Run("required fields", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.CertificateCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-cert" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.Certificate != "-----BEGIN CERTIFICATE-----\n..." {
				t.Errorf("unexpected Certificate: %v", reqBody.Certificate)
			}
			if reqBody.PrivateKey != "-----BEGIN PRIVATE KEY-----\n..." {
				t.Errorf("unexpected PrivateKey: %v", reqBody.PrivateKey)
			}
			json.NewEncoder(w).Encode(schema.CertificateCreateResponse{
				Certificate: schema.Certificate{ID: 1},
			})
		})
		opts := CertificateCreateOpts{
			Name:        "my-cert",
			Certificate: "-----BEGIN CERTIFICATE-----\n...",
			PrivateKey:  "-----BEGIN PRIVATE KEY-----\n...",
		}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("required fields with chain", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.CertificateCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-cert" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.Certificate != "-----BEGIN CERTIFICATE-----\n..." {
				t.Errorf("unexpected Certificate: %v", reqBody.Certificate)
			}
			if reqBody.Chain != "-----BEGIN CHAIN-----\n..." {
				t.Errorf("unexpected Chain: %v", reqBody.Chain)
			}
			if reqBody.PrivateKey != "-----BEGIN PRIVATE KEY-----\n..." {
				t.Errorf("unexpected PrivateKey: %v", reqBody.PrivateKey)
			}
			json.NewEncoder(w).Encode(schema.CertificateCreateResponse{
				Certificate: schema.Certificate{ID: 1},
			})
		})
		opts := CertificateCreateOpts{
			Name:        "my-cert",
			Certificate: "-----BEGIN CERTIFICATE-----\n...",
			Chain:       "-----BEGIN CHAIN-----\n...",
			PrivateKey:  "-----BEGIN PRIVATE KEY-----\n...",
		}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("required fields with labels", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.CertificateCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-cert" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.Certificate != "-----BEGIN CERTIFICATE-----\n..." {
				t.Errorf("unexpected Certificate: %v", reqBody.Certificate)
			}
			if reqBody.Labels == nil || reqBody.Labels["key"] != "value" {
				t.Errorf("no or wrong labels set: %v", reqBody.Labels)
			}
			if reqBody.PrivateKey != "-----BEGIN PRIVATE KEY-----\n..." {
				t.Errorf("unexpected PrivateKey: %v", reqBody.PrivateKey)
			}
			json.NewEncoder(w).Encode(schema.CertificateCreateResponse{
				Certificate: schema.Certificate{ID: 1},
			})
		})
		opts := CertificateCreateOpts{
			Name:        "my-cert",
			Certificate: "-----BEGIN CERTIFICATE-----\n...",
			PrivateKey:  "-----BEGIN PRIVATE KEY-----\n...",
			Labels:      map[string]string{"key": "value"},
		}
		_, _, err := env.Client.Certificate.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestCertificateDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	var (
		ctx        = context.Background()
		certficate = &Certificate{ID: 1}
	)
	_, err := env.Client.Certificate.Delete(ctx, certficate)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCertificateClientUpdate(t *testing.T) {
	var (
		ctx        = context.Background()
		certficate = &Certificate{ID: 1}
	)

	t.Run("update name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "test" {
				t.Errorf("unexpected name: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.CertificateUpdateResponse{
				Certificate: schema.Certificate{
					ID: 1,
				},
			})
		})

		opts := CertificateUpdateOpts{
			Name: "test",
		}
		updatedCertificate, _, err := env.Client.Certificate.Update(ctx, certficate, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedCertificate.ID != 1 {
			t.Errorf("unexpected certficate ID: %v", updatedCertificate.ID)
		}
	})

	t.Run("update labels", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Labels == nil || (*reqBody.Labels)["key"] != "value" {
				t.Errorf("unexpected labels in request: %v", reqBody.Labels)
			}
			json.NewEncoder(w).Encode(schema.CertificateUpdateResponse{
				Certificate: schema.Certificate{
					ID: 1,
				},
			})
		})

		opts := CertificateUpdateOpts{
			Labels: map[string]string{"key": "value"},
		}
		updatedCertificate, _, err := env.Client.Certificate.Update(ctx, certficate, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedCertificate.ID != 1 {
			t.Errorf("unexpected certficate ID: %v", updatedCertificate.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "" {
				t.Errorf("unexpected no name, but got: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.CertificateUpdateResponse{
				Certificate: schema.Certificate{
					ID: 1,
				},
			})
		})

		opts := CertificateUpdateOpts{}
		updatedCertificate, _, err := env.Client.Certificate.Update(ctx, certficate, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedCertificate.ID != 1 {
			t.Errorf("unexpected certficate ID: %v", updatedCertificate.ID)
		}
	})
}
