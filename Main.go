package main

import (
	"fmt"
	"goinvoker/core"
)

func main() {
	var fish = &Fish{
		Name:     "fish",
		SwimFlag: false,
	}
	call(fish)

	var dog = &Dog{
		Name:    "doggy",
		RunFlag: false,
	}
	call(dog)

	fmt.Println(fish)
	fmt.Println(dog)
}

func call(object core.Object) {
	switch object.(type) {
	case Swimable:
		var s = object.(Swimable)
		s.Swim()
	case Runnable:
		var s = object.(Runnable)
		s.Run()
	}
}

type Swimable interface {
	Swim()
}

type Runnable interface {
	Run()
}

type Fish struct {
	Name     string
	SwimFlag bool
}

type Dog struct {
	Name    string
	RunFlag bool
}

func (f *Fish) Swim() {
	f.SwimFlag = true
	fmt.Println("swim")
}

func (d *Dog) Run() {
	d.RunFlag = true
	fmt.Println("run")
}
