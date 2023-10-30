package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/urfave/cli"
)

// Particle function returns the CLI commands for gort particle
func Particle() cli.Command {
	return cli.Command{
		Name:  "particle",
		Usage: "Upload sketches to your Particle Photon or Electron",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"upload"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  gort particle upload <accessToken> <deviceId> <default|tinker-servo|path name> # uploads sketch to Particle Photon or Electron")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "upload":
				if len(c.Args()) < 4 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				accessToken := c.Args()[1]
				deviceId := c.Args()[2]
				fileName := c.Args()[3]
				url := fmt.Sprintf("https://api.particle.io/v1/devices/%v?access_token=%v", deviceId, accessToken)
				extraParams := map[string]string{}
				request, err := newfileUploadRequest(url, extraParams, "file", fileName)
				if err != nil {
					log.Fatal(err)
				}
				client := &http.Client{}
				resp, err := client.Do(request)
				if err != nil {
					log.Fatal(err)
				} else {
					body := &bytes.Buffer{}
					_, err := body.ReadFrom(resp.Body)
					if err != nil {
						log.Fatal(err)
					}
					resp.Body.Close()
					fmt.Println(resp.StatusCode)
					fmt.Println(resp.Header)
					fmt.Println(body)
				}

			}
		},
	}
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	data, fileName, err := openUploadFile(path)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return nil, err
	}
	file := bytes.NewReader(data)
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("PUT", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

func openUploadFile(filePath string) ([]byte, string, error) {
	if filePath == "default" || filePath == "tinker-servo" {
		fileName := "tinker-servo.ino"
		filePath = fmt.Sprintf("support/particle/%v", fileName)
		data, err := Asset(filePath)
		if err != nil {
			return nil, "", err
		}
		return data, fileName, nil
	}

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, "", err
	}
	data := make([]byte, 65535)
	file.Read(data)
	return data, path.Base(filePath), nil
}
