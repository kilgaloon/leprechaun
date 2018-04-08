package client

import (
	"strings"
)
// Context holds variables in present context
type Context struct {
	variables []Variable
}
// Variable is struct definition for variable
type Variable struct {
	name string
	value string
}

// DefineVar defines variable and puts it in present context
func (c *Context) DefineVar(variable string, value string) {
	var v = Variable{name: variable, value: value}
	c.variables = append(c.variables, v)
}
// GetVars returns all defined variables in context
func (c *Context) GetVars() []Variable {
	return c.variables;
}
// GetVar finds var by name and returns its value
func (c *Context) GetVar(name string) Variable {
	for _, value := range c.variables {
		if value.name == name {
			return value;
		}
	}

	return Variable{
		name: name,
		value: name,
	}
}
// Transpile text change variables from context
func (c *Context) Transpile(toCompile string) string {
	var str string

	for _, variable := range c.variables {
		str = strings.Replace(toCompile, "$" + variable.name, variable.value, -1)
	}

	return str;
}

// CurrentContext of client
var CurrentContext = new(Context)