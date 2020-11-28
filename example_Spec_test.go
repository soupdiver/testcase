package testcase_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/testcase"
)

type RoleInterface interface {
	Say() string
}

type MyType struct {
	MyResource RoleInterface
}

func (mt *MyType) MyFunc() {}

func (mt *MyType) IsLower(s string) bool {
	return strings.ToLower(s) == s
}

func (mt *MyType) Fallible() (string, error) {
	return "", nil
}

type MyResourceSupplier struct{}

func (MyResourceSupplier) Say() string {
	return `Hello, world!`
}

func ExampleSpec_myType() {
	var t *testing.T

	// spec do not use any global magic
	// it is just a simple abstraction around testing.T#Context
	// Basically you can easily can run it as you would any other go test
	//   -> `go run ./... -v -run "my/edge/case/nested/block/I/want/to/run/only"`
	//
	spec := testcase.NewSpec(t)

	// when you have no side effects in your testing suite,
	// you can enable parallel execution.
	// You can play parallel even from nested specs to apply parallel testing for that context and below.
	spec.Parallel()
	// or
	spec.NoSideEffect()

	// testcase.variables are thread safe way of setting up complex contexts
	// where some variable need to have different values for edge cases.
	// and I usually work with in-memory implementation for certain shared specs,
	// to make my test coverage run fast and still close to somewhat reality in terms of integration.
	// and to me, it is a necessary thing to have "T#parallel" ContextOption safely available
	var myType = func(t *testcase.T) *MyType {
		return &MyType{}
	}

	spec.Describe(`IsLower`, func(s *testcase.Spec) {
		// it is a convention to me to always make a subject for a certain describe block
		//
		var (
			input   = testcase.Var{Name: `input`}
			subject = func(t *testcase.T) bool {
				return myType(t).IsLower(input.Get(t).(string))
			}
		)

		s.When(`input string has lower case characters`, func(s *testcase.Spec) {
			input.LetValue(s, "all lower case")

			s.Before(func(t *testcase.T) {
				// here you can do setups like cleanup for DB tests
			})

			s.After(func(t *testcase.T) {
				// here you can setup a teardown
			})

			s.Around(func(t *testcase.T) func() {
				// here you can setup things that need teardown
				// such example to me is when I use gomock.Controller and mock setup

				return func() {
					// you can do teardown in this
					// this func will be defered after the test cases
				}
			})

			s.And(`the first character is capitalized`, func(s *testcase.Spec) {
				// you can add more nesting for more concrete specifications,
				// in each nested block, you work on a separate variable stack,
				// so even if you overwrite something here,
				// that has no effect outside of this scope
				input.LetValue(s, "First character is uppercase")

				s.Then(`it will report false`, func(t *testcase.T) {
					require.False(t, subject(t),
						fmt.Sprintf(`it was expected that %q will be reported to be not lowercase`, t.I(`input`)))
				})

			})

			s.Then(`it will return true`, func(t *testcase.T) {
				require.True(t, subject(t),
					fmt.Sprintf(`it was expected that the %q will re reported to be lowercase`, t.I(`input`)))
			})
		})
	})

	spec.Describe(`Fallible`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) (string, error) {
			return myType(t).Fallible()
		}

		var onSuccess = func(t *testcase.T) string {
			someMeaningfulVarName, err := subject(t)
			require.Nil(t, err)
			return someMeaningfulVarName
		}

		s.Then(`it will return an empty string`, func(t *testcase.T) {
			require.Equal(t, "", onSuccess(t))
		})
	})
}
