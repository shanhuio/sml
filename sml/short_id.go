package sml

func shortID(id string) string {
	n := len(id)
	for i, r := range id {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		n = i
		break
	}
	if n > 7 {
		n = 7
	}
	return id[:n]
}
