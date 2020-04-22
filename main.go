package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"

	"golang.org/x/oauth2"

	"github.com/picatz/snyk"

	"github.com/google/go-github/v31/github"
)

func main() {
	ghtoken := os.Getenv("SS_GHTOKEN")
	ghorg := os.Getenv("SS_GHORG")
	snyktoken := os.Getenv("SS_SNYKTOKEN")
	snykorg := os.Getenv("SS_SNYKORG")
	snykintid := os.Getenv("SS_SNYKINTID")

	repos := getGitHubRepos(ghtoken, ghorg)

	// Must be a personal token to import from GitHub; this cannot be a service account.
	projects := getSnykProjects(snyktoken, snykorg)

	missing := compare(repos, projects)
	color.Yellow("Missing projects in Snyk: %s", missing)

	for _, project := range missing {
		color.Green("Attempting to import %s into Snyk...", project)
		split := strings.Split(project, "/")
		createSnykProject(snyktoken, snykorg, snykintid, split[0], split[1], "")
	}
}

func getGitHubRepos(token, org string) []string {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.ListByOrg(ctx, org, nil)
	//repos, _, err := client.Repositories.List(ctx, org, nil)
	if err != nil {
		log.Fatal(err)
	}

	var s []string
	for _, repo := range repos {
		s = append(s, *repo.FullName)
	}
	return s
}

func getSnykProjects(token string, org string) []string {
	client, err := snyk.NewClient(snyk.WithToken(token))
	if err != nil {
		log.Fatal(err)
	}

	projects, err := client.OrganizationProjects(context.TODO(), org)
	if err != nil {
		log.Fatal(err)
	}

	var s []string
	for _, project := range projects {
		s = append(s, strings.Split(project.Name, ":")[0])
	}
	return deduplicate(s)
}

func createSnykProject(token, org, integrationID, owner, name, branch string) {
	client, err := snyk.NewClient(snyk.WithToken(token))
	if err != nil {
		log.Fatal(err)
	}

	cb, err := client.OrganizationImportProject(context.TODO(), org, integrationID, snyk.GitHubImport(owner, name, branch, nil))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cb)
}

func deduplicate(s []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range s {
		if encountered[s[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[s[v]] = true
			// Append to result slice.
			result = append(result, s[v])
		}
	}
	// Return the new slice.
	return result
}

func compare(a, b []string) []string {
	for i := len(a) - 1; i >= 0; i-- {
		for _, v := range b {
			if a[i] == v {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}
	return a
}
