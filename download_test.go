package main

import (
	"reflect"
	"testing"
)

func Test_listAssets(t *testing.T) {
	type args struct {
		release *Release
		archive string
		osType  string
		arch    string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single match",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_darwin_amd64.tar.gz"},
						{Name: "app_v1.0.0_linux_amd64.tar.gz"},
					},
				},
				archive: "tar.gz",
				osType:  "darwin",
				arch:    "amd64",
			},
			want: []string{"app_v1.0.0_darwin_amd64.tar.gz"},
		},
		{
			name: "multiple matches",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_darwin_amd64.tar.gz"},
						{Name: "app_v1.0.0_darwin_amd64_checksums.txt"},
					},
				},
				archive: "tar.gz",
				osType:  "darwin",
				arch:    "amd64",
			},
			want: []string{"app_v1.0.0_darwin_amd64.tar.gz"},
		},
		{
			name: "no match",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_linux_amd64.tar.gz"},
					},
				},
				archive: "zip",
				osType:  "darwin",
				arch:    "amd64",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listAssets(tt.args.release, tt.args.archive, tt.args.osType, tt.args.arch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectAsset(t *testing.T) {
	type args struct {
		release *Release
		archive string
		osType  string
		arch    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Asset
		wantErr bool
	}{
		{
			name: "exact match tar.gz",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_darwin_amd64.tar.gz"},
						{Name: "app_v1.0.0_linux_amd64.tar.gz"},
					},
				},
				archive: "tar.gz",
				osType:  "darwin",
				arch:    "amd64",
			},
			want:    &Asset{Name: "app_v1.0.0_darwin_amd64.tar.gz"},
			wantErr: false,
		},
		{
			name: "exact match zip",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_windows_amd64.zip"},
					},
				},
				archive: "zip",
				osType:  "windows",
				arch:    "amd64",
			},
			want:    &Asset{Name: "app_v1.0.0_windows_amd64.zip"},
			wantErr: false,
		},
		{
			name: "no archive type",
			args: args{
				release: &Release{},
				archive: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no match",
			args: args{
				release: &Release{
					Assets: []Asset{
						{Name: "app_v1.0.0_linux_amd64.tar.gz"},
					},
				},
				archive: "zip",
				osType:  "darwin",
				arch:    "amd64",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectAsset(tt.args.release, tt.args.archive, tt.args.osType, tt.args.arch)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}
