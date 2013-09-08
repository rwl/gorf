// Copyright 2011 John Asmuth. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"fmt"
	//"go/token"
	"code.google.com/p/rog-go/exp/go/token"
	//"go/ast"
	"code.google.com/p/rog-go/exp/go/ast"
	"path/filepath"
	"strings"
	"code.google.com/p/rog-go/exp/go/types"
	"code.google.com/p/rog-go/exp/go/parser"
	//"go/parser"
)

var (
	AllSourceTops = token.NewFileSet()
	AllSources = token.NewFileSet()
	ImportedBy = make(map[string][]string)
	PackageTops = make(map[string]*ast.Package)
	Packages = make(map[string]*ast.Package)
)

func LocalImporter(path string) (pkg *ast.Package) {
	path = filepath.Clean(path)
	
	//fmt.Printf("Importing %s\n", path)
	var ok bool
	var pkgtop *ast.Package
	if pkgtop, ok = PackageTops[path]; !ok {
		pkg = types.DefaultImporter(path)
		return
	}
	if pkg, ok = Packages[path]; ok {
		return
	}
	var sourcefiles []string
	for srcfile := range pkgtop.Files {
		sourcefiles = append(sourcefiles, srcfile)
	}
	//fmt.Printf("Parsing %v\n", sourcefiles)
	dirpkgs, err := parser.ParseFiles(AllSources, sourcefiles, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	pkg = dirpkgs[pkgtop.Name]
	
	Packages[path] = pkg
	
	//fmt.Printf("nil: %v name: %s\n", pkg == nil, pkgtop.Name)
	
	return
}

func ScanAllForImports(dir string) (err error) {
	sw := ScanWalker{}
	filepath.Walk(dir, sw.Walk)
	err = sw.err
	
	return
}

func PreloadImportedBy(path string) {
	for _, ipath := range ImportedBy[path] {
		LocalImporter(ipath)
	}
}

type ScanWalker struct {
	err error
}

func (this *ScanWalker) Walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		this.VisitDir(path, &info)
	} else {
		this.VisitFile(path, &info)
	}
	return err
}

func (s *ScanWalker) VisitDir(path string, f *os.FileInfo) bool {
	err := ScanForImports(path)
	if err != nil {
		s.err = err
	}
	return true
}

func (s *ScanWalker) VisitFile(fpath string, f *os.FileInfo) {
	//if strings.HasSuffix(fpath, ".gorf") || strings.HasSuffix(fpath, ".gorfn") {
	//	os.Remove(fpath)
	//}
}

//Look at the imports, and build up ImportedBy
func ScanForImports(path string) (err error) {
	sourcefiles, _ := filepath.Glob(filepath.Join(path, "*.go"))
	dirpkgs, err := parser.ParseFiles(AllSourceTops, sourcefiles, parser.ImportsOnly)
	
	if err != nil {
		fmt.Println(err)
	}
	
	//take the first non-main. otherwise, main is ok.
	var prime *ast.Package
	for name, pkg := range dirpkgs {
		prime = pkg
		if name != "main" {
			break
		}
	}
	
	if prime == nil {
		return
	}
	
	PackageTops[path] = prime
		
	is := make(ImportScanner)
	
	ast.Walk(is, prime)
	
	if v, ok := is["."]; !v && ok {
		return MakeErr("gorf can not deal with unnamed import in '%s'", path)
	}
	
	for path, _ := range is {
		if strings.HasPrefix(path, ".") {
			return MakeErr("gorf can not deal with relative import in '%s'", path)
		}
	}
	
	for imp := range is {
		ImportedBy[imp] = append(ImportedBy[imp], path)
	}
	
	return
}

type ImportScanner map[string]bool

func (is ImportScanner) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:
		if n.Name != nil && n.Name.Name == "." {
			is["."] = false
		}
		is[string(n.Path.Value)] = true
	}
	return is
}
