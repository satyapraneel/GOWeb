package arrays

func InArray(key string, haystack []string) bool {
	for _, value := range haystack {
		if value == key {
			return true
		}
	}
	return false
}
