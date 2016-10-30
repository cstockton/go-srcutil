package srcutil

import (
	"fmt"
	"go/build"
	"testing"
)

func TestPackage(t *testing.T) {
	ctx := FromWorkDir()

	t.Run("init", func(t *testing.T) {
		buildPkg, err := build.Default.Import(tPkg.ImportPath, defaultToGetwd(), 0)
		tmust(t, err)
		pkg := Package{Package: *buildPkg}
		if pkg.tc != nil {
			t.Errorf("expected nil tc for uninitialized Package")
		}
		pkg.Synopsis()
		exp := fmt.Sprintf("%p", pkg.tc)
		if pkg.tc == nil {
			t.Errorf("expected non-nil tc for initialized Package")
		}
		t.Run("Once", func(t *testing.T) {
			pkg.Synopsis()
			got := fmt.Sprintf("%p", pkg.tc)
			teq(t, exp, got)
		})
	})
	t.Run("Parsing", func(t *testing.T) {
		t.Run("ToAst", func(t *testing.T) {
			pkg, err := ctx.Import(tPkg.ImportPath)
			tmust(t, err)
			fileSet, astPkg, err := pkg.ToAst()
			tmust(t, err)
			teq(t, true, fileSet.Base() > 0)
			teq(t, tPkg.Name, astPkg.Name)

			t.Run("Failure", func(t *testing.T) {
				pkg.Package.Dir = ``
				fileSet, astPkg, err := pkg.ToAst()
				if err == nil {
					t.Errorf("expected error for non-existent import")
				}
				if fileSet != nil {
					t.Errorf("expected nil fileSet for non-existent import")
				}
				if astPkg != nil {
					t.Errorf("expected nil astPkg for non-existent import")
				}
			})
		})
	})
	t.Run("ToTypesInfo", func(t *testing.T) {
		pkg, err := ctx.Import(tPkg.ImportPath)
		tmust(t, err)
		typesInfo, typesPkg, err := pkg.ToInfo()
		tmust(t, err)
		teq(t, tPkg.Name, typesPkg.Name())
		teq(t, true, len(typesInfo.Types) > 0)

		t.Run("Failure", func(t *testing.T) {
			pkg.Package.Dir = ``
			typesInfo, typesPkg, err := pkg.ToInfo()
			if err == nil {
				t.Errorf("expected error for non-existent import")
			}
			if typesPkg != nil {
				t.Errorf("expected nil typesPkg for non-existent import")
			}
			if typesInfo != nil {
				t.Errorf("expected nil typesInfo for non-existent import")
			}
		})
	})
	t.Run("ToTypes", func(t *testing.T) {
		pkg, err := ctx.Import(tPkg.ImportPath)
		tmust(t, err)
		typesPkg, err := pkg.ToTypes()
		tmust(t, err)
		teq(t, tPkg.Name, typesPkg.Name())

		t.Run("Failure", func(t *testing.T) {
			pkg.Package.Dir = ``
			typesPkg, err := pkg.ToTypes()
			if err == nil {
				t.Errorf("expected error for non-existent import")
			}
			if typesPkg != nil {
				t.Errorf("expected nil typesPkg for non-existent import")
			}
		})
	})
	t.Run("ToDoc", func(t *testing.T) {
		pkg, err := ctx.Import(tPkg.ImportPath)
		tmust(t, err)
		docPkg, err := pkg.ToDoc()
		tmust(t, err)
		teq(t, tPkg.Name, docPkg.Name)

		t.Run("Failure", func(t *testing.T) {
			pkg.Package.Dir = ``
			docPkg, err := pkg.ToDoc()
			if err == nil {
				t.Errorf("expected error for non-existent import")
			}
			if docPkg != nil {
				t.Errorf("expected nil docPkg for non-existent import")
			}
		})
	})
}

const (
	TestConstant = "I am for testing."
)

func TestDocMethods(t *testing.T) {
	ctx := FromWorkDir()
	pkg, err := ctx.Import(tPkg.ImportPath)
	tmust(t, err)
	docs := pkg.Docs()

	t.Run("Notes", func(t *testing.T) {
		notes := docs.Notes()
		got, ok := notes["HELLO"]
		teq(t, true, ok)
		teq(t, 2, len(got))
		teq(t, "Note hello 1 for testing.\n", got[0].Body)
		teq(t, "Note hello 2 for testing.\n", got[1].Body)

		got, ok = notes["WORLD"]
		teq(t, true, ok)
		teq(t, 2, len(got))
		teq(t, "Note world 1 for testing.\n", got[0].Body)
		teq(t, "Note world 2 for testing.\n", got[1].Body)
	})
	t.Run("Examples", func(t *testing.T) {
		got := docs.Examples()
		teq(t, true, len(got) > 0)
		teq(t, "PublicStruct", got[0].Name)
	})
	t.Run("Consts", func(t *testing.T) {
		exp := "ConstantOne"
		got := docs.Consts()
		teq(t, true, len(got) > 0)
		teq(t, exp, got[0].Names[0])
	})
	t.Run("Types", func(t *testing.T) {
		types := docs.Types()
		teq(t, true, len(types) > 0)
	})
	t.Run("Methods", func(t *testing.T) {
		methods := docs.Methods()
		teq(t, true, len(methods) > 0)
	})
	t.Run("Vars", func(t *testing.T) {
		vars := docs.Vars()
		teq(t, true, len(vars) > 0)
	})
	t.Run("Funcs", func(t *testing.T) {
		funcs := docs.Funcs()
		teq(t, true, len(funcs) > 0)
	})
}

func TestFiles(t *testing.T) {
	ctx := FromWorkDir()
	pkg, err := ctx.Import(tPkg.ImportPath)
	tmust(t, err)
	files := pkg.Files()
	teq(t, pkg, files.Package)

	t.Run("Names", func(t *testing.T) {
		got := files.Names()
		teq(t, tPkg.Names, got)
	})
	t.Run("Paths", func(t *testing.T) {
		got := files.Paths()
		teq(t, tPkg.paths(tPkg.Names), got)
	})
	t.Run("SourcePaths", func(t *testing.T) {
		got := files.SourcePaths()
		teq(t, tPkg.paths(tPkg.PkgNames), got)
	})
	t.Run("TestPaths", func(t *testing.T) {
		got := files.TestPaths()
		teq(t, tPkg.paths(tPkg.PkgTests), got)
	})
}

func TestFuncs(t *testing.T) {
	ctx := FromWorkDir()
	pkg, err := ctx.Import(tPkg.ImportPath)
	tmust(t, err)

	t.Run("Funcs", func(t *testing.T) {
		ms, err := pkg.MethodSet("PublicStruct")
		tmust(t, err)
		teq(t, ms.Name, "PublicStruct")
		m, ok := ms.Methods["MethodOne"]
		teq(t, true, ok)
		teq(t, "MethodOne", m.Name())
	})
}

func TestMethods(t *testing.T) {
	ctx := FromWorkDir()
	pkg, err := ctx.Import(tPkg.ImportPath)
	tmust(t, err)

	t.Run("Methods", func(t *testing.T) {
		methods := pkg.Methods()
		teq(t, 2, len(methods))
		_, ok := methods["PublicStruct"]
		teq(t, true, ok)
	})
	t.Run("MethodSet", func(t *testing.T) {
		ms, err := pkg.MethodSet("PublicStruct")
		tmust(t, err)
		teq(t, "PublicStruct", ms.Name)
		m, ok := ms.Methods["MethodOne"]
		teq(t, true, ok)
		teq(t, "MethodOne", m.Name())

		t.Run("Names", func(t *testing.T) {
			teq(t, []string{"MethodOne", "MethodOneP", "MethodThree", "MethodThreeP",
				"MethodTwo", "MethodTwoP"}, ms.Names())
		})
	})
}
