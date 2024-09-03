/*
The clienttest package contains helpers for testing.
*/
package testhelpers

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomProjectName() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	length := 20

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[random.Intn(len(charset))]
	}

	return fmt.Sprintf("pr%s", string(randomString))
}
