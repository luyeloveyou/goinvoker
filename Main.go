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
	test := call(fish)
	fish.Name = "fff"
	fmt.Println(test)
	var dog = &Dog{
		Name:    "doggy",
		RunFlag: false,
	}
	dog.Name = "ddd"
	test = call(dog)
	fmt.Println(test)

	fmt.Println(fish)
	fmt.Println(dog)
}

func call(object core.Object) core.Object {
	switch object.(type) {
	case Swimable:
		var s = object.(Swimable)
		s.Swim()
	case Runnable:
		var s = object.(Runnable)
		s.Run()
	}
	return object
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
