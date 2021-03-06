package rest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/fabric8-services/fabric8-auth/resource"

	"github.com/goadesign/goa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbsoluteURLOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	req := &goa.RequestData{
		Request: &http.Request{Host: "api.service.domain.org"},
	}
	// HTTP
	urlStr := AbsoluteURL(req, "/testpath")
	assert.Equal(t, "http://api.service.domain.org/testpath", urlStr)

	// HTTPS
	r, err := http.NewRequest("", "https://api.service.domain.org", nil)
	require.Nil(t, err)
	req = &goa.RequestData{
		Request: r,
	}
	urlStr = AbsoluteURL(req, "/testpath2")
	assert.Equal(t, "https://api.service.domain.org/testpath2", urlStr)
}

func TestAbsoluteURLOKWithProxyForward(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	req := &goa.RequestData{
		Request: &http.Request{Host: "api.service.domain.org"},
	}

	// HTTPS
	r, err := http.NewRequest("", "http://api.service.domain.org", nil)
	require.Nil(t, err)
	r.Header.Set("X-Forwarded-Proto", "https")
	req = &goa.RequestData{
		Request: r,
	}
	urlStr := AbsoluteURL(req, "/testpath2")
	assert.Equal(t, "https://api.service.domain.org/testpath2", urlStr)
}

func TestReplaceDomainPrefixOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	host, err := ReplaceDomainPrefix("api.service.domain.org", "sso")
	assert.Nil(t, err)
	assert.Equal(t, "sso.service.domain.org", host)
}

func TestReplaceDomainPrefixInTooShortHostFails(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	_, err := ReplaceDomainPrefix("org", "sso")
	assert.NotNil(t, err)
}

func TestValidateEmailSuccess(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	isValid, err := ValidateEmail("a@a.com")
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestValidateEmailFail(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	isValid, err := ValidateEmail("a.a@com")
	require.NoError(t, err)
	require.False(t, isValid)
}

func TestAddParamsSuccess(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	testMap := map[string]string{
		"param1": "a",
		"param2": "b",
		"param3": "https://www.redhat.com",
	}
	testHost := "openshift.io"
	generatedURLString, err := AddParams("https://"+testHost, testMap)
	require.NoError(t, err)

	generateURL, err := url.Parse(generatedURLString)
	require.NoError(t, err)

	assert.Equal(t, testHost, generateURL.Host)
	assert.Equal(t, "https", generateURL.Scheme)

	m, _ := url.ParseQuery(generateURL.RawQuery)
	for k, v := range testMap {
		assert.Equal(t, v, m[k][0])
	}
}

func TestAddParamSuccess(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	generatedURL, err := AddParam("https://openshift.io", "param1", "a")
	require.NoError(t, err)
	assert.Equal(t, "https://openshift.io?param1=a", generatedURL)

	generatedURL, err = AddParam(generatedURL, "param2", "abc")
	require.NoError(t, err)
	equalURLs(t, "https://openshift.io?param1=a&param2=abc", generatedURL)
}

func TestAddTrailingSlashToURLSuccess(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	assert.Equal(t, "https://openshift.io/", AddTrailingSlashToURL("https://openshift.io"))
	assert.Equal(t, "https://openshift.io/", AddTrailingSlashToURL("https://openshift.io/"))
	assert.Equal(t, "", AddTrailingSlashToURL(""))
}

// Can't use test.EqualURLs() because of cycle dependency
func equalURLs(t *testing.T, expected string, actual string) {
	expectedURL, err := url.Parse(expected)
	require.Nil(t, err)
	actualURL, err := url.Parse(actual)
	require.Nil(t, err)
	assert.Equal(t, expectedURL.Scheme, actualURL.Scheme)
	assert.Equal(t, expectedURL.Host, actualURL.Host)
	assert.Equal(t, len(expectedURL.Query()), len(actualURL.Query()))
	for name, value := range expectedURL.Query() {
		assert.Equal(t, value, actualURL.Query()[name])
	}
}
