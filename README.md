# argonap

[![CI](https://github.com/ivanklee86/argonap/actions/workflows/ci.yaml/badge.svg)](https://github.com/ivanklee86/argonap/actions/workflows/ci.yaml) [![codecov](https://codecov.io/gh/ivanklee86/argonap/graph/badge.svg?token=KEWN2E756X)](https://codecov.io/gh/ivanklee86/argonap)

CLI to make 🐙 take a quick 💤

aka

ArgoCD [SyncWindows](https://argo-cd.readthedocs.io/en/stable/user-guide/sync_windows/) are great for addressing those whacky situations that somehow pop up in real life:
- Holidays
- Failovers
- Maintainance
- Emergencies where you just want to run lots and lots of `kubectl` commands
- Some (or all) of the above!🤣

`argonap` allows you to create and clear SyncWindows across multiple projects from the comfort of the command line.

## Installation

### Homebrew
```sh
brew tap ivanklee86/homebrew-tap
brew install ivanklee86/tap/argonap
```

### Docker Image
```sh
docker run -it --rm ghcr.io/ivanklee86/argonap:latest
```

### Go
```sh
go install github.com/ivanklee86/argonap@latest
```

## Authentication

`argonap` uses a JWT to authenticate to ArgoCD.  This can be configured in the Helm chartas follows:

```YAML
configs:
  cm:
    accounts.YOUR_ACCOUNT_NAME: apiKey

  rbac:
    policy.csv: |
      p, role:argonap, projects, get, *, allow
      p, argonap, projects, update, *, allow
      g, YOUR_ACCOUNT_NAME, role:argonap
```

A JWT can then generated using the ArgoCD CLI using the following command:
```shell
argocd login # Using username/password or SSO
argocd account generate-token --account YOUR_ACCOUNT
```

`argonap` flags can be set via environment variables with the format `ARGONAP_[flag but replace - with _]` e.g. `ARGONAP_AUTH_TOKEN`.  This allows you to store the auth token securely and pass it to the CLI using your favorite local secrets solution (e.g. [1Password CLI](https://developer.1password.com/docs/cli/secret-references))

## Selection

Projects can be selected by the following CLI options:
- `--name` will cause `argonap` to only make changes to the target AppProject.
- `--label` will only select AppProjects where all labels are matched.  Labels should be in format `key=value` and be supplied multiple times.

Passing no options will run the command on all projects.

## Usage

### `set`

The `set` command takes a file that contain a list SyncWindows to add to projects matching the selection criteria.  The file should be a JSON file containing items that match the [SyncWindows struct](https://pkg.go.dev/github.com/argoproj/argo-cd@v1.8.7/pkg/apis/application/v1alpha1#SyncWindow).

Example file:
```json
[
    {
        "kind": "deny",
        "schedule": "00 3 * * *",
        "duration": "1h",
        "namespaces": ["*"]
    }
]
```

#### Help
```
Set SyncWindows from file

Usage:
  argonap set [flags]

Flags:
      --file string   Path to file with SyncWindows to configure
  -h, --help          help for set

Global Flags:
      --auth-token string       JWT Authentication Token
      --insecure                Don't validate SSL certificate on client request
      --label strings           Labels to filter projects on in format 'key=value'.  Can be used multiple times.
      --name strings            Project names to update.  If specified, label filtering will not apply.  Can be used multiple times.
      --server-address string   ArgoCD server address
      --timeout int             Context timeout in seconds. (default 240)
      --workers int             # of parallel workers. (default 4)
```

### `clear`

The `clear` command removes **all** SyncWindows on projects matching selection criteria.


#### Help
```
Clear SyncWindows on all AppProjects.

Usage:
  argonap clear [flags]

Flags:
  -h, --help   help for clear

Global Flags:
      --auth-token string       JWT Authentication Token
      --insecure                Don't validate SSL certificate on client request
      --label strings           Labels to filter projects on in format 'key=value'.  Can be used multiple times.
      --name strings            Project names to update.  If specified, label filtering will not apply.  Can be used multiple times.
      --server-address string   ArgoCD server address
      --timeout int             Context timeout in seconds. (default 240)
      --workers int             # of parallel workers. (default 4)
```
