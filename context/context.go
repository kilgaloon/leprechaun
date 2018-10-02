package context

import (
	"os"
	"strings"
)

// Context holds variables in present context
type Context struct {
	variables []Variable
}

// Variable is struct definition for variable
type Variable struct {
	name  string
	value string
}

// GetName returns name of variable
func (v Variable) GetName() string {
	return v.name
}

// GetValue returns value of variable
func (v Variable) GetValue() string {
	return v.value
}

// DefineVar defines variable and puts it in present context
func (c *Context) DefineVar(variable string, value string) {
	var v = Variable{name: variable, value: value}
	c.variables = append(c.variables, v)
}

// GetVars returns all defined variables in context
func (c *Context) GetVars() []Variable {
	return c.variables
}

// GetVar finds var by name and returns its value
func (c *Context) GetVar(name string) Variable {
	for _, value := range c.variables {
		if value.name == name {
			return value
		}
	}

	return Variable{
		name:  name,
		value: name,
	}
}

// Transpile text change variables from context
func (c *Context) Transpile(toCompile string) string {
	for _, variable := range c.variables {
		toCompile = strings.Replace(toCompile, "$"+variable.name, variable.value, -1)
	}

	return toCompile
}

//New Create context for agent
func New() *Context {
	context := &Context{}
	// insert environment variables in our context
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		context.DefineVar(pair[0], pair[1])
	}

	return context
}
