package util

func Getip(env, appName *string) ([]string, error) {
	return []string{"10.0.0.11"}, nil
}

func Contains(s *[]string, str *string) bool {
	for _, v := range *s {
		if v == *str {
			return true
		}
	}
	return false
}
