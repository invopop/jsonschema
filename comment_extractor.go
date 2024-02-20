package jsonschema

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	gopath "path"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

// ExtractGoComments will read all the go files contained in the provided path,
// including sub-directories, in order to generate a dictionary of comments
// associated with Types and Fields. The results will be added to the `commentsMap`
// provided in the parameters and expected to be used for Schema "description" fields.
//
// The `go/parser` library is used to extract all the comments and unfortunately doesn't
// have a built-in way to determine the fully qualified name of a package. The `base` paremeter,
// the URL used to import that package, is thus required to be able to match reflected types.
//
// When parsing type comments, we use the `go/doc`'s Synopsis method to extract the first phrase
// only. Field comments, which tend to be much shorter, will include everything.
func ExtractGoComments(base, rootPath string, commentMap map[string]string) error {
	root, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	dict := make(map[string][]*ast.Package)
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// key should consist of a base path + the path relative to the root
		// so that we can treat rootPath of "../" properly
		k := gopath.Join(base, strings.TrimPrefix(path, root))
		// paths may have multiple packages, like for tests
		dict[k] = append(dict[k], maps.Values(d)...)

		return nil
	})
	if err != nil {
		return err
	}

	for pkg, p := range dict {
		for _, f := range p {
			gtxt := ""
			typ := ""
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					typ = x.Name.String()
					if !ast.IsExported(typ) {
						typ = ""
					} else {
						txt := x.Doc.Text()
						if txt == "" && gtxt != "" {
							txt = gtxt
							gtxt = ""
						}
						txt = doc.Synopsis(txt)
						commentMap[fmt.Sprintf("%s.%s", pkg, typ)] = strings.TrimSpace(txt)
					}
				case *ast.Field:
					txt := x.Doc.Text()
					if txt == "" {
						txt = x.Comment.Text()
					}
					if typ != "" && txt != "" {
						for _, n := range x.Names {
							if ast.IsExported(n.String()) {
								k := fmt.Sprintf("%s.%s.%s", pkg, typ, n)
								commentMap[k] = strings.TrimSpace(txt)
							}
						}
					}
				case *ast.GenDecl:
					// remember for the next type
					gtxt = x.Doc.Text()
				}
				return true
			})
		}
	}

	return nil
}
