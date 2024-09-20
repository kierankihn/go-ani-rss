package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func DownloadFileFromUrl(url string, path string, size int64) error {
	// create output file
	outputDir := filepath.Dir(path)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer output.Close()

	// make http request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// progress bar
	bar := progressbar.DefaultBytes(size, "downloading")

	// write into target file
	io.Copy(io.MultiWriter(output, bar), response.Body)

	return nil
}
