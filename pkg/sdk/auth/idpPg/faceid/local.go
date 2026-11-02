package faceid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LocalUniFaceProvider 本地 UniFace 人脸对比提供商
type LocalUniFaceProvider struct {
	Endpoint string
	ApiKey   string
	Client   *http.Client
}

type localUniFaceCompareRequest struct {
	ImageA string `json:"imageA"`
	ImageB string `json:"imageB"`
}

type localUniFaceCompareResponse struct {
	Matched bool    `json:"matched"`
	Score   float64 `json:"score"`
	Reason  string  `json:"reason"`
	Detail  string  `json:"detail"`
}

// NewLocalUniFaceProvider 创建本地 UniFace 提供商
func NewLocalUniFaceProvider(endpoint string, apiKey string) *LocalUniFaceProvider {
	return &LocalUniFaceProvider{
		Endpoint: strings.TrimRight(endpoint, "/"),
		ApiKey:   apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Check 对比两张人脸图片
func (p *LocalUniFaceProvider) Check(base64ImageA string, base64ImageB string) (bool, error) {
	if p.Endpoint == "" {
		return false, fmt.Errorf("Local UniFace endpoint 未配置")
	}

	body, err := json.Marshal(localUniFaceCompareRequest{
		ImageA: base64ImageA,
		ImageB: base64ImageB,
	})
	if err != nil {
		return false, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/compare", p.Endpoint), bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	request.Header.Set("Content-Type", "application/json")
	if p.ApiKey != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.ApiKey))
	}

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return false, fmt.Errorf("Local UniFace compare failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	var compareResponse localUniFaceCompareResponse
	if err = json.Unmarshal(responseBody, &compareResponse); err != nil {
		return false, err
	}

	return compareResponse.Matched, nil
}
