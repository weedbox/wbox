package lib

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type GolangTemplate struct {
	filepath string
	fset     *token.FileSet
	node     *ast.File
}

func OpenGolangTemplate(filepath string) (*GolangTemplate, error) {

	// Open the file
	src, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	gt := &GolangTemplate{
		filepath: filepath,
		fset:     fset,
		node:     node,
	}

	return gt, nil
}

func (t *GolangTemplate) SetPackageName(packageName string) error {
	t.node.Name.Name = ToSnakeCase(packageName)
	return nil
}

func (t *GolangTemplate) SetConstValue(constName string, kind token.Token, value string) error {

	v := value
	if kind == token.STRING {
		v = `"` + value + `"`
	}

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.CONST {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if name.Name == constName {
							if len(valueSpec.Values) > i {
								valueSpec.Values[i] = &ast.BasicLit{
									Kind:  kind,
									Value: v,
								}
							}
						}
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameType(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == source {
						typeSpec.Name.Name = target
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameFunction(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			if decl.Name.Name == source {
				decl.Name.Name = target
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameReceiver(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			if decl.Recv != nil {
				for _, field := range decl.Recv.List {
					if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == source {
						ident.Name = target
					} else if starExpr, ok := field.Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == source {
							ident.Name = target
						}
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameFunctionResult(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			if decl.Type.Results != nil {
				for _, field := range decl.Type.Results.List {
					if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == source {
						ident.Name = target
					} else if starExpr, ok := field.Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == source {
							ident.Name = target
						}
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameVariableType(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if starExpr, ok := valueSpec.Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							if ident.Name == source {
								ident.Name = target
							}
						}
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameAllocationType(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if expr, ok := n.(*ast.CompositeLit); ok {
			if ident, ok := expr.Type.(*ast.Ident); ok {
				if ident.Name == source {
					ident.Name = target
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameFunctionCall(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				if ident.Name == source {
					ident.Name = target
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) RenameFunctionResultInCallExpr(source string, target string) error {

	ast.Inspect(t.node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			for _, arg := range callExpr.Args {
				if funcLit, ok := arg.(*ast.FuncLit); ok {
					if funcLit.Type.Results != nil {
						for _, result := range funcLit.Type.Results.List {
							if starExpr, ok := result.Type.(*ast.StarExpr); ok {
								if ident, ok := starExpr.X.(*ast.Ident); ok {
									if ident.Name == source {
										ident.Name = target
									}
								}
							}

						}
					}
				}
			}
		}

		return true
	})

	return nil
}

func (t *GolangTemplate) Save() error {

	outFile, err := os.Create(t.filepath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := printer.Fprint(outFile, t.fset, t.node); err != nil {
		return err
	}

	return nil
}
