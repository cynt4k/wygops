package array

// ContainsEntry : Check if at least one string is in array
func ContainsEntry(g []string, s ...string) bool {
	for _, entry := range g {
		for _, matchEntry := range s {
			if matchEntry == entry {
				return true
			}
		}
	}
	return false
}
