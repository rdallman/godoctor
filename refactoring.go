package doctor

// This file defines the Refactoring interface, the RefactoringBase struct, and
// several methods common to refactorings based on RefactoringBase, including
// SetSelection, GetLog, and GetResult.

import (
	//"code.google.com/p/go.tools/go/types"
	//"code.google.com/p/go.tools/importer"
	//"fmt"
	"go.tools/importer"
	"go/ast"
	"go/parser"
	"go/token"
)

// The Refactoring interface provides the methods common to all refactorings.
//
// Name returns a human-readable name for the refactoring, properly capitalized
// (e.g., "Rename" or "Extract Function").
//
// TODO: DOCUMENT REMAINING METHODS
type Refactoring interface {
	Name() string
	SetSelection(selection TextSelection) bool
	Run() bool
	GetLog() *Log
	GetResult() (*Log, EditSet)
}

type RefactoringBase struct {
	fset           *token.FileSet
	file           *ast.File
	selectionStart token.Pos
	selectionEnd   token.Pos
	selectedNode   ast.Node
	importer       *importer.Importer
	pkgInfo        *importer.PackageInfo
	log            *Log
	editSet        EditSet
}

// Configures a refactoring by indicating the filename in which text is
// selected and the beginning and end of the selected region.  Internally,
// this configures all of the fields in the RefactoringBase struct.
func (r *RefactoringBase) SetSelection(selection TextSelection) bool {
	r.log = NewLog()

	r.fset = token.NewFileSet()
	r.file = r.parse(selection.filename)
	if r.file == nil {
		return false
	}

	r.selectionStart = lineColToPos(r.fset, r.file,
		selection.startLine, selection.endLine)
	r.selectionEnd = lineColToPos(r.fset, r.file,
		selection.startLine, selection.endLine)

	r.importer = importer.New(new(importer.Config))
	r.pkgInfo = r.importer.CreatePackage(r.file.Name.Name, r.file)
	if r.pkgInfo.Err != nil {
		r.log.Log(FATAL_ERROR, "Analysis error: "+r.pkgInfo.Err.Error())
		return false
	}

	nodes, _ := importer.PathEnclosingInterval(r.file,
		r.selectionStart, r.selectionEnd)
	r.selectedNode = nodes[0]

	r.editSet = NewEditSet()
	return true
}

// Parses the given file, logging errors to the given log, and returning both
// a FileSet and a File
func (r *RefactoringBase) parse(filename string) *ast.File {
	f, err := parser.ParseFile(r.fset, filename, nil, parser.ParseComments)
	if err != nil {
		r.log.Log(FATAL_ERROR, "Error parsing "+filename+": "+
			err.Error())
	}
	return f
}

// Converts a line/column position (where the first character in a file is at
// line 1, column 1) into a token.Pos
func lineColToPos(fset *token.FileSet, astFile *ast.File, line int, column int) token.Pos {
	file := fset.File(astFile.Pos())
	lastLine := -1
	thisColumn := 1
	for i := 0; i < file.Size(); i++ {
		pos := file.Pos(i)
		thisLine := file.Line(pos)
		if thisLine != lastLine {
			thisColumn = 1
		} else {
			thisColumn++
		}
		if thisLine == line && thisColumn == column {
			return pos
		}
		lastLine = line
	}
	return astFile.Pos()
}

func (r *RefactoringBase) GetLog() *Log {
	return r.log
}

func (r *RefactoringBase) GetResult() (*Log, EditSet) {
	return r.log, r.editSet
}
