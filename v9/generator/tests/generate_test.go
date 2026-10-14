package tests

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
				writeCanonicalSpec(t, root, test.spec)
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
	writeCanonicalSpec(t, root, validSpec)
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

func TestGenerationRejectsSpecDigestMismatch(t *testing.T) {
	root := generationFixture(t)
	writeFixture(t, filepath.Join(root, "generator/openapi/spec.json"), validSpec)
	writeFixture(t, filepath.Join(root, "generator/openapi/provenance.json"), `{"spec_version":"fixture-1","sha256":"`+strings.Repeat("0", 64)+`"}`)

	output, err := runGeneration(t, root)
	if err == nil || !strings.Contains(string(output), "does not match SHA-256") {
		t.Fatalf("spec digest mismatch was not rejected: %v, %s", err, output)
	}
	assertPreservedOutput(t, root)
}

func TestGenerationRejectsSpecVersionMismatch(t *testing.T) {
	root := generationFixture(t)
	writeCanonicalSpec(t, root, validSpec)
	digest := sha256.Sum256([]byte(validSpec))
	writeFixture(t, filepath.Join(root, "generator/openapi/provenance.json"), fmt.Sprintf(`{"spec_version":"other-version","sha256":"%x"}`, digest))

	output, err := runGeneration(t, root)
	if err == nil || !strings.Contains(string(output), "version does not match") {
		t.Fatalf("spec version mismatch was not rejected: %v, %s", err, output)
	}
	assertPreservedOutput(t, root)
}

func TestGenerationRequiresSpecProvenance(t *testing.T) {
	root := generationFixture(t)
	writeFixture(t, filepath.Join(root, "generator/openapi/spec.json"), validSpec)
	output, err := runGeneration(t, root)
	if err == nil || !strings.Contains(string(output), "Missing, unreadable or empty spec provenance") {
		t.Fatalf("missing provenance was not rejected: %v, %s", err, output)
	}
	assertPreservedOutput(t, root)
}

func TestInstallSpecPersistsVerifiedCanonicalSource(t *testing.T) {
	root := generationFixture(t)
	archivePath := filepath.Join(root, "release.tar.gz")
	writeSpecArchive(t, archivePath, "release/openapi.json", validSpec)
	digest := sha256.Sum256([]byte(validSpec))

	output, err := runSpecInstall(t, root, archivePath, "release/openapi.json", "fixture-1", hex.EncodeToString(digest[:]))
	if err != nil {
		t.Fatalf("install failed: %v: %s", err, output)
	}
	contents, err := os.ReadFile(filepath.Join(root, "generator/openapi/spec.json"))
	if err != nil || string(contents) != validSpec {
		t.Fatalf("canonical spec was not retained: %v", err)
	}
}

func TestInstallSpecRejectsBadHashOrVersionWithoutReplacingSource(t *testing.T) {
	digest := sha256.Sum256([]byte(validSpec))
	for _, test := range []struct {
		name       string
		version    string
		digest     string
		wantOutput string
	}{
		{name: "bad hash", version: "fixture-1", digest: strings.Repeat("0", 64), wantOutput: "does not match the supplied SHA-256"},
		{name: "wrong version", version: "fixture-2", digest: hex.EncodeToString(digest[:]), wantOutput: "does not match spec version"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := generationFixture(t)
			archivePath := filepath.Join(root, "release.tar.gz")
			writeSpecArchive(t, archivePath, "release/openapi.json", validSpec)
			canonicalPath := filepath.Join(root, "generator/openapi/spec.json")
			writeFixture(t, canonicalPath, sentinel)

			output, err := runSpecInstall(t, root, archivePath, "release/openapi.json", test.version, test.digest)
			if err == nil || !strings.Contains(string(output), test.wantOutput) {
				t.Fatalf("invalid source accepted or wrong failure: %v, %s", err, output)
			}
			contents, readErr := os.ReadFile(canonicalPath)
			if readErr != nil || string(contents) != sentinel {
				t.Fatalf("canonical source changed after failed install: %q, %v", contents, readErr)
			}
		})
	}
}

func TestGenerationPublishesValidOutput(t *testing.T) {
	root := generationFixture(t)
	writeCanonicalSpec(t, root, validSpec)
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
			writeCanonicalSpec(t, root, string(contents))
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

func runSpecInstall(t *testing.T, root, archive, member, version, digest string) ([]byte, error) {
	t.Helper()
	command := exec.CommandContext(t.Context(), "bash", "generator/scripts/install-spec.sh", archive, member, version, digest)
	command.Dir = root
	return command.CombinedOutput()
}

func writeCanonicalSpec(t *testing.T, root, contents string) {
	t.Helper()
	digest := sha256.Sum256([]byte(contents))
	writeFixture(t, filepath.Join(root, "generator/openapi/spec.json"), contents)
	writeFixture(t, filepath.Join(root, "generator/openapi/provenance.json"), fmt.Sprintf(`{"spec_version":"fixture-1","sha256":"%x"}`, digest))
}

func writeSpecArchive(t *testing.T, path, member, contents string) {
	t.Helper()
	var archive bytes.Buffer
	gzipWriter := gzip.NewWriter(&archive)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := tarWriter.WriteHeader(&tar.Header{Name: member, Mode: 0o600, Size: int64(len(contents)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write([]byte(contents)); err != nil {
		t.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, path, archive.String())
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
		if entry.IsDir() {
			if entry.Name() == "tests" || entry.Name() == "openapi" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.Contains(path, string(filepath.Separator)+"tests"+string(filepath.Separator)) {
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
