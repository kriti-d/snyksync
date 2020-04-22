## SnykSync

SnykSync is a simple runner written in Go that quickly compares the repositories in a given GitHub organization to the 
projects tracked in an organization in Snyk. If projects are found to exist in GitHub that are not tracked in
Snyk, the SnykSync tool will attempt to add them. It can be built and executed manually or this repo can be forked
and setup to run using GitHub Actions automatically on a scheduled job.

### Automatic Scanning Via Github Actions

Fork this repository to any org or account that has permission to run GitHub Actions. By default, the workflow
will fire at midnight and noon, UTC time. If you want the script to run more or less often, you can tweak
the `cron` line in the `.github/workflows/go.yml` file to the [cron timing](https://crontab.guru/) of your
preference. 

Before the first scheduled run, you must configure some secrets appropriate to your organization for the application
to work correctly. Open the settings tab for your fork of SnykSync and navigate to `Secrets`. Create each of the
following secrets as appropriate:

* `SS_GHTOKEN`: a personal GitHub access token configured at
[github.com/settings/tokens](https://github.com/settings/tokens) for the `Repo` scope.
* `SS_GHORG`: the name of the organization in GitHub to be scanned. User orgs are not supported at this time.
* `SS_SNYKTOKEN`: a personal Snyk access token configured at [app.snyk.io/account](https://app.snyk.io/account). This
_must_ be a personal token; service accounts are not supported at this time.
* `SS_SNYKORG`: the machine name of the organization in Snyk to add projects to (as seen in the organization URL.)
* `SS_SNYKINTID`: a UUID representing the GitHub integration Snyk will use to set up and scan, found at the very bottom of 
the GitHub integration page for your organization in Snyk.

SnykSync will discover all projects that the user has access to inside the GitHub organization specified. If any of these
projects does not have a corresponding project in Snyk, it will be added. Projects that do not exist in Snyk but do not
have any supported manifest files will still be attempted each time the job runs. This is purposeful as it ensures
projects will be added to Snyk as soon as they're supported.

### Manual Building and Running

If you want to use an automation tool outside of GitHub Actions, simply build the Go application, set environment
variables identical to the secrets in the previous section, and run the application.

```shell script
go build .

export SS_GHTOKEN="14c2ae81763830e98a819373cef15cd61df23c51"
export SS_GHORG="snyk"
export SS_SNYKTOKEN="e4c00eb6-136c-4cac-a037-12e7fcf90de3"
export SS_SNYKORG="my-special-org-name"
export SS_SNYKINTID="cdc0ca60-5e90-40d2-a625-ba4f5e9bd167"

./snyksync
```
