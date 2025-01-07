package acceptance_tests

import "math/rand"

const (
	// charSetAlpha are lower case alphabet letters.
	charSetAlpha = "abcdefghijklmnopqrstuvwxyz"

	// Length of the resource name we wish to generate.
	resourceNameLength = 10
)

func randomString() string {
	result := make([]byte, resourceNameLength)
	for i := 0; i < resourceNameLength; i++ {
		result[i] = charSetAlpha[rand.Intn(len(charSetAlpha))]
	}
	return string(result)
}

func randomStringWithPrefix(prefix string) string {
	return prefix + randomString()
}
