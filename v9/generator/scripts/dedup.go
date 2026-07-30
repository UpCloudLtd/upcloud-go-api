package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
)

// This script runs after oapi-codegen to remove duplicate declarations.
// We generate one model file per API tag, and schemas shared across tags end up
// declared multiple times. This script keeps the first occurrence of each
// type, const, or method and rewrites only files that changed.
const generatedModelsPattern = "./pkg/upcloud/*_models.gen.go"

func main() {
	files, err := filepath.Glob(generatedModelsPattern)
	if err != nil {
		fmt.Printf("find generated model files: %v\n", err)
		os.Exit(1)
	}

	s := make(seen)
	for _, file := range files {
		if err := removeDuplicates(file, s); err != nil {
			fmt.Printf("%s: %v\n", file, err)
			os.Exit(1)
		}
	}
}

// seen tracks package-level names already emitted across files. Package-level
// identifiers share one namespace, so one map covers types, consts, and methods.
type seen map[string]seenEntry

type seenEntry struct {
	file string
	kind string
}

// add records name as seen in file. Returns false (and logs) if it was already seen.
func (s seen) add(kind, name, file string) bool {
	if e, dup := s[name]; dup {
		fmt.Printf("Removing duplicate %s %s from %s (already defined in %s)\n", kind, name, file, e.file)
		return false
	}
	s[name] = seenEntry{file, kind}
	return true
}

func removeDuplicates(file string, s seen) error {
	fileSet := token.NewFileSet()

	parsedFile, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	kept, changed := s.keepUnique(file, parsedFile.Decls)
	if !changed {
		return nil
	}

	parsedFile.Decls = kept
	return writeFormattedGoFile(file, fileSet, parsedFile)
}

func (s seen) keepUnique(file string, decls []ast.Decl) ([]ast.Decl, bool) {
	changed := false
	kept := make([]ast.Decl, 0, len(decls))

	for _, decl := range decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			// Only type and const declarations can collide in our generated model files.
			keep, declChanged := s.keepGenDecl(file, typedDecl)
			if keep {
				kept = append(kept, typedDecl)
			}
			if declChanged {
				changed = true
			}
		case *ast.FuncDecl:
			// Methods such as "NetworkPeeringState.Valid" can also be duplicated
			// when the enum type is generated in more than one tag file.
			name := typedDecl.Name.Name
			if typedDecl.Recv != nil && len(typedDecl.Recv.List) > 0 {
				name = receiverTypeName(typedDecl.Recv.List[0].Type) + "." + name
			}
			if s.add("method", name, file) {
				kept = append(kept, typedDecl)
			} else {
				changed = true
			}
		default:
			kept = append(kept, decl)
		}
	}

	return kept, changed
}

func (s seen) keepGenDecl(file string, decl *ast.GenDecl) (bool, bool) {
	originalSpecCount := len(decl.Specs)

	// A single type/const block can contain multiple names. So we keep the block if any name survives.
	switch decl.Tok {
	case token.TYPE:
		decl.Specs = slices.DeleteFunc(decl.Specs, func(spec ast.Spec) bool {
			ts, ok := spec.(*ast.TypeSpec)
			return ok && !s.add("type", ts.Name.Name, file)
		})
		return len(decl.Specs) > 0, len(decl.Specs) != originalSpecCount
	case token.CONST:
		var changed bool
		decl.Specs, changed = s.keepConstSpecs(file, decl.Specs)
		return len(decl.Specs) > 0, changed || len(decl.Specs) != originalSpecCount
	default:
		return true, false
	}
}

func (s seen) keepConstSpecs(file string, specs []ast.Spec) ([]ast.Spec, bool) {
	kept := make([]ast.Spec, 0, len(specs))
	changed := false

	for _, spec := range specs {
		constSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			kept = append(kept, spec)
			continue
		}

		names, values, namesChanged := s.keepConstNames(file, constSpec)
		if namesChanged {
			changed = true
		}
		if len(names) == 0 {
			continue
		}

		constSpec.Names = names
		constSpec.Values = values
		kept = append(kept, constSpec)
	}

	return kept, changed
}

func (s seen) keepConstNames(file string, spec *ast.ValueSpec) ([]*ast.Ident, []ast.Expr, bool) {
	keptNames := make([]*ast.Ident, 0, len(spec.Names))
	keptValues := make([]ast.Expr, 0, len(spec.Values))
	changed := false

	// Constants can be written as either "A, B = 1, 2" or as an enum-style
	// block where the type/value is inherited from the previous line. Only
	// filter values when each name has its own explicit value.
	explicitValues := len(spec.Values) == len(spec.Names)

	for i, name := range spec.Names {
		if !s.add("const", name.Name, file) {
			changed = true
			continue
		}
		keptNames = append(keptNames, name)
		if explicitValues {
			keptValues = append(keptValues, spec.Values[i])
		}
	}

	if !explicitValues {
		keptValues = spec.Values
	}

	return keptNames, keptValues, changed
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.IndexExpr:
		return receiverTypeName(t.X)
	case *ast.IndexListExpr:
		return receiverTypeName(t.X)
	default:
		return ""
	}
}

func writeFormattedGoFile(file string, fileSet *token.FileSet, parsedFile *ast.File) error {
	var buffer bytes.Buffer
	if err := format.Node(&buffer, fileSet, parsedFile); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	if err := os.WriteFile(file, buffer.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
