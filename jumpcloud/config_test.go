package jumpcloud

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// withStubTokenEndpoint points tokenEndpoint at the given test server for the
// duration of the test and restores it afterwards.
func withStubTokenEndpoint(t *testing.T, srv *httptest.Server) {
	t.Helper()
	prev := tokenEndpoint
	tokenEndpoint = srv.URL
	t.Cleanup(func() { tokenEndpoint = prev })
}

func TestFetchAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, secret, ok := r.BasicAuth()
		if !ok || id != "the-id" || secret != "the-secret" {
			t.Errorf("expected basic auth the-id/the-secret, got id=%q secret=%q ok=%v", id, secret, ok)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %s", err)
		}
		if got := r.Form.Get("grant_type"); got != "client_credentials" {
			t.Errorf("grant_type = %q, want client_credentials", got)
		}
		if got := r.Form.Get("scope"); got != "api" {
			t.Errorf("scope = %q, want api", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"abc123","token_type":"Bearer","expires_in":3600}`))
	}))
	defer srv.Close()
	withStubTokenEndpoint(t, srv)

	token, err := fetchAccessToken("the-id", "the-secret")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if token != "abc123" {
		t.Fatalf("token = %q, want abc123", token)
	}
}

func TestFetchAccessTokenNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
	}))
	defer srv.Close()
	withStubTokenEndpoint(t, srv)

	if _, err := fetchAccessToken("x", "y"); err == nil {
		t.Fatal("expected an error for a 401 response, got nil")
	}
}

func TestFetchAccessTokenEmptyToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token_type":"Bearer"}`))
	}))
	defer srv.Close()
	withStubTokenEndpoint(t, srv)

	if _, err := fetchAccessToken("x", "y"); err == nil {
		t.Fatal("expected an error for an empty access_token, got nil")
	}
}

func TestClientOAuthSetsBearerHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok"}`))
	}))
	defer srv.Close()
	withStubTokenEndpoint(t, srv)

	c := &Config{ClientID: "id", ClientSecret: "secret", OrgID: "org"}
	out, err := c.Client()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	cfg := out.(*jcapiv2.Configuration)
	if got := cfg.DefaultHeader["Authorization"]; got != "Bearer tok" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer tok")
	}
	if _, ok := cfg.DefaultHeader["x-api-key"]; ok {
		t.Error("x-api-key header should not be set in OAuth mode")
	}
	if got := cfg.DefaultHeader["x-org-id"]; got != "org" {
		t.Errorf("x-org-id = %q, want org", got)
	}
}

func TestClientAPIKeyFallback(t *testing.T) {
	c := &Config{APIKey: "legacy-key"}
	out, err := c.Client()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	cfg := out.(*jcapiv2.Configuration)
	if got := cfg.DefaultHeader["x-api-key"]; got != "legacy-key" {
		t.Errorf("x-api-key = %q, want legacy-key", got)
	}
	if _, ok := cfg.DefaultHeader["Authorization"]; ok {
		t.Error("Authorization header should not be set in api_key mode")
	}
}

func TestConvertV2toV1ConfigCarriesBearer(t *testing.T) {
	v2 := jcapiv2.NewConfiguration()
	v2.AddDefaultHeader("Authorization", "Bearer tok")
	v2.AddDefaultHeader("x-org-id", "org")

	v1 := convertV2toV1Config(v2)
	if got := v1.DefaultHeader["Authorization"]; got != "Bearer tok" {
		t.Errorf("v1 Authorization = %q, want %q", got, "Bearer tok")
	}
	if got := v1.DefaultHeader["x-org-id"]; got != "org" {
		t.Errorf("v1 x-org-id = %q, want org", got)
	}
	if _, ok := v1.DefaultHeader["x-api-key"]; ok {
		t.Error("v1 x-api-key should not be set when only Authorization is present")
	}
}

func TestConvertV2toV1ConfigCarriesAPIKey(t *testing.T) {
	v2 := jcapiv2.NewConfiguration()
	v2.AddDefaultHeader("x-api-key", "legacy-key")

	v1 := convertV2toV1Config(v2)
	if got := v1.DefaultHeader["x-api-key"]; got != "legacy-key" {
		t.Errorf("v1 x-api-key = %q, want legacy-key", got)
	}
	if _, ok := v1.DefaultHeader["Authorization"]; ok {
		t.Error("v1 Authorization should not be set when only x-api-key is present")
	}
}
