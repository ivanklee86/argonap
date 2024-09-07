package cli

import (
	"github.com/ivanklee86/argonap/pkg/client"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type WorkerResult struct {
	
}

func worker(argocdClient *client.IArgoCDClient, appProject <-chan *v1alpha1.AppProject) {
	
	
}