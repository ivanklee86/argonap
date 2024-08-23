package cli

import (
	"testing"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMapSubset(t *testing.T) {
	bigMap := make(map[string]string)
	bigMap["a"] = "x"
	bigMap["b"] = "y"
	bigMap["c"] = "z"

	matchingMap := make(map[string]string)
	matchingMap["a"] = "x"
	matchingMap["b"] = "y"

	notMatchingMap := make(map[string]string)
	notMatchingMap["d"] = "u"

	differentValueMap := make(map[string]string)
	differentValueMap["a"] = "c"

	t.Run("Map contains a submap", func(t *testing.T) {
		assert.Equal(t, true, isMapSubset(bigMap, matchingMap))
	})

	t.Run("Map does not contains a submap", func(t *testing.T) {
		assert.Equal(t, false, isMapSubset(bigMap, notMatchingMap))
		assert.Equal(t, false, isMapSubset(bigMap, differentValueMap))
	})
}

func TestProjectFilter(t *testing.T) {
	labelsProjectA := make(map[string]string)
	labelsProjectA["team"] = "Jets"
	labelsProjectA["env"] = "prod"
	labelsProjectA["department"] = "A"

	labelsProjectB := make(map[string]string)
	labelsProjectB["team"] = "Giants"
	labelsProjectB["env"] = "prod"
	labelsProjectB["department"] = "B"

	labelsProjectC := make(map[string]string)
	labelsProjectC["team"] = "Eagles"
	labelsProjectC["env"] = "dev"
	labelsProjectC["department"] = "C"

	projectsArray := []v1alpha1.AppProject{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "ProjectA",
				Labels: labelsProjectA,
			},
			Spec: v1alpha1.AppProjectSpec{
				SyncWindows: v1alpha1.SyncWindows{
					{
						Kind:       "allow",
						Schedule:   "10 1 * * *",
						Duration:   "1h",
						Namespaces: []string{"*"},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "ProjectB",
				Labels: labelsProjectB,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "ProjectC",
				Labels: labelsProjectC,
			},
		},
	}

	projectsList := v1alpha1.AppProjectList{
		Items: projectsArray,
	}

	t.Run("filterProjects returns all projects if no labels are specified", func(t *testing.T) {
		tags := make(map[string]string)
		assert.Len(t, filterProjects(&projectsList, "", tags, false), 3)
	})

	t.Run("filterProjects returns a matching projects if labels are specified", func(t *testing.T) {
		tags := make(map[string]string)
		tags["team"] = "Eagles"
		projectResult := filterProjects(&projectsList, "", tags, false)
		assert.Len(t, projectResult, 1)
		assert.Equal(t, projectResult[0].ObjectMeta.Name, "ProjectC")
	})

	t.Run("filterProjects returns some matching projects if labels are specified", func(t *testing.T) {
		tags := make(map[string]string)
		tags["env"] = "prod"
		assert.Len(t, filterProjects(&projectsList, "", tags, false), 2)
	})

	t.Run("filterProjects returns proejcts with only SyncWindows", func(t *testing.T) {
		tags := make(map[string]string)
		tags["env"] = "prod"
		projectResult := filterProjects(&projectsList, "", tags, true)
		assert.Len(t, projectResult, 1)
		assert.Equal(t, projectResult[0].ObjectMeta.Name, "ProjectA")
	})

	t.Run("filterProjects returns named project if specified", func(t *testing.T) {
		tags := make(map[string]string)
		assert.Len(t, filterProjects(&projectsList, "ProjectA", tags, false), 1)
	})
}
