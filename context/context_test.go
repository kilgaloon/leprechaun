package context

import (
	"testing"
)

var (
	ctx = BuildContext("context")
)

// Test defining variable and getting it back
func TestDefineVarGetVar(t *testing.T) {
	ctx.DefineVar("test_var", "test_value")

	definedVar := ctx.GetVar("test_var")
	if definedVar.value != "test_value" {
		t.Errorf("Expected test_value but got %s", definedVar.value)
	}
}

// Test transpiling
func TestTranspile(t *testing.T) {
	ctx.DefineVar("packageName", "Leprechaun")
	ctx.DefineVar("action", "transpiled")

	stringToTranspile := "This is $packageName, this is string ready to be $action"
	expectedToTranspile := "This is Leprechaun, this is string ready to be transpiled"

	transpiled := ctx.Transpile(stringToTranspile)
	if expectedToTranspile != transpiled {
		t.Errorf("Expected %s got %s", expectedToTranspile, transpiled)
	}
}
