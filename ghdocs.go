package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	personalAccessToken string
	issuesCollection    allIssues
	org                 string
)

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

type allIssues struct {
	issues                                                                      github.IssuesSearchResult
	users                                                                       []github.User
	issuesnewPublic, issues2mPublic, issues6mPublic, issues1yPublic             []github.Issue
	issuesnewOrgMember, issues2mOrgMember, issues6mOrgMember, issues1yOrgMember []github.Issue
	prnewPublic, pr2mPublic, pr6mPublic, pr1yPublic                             []github.Issue
	prnewOrgMember, pr2mOrgMember, pr6mOrgMember, pr1yOrgMember                 []github.Issue
}

func main() {
	org = os.Getenv("GH_ORG")
	personalAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

	if len(personalAccessToken) == 0 {
		log.Fatal("Before you can use this you must set the GITHUB_ACCESS_TOKEN environment variable.")
	}
	if len(org) != 1 {
		log.Fatal("You need to have a single organization name set to GH_ORG environmental variable.")
	}

	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	client := github.NewClient(oauthClient)
	twoMonthsAgo := time.Now().Add(time.Hour * 24 * 30 * 2 * -1)
	sixMonthsAgo := time.Now().Add(time.Hour * 24 * 30 * 6 * -1)
	oneYearAgo := time.Now().Add(time.Hour * 24 * 365 * -1)

	err := populateUsers(org, client)
	if err != nil {
		fmt.Println("Error getting users for " + org)
	}

	err = populateIssues(org, client)
	if err != nil {
		fmt.Println("Error getting issues for " + org)
	}

	totalPR := 0
	totalIssues := 0
	for _, issue := range issuesCollection.issues.Issues {
		if issue.PullRequestLinks == nil {
			totalIssues++
			if issue.UpdatedAt.After(twoMonthsAgo) {
				populateIssueGroup(&issuesCollection.issuesnewOrgMember, &issuesCollection.issuesnewPublic, issue.User, issue)
			} else if issue.UpdatedAt.After(sixMonthsAgo) {
				populateIssueGroup(&issuesCollection.issues2mOrgMember, &issuesCollection.issues2mPublic, issue.User, issue)
			} else if issue.UpdatedAt.After(oneYearAgo) {
				populateIssueGroup(&issuesCollection.issues6mOrgMember, &issuesCollection.issues6mPublic, issue.User, issue)
			} else {
				populateIssueGroup(&issuesCollection.issues1yOrgMember, &issuesCollection.issues1yPublic, issue.User, issue)
			}
		} else {
			totalPR++
			if issue.UpdatedAt.After(twoMonthsAgo) {
				populateIssueGroup(&issuesCollection.prnewOrgMember, &issuesCollection.prnewPublic, issue.User, issue)
			} else if issue.UpdatedAt.After(sixMonthsAgo) {
				populateIssueGroup(&issuesCollection.pr2mOrgMember, &issuesCollection.pr2mPublic, issue.User, issue)
			} else if issue.UpdatedAt.After(oneYearAgo) {
				populateIssueGroup(&issuesCollection.pr6mOrgMember, &issuesCollection.pr6mPublic, issue.User, issue)
			} else {
				populateIssueGroup(&issuesCollection.pr1yOrgMember, &issuesCollection.pr1yPublic, issue.User, issue)
			}
		}
	}

	fmt.Printf("\nSummary\n")
	fmt.Printf("\n\nPull Requests - %d\n", totalPR)
	fmt.Printf("\n  Employee Pull Requests")
	fmt.Printf("\n    New: \t\t%d\n    2-6 months old:\t%d\n    6-12 months old:\t%d\n    Older that 1 year:\t%d", len(issuesCollection.prnewOrgMember), len(issuesCollection.pr2mOrgMember), len(issuesCollection.pr6mOrgMember), len(issuesCollection.pr1yOrgMember))
	fmt.Printf("\n\n  Public Pull Requests")
	fmt.Printf("\n    New: \t\t%d\n    2-6 months old:\t%d\n    6-12 months old:\t%d\n    Older that 1 year:\t%d\n", len(issuesCollection.prnewPublic), len(issuesCollection.pr2mPublic), len(issuesCollection.pr6mPublic), len(issuesCollection.pr1yPublic))

	fmt.Printf("\nIssues - %d\n", totalIssues)
	fmt.Printf("\n  Employee Issues")
	fmt.Printf("\n    New: \t\t%d\n    2-6 months old:\t%d\n    6-12 months old:\t%d\n    Older that 1 year:\t%d", len(issuesCollection.issuesnewOrgMember), len(issuesCollection.issues2mOrgMember), len(issuesCollection.issues6mOrgMember), len(issuesCollection.issues1yOrgMember))
	fmt.Printf("\n\n  Public Issues")
	fmt.Printf("\n    New: \t\t%d\n    2-6 months old:\t%d\n    6-12 months old:\t%d\n    Older that 1 year:\t%d\n", len(issuesCollection.issuesnewPublic), len(issuesCollection.issues2mPublic), len(issuesCollection.issues6mPublic), len(issuesCollection.issues1yPublic))

	fmt.Printf("\nPublic Details - Pull Requests\n")
	printIssues(issuesCollection.pr1yPublic, "1+ year old pull requests")
	printIssues(issuesCollection.pr6mPublic, "6-12 month old pull requests")
	printIssues(issuesCollection.pr2mPublic, "2-6 month old pull requests")
	printIssues(issuesCollection.prnewPublic, "New Pull Requests")

	fmt.Printf("\n\nPublic Details - Issues\n")
	printIssues(issuesCollection.issues1yPublic, "1+ year old issues")
	printIssues(issuesCollection.issues6mPublic, "6-12 month old issues")
	printIssues(issuesCollection.issues2mPublic, "2-6 month old issues")
	printIssues(issuesCollection.issuesnewPublic, "New Issues")
	fmt.Printf("(** indicates that issue or pr is older than 2 days with no comments. ## indicates that issue or pr has no comments in 4 weeks)\n\n")
}

func printIssues(issues []github.Issue, title string) {
	fmt.Printf("  %s\n", title)
	for _, issue := range issues {
		var title string
		if len(*issue.Title) > 60 {
			title = (*issue.Title)[:59] + "..."
		} else {
			title = *issue.Title
		}
		getRepoName(issue)

		fmt.Printf("   %-23s %-19s %4s(%5d) %-62s %s\n", getRepoName(issue), *issue.User.Login, attentionStatus(issue), *issue.Number, title, issue.CreatedAt.Format("Jan 2, 2006"))
	}
	fmt.Printf("\n")
}

func attentionStatus(issue github.Issue) (attentionStatus string) {
	attentionStatus = ""
	if *issue.Comments == 0 && issue.CreatedAt.Before(time.Now().Add(time.Hour*24*2*-1)) {
		attentionStatus = "**"
	}
	if issue.UpdatedAt.Before(time.Now().Add(time.Hour * 24 * 30 * -1)) {
		attentionStatus += "##"
	}

	return
}
func getRepoName(issue github.Issue) string {
	url := *issue.URL
	// fmt.Println(url)
	// fmt.Println(url[30+len(org):])
	// fmt.Println(url[:strings.LastIndex(url, "issues")])
	// fmt.Println(url[30+len(org)])
	repo := (url)[30+len(org) : strings.LastIndex(url, "issues")-1]
	return repo
}
func isUserAnOrgMember(thisuser github.User) bool {
	for _, user := range issuesCollection.users {
		if *thisuser.Login == *user.Login {
			return true
		}
	}
	return false
}

func populateIssueGroup(memberissuelist *[]github.Issue, publicissuelist *[]github.Issue, user *github.User, issue github.Issue) {
	if isUserAnOrgMember(*user) {
		*memberissuelist = append(*memberissuelist, issue)
	} else {
		*publicissuelist = append(*publicissuelist, issue)
	}
}

func populateUsers(org string, client *github.Client) (err error) {
	useropt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{},
	}
	for {
		userSubset, resp, err := client.Organizations.ListMembers(org, useropt)
		fmt.Print(".")
		if err != nil {
			return err
		}
		issuesCollection.users = append(issuesCollection.users, userSubset...)
		if resp.NextPage == 0 {
			break
		}
		useropt.ListOptions.Page = resp.NextPage
	}
	return
}

func populateIssues(org string, client *github.Client) (err error) {
	issueopt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		issuesSubset, resp, err := client.Search.Issues("is:open is:public user:"+org, issueopt)
		fmt.Print("....")
		if err != nil {
			return err
		}
		issuesCollection.issues.Issues = append(issuesCollection.issues.Issues, issuesSubset.Issues...)
		if resp.NextPage == 0 {
			break
		}
		issueopt.ListOptions.Page = resp.NextPage
	}
	return
}
