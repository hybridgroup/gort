package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"bytes"
  "io"
  "log"
  "mime/multipart"
  "net/http"
)

func Spark() cli.Command {
	return cli.Command{
		Name:  "spark",
		Usage: "Upload sketches to your Spark",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"upload"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("  gort spark upload [accessToken] [deviceId] [default|voodoospark|path name] # uploads sketch to Spark")
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
				url := fmt.Sprintf("https://api.spark.io/v1/devices/%v", deviceId)

			  extraParams := map[string]string{}
			  request, err := newfileUploadRequest(url, extraParams, "file", fileName)
			  if err != nil {
			      log.Fatal(err)
			  }
			  request.Header.Set("Authorization", "Bearer " + accessToken)
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
	fmt.Println(path)
	data, err := openUploadFile(path)
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile(paramName, path)
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

func openUploadFile(path string)([]byte, error) {
	filePath := path
	if filePath == "default" || filePath == "voodoospark" {
		filePath = fmt.Sprintf("support/spark/%v.cpp", filePath)
		fmt.Println(filePath)
		data, err := Asset(filePath)
  	if err != nil {
      return nil, err
 		}
		return data, nil					
	} else {
		file, err := os.Open(filePath)
	  defer file.Close()
		if err != nil {
    	return nil, err
		}
	 	data := make([]byte, 4096)
	 	file.Read(data)
		return data, nil
	}
}