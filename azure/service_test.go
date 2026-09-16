package azure

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/context_key"
)

func TestGetApplicableAuthorizationDetails_AuthMethod(t *testing.T) {
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	const sub, tenant, client = "sub", "tenant", "client"

	tests := []struct {
		name   string
		values map[string]string
		want   string
	}{
		{"nothing set", map[string]string{}, "CLI"},
		{"subscription only", map[string]string{auth.SubscriptionID: sub}, "CLI"},
		{"tenant and client, no subscription", map[string]string{auth.TenantID: tenant, auth.ClientID: client}, "CLI"},
		{"secret without subscription", map[string]string{auth.TenantID: tenant, auth.ClientID: client, auth.ClientSecret: "s"}, "CLI"},
		{"secret", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, auth.ClientSecret: "s"}, "Environment"},
		{"managed identity", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client}, "Environment"},
		{"federated token", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, settingFederatedToken: "jwt"}, "OIDC"},
		{"federated token file", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, settingFederatedTokenFile: "/t"}, "OIDC"},
		{"federated token without subscription", map[string]string{auth.TenantID: tenant, auth.ClientID: client, settingFederatedToken: "jwt"}, "CLI"},
		{"federated token without tenant", map[string]string{auth.SubscriptionID: sub, auth.ClientID: client, settingFederatedToken: "jwt"}, "CLI"},
		{"secret wins over federated token", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, auth.ClientSecret: "s", settingFederatedToken: "jwt"}, "Environment"},
		{"certificate wins over federated token", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, auth.CertificatePath: "/c", settingFederatedToken: "jwt"}, "Environment"},
		{"username and password win over federated token", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, auth.Username: "u", auth.Password: "p", settingFederatedToken: "jwt"}, "Environment"},
		{"username without password does not block federated token", map[string]string{auth.SubscriptionID: sub, auth.TenantID: tenant, auth.ClientID: client, auth.Username: "u", settingFederatedToken: "jwt"}, "OIDC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := auth.EnvironmentSettings{Values: tt.values}
			got, _, err := getApplicableAuthorizationDetails(ctx, settings, "MANAGEMENT")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("auth method = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetApplicableAuthorizationDetails_Resource(t *testing.T) {
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	tests := []struct {
		name        string
		environment string
		audience    string
		want        string
		wantErr     bool
	}{
		{"management public", "", "MANAGEMENT", "https://management.azure.com/", false},
		{"default audience is management", "", "", "https://management.azure.com/", false},
		{"graph public", "", "GRAPH", "https://graph.windows.net/", false},
		{"vault public trims trailing slash", "", "VAULT", "https://vault.azure.net", false},
		{"management china", "AZURECHINACLOUD", "MANAGEMENT", "https://management.chinacloudapi.cn/", false},
		{"invalid environment", "NOTACLOUD", "MANAGEMENT", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := auth.EnvironmentSettings{Values: map[string]string{auth.EnvironmentName: tt.environment}}
			_, got, err := getApplicableAuthorizationDetails(ctx, settings, tt.audience)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("resource = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFederatedTokenGetter(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	if err := os.WriteFile(tokenFile, []byte("from-file"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("inline token", func(t *testing.T) {
		got, err := federatedTokenGetter("inline", "")(context.Background())
		if err != nil || got != "inline" {
			t.Fatalf("got %q, %v", got, err)
		}
	})

	t.Run("inline token wins over file", func(t *testing.T) {
		got, err := federatedTokenGetter("inline", tokenFile)(context.Background())
		if err != nil || got != "inline" {
			t.Fatalf("got %q, %v", got, err)
		}
	})

	t.Run("token file", func(t *testing.T) {
		got, err := federatedTokenGetter("", tokenFile)(context.Background())
		if err != nil || got != "from-file" {
			t.Fatalf("got %q, %v", got, err)
		}
	})

	t.Run("token file is re-read on each call", func(t *testing.T) {
		getter := federatedTokenGetter("", tokenFile)
		if _, err := getter(context.Background()); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(tokenFile, []byte("rotated"), 0o600); err != nil {
			t.Fatal(err)
		}
		got, err := getter(context.Background())
		if err != nil || got != "rotated" {
			t.Fatalf("got %q, %v", got, err)
		}
	})

	t.Run("missing token file", func(t *testing.T) {
		missing := filepath.Join(dir, "missing")
		_, err := federatedTokenGetter("", missing)(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), missing) {
			t.Errorf("error %q does not name the file", err)
		}
	})
}

type fakeTokenCredential struct {
	token  string
	err    error
	scopes []string
}

func (f *fakeTokenCredential) GetToken(_ context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	f.scopes = opts.Scopes
	if f.err != nil {
		return azcore.AccessToken{}, f.err
	}
	return azcore.AccessToken{Token: f.token}, nil
}

func TestAzcoreTokenAuthorizer(t *testing.T) {
	newRequest := func() *http.Request {
		req, err := http.NewRequest(http.MethodGet, "https://management.azure.com/subscriptions", nil)
		if err != nil {
			t.Fatal(err)
		}
		return req
	}

	t.Run("sets bearer header with the requested scope", func(t *testing.T) {
		cred := &fakeTokenCredential{token: "abc"}
		a := &azcoreTokenAuthorizer{cred: cred, scope: "https://management.azure.com//.default"}

		req, err := autorest.Prepare(newRequest(), a.WithAuthorization())
		if err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer abc" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer abc")
		}
		if len(cred.scopes) != 1 || cred.scopes[0] != a.scope {
			t.Errorf("scopes = %v, want [%s]", cred.scopes, a.scope)
		}
	})

	t.Run("propagates token error", func(t *testing.T) {
		cred := &fakeTokenCredential{err: errors.New("boom")}
		a := &azcoreTokenAuthorizer{cred: cred, scope: "s"}

		req, err := autorest.Prepare(newRequest(), a.WithAuthorization())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if req.Header.Get("Authorization") != "" {
			t.Error("Authorization header set despite token error")
		}
	})
}
