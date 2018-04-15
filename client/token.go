package client

import (
	"context"
	"fmt"
)

type TokenAuth struct {
	Token    string
	Security bool
	Schema   string
}

func (t *TokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("%s %s", t.Schema, t.Token),
	}, nil
}

func (t *TokenAuth) RequireTransportSecurity() bool {
	return t.Security
}

// fakeOAuth2TokenSource implements a fake oauth2.TokenSource for the purpose of credentials test.
// type FakeOAuth2TokenSource struct {
//   accessToken string
//   Security bool
//   schema string
//   expiry time.Time
// }
//
// func (ts *FakeOAuth2TokenSource) Token() (*oauth2.Token, error) {
//   t := &oauth2.Token{
//     AccessToken: ts.accessToken,
//     Expiry:      time.Now().Add(1 * time.Minute),
//     TokenType:   ts.schema,
//   }
//   return t, nil
// }
//
// func (ts *FakeOAuth2TokenSource) RequireTransportSecurity() bool {
//   return ts.Security
// }
