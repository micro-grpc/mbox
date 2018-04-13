package server

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
)

type contextKey uint

const (
  userKey contextKey = iota
)

// authenticate takes the user from the gRPC metadata and
// adds it into the context values, if available. Otherwise
// an error with gRPC code Unauthenticated is returned.
func Authenticate(ctx context.Context) (context.Context, error) {
  user, ok := ExtractUserFromMD(ctx)
  if !ok {
    return ctx, status.Errorf(codes.Unauthenticated, "request is not authenticated")
  }
  return context.WithValue(ctx, userKey, user), nil
}

// getUser returns the user previously added via authenticate.
func GetUser(ctx context.Context) (string, bool) {
  if user, ok := ctx.Value(userKey).(string); ok && user != "" {
    return user, true
  }
  return "", false
}

// extractUserFromMD extracts the user from gRPC metadata.
func ExtractUserFromMD(ctx context.Context) (string, bool) {
  md, ok := metadata.FromIncomingContext(ctx)
  if !ok {
    return "", false
  }
  values := md["user"]
  if len(values) != 1 || values[0] == "" {
    return "", false
  }
  return values[0], true
}
