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
	// Get environment variables (secrets) to be used.
	ghtoken := os.Getenv("SS_GHTOKEN")
	ghorg := os.Getenv("SS_GHORG")
	snyktoken := os.Getenv("SS_SNYKTOKEN") // Must be a personal token, not a service account.
	snykorg := os.Getenv("SS_SNYKORG")
	snykintid := os.Getenv("SS_SNYKINTID")

	// Get all of our GitHub repos into a simple []string.
	repos := getGitHubRepos(ghtoken, ghorg)

	// Get all of our Snyk projects from a given org into a simple []string.
	projects := getSnykProjects(snyktoken, snykorg)

	// Generate a list of the repos in GitHub that do not have projects in Snyk.
	missing := compare(repos, projects)
	color.Yellow("Missing projects in Snyk: %s", missing)

	// Iterate over each missing item and attempt to import it.
	for _, project := range missing {
		color.Green("Attempting to import %s into Snyk...", project)
		// Split the project names at "/" to get the org and repo names separate.
		split := strings.Split(project, "/")
		// Initiate the import into Snyk.
		createSnykProject(snyktoken, snykorg, snykintid, split[0], split[1], "")
	}
}

// Function to return a list of all repos in GitHub a given user has access to.
func getGitHubRepos(token, org string) []string {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// List all repositories for the authenticated user.
	repos, _, err := client.Repositories.ListByOrg(ctx, org, nil)
	// Comment line above and uncomment line below to use user name instead of org name.
	//repos, _, err := client.Repositories.List(ctx, org, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a slice containing all the full repo names for an org.
	var s []string
	for _, repo := range repos {
		s = append(s, *repo.FullName)
	}
	return s
}

// Function to return a list of all projects in Snyk a given user has access to.
func getSnykProjects(token string, org string) []string {
	client, err := snyk.NewClient(snyk.WithToken(token))
	if err != nil {
		log.Fatal(err)
	}

	// List all projects for the authenticated user.
	projects, err := client.OrganizationProjects(context.TODO(), org)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a slice containing all the full project names for an org.
	var s []string
	for _, project := range projects {
		// Project names contain manifest files after a colon; only grab the project name.
		s = append(s, strings.Split(project.Name, ":")[0])
	}
	// In case there were multiple manifests for a project, there are now multiple items duplicated in
	// our list. The deduplicate() function removes these.
	return deduplicate(s)
}

// Helper function to create one project in Snyk at a time.
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

// Helper function to deduplicate strings in a slice of strings.
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

// Helper function that returns a new slice containing only items that exist in `a` but not in `b`.
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
