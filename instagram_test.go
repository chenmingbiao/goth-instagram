// instagram_test.go

package instagram

import (
	"os"
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := New("client_id", "secret", "http://localhost/callback")

	a.Equal(p.ClientKey, "client_id")
	a.Equal(p.Secret, "secret")
	a.Equal(p.CallbackURL, "http://localhost/callback")
	a.Equal(p.providerName, "instagram")
}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.Implements((*goth.Provider)(nil), New("client_id", "secret", "http://localhost/callback"))
}

func Test_BeginAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := New("client_id", "secret", "http://localhost/callback")
	session, err := p.BeginAuth("state")
	s := session.(*Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "api.instagram.com/oauth/authorize")
	a.Contains(s.AuthURL, "client_id=client_id")
	a.Contains(s.AuthURL, "scope=basic")
	a.Contains(s.AuthURL, "state=state")
}

func Test_FetchUser(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	// This test requires a valid Instagram access token
	// You might want to use a mock HTTP client for this test in a real scenario
	accessToken := os.Getenv("INSTAGRAM_TOKEN")
	if accessToken == "" {
		t.Skip("INSTAGRAM_TOKEN not set. Skipping test.")
	}

	p := New("client_id", "secret", "http://localhost/callback")
	session := &Session{AccessToken: accessToken}

	user, err := p.FetchUser(session)
	a.NoError(err)

	a.NotEmpty(user.Name)
	a.NotEmpty(user.UserID)
	a.Equal("instagram", user.Provider)
}

func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := New("client_id", "secret", "http://localhost/callback")
	session, err := p.UnmarshalSession(`{"AuthURL":"https://api.instagram.com/oauth/authorize","AccessToken":"1234567890"}`)
	a.NoError(err)

	s := session.(*Session)
	a.Equal(s.AuthURL, "https://api.instagram.com/oauth/authorize")
	a.Equal(s.AccessToken, "1234567890")
}
