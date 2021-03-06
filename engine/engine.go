// Copyright 2014 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package engine is the programmatic entrypoint to the Go refactoring engine.
package engine

import (
	"fmt"

	"github.com/godoctor/godoctor/refactoring"
)

// All available refactorings, keyed by a unique, one-word, all-lowercase name
var refactorings map[string]refactoring.Refactoring

// All available refactorings' keys, in the order the refactorings should be
// displayed in a menu presented to the end user
var refactoringsInOrder []string

func init() {
	refactorings = map[string]refactoring.Refactoring{}
	refactoringsInOrder = []string{}

	AddRefactoring("rename", new(refactoring.Rename))
	AddRefactoring("toggle", new(refactoring.ToggleVar))
	AddRefactoring("godoc", new(refactoring.AddGoDoc))
	AddRefactoring("debug", new(refactoring.Debug))
	AddRefactoring("null", new(refactoring.Null))
}

// AllRefactoringNames returns the short names of all refactorings in an
// order suitable for display in a menu.
func AllRefactoringNames() []string {
	return refactoringsInOrder
}

// GetRefactoring returns a Refactoring keyed by the given short name.  The
// short name must be one of the keys in the map returned by AllRefactorings.
func GetRefactoring(shortName string) refactoring.Refactoring {
	return refactorings[shortName]
}

// AddRefactoring allows custom refactorings to be added to the refactoring
// engine.  Invoke this method before starting the command line or protocol
// driver.
func AddRefactoring(shortName string, newRefac refactoring.Refactoring) error {
	if r, ok := refactorings[shortName]; ok {
		return fmt.Errorf("The short name \"%s\" is already "+
			"associated with a refactoring (%s)",
			shortName,
			r.Description().Name)
	}
	refactorings[shortName] = newRefac
	refactoringsInOrder = append(refactoringsInOrder, shortName)
	return nil
}
