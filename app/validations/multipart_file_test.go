package validations

import (
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"github.com/esmailemami/eshop/consts"
)

func TestMultipartFileHeaderSizeValidator(t *testing.T) {
	t.Run("value is not *multipart.FileHeader", func(t *testing.T) {
		err := MultipartFileHeaderSizeValidator(5)("string type")
		if err.Error() != consts.InvalidFileHeader {
			t.Errorf(
				"peyvastFileSizeRule.Validate() wants error: %v got: %v",
				consts.InvalidFileHeader,
				err.Error(),
			)
		}
	})

	t.Run("file is greater than allowed size", func(t *testing.T) {

		fh := multipart.FileHeader{
			Filename: "file.pdf",
			Size:     1000 * 1024 * 1024,
		}

		err := MultipartFileHeaderSizeValidator(5)(&fh)
		if err == nil || !strings.Contains(err.Error(), "یشینه حجم مجاز") {
			t.Errorf(
				"peyvastFileSizeRule.Validate() wants error: %v got: %v",
				"یشینه حجم مجاز",
				err,
			)
		}
	})

	t.Run("ok", func(t *testing.T) {
		fh := multipart.FileHeader{
			Filename: "file.pdf",
			Size:     1 * 1024 * 1024,
		}

		err := MultipartFileHeaderSizeValidator(5)(&fh)
		if err != nil {
			t.Errorf("peyvastFileSizeRule.Validate() wants no error got: %v", err.Error())
		}
	})
}

func TestMultipartFileHeaderMimeTypeValidator(t *testing.T) {

	t.Run("value is not *multipart.FileHeader", func(t *testing.T) {
		err := MultipartFileHeaderMimeTypeValidator("application/pdf")("string type")
		if err.Error() != consts.InvalidFileContentType {
			t.Errorf(
				"peyvastFileMimeTypeRule.Validate() wants error: %v got: %v",
				consts.InvalidFileContentType,
				err.Error(),
			)
		}
	})

	t.Run("mime type in not allowed", func(t *testing.T) {
		fh := multipart.FileHeader{
			Filename: "file.txt",
			Size:     1 * 1024 * 1024,
		}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "not allowed type")

		err := MultipartFileHeaderMimeTypeValidator("application/pdf")(&fh)
		if err == nil || err.Error() != consts.InvalidPeyvastContentType {
			t.Errorf(
				"peyvastFileMimeTypeRule.Validate() wants error: %v got: %v",
				consts.InvalidPeyvastContentType,
				err,
			)
		}
	})

	t.Run("ok", func(t *testing.T) {

		fh := multipart.FileHeader{
			Filename: "file.pdf",
			Size:     1 * 1024 * 1024,
		}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "application/pdf")

		err := MultipartFileHeaderMimeTypeValidator("application/pdf")(&fh)
		if err != nil {
			t.Errorf("peyvastFileMimeTypeRule.Validate() wants no error got: %v", err)
		}
	})
}
