package main

import (
	"cmd/client/internal/compressor"
	"cmd/client/internal/protocol"
	"cmd/client/internal/storage"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// download
func download(currentPath string, paths []string, windowSize int, conn *net.UDPConn, addr *net.UDPAddr) error {
	// Getting full path from every file
	for i := range paths {
		sPath, err := safePath(filepath.Join(currentPath, paths[i]))
		if err != nil {
			return err
		}
		paths[i] = sPath
	}
	cachePath, err := storage.GetShortPath("~/.filedriverCache")
	if err != nil {
		return err
	}
	zipName := createZipFileName()
	err = compressor.ZipFile(cachePath, zipName, paths...)
	if err != nil {
		return err
	}

	chunks, err := storage.NewChunkManager().FragmentData(filepath.Join(cachePath, zipName))
	if err != nil {
		return err
	}
	err = protocol.SendMessage(conn, &protocol.Message{Type: protocol.MsgDownload, SegmentSize: len(chunks)}, addr)
	if err != nil {
		return err
	}

	return protocol.SelectiveRejectSend(chunks, windowSize, conn, addr)
}

func upload(currentPath string, windowSize int, conn *net.UDPConn) error {
	msgServerAns, _, err := protocol.ReceiveMessage(conn)
	if err != nil {
		return err
	}
	// Using protocol to receive all chuks of data from packets
	chunks, err := protocol.SelectiveRejectReceive(windowSize, msgServerAns.SegmentSize, conn)
	if err != nil {
		return err
	}

	sPath, err := safePath(currentPath)
	if err != nil {
		return err
	}

	files, err := storage.NewChunkManager().DefragmentData(chunks, sPath)
	if err != nil {
		return err
	}

	fileWriter := storage.NewFileStorage()
	for _, file := range files {
		path := filepath.Join(currentPath, file.Name+file.Extension)
		sPath, _ := safePath(path)
		fileWriter.WriteToFile(sPath, file.Data, false)
	}
	return nil

}

func pwd(currentPath string) string {
	return filepath.Join(baseDir, currentPath)
}

func mkdir(currentPath string, paths ...string) error {
	var sPath string
	var err error
	for _, path := range paths {
		// In case there is a path expansion else just add the currentPath
		if len(path) >= 2 && path[:2] == "~/" {
			sPath = getBaseShortPath(path)
		} else {
			sPath = filepath.Join(currentPath, path)
			if err != nil {
				return errors.New("mkdir: joining path directories.")
			}
		}
		sPath, err = safePath(sPath)
		err = createDir(sPath)
		if err != nil {
			return errors.New("mkdir: creating directories.")
		}
	}
	return nil
}

func cd(currentPath string, path string) (string, error) {
	// If there is a path with expansion symbol
	if len(path) >= 2 && path[:2] == "~/" {
		return getBaseShortPath(path), nil
	}
	// Join the path in the current path will help to clean the path succesfully
	fullPath := filepath.Join(currentPath, path)
	fullPath = filepath.Clean(fullPath)
	// In case the path does not exist or is not valid is an error
	if !pathIsValid(fullPath) {
		return fullPath, errors.New("cd: no such file or directory:")
	}
	return fullPath, nil
}

func ls(currentPath string, flags []string) (string, error) {
	var allEntries []string
	sPath, err := safePath(currentPath)
	fmt.Println("currentPathLS:", sPath)
	entries, err := os.ReadDir(sPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if len(flags) == 0 && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		allEntries = append(allEntries, entry.Name())
	}

	return strings.Join(allEntries, " "), nil
}

func mv(currentPath string, paths ...string) error {
	fullPath, err := safePath(currentPath)
	if err != nil {
		return err
	}
	destinyDir := filepath.Join(fullPath, paths[len(paths)-1])
	for i := 0; i < len(paths)-1; i++ {
		sPath := filepath.Join(fullPath, paths[i])
		destinyPath := filepath.Join(destinyDir, paths[i])
		err := os.Rename(sPath, destinyPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func rm(currentPath string, paths, flags []string) error {
	fullPath, err := safePath(currentPath)
	if err != nil {
		return err
	}
	flagR, flagF := "", ""
	for _, flag := range flags {
		if flag == "r" {
			flagR = flag
		} else if flag == "f" {
			flagF = flag
		}
	}

	for _, path := range paths {
		rmPath := filepath.Join(fullPath, filepath.Clean(path))
		info, err := os.Stat(rmPath)
		if err != nil {
			return err
		}
		if info.IsDir() && flagR != "" && flagF != "" {
			err = os.RemoveAll(rmPath)
			if err != nil {
				return err
			}
		} else if !info.IsDir() {
			err = os.Remove(rmPath)
			if err != nil {
				return err
			}
		} else {
			return errors.New("rm: cannot remove '" + path + "': Is a directory")
		}
	}
	return nil
}

func pathIsValid(path string) bool {
	fullPath := filepath.Join(baseDir, getBaseShortPath(path))
	_, err := os.Stat(fullPath)
	return os.IsNotExist(err)
}
