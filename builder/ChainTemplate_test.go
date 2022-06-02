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
	lib := NewChainTemplate().Name().Handle().Next(fc).Next(fc)

	rc := NewChainCompile(fc.Clone())
	c1 := NewChainCompile(fc.Clone())
	lt := NewChainCompile(lib.Clone())
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
		fmt.Println("拦截器")
		return 500, nil
	}).Next(rc).Next(c1).Compile()

	caller := lt.Result()
	fmt.Println(caller.Path("lib", "func", "0.0.0").Params(1, 2, 3).Call())
	fmt.Println(string(lib.Description()))
}
