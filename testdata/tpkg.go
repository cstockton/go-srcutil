// Package tpkg is used for testing srcutil. Lorem ipsum dolor sit amet,
// consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et
// dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation
// ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure
// dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
// pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui
// officia deserunt mollit anim id est laborum.
package tpkg

// init gravida malesuada, turpis lacus feugiat quis et diam. Dolor nisl fusce.
func init() {}

// NiladicFunc occaecati accumsan, metus magna sollicitudin, morbi mauris et
// eos quis placerat suspendisse quis, est laoreet sunt vestibulum pharetra
// turpis, etiam ad vestibulum nonummy mus viverra.
func NiladicFunc() string { return `` }

// NiladicVoidFunc occaecati accumsan, metus magna sollicitudin, morbi mauris
// et eos quis placerat suspendisse quis, est laoreet sunt vestibulum pharetra
// turpis, etiam ad vestibulum nonummy mus viverra.
func NiladicVoidFunc() {}

// StringFunc occaecati accumsan, metus magna sollicitudin, morbi mauris et
// eos quis placerat suspendisse quis, est laoreet sunt vestibulum pharetra
// turpis, etiam ad vestibulum nonummy mus viverra.
func StringFunc(str string) string { return `` }

// Pulvinar a habitasse amet illo, iaculis mi condimentum eget id. Consequat
// habitasse erat eros.
const (

	// ConstantOne ipsum non lacus mattis.
	ConstantOne = 42

	// ConstantTwo nulla vel tortor hac leo.
	ConstantTwo = true

	// ConstantThree enectus ante orci turpis leo placerat.
	ConstantThree = "Three"
)

// Lorem ornare accumsan integer, volutpat luctus sed ante malesuada suscipit
// elementum, scelerisque ut non diam pellentesque hymenaeos.
var (

	// VariableOne ipsum non lacus mattis.
	VariableOne = 42

	// VariableTwo nulla vel tortor hac leo.
	VariableTwo = true

	// VariableThree enectus ante orci turpis leo placerat.
	VariableThree = "Three"
)

// PublicStruct proin libero arcu, rerum orci tincidunt, lacus tempor sapien
// platea ullamcorper. Nullam velit, ipsum erat varius nam diam arcu vestibulum.
type PublicStruct struct {
	Name   string
	Number int
}

// MethodOne rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStruct) MethodOne() {}

// MethodTwo rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStruct) MethodTwo() string { return `` }

// MethodThree rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStruct) MethodThree(int, string) string { return `` }

// MethodOneP rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p *PublicStruct) MethodOneP() {}

// MethodTwoP rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p *PublicStruct) MethodTwoP() string { return `` }

// MethodThreeP rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p *PublicStruct) MethodThreeP(int, string) string { return `` }

// PublicStructUnexported proin libero arcu, rerum orci tincidunt, lacus tempor
// sapien platea ullamcorper. Nullam velit, ipsum erat varius nam diam arcu
// vestibulum.
type PublicStructUnexported struct {
	name   string
	number int
}

// MethodOne rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStructUnexported) MethodOne() {}

// MethodTwo rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStructUnexported) MethodTwo() string { return `` }

// MethodThree rutrum convallis lorem lacus, eu fusce mi sapien vitae.
func (p PublicStructUnexported) MethodThree(int, string) string { return `` }

// comment for private method.
func (p PublicStructUnexported) methodOne()                     {}
func (p PublicStructUnexported) methodTwo() string              { return `` }
func (p PublicStructUnexported) methodThree(int, string) string { return `` }
