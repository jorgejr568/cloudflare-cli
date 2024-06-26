package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jorgejr568/cloudflare-cli/internal/utils"
	"io"
	"net/http"
)

type httpCloudflareClient struct {
	client  *http.Client
	apiKey  string
	baseUrl string
}

func (h httpCloudflareClient) acquireRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.apiKey))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (h httpCloudflareClient) acquireResponseError(resp *http.Response, wrap error) error {
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: UNAUTHORIZED. You might need to check your API key", wrap)
	}

	errorMessage, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %s", wrap, err.Error())
	}

	return fmt.Errorf("%w: unexpected status code: %d - %s", wrap, resp.StatusCode, string(errorMessage))
}

func (h httpCloudflareClient) GetZoneByDomain(ctx context.Context, request GetZoneByDomainRequest) (*GetZoneByDomainResponse, error) {
	requestUrl := fmt.Sprintf("%s/client/v4/zones", h.baseUrl)
	req, err := h.acquireRequest(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrZoneListFailed, err.Error())
	}

	q := req.URL.Query()
	q.Add("name", request.Domain)
	req.URL.RawQuery = q.Encode()

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrZoneListFailed, err.Error())
	}

	defer utils.LogErrorIfError(resp.Body.Close())
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrZoneNotFound
		}

		return nil, h.acquireResponseError(resp, ErrZoneListFailed)
	}

	type listZonesResponse struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
	}
	var response listZonesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	for _, zone := range response.Result {
		if zone.Name == request.Domain {
			return &GetZoneByDomainResponse{ZoneID: zone.ID}, nil
		}
	}

	return nil, ErrZoneNotFound
}

func (h httpCloudflareClient) GetZoneRecords(ctx context.Context, request GetZoneRecordsRequest) (*GetZoneRecordsResponse, error) {
	requestUrl := fmt.Sprintf("%s/client/v4/zones/%s/dns_records", h.baseUrl, request.ZoneID)
	req, err := h.acquireRequest(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrZoneRecordsFailed, err.Error())
	}
	q := req.URL.Query()

	q.Add("per_page", "50")
	if request.Name != "" {
		q.Add("name", request.Name)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrZoneRecordsFailed, err.Error())
	}

	defer utils.LogErrorIfError(resp.Body.Close())
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrZoneNotFound
		}

		return nil, h.acquireResponseError(resp, ErrZoneRecordsFailed)
	}

	var response GetZoneRecordsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if request.Type != "" {
		var filteredRecords []ZoneRecord
		for _, record := range response.Records {
			if record.Type == request.Type {
				filteredRecords = append(filteredRecords, record)
			}
		}

		response.Records = filteredRecords
	}

	return &response, nil
}

func (h httpCloudflareClient) AddZoneRecord(ctx context.Context, request AddZoneRecordRequest) (*AddZoneRecordResponse, error) {
	requestUrl := fmt.Sprintf("%s/client/v4/zones/%s/dns_records", h.baseUrl, request.ZoneID)
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(request.Record)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRecordAddFailed, err.Error())
	}

	req, err := h.acquireRequest(ctx, http.MethodPost, requestUrl, &body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRecordAddFailed, err.Error())
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRecordAddFailed, err.Error())
	}

	defer utils.LogErrorIfError(resp.Body.Close())
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrZoneNotFound
		}

		return nil, h.acquireResponseError(resp, ErrRecordAddFailed)
	}

	var response AddZoneRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (h httpCloudflareClient) DeleteZoneRecord(ctx context.Context, request DeleteZoneRecordRequest) error {
	requestUrl := fmt.Sprintf("%s/client/v4/zones/%s/dns_records/%s", h.baseUrl, request.ZoneID, request.RecordID)
	req, err := h.acquireRequest(ctx, http.MethodDelete, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRecordDeleteFailed, err.Error())
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRecordDeleteFailed, err.Error())
	}

	defer utils.LogErrorIfError(resp.Body.Close())
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return ErrRecordNotFound
		}

		return fmt.Errorf("%w: unexpected status code: %d", ErrRecordDeleteFailed, resp.StatusCode)
	}

	return nil
}

func NewHttpCloudflareClient(client *http.Client, apiKey, baseUrl string) CloudflareClient {
	return &httpCloudflareClient{
		client:  client,
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}
