package cli

import (
	"encoding/json"
	"os"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func readSyncWindowsFromFile(path string) ([]v1alpha1.SyncWindow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var syncWindowsToUpdate []v1alpha1.SyncWindow
	err = json.Unmarshal(data, &syncWindowsToUpdate)
	if err != nil {
		return nil, err
	}

	return syncWindowsToUpdate, nil
}
