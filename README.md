<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [testcase](#testcase)
  - [Features](#features)
  - [Guide](#guide)
  - [Official API Documentation](#official-api-documentation)
  - [Getting Started / Example](#getting-started--example)
  - [Modules](#modules)
  - [Summary](#summary)
    - [DRY](#dry)
    - [Modularization](#modularization)
  - [Stability](#stability)
  - [Case Study About `testcase` Package Origin](#case-study-about-testcase-package-origin)
  - [Reference Project](#reference-project)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GoDoc](https://godoc.org/github.com/adamluzsi/testcase?status.png)](https://godoc.org/github.com/adamluzsi/testcase)
[![Build Status](https://travis-ci.org/adamluzsi/testcase.svg?branch=master)](https://travis-ci.org/adamluzsi/testcase)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamluzsi/testcase)](https://goreportcard.com/report/github.com/adamluzsi/testcase)
[![codecov](https://codecov.io/gh/adamluzsi/testcase/branch/master/graph/badge.svg)](https://codecov.io/gh/adamluzsi/testcase)
# testcase

The `testcase` package provides tooling to apply BDD testing conventions.

## Features

- lightweight, it has no dependency, you just import testcase, no clutter
- supports classing flattened testing style, table testing style and nested testing style 
  - [nested testing style](/docs/nested-testing-style.md)
    * allows the usage of an immutable testing subject
    * nesting visualize code complexity
    * testing context building becomes DRY
    * immutable testing subject construction forces disciplined testing conventions
    * supports also the flattened nesting where testing context branches can to a top level fonction  
- DDD based testing composition
  * testing components can be defined and reused across the whole project
  * unit tests can be converted into full-fledged E2E tests for CI/CD pipelines
  * dependency injection during testing becomes a breeze to do
  * reusable testing components increase maintainability aspects
  * reusable testing components decrease required ramp-up time for writing new tests that depend on existing components
- safe parallel testing
  * variables stored per test execution, which prevents leakage and race conditions across different tests
- repeatable pseudo-random fixture generation for test input
- repeatable pseudo-random test order shuffling
  * prevents implicit test dependency on ordering
  * ensures that tests can be added and removed freely without the fear of breaking other tests in the same coverage.
  * flaky tests which depend on test execution order can be noticed at development time

## Guide

`testcase` is a powerful TDD Tooling that requires discipline and understanding about the fundamentals of testing.
If you are looking for a guide that helps streamline your knowledge on the topics,
then please consider read the below listed articles.
 
- [why to test?](/docs/why-to-test.md)
- [nested testing style With testcase.Specs](/docs/nested-testing-style.md)
- [indirection and abstraction and the use in testing](/docs/interface.md)
- [testing doubles, when to use what](/docs/testing-double/README.md)
- [TDD, Role Interface and Contracts](/docs/contracts.md) [WIP-90%]
- [design system with DRY TDD for long term maintainability wins](/docs/spechelper.md) [WIP-1%]

## Official API Documentation

If you already use the framework, and you just want pick an example,
you can go directly to the API documentation that is kept in godoc format.
- [godoc](https://godoc.org/github.com/adamluzsi/testcase)
- [pkg.go.dev](https://pkg.go.dev/github.com/adamluzsi/testcase).

## Getting Started / Example

Examples kept in godoc format.
Every exported functionality aims to have examples provided in the official documentation. 
- [pkg.go.dev](https://pkg.go.dev/github.com/adamluzsi/testcase#pkg-examples)
- [godoc](https://godoc.org/github.com/adamluzsi/testcase#pkg-examples)

A Basic example:

```go
package main

import (
	"testing"

	"github.com/adamluzsi/testcase"
)

func TestMessageWrapper(t *testing.T) {
	s := testcase.NewSpec(t)
	s.NoSideEffect()

	var (
		message        = testcase.Var{Name: `message`}
		messageWrapper = s.Let(`message wrapper`, func(t *testcase.T) interface{} {
			return MessageWrapper{Message: message.Get(t).(string)}
		})
	)

	s.Describe(`#LookupMessage`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) (string, bool) {
			return messageWrapper.Get(t).(MessageWrapper).LookupMessage()
		}

		s.When(`message is empty`, func(s *testcase.Spec) {
			message.LetValue(s, ``)

			s.Then(`it will return with "ok" as false`, func(t *testcase.T) {
				_, ok := subject(t)
				assert.Must(t).True(!ok)
			})
		})

		s.When(`message has content`, func(s *testcase.Spec) {
			message.LetValue(s, fixtures.Random.String())

			s.Then(`it will return with "ok" as true`, func(t *testcase.T) {
				_, ok := subject(t)
				assert.Must(t).True( ok)
			})

			s.Then(`message received back`, func(t *testcase.T) {
				msg, _ := subject(t)
				assert.Must(t).Equal( message.Get(t), msg)
			})
		})
	})
}
```

## Modules
- [httpspec](/httpspec/README.md)
    * spec module helps you create HTTP API Specs.
- [fixtures](/fixtures/README.md)
    * fixtures module helps you create random input values for testing

## Summary

### DRY

`testcase` provides a way to express common Arrange, Act sections for the Asserts with DRY principle in mind.

- First you can define your Act section with a method under test as the subject of your test specification
    * The Act section invokes the method under test with the arranged parameters.
- Then you can build the context of the Act by Arranging the inputs later with humanly explained reasons
    * The Arrange section initializes objects and sets the value of the data that is passed to the method under test.
- And lastly you can define the test expected outcome in an Assert section.
    * The Assert section verifies that the action of the method under test behaves as expected.

Then adding an additional test edge case to the testing suite becomes easier,
as it will have a concrete place where it must be placed.

And if during the creation of the specification, an edge case turns out to be YAGNI,
it can be noted, so visually it will be easier to see what edge case is not specified for the given subject.

The value it gives is that to build test for a certain edge case,
the required mental model size to express the context becomes smaller,
as you only have to focus on one Arrange at a time,
until you fully build the bigger picture.

It also implicitly visualize the required mental model of your production code by the nesting.
[You can read more on that in the nesting section](/docs/nesting.md).

### Modularization

On top of the DRY convention, any time you need to Arrange a common scenario about your projects domain event,
you can modularize these setup blocks in a helper functions.

This helps the readability of the test, while keeping the need of mocks to the minimum as possible for a given test.
As a side effect, integration tests can become low hanging fruit for the project.

e.g.:
```go
package mypkg_test

import (
	"testing"

	"my/project/mypkg"


	"github.com/adamluzsi/testcase"

	. "my/project/testing/pkg"
)

func TestMyTypeMyFunc(t *testing.T) {
	s := testcase.NewSpec(t)

	// high level Arrange helpers from my/project/testing/pkg
	SetupSpec(s)
	GivenWeHaveUser(s, `myuser`)
	// .. other givens

	myType := func() *mypkg.MyType { return &mypkg.MyType{} }

	s.Describe(`#MyFunc`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) { myType().MyFunc(t.I(`myuser`).(*mypkg.User)) } // Act

		s.Then(`edge case description`, func(t *testcase.T) {
			// Assert
			subject(t)
		})
	})
}
```

## Stability

- The package considered stable.
- The package use rolling release conventions.
- No breaking change is planned to the package exported API.
- The package used for production development.
- The package API is only extended if the practical use case proves its necessity.

## [Case Study About `testcase` Package Origin](/docs/history.md)

## Reference Project

- [toggler project, scalable feature toggles on budget for startups](https://github.com/adamluzsi/toggler)
- [frameless project, for a vendor lock free software architecture](https://github.com/adamluzsi/frameless)
- [gorest, a minimalist REST controller for go projects](https://github.com/adamluzsi/gorest)

