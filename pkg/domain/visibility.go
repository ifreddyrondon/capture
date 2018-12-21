package domain

var visibilityTypes = [...]string{"public", "private"}

type Visibility string

func AllowedVisibility(test string) bool {
	if test == "" {
		return false
	}

	for i := range visibilityTypes {
		if visibilityTypes[i] == test {
			return true
		}
	}

	return false
}

var (
	Public  Visibility = "public"
	Private Visibility = "private"
)
