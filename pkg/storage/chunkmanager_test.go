package storage

import "testing"

// TestFilesFragAndUnfrag take 3 files and fragment every file in chunks
// and then convert the chunks on a new file, bassically making a copy of
// the file.
func TestFragAndUnfragFile(t *testing.T) {
	// New instances
	fileWriter := NewFileStorage()
	fragment := NewChunkManager()

	// Fragment multiple files and receiving an array of chunks
	chunks, err := fragment.FragmentData("~/go_logo1.jpg", "~/go_logo2.png", "~/go_logo3.png")
	if err != nil {
		t.Error(err)
	}
	// Now trying to defragment chunks
	files, err := fragment.DefragmentData(chunks, "~/Downloads/")
	if err != nil {
		t.Error(err)
	}
	// Writing a copy of every file
	for _, file := range files {
		fileWriter.WriteToFile(file.Path+"/"+file.Name+"_copy"+file.Extension, file.Data, false)
	}
}
