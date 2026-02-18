package goytmusic

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://music.youtube.com/youtube/v1/"
	defaultUserAgent = "ytmusic/0.0.0"
	defaultOrigin    = "https://music.youtube.com"
)

// A Client abstracts the communication with Innertube API
type Client struct {
	// Base service for all other services to share
	common service

	// Http client to communicate with the API
	client *http.Client
	// Base url for making requests.
	baseURL *url.URL
}

type service struct {
	client *Client
}

// NewClient creats a client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient

	transport := httpClient2.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Ensures any request made by any client has User-Agent and X-Origin defined
	httpClient2.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req = req.Clone(req.Context())
		req.Header.Set("User-Agent", defaultUserAgent)
		req.Header.Set("X-Origin", defaultOrigin)

		return transport.RoundTrip(req)
	})

	c := &Client{client: &httpClient2}

	c.baseURL, _ = url.Parse(defaultBaseURL)

	c.common.client = c

	return c
}

// WithAuthCookie returns a copy of client configured with authentication cookie
func (c *Client) WithAuthCookie(cookie string) *Client {
	sapisid := sapisidFromCookie(cookie)

	c2 := *c
	*c2.client = *c.client
	transport := c2.client.Transport

	c2.client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())

			req.Header.Set("Cookie", cookie)

			if sapisid != "" {
				h := sapisidHash(sapisid)
				req.Header.Set("Authorization", fmt.Sprintf("SAPISIDHASH %s", h))
			}

			return transport.RoundTrip(req) // receiver client transport
		},
	)
	return &c2
}

// create SAPISIDHASH from the sapisid extracted from cookie.
// The Innertube API expects this SAPISIDHASH as auth identity
func sapisidHash(sapisid string) string {
	now := time.Now().Unix()

	concat := fmt.Sprintf("%d %s %s", now, sapisid, defaultOrigin)

	h := sha1.New()
	h.Write([]byte(concat))
	sha1Hash := fmt.Sprintf("%x", h.Sum(nil))

	return fmt.Sprintf("%d_%s", now, sha1Hash)
}

// Extract 3PAPISID from cookie
func sapisidFromCookie(cookie string) string {
	parts := strings.SplitSeq(cookie, ";")
	for part := range parts {
		part = strings.TrimSpace(part)
		if s, found := strings.CutPrefix(part, "__Secure-3PAPISID="); found {
			return s
		}
	}
	return ""
}

// roundTripperFunc creates a RoundTripper (transport).
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
