package idp

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	Id          string
	Username    string
	DisplayName string
	UnionId     string
	Email       string
	Phone       string
	CountryCode string
	AvatarUrl   string
	Extra       map[string]string
}

type ProviderInfo struct {
	Type          string
	SubType       string
	ClientId      string
	ClientSecret  string
	ClientId2     string
	ClientSecret2 string
	AppId         string
	HostUrl       string
	RedirectUrl   string
	DisableSsl    bool
	CodeVerifier  string

	TokenURL    string
	AuthURL     string
	UserInfoURL string
	UserMapping map[string]string

	AppCertificate  string
	RootCertificate string
}

type IdProvider interface {
	SetHttpClient(client *http.Client)
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
}

func GetIdProvider(idpInfo *ProviderInfo, redirectUrl string) (IdProvider, error) {
	switch idpInfo.Type {
	case "GitHub":
		return NewGithubIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Google":
		return NewGoogleIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Twitter":
		provider := NewTwitterIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl)
		provider.CodeVerifier = idpInfo.CodeVerifier
		return provider, nil
	case "LinkedIn":
		return NewLinkedInIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "GitLab":
		return NewGitlabIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Gitee":
		return NewGiteeIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Casdoor":
		return NewCasdoorIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl, idpInfo.HostUrl), nil
	case "Okta":
		return NewOktaIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl, idpInfo.HostUrl), nil
	case "ADFS":
		return NewAdfsIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl, idpInfo.HostUrl), nil
	case "AzureADB2C":
		return NewAzureAdB2cProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl, idpInfo.HostUrl, idpInfo.AppId), nil
	case "Alipay":
		return NewAlipayIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl, idpInfo.AppCertificate, idpInfo.RootCertificate)
	case "MetaMask":
		return NewMetaMaskIdProvider(), nil
	case "Web3Onboard":
		return NewWeb3OnboardIdProvider(), nil
	case "Facebook":
		return NewFacebookIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Telegram":
		return NewTelegramIdProvider(idpInfo.ClientId, idpInfo.ClientSecret, redirectUrl), nil
	case "Custom", "Custom Flexible":
		return NewCustomIdProvider(idpInfo, redirectUrl), nil
	default:
		if isGothSupport(idpInfo.Type) {
			return NewGothIdProvider(idpInfo.Type, idpInfo.ClientId, idpInfo.ClientSecret, idpInfo.ClientId2, idpInfo.ClientSecret2, redirectUrl, idpInfo.HostUrl)
		}
		if strings.HasPrefix(idpInfo.Type, "Custom") {
			return NewCustomIdProvider(idpInfo, redirectUrl), nil
		}
		return nil, fmt.Errorf("OAuth provider type: %s is not supported", idpInfo.Type)
	}
}

var gothList = []string{
	"Apple",
	"AzureAD",
	"Slack",
	"Steam",
	"Line",
	"Amazon",
	"Auth0",
	"BattleNet",
	"Bitbucket",
	"Box",
	"CloudFoundry",
	"Dailymotion",
	"Deezer",
	"DigitalOcean",
	"Discord",
	"Dropbox",
	"EveOnline",
	"Fitbit",
	"Gitea",
	"Heroku",
	"InfluxCloud",
	"Instagram",
	"Intercom",
	"Kakao",
	"Lastfm",
	"Mailru",
	"Meetup",
	"MicrosoftOnline",
	"Naver",
	"Nextcloud",
	"OneDrive",
	"Oura",
	"Patreon",
	"Paypal",
	"SalesForce",
	"Shopify",
	"Soundcloud",
	"Spotify",
	"Strava",
	"Stripe",
	"TikTok",
	"Tumblr",
	"Twitch",
	"Typetalk",
	"Uber",
	"VK",
	"Wepay",
	"Xero",
	"Yahoo",
	"Yammer",
	"Yandex",
	"Zoom",
}

func isGothSupport(provider string) bool {
	for _, value := range gothList {
		if strings.EqualFold(value, provider) {
			return true
		}
	}
	return false
}
