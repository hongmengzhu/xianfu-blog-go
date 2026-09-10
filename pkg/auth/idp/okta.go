package idp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

type OktaIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
	Host   string
}

func NewOktaIdProvider(clientId string, clientSecret string, redirectUrl string, hostUrl string) *OktaIdProvider {
	idp := &OktaIdProvider{}

	config := idp.getConfig(hostUrl, clientId, clientSecret, redirectUrl)
	config.ClientID = clientId
	config.ClientSecret = clientSecret
	config.RedirectURL = redirectUrl
	idp.Config = config
	idp.Host = hostUrl
	return idp
}

func (idp *OktaIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *OktaIdProvider) getConfig(hostUrl string, clientId string, clientSecret string, redirectUrl string) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		TokenURL: fmt.Sprintf("%s/v1/token", hostUrl),
		AuthURL:  fmt.Sprintf("%s/v1/authorize", hostUrl),
	}

	config := &oauth2.Config{
		// openid is required for authentication requests
		// get more details via: https://developer.okta.com/docs/reference/api/oidc/#reserved-scopes
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     endpoint,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	return config
}

type OktaToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

// GetToken use code to get access_token
// get more details via: https://developer.okta.com/docs/reference/api/oidc/#token
func (idp *OktaIdProvider) GetToken(code string) (*oauth2.Token, error) {
	payload := url.Values{}
	payload.Set("code", code)
	payload.Set("grant_type", "authorization_code")
	payload.Set("client_id", idp.Config.ClientID)
	payload.Set("client_secret", idp.Config.ClientSecret)
	payload.Set("redirect_uri", idp.Config.RedirectURL)
	resp, err := idp.Client.PostForm(idp.Config.Endpoint.TokenURL, payload)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pToken := &OktaToken{}
	err = json.Unmarshal(data, pToken)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal token response: %s", err.Error())
	}

	token := &oauth2.Token{
		AccessToken:  pToken.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: pToken.RefreshToken,
		Expiry:       time.Unix(time.Now().Unix()+int64(pToken.ExpiresIn), 0),
	}
	return token, nil
}

type OktaUserInfo struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Picture           string `json:"picture"`
	Sub               string `json:"sub"`
}

// GetUserInfo use token to get user profile
// get more details via: https://developer.okta.com/docs/reference/api/oidc/#userinfo
func (idp *OktaIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/userinfo", idp.Host), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")
	resp, err := idp.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// First unmarshal into a map to capture all claims
	var rawClaims map[string]interface{}
	err = json.Unmarshal(body, &rawClaims)
	if err != nil {
		return nil, err
	}

	var oktaUserInfo OktaUserInfo
	err = json.Unmarshal(body, &oktaUserInfo)
	if err != nil {
		return nil, err
	}

	// Convert raw claims to string map for Extra field
	extra := make(map[string]string)
	for k, v := range rawClaims {
		if v != nil {
			// Convert to string representation
			switch val := v.(type) {
			case string:
				extra[k] = val
			case float64:
				extra[k] = fmt.Sprintf("%v", val)
			case bool:
				extra[k] = fmt.Sprintf("%v", val)
			default:
				// For complex types, marshal to JSON string
				if jsonVal, err := json.Marshal(val); err == nil {
					extra[k] = string(jsonVal)
				}
			}
		}
	}

	userInfo := UserInfo{
		Id:          oktaUserInfo.Sub,
		Username:    oktaUserInfo.PreferredUsername,
		DisplayName: oktaUserInfo.Name,
		Email:       oktaUserInfo.Email,
		AvatarUrl:   oktaUserInfo.Picture,
		Extra:       extra,
	}
	return &userInfo, nil
}
