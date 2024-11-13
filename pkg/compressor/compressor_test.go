package compressor

import (
	"github.com/alexismrosales/FileDriver/pkg/storage"
	"testing"
)

func TestFileCompressor(t *testing.T) {
	downloadPath, err := storage.GetShortPath("~/")
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{"~/go_logo1.jpg", "~/go_logo2.png", "~/go_logo3.png"}
	for i, path := range paths {
		newP, err := storage.GetShortPath(path)
		if err != nil {
			t.Fatal(err)
		}
		paths[i] = newP
	}
	err = ZipFile(downloadPath, "zipexample.zip", paths...)
	if err != nil {
		t.Fatal(err)
	}
}
