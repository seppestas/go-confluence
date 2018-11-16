package confluence

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Wiki struct {
	endPoint   *url.URL
	authMethod AuthMethod
	client     *http.Client
}

func NewWiki(location string, authMethod AuthMethod) (*Wiki, error) {
	u, err := url.ParseRequestURI(location)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	u.Path += "rest/api"

	wiki := new(Wiki)
	wiki.endPoint = u
	wiki.authMethod = authMethod

	wiki.client = &http.Client{}

	return wiki, nil
}

type AuthMethod interface {
	auth(req *http.Request)
}

type basicAuthCallback func() (username, password string)

func (cb basicAuthCallback) auth(req *http.Request) {
	username, password := cb()
	req.SetBasicAuth(username, password)
}

func BasicAuth(username, password string) AuthMethod {
	return basicAuthCallback(func() (string, string) { return username, password })
}

type tokenAuthCallback func() (tokenkey string)

func (cb tokenAuthCallback) auth(req *http.Request) {
	tokenkey := cb()
	c := &http.Cookie{
		Name:     "studio.crowd.tokenkey",
		Value:    tokenkey,
		Path:     "/",
		Domain:   "." + req.URL.Host,
		Secure:   true,
		HttpOnly: true,
	}
	req.AddCookie(c)
}

func TokenAuth(tokenkey string) AuthMethod {
	return tokenAuthCallback(func() string { return tokenkey })
}

func (w *Wiki) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json, */*")
	w.authMethod.auth(req)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	// always read body to give more information about cause of failure
	res, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err2 != nil {
		return nil, err2
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		return res, nil
	case http.StatusNoContent, http.StatusResetContent:
		return res, nil
	case http.StatusUnauthorized:
		return res, fmt.Errorf("Authentication failed.")
	case http.StatusServiceUnavailable:
		return res, fmt.Errorf("Service is not available (%s).", resp.Status)
	case http.StatusInternalServerError:
		return res, fmt.Errorf("Internal server error: %s", resp.Status)
	}

	return res, fmt.Errorf("Unknown response status %s", resp.Status)
}
