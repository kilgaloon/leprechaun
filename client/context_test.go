package client

import "testing"

// Test defining variable and getting it back
func TestDefineVarGetVar(t *testing.T) {
	CurrentContext.DefineVar("test_var", "test_value")

	definedVar := CurrentContext.GetVar("test_var")
	if definedVar.value != "test_value" {
		t.Errorf("Expected test_value but got %s", definedVar.value)
	}
}

// Test transpiling
func TestTranspile(t *testing.T) {
	CurrentContext.DefineVar("packageName", "Leprechaun")
	CurrentContext.DefineVar("action", "transpiled")

	stringToTranspile := "This is $packageName, this is string ready to be $action"
	expectedToTranspile := "This is Leprechaun, this is string ready to be transpiled"

	transpiled := CurrentContext.Transpile(stringToTranspile)
	if expectedToTranspile != transpiled {
		t.Errorf("Expected %s got %s", expectedToTranspile, transpiled)
	}
}
