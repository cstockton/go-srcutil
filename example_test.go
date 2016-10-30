package srcutil_test

import (
	"fmt"
	"go/doc"
	"log"
	"strings"

	"github.com/cstockton/go-srcutil"
)

func Example() {
	pkg, err := srcutil.Import("github.com/cstockton/go-srcutil")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("String:", pkg)
	fmt.Println("Synopsis:", pkg.Doc)

	vars := pkg.Vars()
	for _, v := range vars {
		fmt.Printf("Var: %v %v\n", v.Name(), v.Type())
	}

	// Output:
	// String: srcutil.Package{srcutil}
	// Synopsis: Package srcutil provides utilities for working with Go source code.
	// Var: DefaultContext *srcutil.Context
	// Var: DefaultImportMode go/build.ImportMode
	// Var: DefaultParseMode go/parser.Mode
}

func ExamplePackage_Docs() {
	pkg, err := srcutil.Import("io")
	if err != nil {
		log.Fatal(err)
	}
	docs := pkg.Docs()

	consts := docs.Consts()
	fmt.Printf("// %v", consts[0].Doc)
	fmt.Printf("const(\n  %v\n)\n\n", strings.Join(consts[0].Names, "\n  "))

	vars := docs.Vars()
	for _, v := range vars {
		fmt.Printf("// %v", doc.Synopsis(consts[0].Doc))
		fmt.Printf("var %v\n", v.Names[0])
	}
	fmt.Print("\n")

	types := docs.Types()
	for _, typ := range types {
		if strings.Contains(typ.Name, "Reader") {
			fmt.Printf("// %v\n", doc.Synopsis(typ.Doc))
			for _, f := range typ.Funcs {
				fmt.Printf("// %v\n", doc.Synopsis(f.Doc))
			}
		}
	}

	// Output:
	// // Seek whence values.
	// const(
	//   SeekStart
	//   SeekCurrent
	//   SeekEnd
	// )
	//
	// // Seek whence values.var EOF
	// // Seek whence values.var ErrClosedPipe
	// // Seek whence values.var ErrNoProgress
	// // Seek whence values.var ErrShortBuffer
	// // Seek whence values.var ErrShortWrite
	// // Seek whence values.var ErrUnexpectedEOF
	//
	// // ByteReader is the interface that wraps the ReadByte method.
	// // A LimitedReader reads from R but limits the amount of data returned to just N bytes.
	// // A PipeReader is the read half of a pipe.
	// // Pipe creates a synchronous in-memory pipe.
	// // Reader is the interface that wraps the basic Read method.
	// // LimitReader returns a Reader that reads from r but stops with EOF after n bytes.
	// // MultiReader returns a Reader that's the logical concatenation of the provided input readers.
	// // TeeReader returns a Reader that writes to w what it reads from r.
	// // ReaderAt is the interface that wraps the basic ReadAt method.
	// // ReaderFrom is the interface that wraps the ReadFrom method.
	// // RuneReader is the interface that wraps the ReadRune method.
	// // SectionReader implements Read, Seek, and ReadAt on a section of an underlying ReaderAt.
	// // NewSectionReader returns a SectionReader that reads from r starting at offset off and stops with EOF after n bytes.
}

func ExamplePackage_Methods() {
	pkg, err := srcutil.Import("bytes")
	if err != nil {
		log.Fatal(err)
	}
	pkgMethods := pkg.Methods()

	printer := func(methodSet srcutil.MethodSet) {
		fmt.Printf("type %v (%d methods)\n", methodSet.Name, methodSet.Len())
		for _, name := range methodSet.Names() {
			method := methodSet.Methods[name]
			fmt.Printf("  %v%v\n    returns %v\n", name, method.Params(), method.Results())
		}
	}
	printer(pkgMethods["Reader"])
	printer(pkgMethods["Buffer"])

	// Output:
	// type Reader (11 methods)
	//   Len()
	//     returns (int)
	//   Read(b []byte)
	//     returns (n int, err error)
	//   ReadAt(b []byte, off int64)
	//     returns (n int, err error)
	//   ReadByte()
	//     returns (byte, error)
	//   ReadRune()
	//     returns (ch rune, size int, err error)
	//   Reset(b []byte)
	//     returns ()
	//   Seek(offset int64, whence int)
	//     returns (int64, error)
	//   Size()
	//     returns (int64)
	//   UnreadByte()
	//     returns (error)
	//   UnreadRune()
	//     returns (error)
	//   WriteTo(w io.Writer)
	//     returns (n int64, err error)
	// type Buffer (23 methods)
	//   Bytes()
	//     returns ([]byte)
	//   Cap()
	//     returns (int)
	//   Grow(n int)
	//     returns ()
	//   Len()
	//     returns (int)
	//   Next(n int)
	//     returns ([]byte)
	//   Read(p []byte)
	//     returns (n int, err error)
	//   ReadByte()
	//     returns (byte, error)
	//   ReadBytes(delim byte)
	//     returns (line []byte, err error)
	//   ReadFrom(r io.Reader)
	//     returns (n int64, err error)
	//   ReadRune()
	//     returns (r rune, size int, err error)
	//   ReadString(delim byte)
	//     returns (line string, err error)
	//   Reset()
	//     returns ()
	//   String()
	//     returns (string)
	//   Truncate(n int)
	//     returns ()
	//   UnreadByte()
	//     returns (error)
	//   UnreadRune()
	//     returns (error)
	//   Write(p []byte)
	//     returns (n int, err error)
	//   WriteByte(c byte)
	//     returns (error)
	//   WriteRune(r rune)
	//     returns (n int, err error)
	//   WriteString(s string)
	//     returns (n int, err error)
	//   WriteTo(w io.Writer)
	//     returns (n int64, err error)
	//   grow(n int)
	//     returns (int)
	//   readSlice(delim byte)
	//     returns (line []byte, err error)

}
