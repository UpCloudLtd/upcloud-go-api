package tests

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const validSpec = `{
	"openapi":"3.0.3",
	"info":{"title":"Generation test","version":"fixture-1"},
	"servers":[{"url":"https://api.upcloud.com"}],
	"tags":[{"name":"test"}],
	"paths":{"/test":{"get":{"tags":["test"],"operationId":"getTest",
		"responses":{"200":{"description":"OK","content":{"application/json":{
			"schema":{"type":"object","properties":{"name":{"type":"string"}}}
		}}}}
	}}}
}`

const sentinel = "sentinel\n"

func TestGenerationPreservesOutputOnInvalidInput(t *testing.T) {
	for _, test := range []struct {
		name string
		spec string
	}{
		{name: "missing spec"},
		{name: "empty spec", spec: "\n"},
		{name: "malformed JSON", spec: `{`},
		{name: "not OpenAPI", spec: `{}`},
		{name: "missing tags", spec: `{"openapi":"3.0.3","info":{"title":"Test","version":"1"},"paths":{}}`},
		{name: "unresolved reference", spec: `{"openapi":"3.0.3","info":{"title":"Test","version":"1"},"tags":[{"name":"test"}],"paths":{"/test":{"get":{"tags":["test"],"operationId":"getTest","responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Missing"}}}}}}}}}`},
		{name: "multiple documents", spec: validSpec + "\n" + validSpec},
		{name: "undeclared operation tag", spec: strings.Replace(validSpec, `"tags":[{"name":"test"}]`, `"tags":[{"name":"other"}]`, 1)},
		{name: "untagged operation", spec: strings.Replace(validSpec, `"tags":["test"]`, `"tags":[]`, 1)},
		{name: "multiple operation tags", spec: strings.Replace(strings.Replace(validSpec, `"tags":[{"name":"test"}]`, `"tags":[{"name":"test"},{"name":"other"}]`, 1), `"tags":["test"]`, `"tags":["test","other"]`, 1)},
		{name: "duplicate tags", spec: strings.Replace(validSpec, `[{"name":"test"}]`, `[{"name":"test"},{"name":"test"}]`, 1)},
		{name: "unsafe tag filename", spec: strings.ReplaceAll(validSpec, `"test"`, `"../outside"`)},
		{name: "tag trailing newline", spec: strings.ReplaceAll(validSpec, `"test"`, `"test\n"`)},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := generationFixture(t)
			if test.spec != "" {
				writeFixture(t, filepath.Join(root, "generator/spec.json"), test.spec)
			}
			output, err := runGeneration(t, root)
			if err == nil {
				t.Fatalf("invalid input accepted: %s", output)
			}
			assertPreservedOutput(t, root)
		})
	}
}

func TestGenerationPreservesOutputOnTemplateFailure(t *testing.T) {
	root := generationFixture(t)
	writeFixture(t, filepath.Join(root, "generator/spec.json"), validSpec)
	writeFixture(t, filepath.Join(root, "generator/config/templates/client.tmpl"), "{{invalid template syntax")
	output, err := runGeneration(t, root)
	if err == nil {
		t.Fatalf("broken template accepted: %s", output)
	}
	if !strings.Contains(string(output), "Generating test") {
		t.Fatalf("test did not reach tag generation: %s", output)
	}
	assertPreservedOutput(t, root)
}

func TestGenerationPublishesValidOutput(t *testing.T) {
	root := generationFixture(t)
	writeFixture(t, filepath.Join(root, "generator/spec.json"), validSpec)
	output, err := runGeneration(t, root)
	if err != nil {
		t.Fatalf("valid generation failed: %v: %s", err, output)
	}
	for _, name := range []string{"server_urls.gen.go", "test_models.gen.go", "test_client.gen.go"} {
		path := filepath.Join(root, "pkg/upcloud", name)
		if _, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors); parseErr != nil {
			t.Fatalf("invalid generated Go in %s: %v", name, parseErr)
		}
	}
	for _, name := range []string{"old_models.gen.go", "old_client.gen.go"} {
		if _, statErr := os.Stat(filepath.Join(root, "pkg/upcloud", name)); !os.IsNotExist(statErr) {
			t.Fatalf("stale output %s retained: %v", name, statErr)
		}
	}
	contents, err := os.ReadFile(filepath.Join(root, "pkg/upcloud/handwritten.go"))
	if err != nil || string(contents) != sentinel {
		t.Fatalf("handwritten file changed: %q, %v", contents, err)
	}
}

func TestGenerationRejectsInvalidInputShape(t *testing.T) {
	for _, test := range []struct {
		name  string
		key   string
		value any
	}{
		{name: "unsupported version", key: "openapi", value: "2.0"},
		{name: "missing info", key: "info", value: nil},
		{name: "empty paths", key: "paths", value: map[string]any{}},
		{name: "wrong paths type", key: "paths", value: []any{}},
		{name: "wrong tags type", key: "tags", value: "test"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var document map[string]any
			if err := json.Unmarshal([]byte(validSpec), &document); err != nil {
				t.Fatal(err)
			}
			document[test.key] = test.value
			contents, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			root := generationFixture(t)
			writeFixture(t, filepath.Join(root, "generator/spec.json"), string(contents))
			output, err := runGeneration(t, root)
			if err == nil || !strings.Contains(string(output), "Invalid OpenAPI input") {
				t.Fatalf("invalid shape accepted or wrong failure: %v, %s", err, output)
			}
			assertPreservedOutput(t, root)
		})
	}
}

func runGeneration(t *testing.T, root string) ([]byte, error) {
	t.Helper()
	staging := filepath.Join(root, "staging")
	if err := os.MkdirAll(staging, 0o755); err != nil {
		t.Fatal(err)
	}
	command := exec.CommandContext(t.Context(), "bash", "generator/scripts/generate.sh")
	command.Dir = root
	command.Env = append(os.Environ(), "TMPDIR="+staging)
	output, err := command.CombinedOutput()
	entries, readErr := os.ReadDir(staging)
	if readErr != nil || len(entries) != 0 {
		t.Fatalf("temporary files leaked: %v, %v; output: %s", entries, readErr, output)
	}
	return output, err
}

func assertPreservedOutput(t *testing.T, root string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, "pkg/upcloud"))
	if err != nil || len(entries) != 4 {
		t.Fatalf("unexpected output files: %v, %v", entries, err)
	}
	for _, name := range []string{"server_urls.gen.go", "old_models.gen.go", "old_client.gen.go", "handwritten.go"} {
		contents, readErr := os.ReadFile(filepath.Join(root, "pkg/upcloud", name))
		if readErr != nil || string(contents) != sentinel {
			t.Fatalf("%s changed on failure: %q, %v", name, contents, readErr)
		}
	}
}

func generationFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	sdkRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"go.mod", "go.sum"} {
		contents, readErr := os.ReadFile(filepath.Join(sdkRoot, name))
		if readErr != nil {
			t.Fatal(readErr)
		}
		writeFixture(t, filepath.Join(root, name), string(contents))
	}
	err = filepath.WalkDir(filepath.Join(sdkRoot, "generator"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || strings.Contains(path, string(filepath.Separator)+"tests"+string(filepath.Separator)) {
			return nil
		}
		relative, relativeErr := filepath.Rel(sdkRoot, path)
		if relativeErr != nil {
			return relativeErr
		}
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		writeFixture(t, filepath.Join(root, relative), string(contents))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"server_urls.gen.go", "old_models.gen.go", "old_client.gen.go", "handwritten.go"} {
		writeFixture(t, filepath.Join(root, "pkg/upcloud", name), sentinel)
	}
	return root
}

func writeFixture(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
