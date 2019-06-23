package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func testFilePath(parts ...string) string {
	_, filename, _, _ := runtime.Caller(0)
	path := []string{filepath.Dir(filename)}
	path = append(path, parts...)
	return filepath.Join(path...)
}

func TestUnmarshallResponse(t *testing.T) {
	jf := testFilePath("test-data", "response_OK.json")
	f, err := os.Open(jf)
	if err != nil {
		t.Fatalf("can't open test file: '%s'", jf)
	}
	defer f.Close()
	d := json.NewDecoder(f)

	var resp Response
	err = d.Decode(&resp)
	if err != nil {
		t.Fatalf("can't unmarshall response: %s", err)
	}
	if l := len(resp.Data); l != 2 {
		t.Errorf("expected %d data objects, got %d", 2, l)
	}
}
