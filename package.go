package srcutil

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// A Package here has the same meaning as in Go. It embeds a build.Package and
// provides methods to centralize some of the common operations for working with
// Go source code and makes it easy to create some of the lower level compiler
// toolchain primitives like ast, types and parser packages. It should be safe
// for concurrent use from multiple Goroutines.
//
// You should not create Package values with composite literals, instead use one
// of the functions in this package so it may be initialized safely.
type Package struct {
	build.Package
	once sync.Once
	tc   *toolchain
}

// Import is shorthand for FromWorkDir().Import("pkgname").
func Import(pkgName string) (*Package, error) {
	return FromWorkDir().Import(pkgName)
}

type toolchain struct {
	fileSet   *token.FileSet
	astPkg    *ast.Package
	docPkg    *doc.Package
	typesPkg  *types.Package
	typesInfo *types.Info
}

// Synopsis implements fmt.Stringer.
func (p *Package) Synopsis() string {
	p.init()
	return doc.Synopsis(p.tc.docPkg.Doc)
}

// String implements fmt.Stringer.
func (p *Package) String() string {
	p.init()
	return fmt.Sprintf("srcutil.Package{%s}", p.Name)
}

// ToAst provides access to an associated pair of *token.FileSet and
// ast.Package. A new pair is created each call and a nil pointer will be
// returned when error is non-nil.
func (p *Package) ToAst() (*token.FileSet, *ast.Package, error) {
	tc, err := p.toToolchain(nil)
	if err != nil {
		return nil, nil, err
	}
	return tc.fileSet, tc.astPkg, nil
}

// ToDoc provides access to a *doc.Package. A new *doc.Package will be created
// each call and a nil pointer will be returned when error is non-nil.
func (p *Package) ToDoc() (*doc.Package, error) {
	tc, err := p.toToolchain(nil)
	if err != nil {
		return nil, err
	}
	return tc.docPkg, nil
}

func (p *Package) typesInfo() *types.Info {
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	return info
}

// ToTypes provides access to a *types.Package. A new *types.Package will be
// created each call and a nil pointer will be returned when error is non-nil.
func (p *Package) ToTypes() (*types.Package, error) {
	tc, err := p.toToolchain(nil)
	if err != nil {
		return nil, err
	}
	return tc.typesPkg, nil
}

// ToInfo is like ToTypes but also returns a *types.Info that contains all the
// Info maps declared and ready to query.
func (p *Package) ToInfo() (*types.Info, *types.Package, error) {
	tc, err := p.toToolchain(p.typesInfo())
	if err != nil {
		return nil, nil, err
	}
	return tc.typesInfo, tc.typesPkg, nil
}

// Docs groups the documentation related methods.
type Docs struct {
	Package *Package
}

// Docs returns a Docs struct to perform common operations related to
// documentation using the go/doc
func (p *Package) Docs() Docs {
	p.init()
	return Docs{p}
}

func (d *Docs) indirectValues(s []*doc.Value) []doc.Value {
	out := make([]doc.Value, len(s))
	for i := range s {
		out[i] = *s[i]
	}
	return out
}

// Examples returns a slice of doc.Example for each declared Go example.
func (d *Docs) Examples() []doc.Example {
	astFiles := d.Package.astFiles(d.Package.tc.astPkg)
	s := doc.Examples(astFiles...)
	out := make([]doc.Example, len(s))
	for i := range s {
		out[i] = *s[i]
	}
	return out
}

// Notes returns all marked comments starting with "MARKER(uid): note body."
// as described in the go/doc package. I.E.:
//   // TODO(cstockton): Fix this.
//   // BUG(cstockton): Broken.
func (d *Docs) Notes() map[string][]doc.Note {
	m := d.Package.tc.docPkg.Notes
	out := make(map[string][]doc.Note)
	for k, ns := range m {
		out[k] = make([]doc.Note, len(ns))
		for i := range ns {
			out[k][i] = *ns[i]
		}
	}
	return out
}

// Consts returns declared constants in the go/doc package style, which
// groups by the entire const ( Const1 = 1, Const2 = .. ) blocks.
func (d *Docs) Consts() []doc.Value {
	return d.indirectValues(d.Package.tc.docPkg.Consts)
}

// Types returns a slice of doc.Type representing exported functions.
func (d *Docs) Types() []doc.Type {
	s := d.Package.tc.docPkg.Types
	out := make([]doc.Type, len(s))
	for i := range s {
		out[i] = *s[i]
	}
	return out
}

// Methods returns declared methods of doc.Func types grouped in a map of
// string type names.
func (d *Docs) Methods() map[string][]doc.Func {
	types := d.Types()
	out := make(map[string][]doc.Func)
	for _, typ := range types {
		out[typ.Name] = make([]doc.Func, len(typ.Funcs))
		for i := range typ.Funcs {
			out[typ.Name][i] = *typ.Funcs[i]
		}
	}
	return out
}

// Vars returns declared variables in the go/doc package style, which groups
// the by the var ( Var1 = 1, Var2 = .. ) blocks.
func (d *Docs) Vars() []doc.Value {
	return d.indirectValues(d.Package.tc.docPkg.Vars)
}

// Funcs returns a slice of doc.Func representing exported functions.
func (d *Docs) Funcs() []doc.Func {
	s := d.Package.tc.docPkg.Funcs
	out := make([]doc.Func, len(s))
	for i := range s {
		out[i] = *s[i]
	}
	return out
}

// Function groups a types.Func and types.Signature, it will never be part of a
// method so Recv() will always be nul.
type Function struct {
	*types.Func
	*types.Signature
}

// String implements fmt.Stringer.
func (f Function) String() string {
	return f.Func.String()
}

// NewFunction returns a Function, typeFunc must not be nil.
func NewFunction(typeFunc *types.Func) Function {
	// funcs always have signatures
	return Function{typeFunc, typeFunc.Type().(*types.Signature)}
}

// Functions returns all the packages named functions from the packages outer
// most scope.
func (p *Package) Functions() []Function {
	p.init()
	var funcs []Function
	scope := p.tc.typesPkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() || isTest(name, "Test") || isTest(name, "Example") {
			continue
		}
		asFunc, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		funcs = append(funcs, NewFunction(asFunc))
	}
	return funcs
}

// Var groups a types.Func and types.Signature.
type Var struct {
	*types.Func
	*types.Signature
}

// MethodSet represents a set of methods belonging to a named type.
type MethodSet struct {
	Name    string
	Obj     types.Object
	Methods map[string]Function
}

// NewMethodSet returns a initialized MethodSet.
func NewMethodSet(name string, obj types.Object) MethodSet {
	return MethodSet{
		Name:    name,
		Obj:     obj,
		Methods: make(map[string]Function),
	}
}

// Names returns the names of the methods in this MethodSet.
func (m MethodSet) Names() []string {
	var names []string
	for k := range m.Methods {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Len returns the current number of Method's within this MethodSet.
func (m MethodSet) Len() int {
	return len(m.Methods)
}

// Methods returns a map keyed off of the name type with a value of MethodSet.
// Only types with at least one method are included.
func (p *Package) Methods() map[string]MethodSet {
	p.init()
	methods := make(map[string]MethodSet)
	scope := p.tc.typesPkg.Scope()

	for _, name := range scope.Names() {
		methodSet, err := p.MethodSet(name)
		if err != nil {
			continue
		}
		if methodSet.Len() > 0 {
			methods[name] = methodSet
		}
	}
	return methods
}

// MethodSet returns the set of methods for the given name.
func (p *Package) MethodSet(name string) (MethodSet, error) {
	p.init()
	obj := p.tc.typesPkg.Scope().Lookup(name)
	if obj == nil {
		return MethodSet{}, fmt.Errorf("named type was not found")
	}
	if !obj.Exported() || isTest(name, "Test") || isTest(name, "Example") {
		return MethodSet{}, fmt.Errorf("named type was not exported")
	}

	typ := obj.Type()
	ms := NewMethodSet(name, obj)
	for _, t := range []types.Type{typ, types.NewPointer(typ)} {
		mset := types.NewMethodSet(t)
		for i := 0; i < mset.Len(); i++ {
			z := mset.At(i)
			f, ok := z.Obj().(*types.Func)
			if !ok {
				continue // must be *Var field selection
			}
			ms.Methods[f.Name()] = NewFunction(f)
		}
	}
	return ms, nil
}

// init is called for you by all functions and methods that return a Package
// type, init will be ran only once within a sync.Once, multiple calls are safe.
func (p *Package) init() (err error) {
	p.once.Do(func() {
		p.tc, err = p.toToolchain(p.typesInfo())
	})
	return err
}

// toToolchain is used to initialize the package for usage.
func (p *Package) toToolchain(typesInfo *types.Info) (*toolchain, error) {
	tc := &toolchain{}
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, p.Dir, nil, DefaultParseMode)
	if err != nil {
		return tc, err
	}

	// @TODO I'm not sure the best way to Copy() an ast. There may be a utility
	// func somewhere but I couldn't find it, this will have to do for now.
	//   -> doc.New takes ownership of the AST pkg and may edit or overwrite it.
	docFileSet := token.NewFileSet()
	docPkgs, err := parser.ParseDir(docFileSet, p.Dir, nil, DefaultParseMode)
	if err != nil {
		return tc, err
	}

	astPkg, okAst := pkgs[p.Name]
	docAstPkg, okDoc := docPkgs[p.Name]
	if !okAst || !okDoc {
		return tc, fmt.Errorf(
			`unable to find pkg "%s" in the "%s" directory`, p.Name, p.Dir)
	}
	docPkg := doc.New(docAstPkg, p.Dir, doc.Mode(0))
	astFiles := p.astFiles(astPkg)
	conf := types.Config{Importer: importer.Default()}
	typesPkg, err := conf.Check(p.Name, fileSet, astFiles, typesInfo)
	if err != nil {
		return tc, err
	}

	tc.fileSet, tc.astPkg, tc.docPkg, tc.typesPkg, tc.typesInfo =
		fileSet, astPkg, docPkg, typesPkg, typesInfo
	return tc, nil
}

func (p *Package) astFiles(astPkg *ast.Package) (out []*ast.File) {
	for key := range astPkg.Files {
		out = append(out, astPkg.Files[key])
	}
	return
}

// Exact check for a test func string from:
//   https://golang.org/src/cmd/go/test.go
//
// isTest tells whether name looks like a test (or benchmark, according to prefix).
// It is a Test (say) if there is a character after Test that is not a lower-case letter.
// We don't want TesticularCancer.
//               ^
//     lol      /
//             /
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

// Files groups operations on a packages files.
type Files struct {
	Package *Package
}

// Files returns a sorted slice of full file paths for this package.
func (p *Package) Files() Files {
	return Files{p}
}

// Names returns a sorted slice of file names for this package.
func (pf *Files) Names() []string {
	pf.Package.init()
	s := pf.Package.tc.docPkg.Filenames
	out := make([]string, len(s))
	for i := range s {
		out[i] = filepath.Base(s[i])
	}
	return out
}

// Paths returns a sorted slice of full file paths for this package.
func (pf *Files) Paths() []string {
	pf.Package.init()
	s := pf.Package.tc.docPkg.Filenames
	out := make([]string, len(s))
	copy(out, s)
	return out
}

// SourcePaths is like FilePaths but will include all files found by the build
// importer, .cc, .m, .s, etc while excluding test files.
func (pf *Files) SourcePaths() []string {
	var names []string
	p := pf.Package
	names = append(names, append(p.GoFiles, p.CgoFiles...)...)
	names = append(names, append(p.CXXFiles, p.MFiles...)...)
	names = append(names, append(p.SFiles, p.SwigFiles...)...)
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(p.Dir, name)
	}
	sort.Strings(paths)
	return paths
}

// TestPaths is like FilePaths but will include only test files.
func (pf *Files) TestPaths() []string {
	var names []string
	p := pf.Package
	names = append(names, append(p.TestGoFiles, p.XTestGoFiles...)...)
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(p.Dir, name)
	}
	sort.Strings(paths)
	return paths
}
