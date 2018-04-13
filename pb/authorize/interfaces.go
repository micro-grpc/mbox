package mbox_authorize

type AuthorizeService interface {
  ParseToken(token string) (*AuthorizeResponse, error)
}
