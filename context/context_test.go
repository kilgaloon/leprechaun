package context_test

import (
	"testing"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	fakeClient = agent.New("test", cfgWrap.New("test", *path))
	ctx        = context.New()
)

// Test defining variable and getting it back
func TestDefineVarGetVar(t *testing.T) {
	ctx.DefineVar("test_var", "test_value")

	if len(ctx.GetVars()) < 1 {
		t.Fail()
	}

	definedVar := ctx.GetVar("test_var")
	if definedVar.GetValue() != "test_value" {
		t.Errorf("Expected test_value but got %s", definedVar.GetValue())
	}

	undefinedVar := ctx.GetVar("test_not_var")
	if undefinedVar.GetName() != "test_not_var" || undefinedVar.GetValue() != "test_not_var" {
		t.Fail()
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
