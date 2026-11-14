package device42

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressCommaSeparatedOrderDiff suppresses diffs where the only difference is
// the ordering of a comma-separated list (e.g. tags/groups). Device42 returns
// such lists in an arbitrary order, so "a,b" and "b,a" are semantically equal.
func suppressCommaSeparatedOrderDiff(k, old, new string, d *schema.ResourceData) bool {
	if old == new {
		return true
	}
	normalize := func(s string) string {
		if s == "" {
			return ""
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		sort.Strings(parts)
		return strings.Join(parts, ",")
	}
	return normalize(old) == normalize(new)
}

func apiDevice42Get(client *resty.Client, path string, result interface{}) (*resty.Response, error) {
	log.Printf("[DEBUG] apiGet - Calling API on path %s", path)

	resp, err := client.R().
		SetResult(result).
		Get(path)

	if err != nil {
		log.Printf("[WARN] apiGet - Error in GET request for path %s", path)
		return nil, err
	}

	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		log.Printf("[WARN] apiGet - GET %s returned HTTP %d", path, resp.StatusCode())
		return nil, fmt.Errorf("apiGet - HTTP %d: %s", resp.StatusCode(), resp.String())
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

	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		log.Printf("[WARN] apiDelete - DELETE %s returned HTTP %d", path, resp.StatusCode())
		return nil, fmt.Errorf("apiDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	return resp, nil
}
