package examples_test

import (
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/adamluzsi/testcase/assert"
)

func IfSubject(condition bool) string {
	if condition {
		return `A`
	} else {
		return `B`
	}
}

func TestIfSubject(t *testing.T) {
	s := testcase.NewSpec(t)

	var (
		condition = testcase.Var{Name: `condition`}
		subject   = func(t *testcase.T) string {
			return IfSubject(condition.Get(t).(bool))
		}
	)

	s.When(`condition described`, func(s *testcase.Spec) {
		condition.LetValue(s, true)

		s.Then(`it will return ...`, func(t *testcase.T) {
			assert.Must(t).Equal(`A`, subject(t))
		})
	})

	s.When(`condition opposite described`, func(s *testcase.Spec) {
		condition.LetValue(s, false)

		s.Then(`it will return ...`, func(t *testcase.T) {
			assert.Must(t).Equal(`B`, subject(t))
		})
	})
}
