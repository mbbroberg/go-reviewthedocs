package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	personalAccessToken string
	// issuesCollection    allIssues
	org string
)

// TokenSource is an encapsulation of the AccessToken string
type TokenSource struct {
	AccessToken string
}

// Token authenticates via oauth
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	org = os.Getenv("GH_ORG")
	personalAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

	if len(personalAccessToken) == 0 {
		log.Fatal("Before you can use this you must set the GITHUB_ACCESS_TOKEN environment variable.")
	}
	if len(org) < 1 {
		log.Fatal("You need to have a single organization name set to GH_ORG environmental variable.")
	}

	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient) // authenticated to GitHub here

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(org, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		fmt.Println(*repo.FullName)
	}
}
