package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyFile(sourcePath string, destPath string) (err error) {
	var stats os.FileInfo
	stats, err = os.Stat(sourcePath)
	if err != nil {
		return
	}

	var blob []byte
	blob, err = ioutil.ReadFile(sourcePath)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(destPath, blob, stats.Mode())
	return err
}

func downloadFromUrl(dirName string, url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(dirName + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName
}

func gortDirName() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + string(os.PathSeparator) + "gort"
}

func createGortDirectory() (string, error) {
	dirName := gortDirName()
	fileExists, err := exists(dirName)
	if fileExists {
		fmt.Println("Gort lives")
	} else {
		os.Mkdir(dirName, 0777)
	}
	return dirName, err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func supportDir(support string) (dir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	dir = usr.HomeDir + "/.gort/support/" + support
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
		}
	}
	return
}
