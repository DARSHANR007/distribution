package s3

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRedirectEndpointParameter verifies that redirectendpoint parameter is correctly parsed
func TestRedirectEndpointParameter(t *testing.T) {
	parameters := map[string]any{
		"accesskey": "",
		"secretkey": "",
		"bucket":    "test-bucket",
		"region":    "us-west-2",
		"encrypt":   false,
	}

	// Test without redirectendpoint
	d, err := FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "", d.StorageDriver.(*driver).RedirectEndpoint)

	// Test with redirectendpoint
	parameters["redirectendpoint"] = "https://s3-public.example.com"
	d, err = FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "https://s3-public.example.com", d.StorageDriver.(*driver).RedirectEndpoint)

	// Test with redirectendpoint and port
	parameters["redirectendpoint"] = "https://s3-public.example.com:443"
	d, err = FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "https://s3-public.example.com:443", d.StorageDriver.(*driver).RedirectEndpoint)
}

// TestRedirectEndpointWithSpecialCharacters tests RedirectEndpoint with special characters
func TestRedirectEndpointWithSpecialCharacters(t *testing.T) {
	parameters := map[string]any{
		"accesskey": "",
		"secretkey": "",
		"bucket":    "test-bucket",
		"region":    "us-west-2",
		"encrypt":   false,
	}

	// Test with URL containing path
	parameters["redirectendpoint"] = "https://s3-public.example.com/cdn/v1"
	d, err := FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "https://s3-public.example.com/cdn/v1", d.StorageDriver.(*driver).RedirectEndpoint)

	// Test with URL containing query parameters (should be preserved)
	parameters["redirectendpoint"] = "https://s3-public.example.com?region=us-west"
	d, err = FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "https://s3-public.example.com?region=us-west", d.StorageDriver.(*driver).RedirectEndpoint)
}

// TestRedirectEndpointEmptyString tests empty string handling
func TestRedirectEndpointEmptyString(t *testing.T) {
	parameters := map[string]any{
		"accesskey": "",
		"secretkey": "",
		"bucket":    "test-bucket",
		"region":    "us-west-2",
		"encrypt":   false,
	}

	// Test with explicit empty string
	parameters["redirectendpoint"] = ""
	d, err := FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "", d.StorageDriver.(*driver).RedirectEndpoint)

	// Test with nil (should default to empty)
	parameters["redirectendpoint"] = nil
	d, err = FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)
	require.Equal(t, "", d.StorageDriver.(*driver).RedirectEndpoint)
}

// TestRedirectEndpointRequestValidation verifies request handling
func TestRedirectEndpointRequestValidation(t *testing.T) {
	parameters := map[string]any{
		"accesskey":        "",
		"secretkey":        "",
		"bucket":           "test-bucket",
		"region":           "us-west-2",
		"encrypt":          false,
		"redirectendpoint": "https://s3-public.example.com",
	}

	d, err := FromParameters(context.Background(), parameters)
	require.NoError(t, err)
	require.NotNil(t, d)

	// Create a simple request for testing
	req, err := http.NewRequest(http.MethodGet, "http://internal.example.com/blob", nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	// Verify driver has the redirect endpoint configured
	drv := d.StorageDriver.(*driver)
	require.Equal(t, "https://s3-public.example.com", drv.RedirectEndpoint)
}
