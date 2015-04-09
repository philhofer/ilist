package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
)

type T struct {
	next *T
	prev *T
}

var tmplFile *ast.File
var fileset *token.FileSet

func parseTemplate() error {
	fileset = token.NewFileSet()
	_, thisfile, _, _ := runtime.Caller(0)
	listfile := filepath.Join(filepath.Dir(thisfile), "list.go")
	var err error
	tmplFile, err = parser.ParseFile(fileset, listfile, nil, parser.ParseComments)
	return err
}

type idVisitor struct {
	oldName string
	newName string
}

func (i *idVisitor) Visit(node ast.Node) ast.Visitor {
	if id, ok := node.(*ast.Ident); ok {
		if id.Name == i.oldName {
			id.Name = i.newName
		}
	}
	return i
}

func renameID(name string, fs *ast.File) {
	for _, decl := range fs.Decls {
		ast.Walk(&idVisitor{
			oldName: "T",
			newName: name,
		}, decl)
		ast.Walk(&idVisitor{
			oldName: "TDequeue",
			newName: name + "List",
		}, decl)
	}
}

type structVisitor struct {
	name  string
	fail  error
	found bool
}

func (s *structVisitor) Visit(node ast.Node) ast.Visitor {
	if ts, ok := node.(*ast.TypeSpec); ok && ts.Name.Name == s.name {
		sd, ok := ts.Type.(*ast.StructType)
		if !ok {
			s.fail = fmt.Errorf("type %q is not a struct!", node)
		} else {
			ptrToSelf, _ := parser.ParseExpr("*" + s.name)
			needsFields(sd, ptrToSelf)
		}
		s.found = true
		return nil
	}
	return s
}

func needsFields(st *ast.StructType, typ ast.Expr) {
	foundnext := false
	foundprev := false
	for _, f := range st.Fields.List {
		switch len(f.Names) {
		case 2:
			if (f.Names[0].Name == "next" && f.Names[1].Name == "prev") ||
				(f.Names[0].Name == "prev" && f.Names[1].Name == "next") {
				return
			}
		case 1:
			switch f.Names[0].Name {
			case "next":
				foundnext = true
			case "prev":
				foundprev = true
			}
		default:
			continue
		}
	}
	if !foundnext {
		st.Fields.List = append(st.Fields.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("next")},
			Type:  typ,
		})
	}
	if !foundprev {
		st.Fields.List = append(st.Fields.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("prev")},
			Type:  typ,
		})
	}
}

func AddFieldsTo(filename string, typename string) error {
	fmt.Printf("Adding impl to type %q in file %q\n", typename, filename)
	fset := token.NewFileSet()
	fileast, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	tmplFile.Name = fileast.Name
	renameID(typename, tmplFile)
	visitor := &structVisitor{name: typename}
	for _, decl := range fileast.Decls {
		ast.Walk(visitor, decl)
		if visitor.fail != nil {
			return visitor.fail
		}
	}
	if !visitor.found {
		return fmt.Errorf("no type named %q found", typename)
	}
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	bwr := bufio.NewWriter(outfile)
	err = printer.Fprint(bwr, fset, fileast)
	if err == nil {
		err = bwr.Flush()
	}
	outfile.Close()
	return err
}

type pkgRenamer string

func (p pkgRenamer) Visit(node ast.Node) ast.Visitor {
	if pn, ok := node.(*ast.Package); ok {
		pn.Name = string(p)
		return nil
	}
	return p
}

func WriteImplFileTo(file string) error {
	ofile, err := os.Create(file)
	if err != nil {
		return err
	}
	bwr := bufio.NewWriter(ofile)
	err = printer.Fprint(bwr, fileset, tmplFile)
	if err == nil {
		err = bwr.Flush()
	}
	ofile.Close()
	return err
}
