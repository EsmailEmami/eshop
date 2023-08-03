package validations

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/esmailemami/eshop/consts"
)

// MultipartFileHeaderSizeValidator validates the maximum allowed size for the uploaded file.
func MultipartFileHeaderSizeValidator(maxAllowedMB int64) func(value interface{}) error {
	return func(value interface{}) error {
		v, ok := value.(*multipart.FileHeader)
		if !ok {
			return errors.New(consts.InvalidFileHeader)
		}
		if v.Size > (maxAllowedMB * 1024 * 1024) {
			return fmt.Errorf("The maximum allowed attachment file size is %d megabytes.", maxAllowedMB)
		}
		return nil
	}
}

// MultipartFileHeaderMimeTypeValidator validates the uploaded file's type against an allowed list.
func MultipartFileHeaderMimeTypeValidator(allowedTypes ...string) func(value interface{}) error {
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
