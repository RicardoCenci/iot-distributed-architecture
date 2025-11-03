package config

import (
	"slices"
	"strings"
)

type dotNotationMap struct {
	m map[string]interface{}
}

func newDotNotationMap() *dotNotationMap {
	return &dotNotationMap{
		m: make(map[string]interface{}),
	}
}

func (m *dotNotationMap) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	slices.Reverse(keys)
	stack := newStack(keys)

	return m.getRecursive(stack, m.m)
}

func (m *dotNotationMap) getRecursive(stack *stack, currentValue interface{}) interface{} {
	if stack.IsEmpty() {
		return currentValue
	}

	currentKey := stack.Pop()

	if _, ok := currentValue.(map[string]interface{})[currentKey]; ok {
		return m.getRecursive(stack, currentValue.(map[string]interface{})[currentKey])
	}

	return nil
}

func (m *dotNotationMap) Set(key string, value interface{}) {

	keys := strings.Split(key, ".")
	slices.Reverse(keys)
	stack := newStack(keys)

	m.setRecursive(stack, m.m, value)
}

func (m *dotNotationMap) setRecursive(stack *stack, currentMap map[string]interface{}, value interface{}) {

	currentKey := stack.Pop()

	if currentKey == "" {
		return
	}

	if _, ok := currentMap[currentKey]; !ok {
		currentMap[currentKey] = make(map[string]interface{})
	}

	if stack.IsEmpty() {
		currentMap[currentKey] = value
		return
	}

	m.setRecursive(
		stack,
		currentMap[currentKey].(map[string]interface{}),
		value,
	)
}

func (m *dotNotationMap) GetAllAsMap() map[string]interface{} {
	return m.m
}
