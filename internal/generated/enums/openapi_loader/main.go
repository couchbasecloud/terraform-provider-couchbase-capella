// Fetches the OpenAPI spec from the Couchbase docs page and writes it to disk.
// The spec is embedded as JSON in a <script> tag; we parse the HTML, extract
// the JSON object, and write it as-is. JSON is a strict subset of YAML, so
// kin-openapi (used by gen-api and gen-enums) loads the result without
// conversion.
//
// Fetch logic is kept in sync with internal/docs/openapi_loader.go.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const specURL = "https://docs.couchbase.com/cloud/management-api-reference/index.html"

func main() {
	data, err := fetch()
	if err != nil {
		log.Fatalf("fetch failed: %v", err)
	}

	f, err := os.OpenFile("openapi.generated.yaml", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		log.Fatalf("open failed: %v", err)
	}
	n, err := f.Write(data)
	if cerr := f.Close(); cerr != nil && err == nil {
		err = cerr
	}
	if err != nil {
		log.Fatalf("write failed: %v", err)
	}

	fmt.Fprintf(os.Stderr, "wrote %d bytes to openapi.generated.yaml\n", n)
}

func fetch() ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(specURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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
