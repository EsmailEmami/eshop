package str

import (
	"strings"

	"github.com/esmailemami/eshop/services/numeric"
)

// تبدیل کاراکتر های عربی به فارسی
// این تابع اعداد را به اعداد انگلیسی تبدیل میکنید
func ArToFa(ar string) string {
	s := strings.Trim(ar, " ")
	s = strings.ReplaceAll(s, "ي", "ی")
	s = strings.ReplaceAll(s, "ك", "ک")
	s = strings.ReplaceAll(s, "د", "د")

	s = numeric.TransformFa2En(s)

	return s
}
