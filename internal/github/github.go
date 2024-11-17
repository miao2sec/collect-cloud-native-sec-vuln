package github

import (
	"context"
	"github.com/google/go-github/v66/github"
	"github.com/rs/zerolog/log"
	"github.com/y4ney/cloud-native-security-vuln/internal/component"
	"github.com/y4ney/cloud-native-security-vuln/internal/config"
	"golang.org/x/oauth2"
	"net/http"
)

type Client struct {
	Token  string
	Ctx    context.Context
	Client *github.Client
}

type ClientFunc func(*Client)

// WithToken 配置 Token
func WithToken(token string) ClientFunc {
	return func(c *Client) { c.Token = token }
}

func NewClient(opts ...ClientFunc) *Client {
	var (
		c  = &Client{Ctx: context.Background()}
		tc = &http.Client{}
	)

	for _, opt := range opts {
		opt(c)
	}

	if c.Token == "" || c.Token == config.Token {
		log.Debug().Msg("token is empty")
	} else {
		tc = oauth2.NewClient(c.Ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token}))
	}

	c.Client = github.NewClient(tc)
	return c
}

func (c *Client) GetAdvisories(component *component.Component) ([]*github.SecurityAdvisory, error) {
	advisories, _, err := c.Client.SecurityAdvisories.ListRepositorySecurityAdvisories(c.Ctx,
		component.Owner, component.Repo, nil)
	if err != nil {
		return nil, err
	}
	return advisories, nil
}
