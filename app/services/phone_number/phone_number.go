package phone_number

import "regexp"

func IsMobileNumber(number string) bool {
	r := regexp.MustCompile(`^09(\d){9}$`)
	return r.MatchString(number)
}
