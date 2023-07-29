package validations

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/esmailemami/eshop/consts"
)

// MultipartFileHeaderSizeValidator بیشینه مجاز برای فایل آپلود شده را اعتباری سنجی میکند
func MultipartFileHeaderSizeValidator(maxAllowedMB int64) func(value interface{}) error {
	return func(value interface{}) error {
		v, ok := value.(*multipart.FileHeader)
		if !ok {
			return errors.New(consts.InvalidFileHeader)
		}
		if v.Size > (maxAllowedMB * 1024 * 1024) {
			return fmt.Errorf("بیشینه حجم مجاز برای فایل پیوست %d مگابایت است.", maxAllowedMB)
		}
		return nil
	}
}

// MultipartFileHeaderMimeTypeValidator نوع فایل آپلود شده را اعتبار سنجی میکند که از لیست مجاز باشد
func MultipartFileHeaderMimeTypeValidator(allowedTypes []string) func(value interface{}) error {
	return func(value interface{}) error {
		v, ok := value.(*multipart.FileHeader)
		if !ok {
			return errors.New(consts.InvalidFileContentType)
		}
		ct := v.Header.Get("Content-Type")
		for _, mt := range allowedTypes {
			if ct == mt {
				return nil
			}
		}
		return errors.New(consts.InvalidFileContentType)
	}
}
