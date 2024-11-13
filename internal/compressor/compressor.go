package compressor

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ZipFile(destinyPath string, zipName string, paths ...string) error {
	// Creating a new zip file
	file, err := os.Create(filepath.Join(destinyPath, zipName))
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	for _, path := range paths {
		err = addFileToZip(zipWriter, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(writer *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	//Create a header of the file inside the ZIP
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// Crear a writing for the file in the ZIP
	writerFile, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	//Copy content to the zip
	_, err = io.Copy(writerFile, file)
	return err
}
