package cli

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func IsMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}

func filterProjects(appProjects *v1alpha1.AppProjectList, labels map[string]string) []v1alpha1.AppProject {
	matchingProjects := []v1alpha1.AppProject{}

	for _, appProject := range appProjects.Items {
		if IsMapSubset(appProject.ObjectMeta.Labels, labels) {
			matchingProjects = append(matchingProjects, appProject)
		}
	}

	return matchingProjects
}
