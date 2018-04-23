// Package handler - {{ .RelativeName }} service handler
{{ comment .Copyright }}
{{if .Licenses}}{{ comment .Licenses }}{{end}}

package handler

import (
  "{{ .ProjectDir }}/pb/{{ .PackageName }}"
  "golang.org/x/net/context"
  log "github.com/sirupsen/logrus"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc"
  "github.com/jmoiron/sqlx"
  "fmt"
)
  //"github.com/micro-grpc/mbox/services/authorize"
// ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

// {{ .RelativeName }}Service - service structure
type {{ .RelativeName }}Service struct {
  Debug bool
  Verbose int
  Version string
  Data string
  {{ .PackageName }}.{{ .RelativeName }}ServiceServer
}

// New{{ .RelativeName }}Service - new service
func New{{ .RelativeName }}Service() *{{ .RelativeName }}Service {
  s := &{{ .RelativeName }}Service{}
  return s
}

// Register - register service
func (s *{{ .RelativeName }}Service) Register(srv *grpc.Server)  {
  {{ .PackageName }}.Register{{ .RelativeName }}ServiceServer(srv, s)
}

// Read - load ine row
func (s *{{ .RelativeName }}Service) Read(ctx context.Context, in *{{ .PackageName }}.Query) (res *{{ .PackageName }}.Response, err error)  {
  //log.Println("Receive message", in.GetId(), "debug:", s.Debug)

  header := metadata.Pairs("ver", s.Version)
  grpc.SendHeader(ctx, header)

  trailer := metadata.Pairs("mes", "footer message")
  grpc.SetTrailer(ctx, trailer)

  //if user, ok := authorize.GetUser(ctx); ok {
  //  log.Println("USER:", user)
  //}
  //if organization, ok := authorize.GetTokenParams(ctx, "organization"); ok {
  //  organizationID, _ := authorize.GetTokenParams(ctx, "organization_id")
  //  log.Println("organization ID:", organizationID, "name:", organization)
  //}
  //if olala, ok := ctx.Value("olala").(string); ok {
  //  fmt.Println("olala:", olala)
  //}

  db, err := GetDB(ctx)
  if err != nil {
    log.Errorln("[DB:connect]", err.Error())
    return &{{ .PackageName }}.Response{}, err
  }
  ql := {{ .PackageName }}.NewQuery{{ .RelativeName }}(s.Verbose)
  ql.SetDB(db)

  //row, err := ql.Read(in)
  //if err != nil {
  //  log.Errorln("[Read]", err.Error())
  //}

  //return row, err
  return &{{ .PackageName }}.Response{}, nil
}

// GetConnect - return connection string from context
func GetConnect(ctx context.Context) (string, error) {
  conn, ok := ctx.Value("connectDB").(string)

  if !ok {
    return "", fmt.Errorf("is not connect string")
  }
  return conn, nil
}

// GetDB - return database connection
func GetDB(ctx context.Context) (*sqlx.DB, error) {
  db, ok := ctx.Value("DB").(*sqlx.DB)
  if !ok {
    return nil, fmt.Errorf("is not connect string")
  }
  return db, nil
}
