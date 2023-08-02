package device42

import (
	"log"

	"github.com/go-resty/resty/v2"
)

func apiDevice42Get(client *resty.Client, path string, result interface{}) (*resty.Response, error) {
	log.Printf("[DEBUG] apiGet - Calling API on path %s", path)

	resp, err := client.R().
		SetResult(result).
		Get(path)

	if err != nil {
		log.Printf("[WARN] apiGet - Error in GET request for path %s", path)
		return nil, err
	}

	return resp, nil
}

func apiDevice42Post(client *resty.Client, path string, formData map[string]string, result interface{}) (*resty.Response, error) {
	log.Printf("[DEBUG] apiPost - Calling API on path %s", path)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(result).
		Post(path)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func apiDevice42Delete(client *resty.Client, path string, result interface{}) (*resty.Response, error) {
	log.Printf("[DEBUG] apiDelete - Calling API on path %s", path)

	resp, err := client.R().
		SetResult(result).
		Delete(path)

	if err != nil {
		log.Printf("[WARN] apiDelete - Error in DELETE request for path %s", path)
		return nil, err
	}

	return resp, nil
}
