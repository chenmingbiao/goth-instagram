// session_test.go

package instagram

import (
	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Implements_Session(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	s := &Session{}
	a.Implements((*goth.Session)(nil), s)
}

func Test_GetAuthURL(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	s := &Session{}

	_, err := s.GetAuthURL()
	a.Error(err)

	s.AuthURL = "/foo"

	url, _ := s.GetAuthURL()
	a.Equal(url, "/foo")
}

func Test_ToJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	s := &Session{}

	data := s.Marshal()
	a.Equal(data, `{"AuthURL":"","AccessToken":"","RefreshToken":"","ExpiresAt":"0001-01-01T00:00:00Z"}`)
}

func Test_String(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	s := &Session{}

	a.Equal(s.String(), s.Marshal())
}

func Test_UnmarshalSession(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := Provider{}
	s, err := p.UnmarshalSession(`{"AuthURL":"https://instagram.com/auth","AccessToken":"1234567890"}`)
	a.NoError(err)
	session := s.(*Session)
	a.Equal(session.AuthURL, "https://instagram.com/auth")
	a.Equal(session.AccessToken, "1234567890")
}

func Test_Authorize(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := New("cid", "csecret", "http://localhost/callback")
	s := &Session{}

	// This test would typically use a mock OAuth2 config
	// For simplicity, we're just testing the error case here
	_, err := s.Authorize(p, nil)
	a.Error(err)

	// In a real test, you'd mock the OAuth2 config and test both success and failure cases
}
