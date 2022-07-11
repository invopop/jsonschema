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
)

func handleType(expr ast.Expr, breadcrumb string, comments map[string]string) {
	switch t := expr.(type) {
	case *ast.StructType:
		for _, field := range t.Fields.List {
			for _, name := range field.Names {
				if !ast.IsExported(name.Name) {
					continue
				}

				b := fmt.Sprintf("%s.%s", breadcrumb, name.Name)
				comments[b] = strings.TrimSpace(field.Doc.Text())
				handleType(field.Type, b, comments)
			}
		}
	case *ast.ArrayType:
		handleType(t.Elt, fmt.Sprintf("%s.[]", breadcrumb), comments)
	case *ast.MapType:
		handleType(t.Key, fmt.Sprintf("%s.[key]", breadcrumb), comments)
		handleType(t.Value, fmt.Sprintf("%s.[value]", breadcrumb), comments)
	case *ast.StarExpr:
		handleType(t.X, breadcrumb, comments)
	}
}

func getPackages(base, path string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	dict := make(map[string]*ast.Package)
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			for pkgName, v := range d {
				k := gopath.Join(base, gopath.Dir(path), pkgName)
				dict[k] = v
			}
		}
		return nil
	})
	return dict, err
}

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
func ExtractGoComments(base, path string, commentMap map[string]string) error {
	pkgs, err := getPackages(base, path)
	if err != nil {
		return err
	}

	for qualifiedName, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if d, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range d.Specs {
						if s, ok := spec.(*ast.TypeSpec); ok {
							if !ast.IsExported(s.Name.Name) {
								continue
							}

							breadcrumb := fmt.Sprintf("%s.%s", qualifiedName, s.Name.Name)
							txt := s.Doc.Text()
							if txt == "" {
								txt = d.Doc.Text()
							}
							commentMap[breadcrumb] = strings.TrimSpace(doc.Synopsis(txt))
							handleType(s.Type, breadcrumb, commentMap)
						}
					}
				}
			}
		}
	}
	return nil
}
