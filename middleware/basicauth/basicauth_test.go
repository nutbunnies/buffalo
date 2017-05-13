package basicauth_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/basicauth"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, nil)
	}
	auth := func(c buffalo.Context, u, p string) bool {
		return u == "tester" && p == "pass123"
	}
	a := buffalo.Automatic(buffalo.Options{})
	a.GET("/", basicauth.BasicAuth(auth)(h))
	return a
}

func TestBasicAuth(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())

	authfail := "invalid basic auth"

	// missing authorization
	res := w.Request("/").Get()
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), authfail)

	// bad cred tokens
	req := w.Request("/")
	req.Headers["Authorization"] = "badcreds"
	res = req.Get()
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), authfail)

	// bad cred values
	req = w.Request("/")
	req.Headers["Authorization"] = "badcreds:badpass"
	res = req.Get()
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), authfail)

	creds := base64.StdEncoding.EncodeToString([]byte("tester:pass123"))

	// valid cred values
	req = w.Request("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(200, res.Code)
}