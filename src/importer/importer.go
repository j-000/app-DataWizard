package importer

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/j-000/app-DataWizard.git/utils"
)

/*
	Main functionality is to query data from the ATS source.

	Data can be in XML or JSON format.
	Querying can happen over HTTPS or SFTP

	HTTP + XML
	HTTP + JSON
	SFTP + XML
	SFTP + JSON
*/

type RequestParams struct {
	Url     string
	Headers map[string]string
}

// TODO Handle SFTP imports

func HTTP(args RequestParams, contentType string) (string, error) {
	switch contentType {
	case "xml":
		args.Headers["Accept"] = "application/xml"
	case "json":
		args.Headers["Accept"] = "application/json"
	default:
		return "", fmt.Errorf("contentType must be xml or json")
	}
	data, err := httpGet(args)
	if err != nil {
		return "", err
	}
	return data, nil
}

func httpGet(args RequestParams) (string, error) {
	log.Printf("[httpGet] GET url=%s\n", args.Url)
	request, err := http.NewRequest("GET", args.Url, nil)
	if err != nil {
		log.Printf("[httpGet] Error=%s \n", err.Error())
		return "", err
	}
	// Set default headers
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	// Set custom headers
	for key, value := range args.Headers {
		request.Header.Add(key, value)
	}
	log.Printf("[httpGet] Headers = %s\n", request.Header)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("[httpGet] Error: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	log.Printf("[httpGet] response statuscode= %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code = %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[httpGet] Error: %s", err.Error())
		return "", err
	}
	data := string(body)
	log.Printf("[httpGet] response data size=%d", len(data))
	log.Printf("[httpGet] response data (100 chars crop.)=%s...\n\n", utils.Substring(data, 100))
	return data, nil
}
