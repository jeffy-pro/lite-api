package hotelbeds

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"lite-api/internal/client"
	liteapierrors "lite-api/internal/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.nhat.io/clock"
)

//go:embed testdata/hotelbeds_response.json
var hotelbedsResponse []byte

func TestHotelBeds_Search(t *testing.T) {
	staticClock := clock.Fix(time.Date(2024, 7, 12, 11, 4, 5, 0, time.UTC))
	apiKey, secret := "12345", "6789"

	t.Run("client error", func(t *testing.T) {
		hotelBedsCli := NewHotelBeds("0.0.0.0", apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.Error(t, err)
		require.Zero(t, res)
	})

	t.Run("invalid request", func(t *testing.T) {
		hotelBedsCli := NewHotelBeds("http://///invalid-url", apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.Error(t, err)
		require.Zero(t, res)
	})

	t.Run("invalid request", func(t *testing.T) {
		hotelBedsCli := NewHotelBeds("http://///invalid-url", apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.Error(t, err)
		require.Zero(t, res)
	})

	t.Run("test signature", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != hotelsEndpoint {
				t.Fail()
				return
			}

			in := fmt.Sprintf("%s%s%d", apiKey, secret, staticClock.Now().Unix())
			hash := sha256.New()
			hash.Write([]byte(in))
			out := hex.EncodeToString(hash.Sum(nil))
			reqHeader := r.Header.Get(headerXSignature)

			require.Equal(t, out, reqHeader)

			_, _ = fmt.Fprintln(w, `{"message": "hello, world"}`)
		}))
		defer mockServer.Close()
		hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.NoError(t, err)
		require.Zero(t, res)
	})

	t.Run("verify headers", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != hotelsEndpoint {
				t.Fail()
				return
			}

			in := fmt.Sprintf("%s%s%d", apiKey, secret, staticClock.Now().Unix())
			hash := sha256.New()
			hash.Write([]byte(in))
			out := hex.EncodeToString(hash.Sum(nil))
			reqHeader := r.Header.Get(headerXSignature)

			require.Equal(t, out, reqHeader)
			require.Equal(t, gzipEncoding, r.Header.Get(headerAcceptEncoding))
			require.Equal(t, applicationJSON, r.Header.Get(headerAccept))
			require.Equal(t, apiKey, r.Header.Get(headerApiKey))
			require.Equal(t, applicationJSON, r.Header.Get(headerContentType))

			_, _ = fmt.Fprintln(w, `{"message": "hello, world"}`)
		}))
		defer mockServer.Close()
		hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.NoError(t, err)
		require.Zero(t, res)
	})

	t.Run("error status codes", func(t *testing.T) {
		t.Run("no body", func(t *testing.T) {
			noBodyStatusCodes := []int{
				http.StatusPaymentRequired, http.StatusNotAcceptable, http.StatusConflict, http.StatusGone,
				http.StatusUnsupportedMediaType, http.StatusBadGateway, http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			}

			noBodyStatusTestHelper := func(t *testing.T, statusCode int) {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != hotelsEndpoint {
						t.Fail()
						return
					}
					w.WriteHeader(statusCode)
					_, _ = w.Write([]byte(`don't care'`))
				}))
				defer mockServer.Close()

				hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
				res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
				expectedErr := liteapierrors.NewAPIErr(http.StatusText(statusCode), "internal server error")
				require.ErrorContains(t, err, expectedErr.Error())
				require.Zero(t, res)
			}

			for _, statusCode := range noBodyStatusCodes {
				noBodyStatusTestHelper(t, statusCode)
			}
		})

		t.Run("simple error", func(t *testing.T) {
			simpleErrStatusCodes := []int{
				http.StatusUnauthorized, http.StatusForbidden, http.StatusTooManyRequests,
			}

			simpleErrStatusTestHelper := func(t *testing.T, statusCode int) {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != hotelsEndpoint {
						t.Fail()
						return
					}
					w.WriteHeader(statusCode)
					_, _ = w.Write([]byte(`{"error": "something went wrong"}`))
				}))
				defer mockServer.Close()

				hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
				res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
				expectedErr := liteapierrors.NewAPIErr(http.StatusText(statusCode), "something went wrong")
				require.ErrorContains(t, err, expectedErr.Error())
				require.Zero(t, res)
			}

			for _, statusCode := range simpleErrStatusCodes {
				simpleErrStatusTestHelper(t, statusCode)
			}
		})

		t.Run("larger error", func(t *testing.T) {
			simpleErrStatusCodes := []int{
				http.StatusBadRequest, http.StatusNotFound,
			}

			largerErrStatusTestHelper := func(t *testing.T, statusCode int) {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != hotelsEndpoint {
						t.Fail()
						return
					}
					w.WriteHeader(statusCode)
					_, _ = w.Write([]byte(`
								{
								"auditData": {},
								"error": {
									"code": "ERR CODE", 
									"message": "detailed message"
									}
								}`))
				}))
				defer mockServer.Close()

				hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
				res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
				expectedErr := liteapierrors.NewAPIErr("ERR CODE", "detailed message")
				require.ErrorContains(t, err, expectedErr.Error())
				require.Zero(t, res)
			}

			for _, statusCode := range simpleErrStatusCodes {
				largerErrStatusTestHelper(t, statusCode)
			}
		})
	})

	t.Run("error decoding response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != hotelsEndpoint {
				t.Fail()
				return
			}

			_, _ = fmt.Fprintln(w, `{`)
		}))
		defer mockServer.Close()
		hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.Error(t, err)
		require.Zero(t, res)
	})

	t.Run("client success", func(t *testing.T) {
		var upstreamResp client.SearchResponse
		require.NoError(t, json.Unmarshal(hotelbedsResponse, &upstreamResp))

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != hotelsEndpoint {
				t.Fail()
				return
			}

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Fatal("Expected http.ResponseWriter to be a http.Flusher")
			}

			_, _ = fmt.Fprintln(w, string(hotelbedsResponse))
			flusher.Flush()
		}))

		defer mockServer.Close()

		hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.NoError(t, err)
		require.Equal(t, upstreamResp, res)
	})

	t.Run("client success - gzip response", func(t *testing.T) {
		var upstreamResp client.SearchResponse
		require.NoError(t, json.Unmarshal(hotelbedsResponse, &upstreamResp))

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != hotelsEndpoint {
				t.Fail()
				return
			}

			gz := gzip.NewWriter(w)
			defer gz.Close()

			w.Header().Set(headerContentEncoding, gzipEncoding)
			w.Header().Set(headerContentType, applicationJSON)

			_, err := gz.Write(hotelbedsResponse)
			require.NoError(t, err)
		}))

		defer mockServer.Close()

		hotelBedsCli := NewHotelBeds(mockServer.URL, apiKey, secret, staticClock)
		res, err := hotelBedsCli.Search(context.Background(), client.SearchRequest{})
		require.NoError(t, err)
		require.Equal(t, upstreamResp, res)
	})
}
