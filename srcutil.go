// Package srcutil provides utilities for working with Go source code. The Go
// standard library provides a powerful suite of packages "go/{ast,doc,...}"
// which are used by the Go tool chain to compile Go programs. As you initially
// try to find your way around you hit a small dependency barrier and have to
// learn a small portion of each package. There is a fantastic write up and
// collection of examples that I used to learn (or shamelessly copy pasta'd)
// while creating this package, currently maintained by:
//
//   Alan Donovan (https://github.com/golang/example/tree/master/gotypes)
//
// In the mean time this package can help you get started with some common use
// cases.
package srcutil

import (
	"fmt"
	"go/build"
	"go/parser"
	"os"
	"path/filepath"
)

var (

	// DefaultContext is a Context configured as if you were in the directory of
	// the source code calling this library and working normally. It uses the
	// currently configured GOROOT, GOPATH and the source dir is set to your
	// current working directory.
	DefaultContext = FromWorkDir()

	// DefaultParseMode adds ParseComments to the default parser.Mode used in the
	// go/parser package.
	DefaultParseMode = parser.ParseComments

	// DefaultImportMode is the same as the default build.ImportMode used in the
	// go/build package.
	DefaultImportMode = build.ImportMode(0)
)

// FromDir returns a Context configured with the SourceDir set to the given dir.
// It uses the default build.ImportMode and build.Default for build.Context and
// your GOROOT and GOPATH from the environment to determine the location of
// packages and the standard library.
func FromDir(fromDir string) *Context {
	return &Context{Context: build.Default, SourceDir: fromDir}
}

// FromStandard returns a Context configured to only contain the Go standard
// library, to do this it simply excludes your GOPATH and sets the SourceDir
// to the GOROOT/src.
func FromStandard() *Context {
	ctx := FromDir(``)
	ctx.GOPATH = ``
	ctx.SourceDir = filepath.Join(ctx.GOROOT, "src")
	return ctx
}

// FromWorkDir is like FromDir except it sets the Source dir to your working
// directory. It is used as the DefaultContext.
func FromWorkDir() *Context {
	return FromDir(defaultToGetwd())
}

// Context is not something you need to interact with for common use cases,
// instead calling the top level functions that return Package types directly
// which will use the DefaultContext.
//
// This structure is like build.Context except it includes the ImportMode and a
// SourceDir that will default to the current working directory. It sits at the
// top of this packages dependency hierarchy as it loads the top level useful
// object, Package.
type Context struct {
	build.Context

	// SourceDir defines where the code lives relative for operations you perform.
	// For example if you call Context.Import("reflect") it will attempt to import
	// that as if the calling package was SourceDir. So if a vendor/ directory
	// existing within SourceDir with a reflect package that would be imported
	// instead.
	SourceDir string
}

// String implements fmt.Stringer.
func (c *Context) String() string {
	return fmt.Sprintf("Context(%s -> %s)", c.SourceDir, c.GOROOT)
}

func defaultToGetwd(srcDirs ...string) string {
	for _, srcDir := range srcDirs {
		if len(srcDir) > 0 {
			return srcDir
		}
	}
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	return dir
}

// Import will behave just like using a import declaration from code residing
// within the SourceDir of this Context. If you did not explicitly set the
// SourceDir it will use your current working directory.
func (c *Context) Import(pkgName string) (*Package, error) {
	buildPkg, err := c.Context.Import(pkgName, defaultToGetwd(c.SourceDir), DefaultImportMode)
	if err != nil {
		return nil, err
	}
	return &Package{Package: *buildPkg}, nil
}
