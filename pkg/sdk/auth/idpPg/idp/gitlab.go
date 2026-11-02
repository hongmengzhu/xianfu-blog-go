package idp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

type GitlabIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewGitlabIdProvider(clientId string, clientSecret string, redirectUrl string) *GitlabIdProvider {
	idp := &GitlabIdProvider{}

	config := idp.getConfig(clientId, clientSecret, redirectUrl)
	idp.Config = config

	return idp
}

func (idp *GitlabIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

// getConfig return a point of Config, which describes a typical 3-legged OAuth2 flow
func (idp *GitlabIdProvider) getConfig(clientId string, clientSecret string, redirectUrl string) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		TokenURL: "https://gitlab.com/oauth/token",
	}

	config := &oauth2.Config{
		Scopes:       []string{"read_user+profile"},
		Endpoint:     endpoint,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	return config
}

type GitlabProviderToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

// GetToken use code get access_token (*operation of getting code ought to be done in front)
// get more detail via: https://docs.gitlab.com/ee/api/oauth2.html
func (idp *GitlabIdProvider) GetToken(code string) (*oauth2.Token, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", idp.Config.ClientID)
	params.Add("client_secret", idp.Config.ClientSecret)
	params.Add("code", code)
	params.Add("redirect_uri", idp.Config.RedirectURL)

	accessTokenUrl := fmt.Sprintf("%s?%s", idp.Config.Endpoint.TokenURL, params.Encode())
	resp, err := idp.Client.Post(accessTokenUrl, "application/json;charset=UTF-8", nil)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gtoken := &GitlabProviderToken{}
	if err = json.Unmarshal(data, gtoken); err != nil {
		return nil, err
	}

	// gtoken.ExpiresIn always returns 0, so we set Expiry=7200 to avoid verification errors.
	token := &oauth2.Token{
		AccessToken:  gtoken.AccessToken,
		TokenType:    gtoken.TokenType,
		RefreshToken: gtoken.RefreshToken,
		Expiry:       time.Unix(time.Now().Unix()+int64(7200), 0),
	}

	return token, nil
}

type GitlabUserInfo struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
	Username        string      `json:"username"`
	State           string      `json:"state"`
	AvatarUrl       string      `json:"avatar_url"`
	WebUrl          string      `json:"web_url"`
	CreatedAt       time.Time   `json:"created_at"`
	Bio             string      `json:"bio"`
	BioHtml         string      `json:"bio_html"`
	Location        string      `json:"location"`
	PublicEmail     string      `json:"public_email"`
	Skype           string      `json:"skype"`
	Linkedin        string      `json:"linkedin"`
	Twitter         string      `json:"twitter"`
	WebsiteUrl      string      `json:"website_url"`
	Organization    string      `json:"organization"`
	JobTitle        string      `json:"job_title"`
	Pronouns        interface{} `json:"pronouns"`
	Bot             bool        `json:"bot"`
	WorkInformation interface{} `json:"work_information"`
	Followers       int         `json:"followers"`
	Following       int         `json:"following"`
	LastSignInAt    time.Time   `json:"last_sign_in_at"`
	ConfirmedAt     time.Time   `json:"confirmed_at"`
	LastActivityOn  string      `json:"last_activity_on"`
	Email           string      `json:"email"`
	ThemeId         int         `json:"theme_id"`
	ColorSchemeId   int         `json:"color_scheme_id"`
	ProjectsLimit   int         `json:"projects_limit"`
	CurrentSignInAt time.Time   `json:"current_sign_in_at"`
	Identities      []struct {
		Provider       string      `json:"provider"`
		ExternUid      string      `json:"extern_uid"`
		SamlProviderId interface{} `json:"saml_provider_id"`
	} `json:"identities"`
	CanCreateGroup                 bool        `json:"can_create_group"`
	CanCreateProject               bool        `json:"can_create_project"`
	TwoFactorEnabled               bool        `json:"two_factor_enabled"`
	External                       bool        `json:"external"`
	PrivateProfile                 bool        `json:"private_profile"`
	CommitEmail                    string      `json:"commit_email"`
	SharedRunnersMinutesLimit      interface{} `json:"shared_runners_minutes_limit"`
	ExtraSharedRunnersMinutesLimit interface{} `json:"extra_shared_runners_minutes_limit"`
}

// GetUserInfo use GitlabProviderToken gotten before return GitlabUserInfo
func (idp *GitlabIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	resp, err := idp.Client.Get("https://gitlab.com/api/v4/user?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	guser := GitlabUserInfo{}
	if err = json.Unmarshal(data, &guser); err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		Id:          strconv.Itoa(guser.Id),
		Username:    guser.Username,
		DisplayName: guser.Name,
		AvatarUrl:   guser.AvatarUrl,
		Email:       guser.Email,
	}
	return &userInfo, nil
}
