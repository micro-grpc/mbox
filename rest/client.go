package rest

import (
  "fmt"
  "gopkg.in/resty.v1"
)

type Client struct {
  URL      string
  Debug    bool
  Verbose  int
  Test     bool
  Host     string
  Protocol string
  User     string
  Token    string
  C        *resty.Client
}

func (client *Client) Init() {
  client.C = resty.New()
  client.URL = fmt.Sprintf("%s://%s/api/", client.Protocol, client.Host)
  client.C.SetBasicAuth(client.User, client.Token)
  client.C.SetDebug(client.Debug)
  client.C.SetHostURL(client.URL)
  client.C.SetRESTMode()
}

func (client *Client) PUT(url string, data []byte) error {
  c := resty.New()
  c.SetBasicAuth(client.User, client.Token)
  c.SetDebug(client.Debug)
  resp, err := c.R().
    SetHeader("Content-Type", "application/json").
    SetBody(data).
    // SetResult(&AuthSuccess{}).    // or SetResult(AuthSuccess{}).
    Put(fmt.Sprintf("%s://%s/%s", client.Protocol, client.Host, url))
  if err != nil {
    return err
  }
  if resp.StatusCode() != 200 {
    //if resp.Body() != nil {
      // fmt.Printf("\nResponse Body: %v", resp)
      //return nil
    //}
  //noinspection ALL
  return fmt.Errorf("Error status code: %d", resp.StatusCode())
  }
  return nil
}
