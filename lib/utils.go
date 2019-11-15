package lib

func ContainsString(target string, source []string) bool {
	for _, elem := range source {
		if elem == target {
			return true
		}
	}

	return false
}
