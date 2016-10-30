package tpkg

// Pulvinar a habitasse amet illo, iaculis mi condimentum eget id. Consequat
// habitasse erat eros.
const (

	// constantOne ipsum non lacus mattis.
	constantOne = 42

	// ConstantTwo nulla vel tortor hac leo.
	constantTwo = true

	// ConstantThree enectus ante orci turpis leo placerat.
	constantThree = "Three"
)

// Lorem ornare accumsan integer, volutpat luctus sed ante malesuada suscipit
// elementum, scelerisque ut non diam pellentesque hymenaeos.
var (

	// VariableOne ipsum non lacus mattis.
	variableOne = 42

	// VariableTwo nulla vel tortor hac leo.
	variableTwo = true

	// VariableThree enectus ante orci turpis leo placerat.
	variableThree = "Three"
)

// comment for private func.
func funcOne()                     {}
func funcTwo() string              { return `` }
func funcThree(int, string) string { return `` }

// privateStruct Consectetuer metus blandit est, auctor laoreet, enim leo id
// ante. Suspendisse tincidunt quam ipsum lacinia dui, a ac in enim sed.
type privateStruct struct {
	unexportedStr string
	unexportedInt int
}

// comment for privateStruct method.
func (p privateStruct) funcOne()                       {}
func (p privateStruct) funcTwo() string                { return `` }
func (p privateStruct) funcThree(int, string) string   { return `` }
func (p *privateStruct) funcOneP()                     {}
func (p *privateStruct) funcTwoP() string              { return `` }
func (p *privateStruct) funcThreeP(int, string) string { return `` }

// privateStruct Consectetuer metus blandit est, auctor laoreet, enim leo id
// ante. Suspendisse tincidunt quam ipsum lacinia dui, a ac in enim sed.
type privateStructExported struct {
	ExportedStr string
	ExportedInt int
}

// comment for privateStruct method.
func (p privateStructExported) FuncOne()                       {}
func (p privateStructExported) funcTwo() string                { return `` }
func (p privateStructExported) funcThree(int, string) string   { return `` }
func (p *privateStructExported) FuncOneP()                     {}
func (p *privateStructExported) funcTwoP() string              { return `` }
func (p *privateStructExported) funcThreeP(int, string) string { return `` }

// HELLO(cstockton): Note hello 1 for testing.
// HELLO(chris): Note hello 2 for testing.
// WORLD(cstockton): Note world 1 for testing.
// WORLD(chris): Note world 2 for testing.
