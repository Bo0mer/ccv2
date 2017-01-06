// Package ccv2 implements read-only Cloud Foundry Cloud Controller API client,
// targeting version 2 of the API.
// For more details about the API, see https://apidocs.cloudfoundry.org/.
//
// Note that this is only a Cloud Controller client, thus it does not deal
// with authentication and authorization. It is responsiblity of the client
// the provided an authenticated HTTP client.
//
// Example usage:
//   apiURL, _ := url.Parse("https://api.bosh-lite.com")
//   cf := &ccv2.Client{
//     API:        apiURL,
//     HTTPClient: http.DefaultClient,
//   }
//
//   ctx := context.Background()
//   infoCtx, cancel := context.WithTimeout(ctx, time.Second*5)
//   defer cancel()
//   info, err := cf.Info(infoCtx)
//   if err != nil {
//     log.Fatalf("error fetching info: %v\n", err)
//   }
//
//   authConfig := &oauth2.Config{
//     ClientID: "cf",
//     Scopes:   []string{""},
//     Endpoint: oauth2.Endpoint{
//       AuthURL:  info.TokenEndpoint + "/oauth/auth",
//       TokenURL: info.TokenEndpoint + "/oauth/token",
//     },
//   }
//
//   tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
//   defer cancel()
//   token, err := authConfig.PasswordCredentialsToken(tCtx, "admin", "admin")
//   if err != nil {
//     log.Fatalf("error fetching token: %v\n", err)
//   }
//
//   cf = &ccv2.Client{
//     API:        apiURL,
//     HTTPClient: authConfig.Client(ctx, token),
//   }
//   // Use cf as authenticated on behalf of admin:admin.
package ccv2
