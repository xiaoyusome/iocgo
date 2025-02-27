package iocgo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var log string

func Println(a ...interface{}) {
	msg := fmt.Sprintln(a...)
	log += msg
	fmt.Println(msg)
}

type Fooer interface {
	Foo(int)
}
type Foo struct {
}

func NewFoo() *Foo {
	return &Foo{}
}
func (Foo) Foo(i int) {
	Println("foo:", i)
}

type Barer interface {
	Bar(string)
}
type Bar struct {
}

func NewBar() *Bar {
	return &Bar{}
}
func (Bar) Bar(s string) {
	Println("bar:", s)
}

type Foobarer interface {
	Say(int, string)
}
type Foobar struct {
	foo Fooer
	bar Barer
	msg string
}

func NewFoobar(f Fooer, b Barer) Foobarer {
	return &Foobar{
		foo: f,
		bar: b,
	}
}
func NewFoobarWithMsg(f Fooer, b Barer, msg string) Foobarer {
	//Println("NewFoobarWithMsg~~~~~~~")
	return &Foobar{
		foo: f,
		bar: b,
		msg: msg,
	}
}
func (f *Foobar) Say(i int, s string) {
	if f.foo != nil {
		f.foo.Foo(i)
	}
	if f.bar != nil {
		f.bar.Bar(s)
	}
	f.msg += fmt.Sprintf("[%d,%s]", i, s)
	Println("Foobar msg:", f.msg)
}
func TestContainer_SimpleRegister(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar)
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} })
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))

}
func TestContainer_RegisterParameters(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobarWithMsg, Parameters(map[int]interface{}{2: "studyzy"}))
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} })
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
	assert.True(t, strings.Contains(log, "studyzy"))
}

type Baz struct{}

func (Baz) Bar(s string) {
	Println("baz:", s)
}
func TestContainer_RegisterDefault(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar)
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} })
	container.Register(func() Barer { return &Baz{} }, Default())
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "baz:"))
}

func TestContainer_RegisterOptional(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar, Optional(0))
	container.Register(func() Barer { return &Bar{} })
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, !strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
}
func TestContainer_RegisterInterface(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar)
	var f Fooer
	err := container.Register(NewFoo, Interface(&f))
	assert.Nil(t, err)
	var b Barer
	err = container.Register(NewBar, Interface(&b))
	assert.Nil(t, err)
	var fb Foobarer
	err = container.Resolve(&fb)
	assert.Nil(t, err)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
}
func TestContainer_RegisterLifestyleTransient(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobarWithMsg, Parameters(map[int]interface{}{2: "studyzy"}), Lifestyle(true))
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} })
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
	assert.True(t, strings.Contains(log, "123"))
	log = ""
	var fb2 Foobarer
	container.Resolve(&fb2) //resolve a new instance since Lifestyle(transient=true)
	fb2.Say(456, "Hi")
	t.Log(log)
	assert.True(t, strings.Contains(log, "Hi"))
	assert.True(t, !strings.Contains(log, "123"))

}
func TestContainer_RegisterDependsOn(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar, DependsOn(map[int]string{1: "bar"})) //depend on "bar" name Barer
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} }, Name("bar"))
	container.Register(func() Barer { return &Baz{} }, Name("baz"))
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
}
func TestContainer_RegisterInstance(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobar)
	container.Register(func() Fooer { return &Foo{} })
	b := &Bar{}
	var bar Barer
	container.RegisterInstance(&bar, b, Default()) // register interface -> instance
	var fb Foobarer
	container.Resolve(&fb)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
}
func TestContainer_Resolve(t *testing.T) {
	log = ""
	container := NewContainer()
	container.Register(NewFoobarWithMsg, Parameters(map[int]interface{}{2: "studyzy"}))                  //default Foobar register
	container.Register(NewFoobarWithMsg, Parameters(map[int]interface{}{2: "Devin"}), Name("instance2")) //named Foobar register
	container.Register(func() Fooer { return &Foo{} })
	container.Register(func() Barer { return &Bar{} })
	var fb Foobarer
	container.Resolve(&fb, Arguments(map[int]interface{}{2: "arg2"})) //resolve use new argument to replace register parameters
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
	assert.True(t, strings.Contains(log, "arg2"))
	assert.False(t, strings.Contains(log, "studyzy"))
	log = ""
	var fb2 Foobarer
	err := container.Resolve(&fb2, ResolveName("instance2")) //resolve by name
	assert.Nil(t, err)
	fb2.Say(456, "New World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "Devin"))

}

func TestContainer_Call(t *testing.T) {
	log = ""

	Register(NewFoobar)
	Register(func() Fooer { return &Foo{} })
	Register(func() Barer { return &Bar{} })

	Call(SayHi1)
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
	log = ""
	Call(SayHi2, CallArguments(map[int]interface{}{2: "Devin"}))
	assert.True(t, strings.Contains(log, "Devin"))
	log = ""
	Register(func() Barer { return &Baz{} }, Name("baz"))
	Call(SayHi1, CallDependsOn(map[int]string{1: "baz"}))
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "baz:"))

}
func SayHi1(f Fooer, b Barer) {
	f.Foo(1234)
	b.Bar("hi")
}

func SayHi2(f Fooer, b Barer, hi string) {
	f.Foo(len(hi))
	b.Bar(hi)
	Println("SayHi")
}

func TestContainer_Reset(t *testing.T) {
	log = ""
	Register(NewFoobar)
	Register(func() Fooer { return &Foo{} })
	Register(func() Barer { return &Bar{} })
	var fb Foobarer
	err := Resolve(&fb)
	assert.Nil(t, err)
	fb.Say(123, "Hello World")
	t.Log(log)
	assert.True(t, strings.Contains(log, "foo:"))
	assert.True(t, strings.Contains(log, "bar:"))
	Reset()
	err = Resolve(&fb)
	assert.NotNil(t, err)
}
