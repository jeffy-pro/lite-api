package hotelbeds

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.nhat.io/clock"
	"lite-api/internal/client"
	liteapierrors "lite-api/internal/errors"
	"net/http"
	"time"
)

const (
	hotelsEndpoint       = "/hotel-api/1.0/hotels"
	headerXSignature     = "X-Signature"
	headerApiKey         = "Api-key"
	headerAccept         = "Accept"
	headerAcceptEncoding = "Accept-Encoding"
	headerContentType    = "Content-Type"
	applicationJSON      = "application/json"
	gzipAccept           = "gzip"
)

type HotelBeds struct {
	clock  clock.Clock
	cli    *http.Client
	apiKey string
	secret string
	host   string
}

func NewHotelBeds(host, apiKey, secret string, clock clock.Clock) *HotelBeds {
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	return &HotelBeds{
		cli:    &http.Client{Timeout: time.Second * 5},
		apiKey: apiKey,
		secret: secret,
		host:   host,
		clock:  clock,
	}
}

func (h *HotelBeds) Search(ctx context.Context, searchReq client.SearchRequest) (client.SearchResponse, error) {
	reqbody, err := json.Marshal(&searchReq)
	if err != nil {
		return client.SearchResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", h.host,
		hotelsEndpoint), bytes.NewBuffer(reqbody))

	if err != nil {
		return client.SearchResponse{}, err
	}

	req.Header.Add(headerApiKey, h.apiKey)
	req.Header.Add(headerAccept, applicationJSON)
	req.Header.Add(headerAcceptEncoding, gzipAccept)
	req.Header.Add(headerContentType, applicationJSON)

	signature := h.sign()
	req.Header.Add(headerXSignature, signature)

	resp, err := h.cli.Do(req)
	if err != nil {
		return client.SearchResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return client.SearchResponse{}, h.handleErrors(resp)
	}

	var searchResp client.SearchResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&searchResp)
	if err != nil {
		return client.SearchResponse{}, err
	}

	return searchResp, nil
}

func (h *HotelBeds) sign() string {
	// Begin Signature Assembly
	assemble := fmt.Sprintf("%s%s%d", h.apiKey, h.secret, h.clock.Now().Unix())

	// Begin SHA-256 Encryption
	hash := sha256.New()
	hash.Write([]byte(assemble))
	return hex.EncodeToString(hash.Sum(nil))
}

func (h *HotelBeds) handleErrors(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusTooManyRequests:
		simpleErr := client.SimpleError{}
		err := decoder.Decode(&simpleErr)
		if err != nil {
			return fmt.Errorf("error decoding error message: %w", err)
		}

		return liteapierrors.NewAPIErr(http.StatusText(resp.StatusCode), simpleErr.Error)
	case http.StatusPaymentRequired, http.StatusNotAcceptable, http.StatusConflict, http.StatusGone,
		http.StatusUnsupportedMediaType, http.StatusBadGateway, http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:

		return liteapierrors.NewAPIErr(http.StatusText(resp.StatusCode), "internal server error")
	default:
		searchResp := client.SearchResponse{}
		err := decoder.Decode(&searchResp)
		if err != nil {
			return fmt.Errorf("error decoding error message: %w", err)
		}

		return liteapierrors.NewAPIErr(searchResp.Error.Code, searchResp.Error.Message)
	}
}
