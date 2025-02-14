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

func (t *GolangTemplate) SetConstValue(constName string, kind token.Token, value string) *ast.File {

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
