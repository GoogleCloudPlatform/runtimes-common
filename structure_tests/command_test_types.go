package main

import (
	"strings"
)

type CommandTest interface {
	GetName()	string
	GetCommand()	string
	GetFlags()	string
	GetExpectedOutput() []string
	GetExcludedOutput() []string
	GetExpectedError()  []string
	GetExcludedError()  []string // excluded error from running command
}

type CommandTestv0 struct {
	Name           string
	Command        string
	Flags          string
	ExpectedOutput []string
	ExcludedOutput []string
	ExpectedError  []string
	ExcludedError  []string // excluded error from running command
}

func (t CommandTestv0) GetName() string {
	return t.Name
}

func (t CommandTestv0) GetCommand() string {
	return t.Command
}

func (t CommandTestv0) GetFlags() string {
	return t.Flags
}

func (t CommandTestv0) GetExpectedOutput() []string {
	return t.ExpectedOutput
}

func (t CommandTestv0) GetExcludedOutput() []string {
	return t.ExcludedOutput
}

func (t CommandTestv0) GetExpectedError() []string {
	return t.ExpectedError
}

func (t CommandTestv0) GetExcludedError() []string {
	return t.ExcludedError
}

type CommandTestv1 struct {
	CommandTest
	Name           string
	Command        string
	Flags          []string
	ExpectedOutput []string
	ExcludedOutput []string
	ExpectedError  []string
	ExcludedError  []string // excluded error from running command
}

func (t CommandTestv1) GetName() string {
	return t.Name
}

func (t CommandTestv1) GetCommand() string {
	return t.Command
}

func (t CommandTestv1) GetFlags() string {
	return strings.Join(t.Flags, " ")
}

func (t CommandTestv1) GetExpectedOutput() []string {
	return t.ExpectedOutput
}

func (t CommandTestv1) GetExcludedOutput() []string {
	return t.ExcludedOutput
}

func (t CommandTestv1) GetExpectedError() []string {
	return t.ExpectedError
}

func (t CommandTestv1) GetExcludedError() []string {
	return t.ExcludedError
}