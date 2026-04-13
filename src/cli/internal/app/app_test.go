package app

import (
	"strings"
	"testing"

	"github.com/ShiunduZachariah/azopscli/internal/config"
)

func TestDetermineCredentialStrategy(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		want    string
		wantErr string
	}{
		{
			name: "auto falls back to default chain",
			cfg: config.Config{
				AuthMode: "auto",
			},
			want: authModeAuto,
		},
		{
			name: "auto selects service principal when credentials exist",
			cfg: config.Config{
				AuthMode:     "auto",
				TenantID:     "tenant-123",
				ClientID:     "client-123",
				ClientSecret: "secret-123",
			},
			want: authModeServicePrincipal,
		},
		{
			name: "service principal requires full credentials",
			cfg: config.Config{
				AuthMode: "service-principal",
				TenantID: "tenant-123",
			},
			wantErr: "service principal auth requires tenant id, client id, and client secret",
		},
		{
			name: "service principal accepts full credentials",
			cfg: config.Config{
				AuthMode:     "service-principal",
				TenantID:     "tenant-123",
				ClientID:     "client-123",
				ClientSecret: "secret-123",
			},
			want: authModeServicePrincipal,
		},
		{
			name: "unsupported auth mode errors",
			cfg: config.Config{
				AuthMode: "mystery-mode",
			},
			wantErr: "unsupported auth mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.Normalize()

			got, err := determineCredentialStrategy(tt.cfg)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tt.wantErr)
				}
				if got != "" {
					t.Fatalf("expected empty strategy when error occurs, got %q", got)
				}
				if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("determine strategy: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected strategy: got %q want %q", got, tt.want)
			}
		})
	}
}
