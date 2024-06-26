package config

import (
	mock_config "github.com/jorgejr568/cloudflare-cli/internal/config/mocks"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestNewComposite(t *testing.T) {

	tests := []struct {
		name    string
		want    Composite
		wantErr bool
	}{
		{
			name: "success",
			want: &composite{
				entries: []Composite{
					func() Composite {
						c, err := LoadLocalConfig()
						if err != nil {
							t.Errorf("NewComposite() error = %v, wantErr %v", err, false)
						}
						return c
					}(),
					envConfig{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewComposite()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewComposite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComposite() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_composite_CloudflareAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		entries []Composite
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return empty if it has no sources",
			fields: fields{
				entries: []Composite{},
			},
			want: "",
		},
		{
			name: "should return empty if all sources return empty",
			fields: fields{
				entries: []Composite{
					func() Composite {
						m := mock_config.NewMockComposite(ctrl)
						m.EXPECT().CloudflareAPIKey().
							Return("").
							Times(1)
						return m
					}(),
					func() Composite {
						m := mock_config.NewMockComposite(ctrl)
						m.EXPECT().CloudflareAPIKey().
							Return("").
							Times(1)
						return m
					}(),
				},
			},
			want: "",
		},
		{
			name: "should return the first non-empty value",
			fields: fields{
				entries: []Composite{
					func() Composite {
						m := mock_config.NewMockComposite(ctrl)
						m.EXPECT().CloudflareAPIKey().
							Return("").
							Times(1)
						return m
					}(),
					func() Composite {
						m := mock_config.NewMockComposite(ctrl)
						m.EXPECT().CloudflareAPIKey().
							Return("foo").
							Times(1)
						return m
					}(),
					func() Composite {
						m := mock_config.NewMockComposite(ctrl)
						m.EXPECT().CloudflareAPIKey().
							Return("bar").
							Times(0)
						return m
					}(),
				},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &composite{
				entries: tt.fields.entries,
			}
			if got := c.CloudflareAPIKey(); got != tt.want {
				t.Errorf("CloudflareAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
