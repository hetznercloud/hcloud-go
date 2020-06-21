package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestCertificateCreateOptsValidate(t *testing.T) {
	testCases := map[string]struct {
		Opts  CertificateCreateOpts
		Valid bool
	}{
		"empty": {
			Opts:  CertificateCreateOpts{},
			Valid: false,
		},
		"all set": {
			Opts: CertificateCreateOpts{
				Name:        "name",
				Certificate: "cert",
				PrivateKey:  "key",
				Labels:      map[string]string{},
			},
			Valid: true,
		},
		"no name": {
			Opts: CertificateCreateOpts{
				Certificate: "cert",
				PrivateKey:  "key",
				Labels:      map[string]string{},
			},
			Valid: false,
		},
		"no certificate": {
			Opts: CertificateCreateOpts{
				Name:       "name",
				PrivateKey: "key",
				Labels:     map[string]string{},
			},
			Valid: false,
		},
		"no private key": {
			Opts: CertificateCreateOpts{
				Name:        "name",
				Certificate: "cert",
				Labels:      map[string]string{},
			},
			Valid: false,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := testCase.Opts.Validate()
			if err == nil && !testCase.Valid || err != nil && testCase.Valid {
				t.FailNow()
			}
		})
	}
}

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
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.CertificateCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.CertificateCreateRequest{
			Name:        "my-cert",
			Certificate: "-----BEGIN CERTIFICATE-----\n...",
			PrivateKey:  "-----BEGIN PRIVATE KEY-----\n...",
			Labels: func() *map[string]string {
				labels := map[string]string{"key": "value"}
				return &labels
			}(),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.CertificateCreateResponse{
			Certificate: schema.Certificate{ID: 1},
		})
	})

	ctx := context.Background()

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
}

func TestCertificateCreateValidation(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()
	opts := CertificateCreateOpts{}
	_, _, err := env.Client.Certificate.Create(ctx, opts)
	if err == nil || err.Error() != "missing name" {
		t.Fatalf("unexpected error: %v", err)
	}
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
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/certificates/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("expected PUT")
		}
		var reqBody schema.CertificateUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.CertificateUpdateRequest{
			Name: String("test"),
			Labels: func() *map[string]string {
				labels := map[string]string{"key": "value"}
				return &labels
			}(),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.CertificateUpdateResponse{
			Certificate: schema.Certificate{
				ID: 1,
			},
		})
	})

	var (
		ctx        = context.Background()
		certficate = &Certificate{ID: 1}
	)

	opts := CertificateUpdateOpts{
		Name:   "test",
		Labels: map[string]string{"key": "value"},
	}
	updatedCertificate, _, err := env.Client.Certificate.Update(ctx, certficate, opts)
	if err != nil {
		t.Fatal(err)
	}
	if updatedCertificate.ID != 1 {
		t.Errorf("unexpected certficate ID: %v", updatedCertificate.ID)
	}
}
