package jsonschema

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	gopath "path"
	"path/filepath"
	"strings"
)

type breadcrumb []string

func (b breadcrumb) With(breadcrumb string) breadcrumb {
	return append(b, breadcrumb)
}

func (b breadcrumb) Field(fieldName string) breadcrumb {
	return b.With(fieldName)
}

func (b breadcrumb) SliceElem() breadcrumb {
	return b.With("[]")
}

func (b breadcrumb) MapKey() breadcrumb {
	return b.With("[key]")
}

func (b breadcrumb) MapElem() breadcrumb {
	return b.With("[value]")
}

func (b breadcrumb) String() string {
	return strings.Join(b, ".")
}

func handleType(expr ast.Expr, breadcrumb breadcrumb, comments map[string]string) {
	switch t := expr.(type) {
	case *ast.StructType:
		for _, field := range t.Fields.List {
			for _, name := range field.Names {
				if !ast.IsExported(name.Name) {
					continue
				}

				b := breadcrumb.Field(name.Name)
				txt := field.Doc.Text()
				if txt == "" {
					txt = field.Comment.Text()
				}
				comments[b.String()] = strings.TrimSpace(txt)
				handleType(field.Type, b, comments)
			}
		}
	case *ast.ArrayType:
		handleType(t.Elt, breadcrumb.SliceElem(), comments)
	case *ast.MapType:
		handleType(t.Key, breadcrumb.MapKey(), comments)
		handleType(t.Value, breadcrumb.MapElem(), comments)
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
		rootBreadcrumb := breadcrumb{qualifiedName}
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if d, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range d.Specs {
						if s, ok := spec.(*ast.TypeSpec); ok {
							if !ast.IsExported(s.Name.Name) {
								continue
							}

							breadcrumb := rootBreadcrumb.With(s.Name.Name)
							txt := s.Doc.Text()
							if txt == "" {
								txt = d.Doc.Text()
							}
							commentMap[breadcrumb.String()] = strings.TrimSpace(doc.Synopsis(txt))
							handleType(s.Type, breadcrumb, commentMap)
						}
					}
				}
			}
		}
	}
	return nil
}
