package cloudflare

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockServerConfig struct {
	path       string
	method     string
	statusCode int
	response   interface{}
	assert     func(r *http.Request)
}

func newMockServer(configs ...mockServerConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, config := range configs {
			if r.URL.Path == config.path && r.Method == config.method {
				if config.assert != nil {
					config.assert(r)
				}

				w.WriteHeader(config.statusCode)
				if config.response != nil {
					_ = json.NewEncoder(w).Encode(config.response)
				}
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestNewHttpCloudflareClient(t *testing.T) {
	type args struct {
		client  *http.Client
		apiKey  string
		baseUrl string
	}
	tests := []struct {
		name string
		args args
		want CloudflareClient
	}{
		{
			name: "success",
			args: args{
				client:  &http.Client{},
				apiKey:  "api-key",
				baseUrl: "http://localhost",
			},
			want: &httpCloudflareClient{
				client:  &http.Client{},
				apiKey:  "api-key",
				baseUrl: "http://localhost",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHttpCloudflareClient(tt.args.client, tt.args.apiKey, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpCloudflareClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpCloudflareClient_AddZoneRecord(t *testing.T) {
	type fields struct {
		server func() *httptest.Server
		apiKey string
	}

	type args struct {
		ctx     context.Context
		request AddZoneRecordRequest
	}
	defaultArgs := args{
		ctx: context.Background(),
		request: AddZoneRecordRequest{
			ZoneID: "zone-id",
			Record: ZoneRecordRequest{
				Type:    "A",
				Name:    "app.example.com",
				Content: "1.1.1.1",
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AddZoneRecordResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an unauthorized error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records",
						method:     http.MethodPost,
						statusCode: http.StatusForbidden,
						response:   nil,
						assert: func(r *http.Request) {
							assert.Equal(t, r.Header.Get("Authorization"), "Bearer api-key")
						},
					})
				},
				apiKey: "api-key",
			},
			args: defaultArgs,
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRecordAddFailed)
			},
		},
		{
			name: "it should return a not found error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records",
						method:     http.MethodPost,
						statusCode: http.StatusNotFound,
						response:   nil,
					})
				},
			},
			args: defaultArgs,
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrZoneNotFound)
			},
		},
		{
			name: "it should return an unexpected status code error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records",
						method:     http.MethodPost,
						statusCode: http.StatusInternalServerError,
						response:   nil,
					})
				},
			},
			args: defaultArgs,
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRecordAddFailed)
			},
		},
		{
			name: "it should add a record",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/zone-id/dns_records",
						method: http.MethodPost,
						response: map[string]interface{}{
							"result": map[string]interface{}{
								"id": "record-id",
							},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: defaultArgs,
			want: &AddZoneRecordResponse{
				Record: ZoneRecord{
					ID: "record-id",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.fields.server()
			h := httpCloudflareClient{
				client:  server.Client(),
				apiKey:  tt.fields.apiKey,
				baseUrl: server.URL,
			}
			got, err := h.AddZoneRecord(tt.args.ctx, tt.args.request)
			if tt.wantErr(t, err) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddZoneRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpCloudflareClient_DeleteZoneRecord(t *testing.T) {
	type fields struct {
		server func() *httptest.Server
		apiKey string
	}
	type args struct {
		ctx     context.Context
		request DeleteZoneRecordRequest
	}
	defaultArgs := args{
		ctx: context.Background(),
		request: DeleteZoneRecordRequest{
			ZoneID:   "zone-id",
			RecordID: "record-id",
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a forbidden error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/zone-id/dns_records/record-id",
						method: http.MethodDelete,
						assert: func(r *http.Request) {
							assert.Equal(t, r.Header.Get("Authorization"), "Bearer api-key")
						},
						statusCode: http.StatusForbidden,
					})
				},
				apiKey: "api-key",
			},
			args: defaultArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRecordDeleteFailed)
			},
		},
		{
			name: "it should return a not found error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records/record-id",
						method:     http.MethodDelete,
						statusCode: http.StatusNotFound,
					})
				},
			},
			args: defaultArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRecordNotFound)
			},
		},
		{
			name: "it should return an unexpected status code error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records/record-id",
						method:     http.MethodDelete,
						statusCode: http.StatusInternalServerError,
					})
				},
			},
			args: defaultArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRecordDeleteFailed)
			},
		},
		{
			name: "it should delete the record",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/zone-id/dns_records/record-id",
						method:     http.MethodDelete,
						statusCode: http.StatusOK,
					})
				},
			},
			args:    defaultArgs,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.fields.server()
			h := httpCloudflareClient{
				client:  server.Client(),
				apiKey:  tt.fields.apiKey,
				baseUrl: server.URL,
			}
			if err := h.DeleteZoneRecord(tt.args.ctx, tt.args.request); tt.wantErr(t, err) {
				return
			}
		})
	}
}

func Test_httpCloudflareClient_GetZoneByDomain(t *testing.T) {
	type fields struct {
		server func() *httptest.Server
		apiKey string
	}
	type args struct {
		ctx     context.Context
		request GetZoneByDomainRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetZoneByDomainResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a forbidden error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones",
						method: http.MethodGet,
						assert: func(r *http.Request) {
							assert.Equal(t, r.Header.Get("Authorization"), "Bearer api-key")
							assert.Equal(t, r.URL.Query().Get("name"), "example.com")
						},
						statusCode: http.StatusForbidden,
						response:   nil,
					})
				},
				apiKey: "api-key",
			},
			args: args{
				ctx: context.Background(),
				request: GetZoneByDomainRequest{
					Domain: "example.com",
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrZoneListFailed)
			},
		},
		{
			name: "it should return a not found error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones",
						method: http.MethodGet,
						response: map[string]interface{}{
							"result": []map[string]interface{}{},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: args{
				ctx: context.Background(),
				request: GetZoneByDomainRequest{
					Domain: "example.com",
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrZoneNotFound)
			},
		},
		{
			name: "it should return a zone",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones",
						method: http.MethodGet,
						response: map[string]interface{}{
							"result": []map[string]interface{}{
								{
									"id":   "zone-id",
									"name": "example.com",
								},
							},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: args{
				ctx: context.Background(),
				request: GetZoneByDomainRequest{
					Domain: "example.com",
				},
			},
			want: &GetZoneByDomainResponse{
				ZoneID: "zone-id",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.fields.server()
			h := httpCloudflareClient{
				client:  server.Client(),
				apiKey:  tt.fields.apiKey,
				baseUrl: server.URL,
			}
			got, err := h.GetZoneByDomain(tt.args.ctx, tt.args.request)
			if tt.wantErr(t, err) {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZoneByDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpCloudflareClient_GetZoneRecords(t *testing.T) {
	type fields struct {
		server func() *httptest.Server
		apiKey string
	}
	type args struct {
		ctx     context.Context
		request GetZoneRecordsRequest
	}
	defaultArgs := args{
		ctx: context.Background(),
		request: GetZoneRecordsRequest{
			ZoneID: "mock-zone-id",
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetZoneRecordsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a forbidden error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/mock-zone-id/dns_records",
						method: http.MethodGet,
						assert: func(r *http.Request) {
							assert.Equal(t, r.Header.Get("Authorization"), "Bearer api-key")
							assert.Equal(t, r.URL.Query().Get("per_page"), "50")
						},
						statusCode: http.StatusForbidden,
					})
				},
				apiKey: "api-key",
			},
			args: defaultArgs,
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrZoneRecordsFailed)
			},
		},
		{
			name: "it should return a not found error",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:       "/client/v4/zones/mock-zone-id/dns_records",
						statusCode: http.StatusNotFound,
					})
				},
			},
			args: defaultArgs,
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrZoneNotFound)
			},
		},
		{
			name: "it should filter by name",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/mock-zone-id/dns_records",
						method: http.MethodGet,
						assert: func(r *http.Request) {
							assert.Equal(t, r.URL.Query().Get("name"), "app.example.com")
						},
						response: map[string]interface{}{
							"result": []map[string]interface{}{},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: func() args {
				args := defaultArgs
				args.request.Name = "app.example.com"
				return args
			}(),
			want: &GetZoneRecordsResponse{
				Records: []ZoneRecord{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should filter by type",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/mock-zone-id/dns_records",
						method: http.MethodGet,
						response: map[string]interface{}{
							"result": []map[string]interface{}{
								{
									"id":   "record-id",
									"type": "A",
								},
								{
									"id":   "record-id-2",
									"type": "CNAME",
								},
							},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: func() args {
				args := defaultArgs
				args.request.Type = "A"
				return args
			}(),
			want: &GetZoneRecordsResponse{
				Records: []ZoneRecord{
					{
						ID:   "record-id",
						Type: "A",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return all records",
			fields: fields{
				server: func() *httptest.Server {
					return newMockServer(mockServerConfig{
						path:   "/client/v4/zones/mock-zone-id/dns_records",
						method: http.MethodGet,
						response: map[string]interface{}{
							"result": []map[string]interface{}{
								{
									"id":   "record-id",
									"type": "A",
								},
								{
									"id":   "record-id-2",
									"type": "CNAME",
								},
							},
						},
						statusCode: http.StatusOK,
					})
				},
			},
			args: defaultArgs,
			want: &GetZoneRecordsResponse{
				Records: []ZoneRecord{
					{
						ID:   "record-id",
						Type: "A",
					},
					{
						ID:   "record-id-2",
						Type: "CNAME",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.fields.server()
			h := httpCloudflareClient{
				client:  server.Client(),
				apiKey:  tt.fields.apiKey,
				baseUrl: server.URL,
			}
			got, err := h.GetZoneRecords(tt.args.ctx, tt.args.request)
			if tt.wantErr(t, err) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZoneRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}
