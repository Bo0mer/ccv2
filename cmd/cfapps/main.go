// cfapps is a simple command line utility that prints one's Cloud Foundry
// organizations and applications.
//
// It's purpose is just to demonstrate how to use package ccv2.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Bo0mer/ccv2"
	"golang.org/x/oauth2"
)

var api string
var username string
var password string

func init() {
	flag.StringVar(&api, "api", "https://api.bosh-lite.com", "URL of the CC.")
	flag.StringVar(&username, "username", "admin", "Username.")
	flag.StringVar(&password, "password", "admin", "Password.")
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	apiURL, err := url.Parse(api)
	if err != nil {
		log.Fatalf("error parsing api url: %v\n", err)
	}

	cf := &ccv2.Client{
		API:        apiURL,
		HTTPClient: http.DefaultClient,
	}
	ctx := context.Background()
	infoCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	info, err := cf.Info(infoCtx)
	if err != nil {
		log.Fatalf("error fetching info: %v\n", err)
	}

	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  info.TokenEndpoint + "/oauth/auth",
			TokenURL: info.TokenEndpoint + "/oauth/token",
		},
	}
	tokenCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	token, err := authConfig.PasswordCredentialsToken(tokenCtx, username, password)
	if err != nil {
		log.Fatalf("error fetching oauth2 token: %v\n", err)
	}
	cf = &ccv2.Client{
		API:        apiURL,
		HTTPClient: authConfig.Client(ctx, token),
	}

	orgsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	orgs, err := cf.Organizations(orgsCtx)
	if err != nil {
		log.Fatalf("error fetching organizations: %v\n", err)
	}
	fmt.Printf("===== Organizations =====\n")
	for _, org := range orgs {
		fmt.Printf("%s\n", org.Entity.Name)
	}

	appsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	apps, err := cf.Applications(appsCtx)
	if err != nil {
		log.Fatalf("error fetching applications: %v\n", err)
	}
	fmt.Printf("===== Applications ===== \n")
	for _, app := range apps {
		fmt.Printf("%s\n", app.Entity.Name)
	}
}
