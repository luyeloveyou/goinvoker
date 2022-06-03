package builder

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewChainTemplate(t *testing.T) {
	template := NewChainTemplate()
	if template == nil {
		t.Errorf("生成的模板为nil")
	}
}

func TestChainTemplate_Name(t *testing.T) {
	template := NewChainTemplate()
	template.Name()
	if !reflect.DeepEqual(template.Description(), []rune{NameRouter}) {
		t.Errorf("expected: %v, got: %v", string(NameRouter), string(template.Description()))
	}
}

func TestChainTemplate_Version(t *testing.T) {
	template := NewChainTemplate()
	template.Version()
	if !reflect.DeepEqual(template.Description(), []rune{VersionRouter}) {
		t.Errorf("expected: %v, got: %v", string(VersionRouter), string(template.Description()))
	}
}
func TestChainTemplate_Build(t *testing.T) {
	template := NewChainTemplate()
	template.Name().Version().Handle()
	template.Fill("test", "1.0.0", HandlerHelper(func(id uint64, result any, params []any) (any, error) {
		fmt.Println(id, result, params)
		return 1, nil
	})).Build()
	caller := template.Result()
	caller.Path("test", "1.0.0").Result(1).Call()
	caller.Path("test", "1.0.0").Call()
}

func TestChainCompile_Compile(t *testing.T) {
	fc := NewChainTemplate().Name().Version().Handle()
	lib := NewChainTemplate().Name().Handle().Next(fc, -1)
	invoker := NewChainTemplate().Name().Handle().Next(lib, 1)

	rc := NewChainCompile(fc.Clone())
	c1 := NewChainCompile(fc.Clone())
	lt := NewChainCompile(lib.Clone())
	it := NewChainCompile(invoker.Clone())
	ct := NewChainCompile(fc.Clone())

	ct.Name("func").Version("0.0.0")
	rc.Resolve(ct).Handle(func(id uint64, result any, params []any) (any, error) {
		fmt.Println(id, result, params)
		return 100, nil
	})
	c1.Resolve(ct).Handle(func(id uint64, result any, params []any) (any, error) {
		fmt.Println(id, result, params)
		return 2, nil
	})
	lt.Name("lib").Handle(func(id uint64, result any, params []any) (any, error) {
		fmt.Println("lib拦截器")
		fmt.Println(id, result, params)
		return 500, nil
	}).Next(rc).Next(c1)
	it.Name("invoker").Handle(func(id uint64, result any, params []any) (any, error) {
		fmt.Println("invoker拦截器")
		fmt.Println(id, result, params)
		return 600, nil
	}).Next(lt).Compile()

	caller := it.Result()
	fmt.Println(caller.Path("invoker", "lib", "func", "0.0.0").Params(1, 2, 3).Call())
	fmt.Println(string(invoker.Description()))
}

func TestNewChainCompile(t *testing.T) {
	fc := NewChainTemplate().Name().Version().Handle()
	lib := NewChainTemplate().Name().Handle().Next(fc, -1)
	lt := NewChainCompile(lib.Clone())
	fmt.Printf("%#v\n", lt)
}
