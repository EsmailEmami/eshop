package file

import (
	"os"
	"testing"
)

func TestGetFileExetension(t *testing.T) {
	type testcase struct {
		filename string
		wanted   string
	}
	tests := []testcase{
		{"myfile.txt", "txt"},
		{"myfile.tar.gz", "gz"},
		{"myfile", ""},
	}
	for _, tc := range tests {
		ext := GetFileExetension(tc.filename)
		if ext != tc.wanted {
			t.Errorf("GetFileExetension() wants: %v got: %v", tc.wanted, ext)
		}
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	t.Run("no same file", func(t *testing.T) {
		dir := os.TempDir() + "/sg"
		os.RemoveAll(dir)
		got := GenerateUniqueFilename(dir, "image.jpg", 1)

		if got != "image.jpg" {
			t.Errorf("GenerateUniqueFilename() want: %v got: %v", "image.jpg", got)
		}
	})
	t.Run("file with same name exists", func(t *testing.T) {
		dir := os.TempDir() + "/sg"
		os.RemoveAll(dir)

		_ = os.MkdirAll(dir, os.ModePerm)
		_, _ = os.Create(dir + "/image.jpg")

		got := GenerateUniqueFilename(dir, "image.jpg", 1)

		if got != "image-001.jpg" {
			t.Errorf("GenerateUniqueFilename() want: %v got: %v", "image-001.jpg", got)
		}
	})

	t.Run("some files with same name exist", func(t *testing.T) {
		dir := os.TempDir() + "/sg"
		os.RemoveAll(dir)

		_ = os.MkdirAll(dir, os.ModePerm)
		_, _ = os.Create(dir + "/image.jpg")
		_, _ = os.Create(dir + "/image-001.jpg")
		_, _ = os.Create(dir + "/image-002.jpg")

		got := GenerateUniqueFilename(dir, "image.jpg", 1)

		if got != "image-003.jpg" {
			t.Errorf("GenerateUniqueFilename() want: %v got: %v", "image-003.jpg", got)
		}
	})
}
