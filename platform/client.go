package platform

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// SessionCookieToken mattermost cookie field
	SessionCookieToken = "MMAUTHTOKEN"
	// HeaderETAGServer Header Field for ETag
	HeaderETAGServer = "ETag"
	// HeaderETAGClient Header Field for ETag
	HeaderETAGClient = "If-None-Match"
	// HeaderToken Header Field for Token
	HeaderToken = "token"
	// HeaderAuth Header Field for Auth
	HeaderAuth = "Authorization"
	// HeaderBearer prefix for auth field
	HeaderBearer = "BEARER"
	// HeaderRequestID Header Field for Request ID
	HeaderRequestID = "X-Request-ID"
	// APIURLSuffixV3 holds api suffice for api
	APIURLSuffixV3 = "/api/v3"
	// APIURLSuffix holds the current suffix
	APIURLSuffix = APIURLSuffixV3
)

// Client for Mattermost
type Client struct {
	URL           string       // The location of the server like "http://localhost:8065"
	APIUrl        string       // The api location of the server like "http://localhost:8065/api/v3"
	HTTPClient    *http.Client // The http client
	AuthToken     string
	AuthType      string
	Teams         map[string]*Team
	Channels      map[string][]*Channel
	RequestID     string
	Etag          string
	ServerVersion string
	User          *User
}

// Result holds requets results
type Result struct {
	RequestID string
	Etag      string
	Data      interface{}
}

// NewClient with default settings
func NewClient(url string) *Client {
	return &Client{url,
		url + APIURLSuffix,
		&http.Client{},
		"", "",
		make(map[string]*Team),
		make(map[string][]*Channel),
		"", "", "", nil}
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func getCookie(name string, resp *http.Response) *http.Cookie {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// Login with username and passowrd
func (c *Client) Login(loginID string, password string) (*Result, *ClientError) {
	m := make(map[string]string)
	m["login_id"] = loginID
	m["password"] = password
	return c.login(m)
}

func (c *Client) login(m map[string]string) (*Result, *ClientError) {
	r, err := c.DoAPIPost("/users/login", MapToJSON(m))
	if err != nil {
		return nil, err
	}
	c.AuthToken = r.Header.Get(HeaderToken)
	c.AuthType = HeaderBearer
	sessionToken := getCookie(SessionCookieToken, r)

	if c.AuthToken != sessionToken.Value {
		return nil, NewClientError("/users/login", "model.client.login.app_error", nil, "")
	}

	defer closeBody(r)
	c.User = UserFromJSON(r.Body)
	return &Result{r.Header.Get(HeaderRequestID),
		r.Header.Get(HeaderETAGServer),
		c.User}, nil
}

// DoAPIPost performs client api action
func (c *Client) DoAPIPost(url string, data string) (*http.Response, *ClientError) {
	rq, _ := http.NewRequest("POST", c.APIUrl+url, strings.NewReader(data))

	if len(c.AuthToken) > 0 {
		rq.Header.Set(HeaderAuth, c.AuthType+" "+c.AuthToken)
	}

	if rp, err := c.HTTPClient.Do(rq); err != nil {
		return nil, NewClientError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return nil, ClientErrorFromJSON(rp.Body)
	} else {
		return rp, nil
	}
}

// DoPost perform post action on url
func (c *Client) DoPost(url, data, contentType string) (*http.Response, *ClientError) {
	rq, _ := http.NewRequest("POST", c.URL+url, strings.NewReader(data))
	rq.Header.Set("Content-Type", contentType)

	if rp, err := c.HTTPClient.Do(rq); err != nil {
		return nil, NewClientError(url, "client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return nil, ClientErrorFromJSON(rp.Body)
	} else {
		return rp, nil
	}
}

// DoAPIGet performs client api get action
func (c *Client) DoAPIGet(url string, data string, etag string) (*http.Response, *ClientError) {
	rq, _ := http.NewRequest("GET", c.APIUrl+url, strings.NewReader(data))

	if len(etag) > 0 {
		rq.Header.Set(HeaderETAGClient, etag)
	}

	if len(c.AuthToken) > 0 {
		rq.Header.Set(HeaderAuth, c.AuthType+" "+c.AuthToken)
	}

	rp, err := c.HTTPClient.Do(rq)
	if err != nil {
		return nil, NewClientError(url, "model.client.connecting.app_error", nil, err.Error())
	} else if rp.StatusCode == 304 {
		return rp, nil
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return rp, ClientErrorFromJSON(rp.Body)
	}
	return rp, nil
}

// GetAllTeams returns a map of all teams using team ids as the key.
func (c *Client) GetAllTeams() (*Result, *ClientError) {
	r, err := c.DoAPIGet("/teams/all", "", "")
	if err != nil {
		return nil, err
	}
	defer closeBody(r)
	c.Teams = TeamMapFromJSON(r.Body)
	return &Result{r.Header.Get(HeaderRequestID),
		r.Header.Get(HeaderETAGServer), c.Teams}, nil
}

// GetTeamRoute concat route with team id
func (c *Client) GetTeamRoute(name string) (string, error) {
	id, err := c.GetTeamID(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/teams/%v", id), nil
}

// GetTeamID from pre fetched list of teams
func (c *Client) GetTeamID(name string) (string, error) {
	for _, t := range c.Teams {
		if t.Name == name {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("Failed to find team id for %q", name)
}

// GetChannels from api returns []*channels
func (c *Client) GetChannels(team string, etag string) ([]*Channel, *ClientError) {
	if v, ok := c.Channels[team]; ok {
		return v, nil
	}
	route, rErr := c.GetTeamRoute(team)
	if rErr != nil {
		return nil, NewClientError("Team Lookup", "client.channels.get", nil, rErr.Error())
	}
	r, err := c.DoAPIGet(route+"/channels/", "", etag)
	if err != nil {
		return nil, err
	}
	defer closeBody(r)
	c.Channels[team] = ChannelListFromJSON(r.Body).Channels
	return c.Channels[team], nil
}

// FindChannelIDByName will take a name of a channel and return id
func (c *Client) FindChannelIDByName(name, team string) (string, error) {
	channels, err := c.GetChannels(team, "")
	if err != nil {
		return "", fmt.Errorf(err.Message)
	}
	for _, c := range channels {
		if c.Name == name {
			return c.ID, nil
		}
	}
	return "", fmt.Errorf("Channel %q not found for team %q", name, team)
}

// CreateIncomingWebhook create incoming webhook
func (c *Client) CreateIncomingWebhook(team string, hook *IncomingWebhook) (*Result, *ClientError) {
	route, err := c.GetTeamRoute(team)
	if err != nil {
		return nil, NewClientError("Team Lookup", "client.create.incoming.app_error", nil, err.Error())
	}
	r, cErr := c.DoAPIPost(route+"/hooks/incoming/create", hook.ToJSON())
	if cErr != nil {
		return nil, cErr
	}
	defer closeBody(r)
	return &Result{r.Header.Get(HeaderRequestID),
		r.Header.Get(HeaderETAGServer),
		IncomingWebhookFromJSON(r.Body)}, nil
}

// CreateOutgoingWebhook creates outgoing webhook
func (c *Client) CreateOutgoingWebhook(team string, hook *OutgoingWebhook) (*Result, *ClientError) {
	route, err := c.GetTeamRoute(team)
	if err != nil {
		return nil, NewClientError("Team Lookup", "client.create.incoming.app_error", nil, err.Error())
	}
	r, cErr := c.DoAPIPost(route+"/hooks/outgoing/create", hook.ToJSON())
	if cErr != nil {
		return nil, cErr
	}
	defer closeBody(r)
	return &Result{r.Header.Get(HeaderRequestID),
		r.Header.Get(HeaderETAGServer),
		OutgoingWebhookFromJSON(r.Body)}, nil
}

// PostToWebhook sent payload to webhook
func (c *Client) PostToWebhook(id, payload string) (*Result, *ClientError) {
	r, err := c.DoPost("/hooks/"+id, payload, "application/json")
	if err != nil {
		return nil, err
	}
	defer closeBody(r)
	return &Result{r.Header.Get(HeaderRequestID),
		r.Header.Get(HeaderETAGServer), nil}, nil
}

//
// func (c *Client) DeleteIncomingWebhook(id string) (*Result, *ClientError) {
// 	data := make(map[string]string)
// 	data["id"] = id
// 	if r, err := c.DoApiPost(c.GetTeamRoute()+"/hooks/incoming/delete", MapToJson(data)); err != nil {
// 		return nil, err
// 	} else {
// 		defer closeBody(r)
// 		return &Result{r.Header.Get(HEADER_REQUEST_ID),
// 			r.Header.Get(HEADER_ETAG_SERVER), MapFromJson(r.Body)}, nil
// 	}
// }
//
// func (c *Client) DeleteOutgoingWebhook(id string) (*Result, *ClientError) {
// 	data := make(map[string]string)
// 	data["id"] = id
// 	if r, err := c.DoApiPost(c.GetTeamRoute()+"/hooks/outgoing/delete", MapToJson(data)); err != nil {
// 		return nil, err
// 	} else {
// 		defer closeBody(r)
// 		return &Result{r.Header.Get(HEADER_REQUEST_ID),
// 			r.Header.Get(HEADER_ETAG_SERVER), MapFromJson(r.Body)}, nil
// 	}
// }
