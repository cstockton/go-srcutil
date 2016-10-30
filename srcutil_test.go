package srcutil

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func tmust(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func teq(t *testing.T, exp, got interface{}) {
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("DeepEqual failed:\n  exp: %#v\n  got: %#v", exp, got)
	}
}

type testPackage struct {
	Name       string
	Path       string
	ImportPath string
	Names      []string
	PkgNames   []string
	PkgTests   []string
	Types      map[string]bool
	Funcs      map[string]bool
	Vars       map[string]bool
	Consts     map[string]bool
	DocFuncs   map[string]bool
	DocVars    map[string]bool
	DocMethods map[string]bool
}

func newTestPackage() *testPackage {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tPkg := &testPackage{
		Name:       "tpkg",
		Path:       filepath.Join(cwd, "testdata"),
		ImportPath: "github.com/cstockton/go-srcutil/testdata",
		Names: []string{
			"tpkg.go", "tpkg_private.go", "tpkg_test.go"},
		PkgNames: []string{
			"tpkg.go", "tpkg_private.go"},
		PkgTests: []string{
			"tpkg_example_test.go", "tpkg_test.go"},
		Types: map[string]bool{
			"PublicStruct": true, "PublicStructUnexported": true,
			"privateStruct": true, "privateStructExported": true},
		Funcs: map[string]bool{
			"funcOne": true, "funcTwo": true, "funcThree": true, "StringFunc": true,
			"NiladicFunc": true, "NiladicVoidFunc": true},
		Vars: map[string]bool{
			"variableOne": true, "variableTwo": true, "variableThree": true,
			"VariableOne": true, "VariableTwo": true, "VariableThree": true},
		Consts: map[string]bool{
			"constantOne": true, "constantTwo": true, "constantThree": true,
			"ConstantOne": true, "ConstantTwo": true, "ConstantThree": true},
	}
	return tPkg
}

func (p *testPackage) paths(s []string) []string {
	paths := make([]string, len(s))
	for i, name := range s {
		paths[i] = filepath.Join(p.Path, name)
	}
	sort.Strings(paths)
	return paths
}

var tPkg = newTestPackage()

func TestContext(t *testing.T) {
	cwd, err := os.Getwd()
	tmust(t, err)

	t.Run("Creation", func(t *testing.T) {
		t.Run("FromWorkDir", func(t *testing.T) {
			ctx := FromWorkDir()
			teq(t, cwd, ctx.SourceDir)
		})
		t.Run("FromStandard", func(t *testing.T) {
			ctx := FromStandard()
			teq(t, ``, ctx.GOPATH)
			teq(t, ctx.GOROOT+"/src", ctx.SourceDir)
		})
		t.Run("FromDir", func(t *testing.T) {
			ctx := FromDir(`.`)
			teq(t, `.`, ctx.SourceDir)
		})
	})
}

func TestImport(t *testing.T) {
	ctx := FromWorkDir()

	t.Run("FromGOPATH", func(t *testing.T) {
		pkg, err := ctx.Import(tPkg.ImportPath)
		tmust(t, err)
		teq(t, tPkg.Name, pkg.Name)
	})
	t.Run("FromStandardLibrary", func(t *testing.T) {
		pkg, err := ctx.Import("reflect")
		tmust(t, err)
		teq(t, "reflect", pkg.Name)
	})
	t.Run("Failure", func(t *testing.T) {
		pkg, err := ctx.Import("thislibrarydoesntexist")
		if err == nil {
			t.Errorf("expected error for non-existent import")
		}
		if pkg != nil {
			t.Errorf("expected nil pkg for non-existent import")
		}
	})
}
