package lib

import (
  "strings"
  "encoding/base64"
  "fmt"
)

// type AuthData struct {
//   User string
//   Data map[string]string
// }
//
// type ParseTokenInterface interface {
//   ParseToken(token string) (AuthData, error)
// }
//
// type DefaultParseToken func(token string) (AuthData, error)
//
// type ParseTokenType struct {
//   Secret string
// }

// ParseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func ParseBasicAuth(auth string) (username, password string, ok bool) {
  const prefix = "Basic "
  src := ""
  if !strings.HasPrefix(auth, prefix) {
    src = auth
  } else {
    src = auth[len(prefix):]
  }
  c, err := base64.StdEncoding.DecodeString(src)
  if err != nil {
    return
  }
  cs := string(c)
  s := strings.IndexByte(cs, ':')
  if s < 0 {
    return
  }
  return cs[:s], cs[s+1:], true
}

// SetBasicAuth - sets the request's Authorization header to use HTTP
// Basic Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password
// are not encrypted.
func SetBasicAuth(username, password string) string {
  return fmt.Sprintf("Authorization Basic %v", BasicAuth(username, password))
}

// BasicAuth - see 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
  auth := fmt.Sprintf("%s:%s", username, password)
  return base64.StdEncoding.EncodeToString([]byte(auth))
}


