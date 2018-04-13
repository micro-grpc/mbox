package authorize

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro-grpc/mbox/lib"
	pbauthorize "github.com/micro-grpc/mbox/pb/authorize"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const tokenInfoKey string = "tokenInfo"

type Authorize struct {
	Secret     string
	Security   bool
	Params     []string
	Prefix     string
	Check      Check
	CheckBasic CheckBasic
}

type Check func(token string) error
type CheckBasic func(username string, password string) error

// ParseToken
func (a *Authorize) ParseToken(t string) (*pbauthorize.AuthorizeResponse, error) {
	err := a.Check(t)
	if err != nil {
		return nil, err
	}
	log.Debugln("[ParseToken] Token:", t)
	log.Debugln("[ParseToken] Secret:", a.Secret)
	// claims := jwt.StandardClaims{}
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, fmt.Errorf("Error 401 Unauthorized")
	}
	log.Debugln("[ParseToken] ---> ParseWithClaims: token", token.Claims)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwtauth: expecting jwt.MapClaims")
	}

	res := &pbauthorize.AuthorizeResponse{}
	if user, ok := claims["user"]; ok {
		res.User = user.(string)
	}
	if userId, ok := claims["user_id"]; ok {
		res.UserId = userId.(string)
	}
	if len(a.Params) > 0 {
		for _, k := range a.Params {
			log.Println("key:", k)
			if val, ok := claims[k]; ok {
				param := pbauthorize.Param{Name: k, Value: val.(string)}
				res.Params = append(res.Params, &param)
			}
		}
	}

	return res, nil
}

func (a *Authorize) ParseBasic(t string) (*pbauthorize.AuthorizeResponse, error) {
	username, password, ok := lib.ParseBasicAuth(t)
	if !ok {
		return nil, fmt.Errorf("invalid basic auth")
	}
	err := a.CheckBasic(username, password)
	if err != nil {
		return nil, err
	}

	res := &pbauthorize.AuthorizeResponse{}
	res.User = username
	res.UserId = username

	return res, nil
}

// NewAuthorize
func NewAuthorize(secret string, security bool, Check Check, CheckBasic CheckBasic) *Authorize {
	authorize := &Authorize{Secret: secret, Security: security}
	authorize.Params = []string{}
	authorize.Check = Check
	authorize.CheckBasic = CheckBasic
	return authorize
}

// GetUser - return authorize user
func GetUser(ctx context.Context) (string, bool) {
	// log.Println(ctx.Value(tokenInfoKey))
	if tokenInfo, ok := ctx.Value(tokenInfoKey).(*pbauthorize.AuthorizeResponse); ok {
		if len(tokenInfo.User) > 0 {
			return tokenInfo.User, true
		} else if len(tokenInfo.UserId) > 0 {
			return tokenInfo.UserId, true
		}
	}
	return "", false
}

// GetUserId - return authorize user id
func GetUserId(ctx context.Context) (string, bool) {
	if tokenInfo, ok := ctx.Value(tokenInfoKey).(*pbauthorize.AuthorizeResponse); ok {
		if len(tokenInfo.UserId) > 0 {
			return tokenInfo.UserId, true
		}
	}
	return "", false
}

// GetTokenParams - return custom param
func GetTokenParams(ctx context.Context, name string) (string, bool) {
	if tokenInfo, ok := ctx.Value(tokenInfoKey).(*pbauthorize.AuthorizeResponse); ok {
		for _, param := range tokenInfo.Params {
			if param.Name == name {
				return param.Value, true
			}
		}
	}
	return "", false
}
