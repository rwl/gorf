// Copyright 2011 John Asmuth. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"path/filepath"
	"os"
	//"go/ast"
	"code.google.com/p/rog-go/exp/go/ast"
	//"go/printer"
	"code.google.com/p/rog-go/exp/go/printer"
)

func Copy(srcpath, dstpath string) (err error) {
	var srcFile *os.File
	srcFile, err = os.Open(srcpath)
	if err != nil {
		return
	}

	var dstFile *os.File
	dstFile, err = os.Create(dstpath)
	if err != nil {
		return
	}

	io.Copy(dstFile, srcFile)

	dstFile.Close()
	srcFile.Close()

	return
}

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath);
	return err == nil
}

func BackupSource(fpath string) (err error) {
	dir, name := filepath.Split(fpath)
	backup := filepath.Join(dir, "."+name+".0.gorf")
	if !FileExists(backup) {
		err = Copy(fpath, backup)
	}
	return
}

func Touch(fpath string) (err error) {
	f, err := os.Create(fpath)
	f.Close()
	return
}

func MoveSource(oldpath, newpath string) (err error) {
	fmt.Printf("Moving %s to %s\n", oldpath, newpath)
	
	ndir, nfile := filepath.Split(newpath)
	nmarker := filepath.Join(ndir, "."+nfile+".0.gorfn")
	
	if FileExists(newpath) {
		BackupSource(newpath)
	}
	
	odir, ofile := filepath.Split(newpath)
	if !FileExists(filepath.Join(odir, "."+ofile+".0.gorfn")) {
		err = BackupSource(oldpath)
		if err != nil {
			return
		}
	}
	
	err = os.MkdirAll(ndir, 0755)
	if err != nil {
		return
	}
	
	err = Touch(nmarker)
	if err != nil {
		return
	}
	
	err = Copy(oldpath, newpath)
	if err != nil {
		return
	}
	
	err = os.Remove(oldpath)
	
	return
}

func NewSource(fpath string, file *ast.File) (err error) {
	fmt.Printf("Creating %s\n", fpath)
	
	dir, name := filepath.Split(fpath)
	
	err = Touch(filepath.Join(dir, "."+name+".0.gorfn"))
	if err != nil {
		return
	}
	
	var out io.Writer
	out, err = os.Create(fpath)
	if err != nil {
		return
	}
	
	err = printer.Fprint(out, AllSources, file)
	
	return
}

func RewriteSource(fpath string, file *ast.File) (err error) {
	fmt.Printf("Rewriting %s\n", fpath)

	err = BackupSource(fpath)
	if err != nil {
		return
	}

	var out io.Writer
	out, err = os.Create(fpath)
	if err != nil {
		return
	}
	
	err = printer.Fprint(out, AllSources, file)
	
	return
}
