package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

const (
	authURL      = "https://api.instagram.com/oauth/authorize"
	tokenURL     = "https://api.instagram.com/oauth/access_token"
	userAPIURL   = "https://graph.instagram.com/me"
	providerName = "instagram"
)

// Provider is the implementation of `goth.Provider` for accessing Instagram.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// New creates a new Instagram provider, and sets up important connection details.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: providerName,
		HTTPClient:   &http.Client{},
	}

	if len(scopes) == 0 {
		scopes = []string{"user_profile"}
	}

	p.config = &oauth2.Config{
		ClientID:     clientKey,
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	return p
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to set the name of the provider (needed to satisfy the provider interface)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// Client returns an HTTP client using the provided token.
func (p *Provider) Client(token *oauth2.Token) *http.Client {
	return p.config.Client(oauth2.NoContext, token)
}

// BeginAuth asks Instagram for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	return &Session{
		AuthURL: url,
	}, nil
}

// FetchUser will go to Instagram and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken: sess.AccessToken,
		Provider:    p.Name(),
	}

	if user.AccessToken == "" {
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	response, err := p.Client(sess.Token()).Get(userAPIURL + "?fields=id,username")
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("%s responded with a %d trying to fetch user information", p.providerName, response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(bits, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Session stores data during the auth process with Instagram.
type Session struct {
	AuthURL     string
	AccessToken string
	Token       *oauth2.Token
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Instagram provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New("an AuthURL has not been set")
	}
	return s.AuthURL, nil
}

// Authorize the session with Instagram and return the access token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(oauth2.NoContext, params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.Token = token
	return token.AccessToken, nil
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}
	err := json.Unmarshal([]byte(data), s)
	return s, err
}
