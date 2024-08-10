package main

import (
	"context"
	"fmt"

	argocd_client "github.com/ivanklee/argocd_client/pkg/client"
)

func main() {
	argocdClientOptions := argocd_client.ArgoCDClientOptions{
        ServerAddr: "localhost:8080",
        Insecure: true,
        AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhdXRvbWF0aW9uOmFwaUtleSIsIm5iZiI6MTcyMzI3MTg4OCwiaWF0IjoxNzIzMjcxODg4LCJqdGkiOiI5ZmQ4YWIwMi02YzhkLTRmY2EtODhjMi1lYWRlZGZiNjFmYzUifQ.fBpbZ7-RmOfINHQehl8BJuMWa87M5j46hFbaTbiIqiE",
    }

    client, err := argocd_client.NewArgoCDClient(&argocdClientOptions)
    if err != nil {
        panic(err)
    }

    project, err := client.CreateProject(context.Background(), "test")
    if err != nil {
        panic(err)
    }

    fmt.Sprint("%s", project)
}
