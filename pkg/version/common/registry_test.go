package common

import (
	"errors"
	"testing"
)

// Mock registry client for testing
type mockRegistryClient struct {
	versions []string
	err      error
}

func (m *mockRegistryClient) GetVersions(packageName string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.versions, nil
}

func (m *mockRegistryClient) GetPackageInfo(packageName string) (PackageInfo, error) {
	return PackageInfo{
		Name:        packageName,
		Description: "Test package",
	}, nil
}

func TestGetLatestVersion(t *testing.T) {
	tests := []struct {
		name    string
		client  RegistryClient
		pkg     string
		want    string
		wantErr bool
	}{
		{
			name:    "successful retrieval",
			client:  &mockRegistryClient{versions: []string{"1.0.0", "1.1.0", "2.0.0"}},
			pkg:     "test-package",
			want:    "2.0.0",
			wantErr: false,
		},
		{
			name:    "no versions found",
			client:  &mockRegistryClient{versions: []string{}},
			pkg:     "test-package",
			wantErr: true,
		},
		{
			name:    "client error",
			client:  &mockRegistryClient{err: errors.New("network error")},
			pkg:     "test-package",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestVersion(tt.client, tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionHistory(t *testing.T) {
	client := &mockRegistryClient{versions: []string{"1.0.0", "1.1.0", "2.0.0", "2.1.0"}}

	tests := []struct {
		name    string
		pkg     string
		limit   int
		wantLen int
		wantErr bool
	}{
		{
			name:    "no limit",
			pkg:     "test-package",
			limit:   0,
			wantLen: 4,
			wantErr: false,
		},
		{
			name:    "with limit",
			pkg:     "test-package",
			limit:   2,
			wantLen: 2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVersionHistory(client, tt.pkg, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersionHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("GetVersionHistory() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{"equal versions", "1.0.0", "1.0.0", 0},
		{"v1 greater", "2.0.0", "1.0.0", 1},
		{"v1 lesser", "1.0.0", "2.0.0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("compareVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
