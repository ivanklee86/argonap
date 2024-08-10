package argocd_client

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IArgoCDClient interface {
	CreateProject(context.Context, string) (*v1alpha1.AppProject, error)
}

type ArgoCDClientOptions struct {
	ServerAddr string
	Insecure   bool
	AuthToken  string
}

type ArgoCDClient struct {
	ProjectClient project.ProjectServiceClient
}

func NewArgoCDClient(aco *ArgoCDClientOptions) (IArgoCDClient, error) {
	apiClient, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: aco.ServerAddr,
		Insecure: aco.Insecure,
		AuthToken: aco.AuthToken,
	})
	if err != nil {
		return nil, err
	}

	_, projectClient, err := apiClient.NewProjectClient()
	if err != nil {
		return nil, err
	}

	return &ArgoCDClient{
		ProjectClient: projectClient,
	}, nil
}

func (c *ArgoCDClient) CreateProject(ctx context.Context, name string) (*v1alpha1.AppProject, error) {
	return c.ProjectClient.Create(ctx, &project.ProjectCreateRequest{
		Project: &v1alpha1.AppProject{
			ObjectMeta: v1.ObjectMeta{
				Name: name,
			},
		},
	})
}