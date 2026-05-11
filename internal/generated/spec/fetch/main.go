// Fetches the OpenAPI spec from OPENAPI_SPEC_URL and writes it to disk.
// The spec is embedded as JSON in a <script> tag on the Couchbase docs page;
// we parse the HTML, extract the JSON object, and write it as-is. JSON is a
// strict subset of YAML, so kin-openapi (used by gen-api and gen-enums) loads
// the result without conversion.
//
// Fetch logic is kept in sync with internal/docs/openapi_loader.go.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const defaultURL = "https://docs.couchbase.com/cloud/management-api-reference/index.html"

func main() {
	out := flag.String("o", "openapi.generated.yaml", "output path")
	flag.Parse()

	url := os.Getenv("OPENAPI_SPEC_URL")
	if url == "" {
		url = defaultURL
	}

	data, err := fetch(url)
	if err != nil {
		log.Fatalf("fetch %s: %v", url, err)
	}

	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatalf("write %s: %v", *out, err)
	}

	fmt.Fprintf(os.Stderr, "wrote %d bytes to %s\n", len(data), *out)
}

func fetch(url string) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(http.DetectContentType(data), "text/html") {
		return extractEmbeddedSpec(data)
	}
	return data, nil
}

func extractEmbeddedSpec(htmlData []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var find func(*html.Node) string
	find = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
			if idx := strings.Index(n.FirstChild.Data, `{"openapi":"3.0`); idx != -1 {
				spec, _ := extractJSONObject(n.FirstChild.Data[idx:])
				return spec
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if spec := find(c); spec != "" {
				return spec
			}
		}
		return ""
	}

	if spec := find(doc); spec != "" {
		return []byte(spec), nil
	}
	return nil, fmt.Errorf("no embedded OpenAPI spec found in HTML")
}

func extractJSONObject(s string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(s))
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return "", err
	}
	return string(raw), nil
}
