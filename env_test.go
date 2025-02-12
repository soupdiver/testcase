package testcase_test

import (
	"os"
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/adamluzsi/testcase/fixtures"
	"github.com/adamluzsi/testcase/internal"
)

func TestEnvVarHelpers(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Describe(`#SetEnv`, func(s *testcase.Spec) {
		var (
			tb = s.Let(`TB`, func(t *testcase.T) interface{} {
				return &internal.RecorderTB{TB: &internal.StubTB{}}
			})
			tbCleanupNow = func(t *testcase.T) { tb.Get(t).(*internal.RecorderTB).CleanupNow() }
			key          = s.LetValue(`key`, `TESTING_DATA_`+fixtures.Random.String())
			value        = s.LetValue(`value`, fixtures.Random.String())
			subject      = func(t *testcase.T) {
				testcase.SetEnv(tb.Get(t).(testing.TB), key.Get(t).(string), value.Get(t).(string))
			}
		)

		s.After(func(t *testcase.T) {
			t.Must.Nil(os.Unsetenv(key.Get(t).(string)))
		})

		s.When(`environment key is invalid`, func(s *testcase.Spec) {
			key.LetValue(s, ``)

			s.Then(`it will return with error`, func(t *testcase.T) {
				var finished bool
				internal.RecoverExceptGoexit(func() {
					subject(t)
					finished = true
				})
				t.Must.True(!finished)
				t.Must.True(tb.Get(t).(*internal.RecorderTB).IsFailed)
			})
		})

		s.When(`no environment variable is set before the call`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				t.Must.Nil(os.Unsetenv(key.Get(t).(string)))
			})

			s.Then(`value will be set`, func(t *testcase.T) {
				subject(t)

				v, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(ok)
				t.Must.Equal(v, value.Get(t))
			})

			s.Then(`value will be unset after Cleanup`, func(t *testcase.T) {
				subject(t)
				tbCleanupNow(t)

				_, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(!ok)
			})
		})

		s.When(`environment variable already had a value`, func(s *testcase.Spec) {
			originalValue := s.LetValue(`original value`, fixtures.Random.String())

			s.Before(func(t *testcase.T) {
				t.Must.Nil(os.Setenv(key.Get(t).(string), originalValue.Get(t).(string)))
			})

			s.Then(`new value will be set`, func(t *testcase.T) {
				subject(t)

				v, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(ok)
				t.Must.Equal(v, value.Get(t))
			})

			s.Then(`old value will be restored on Cleanup`, func(t *testcase.T) {
				subject(t)
				tbCleanupNow(t)

				v, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(ok)
				t.Must.Equal(v, originalValue.Get(t))
			})
		})
	})

	s.Describe(`#UnsetEnv`, func(s *testcase.Spec) {
		var (
			tb           = s.Let(`TB`, func(t *testcase.T) interface{} { return &internal.RecorderTB{} })
			tbCleanupNow = func(t *testcase.T) { tb.Get(t).(*internal.RecorderTB).CleanupNow() }
			key          = s.LetValue(`key`, `TESTING_DATA_`+fixtures.Random.String())
			subject      = func(t *testcase.T) {
				testcase.UnsetEnv(tb.Get(t).(testing.TB), key.Get(t).(string))
			}
		)

		s.After(func(t *testcase.T) {
			t.Must.Nil(os.Unsetenv(key.Get(t).(string)))
		})

		s.When(`no environment variable is set before the call`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				t.Must.Nil(os.Unsetenv(key.Get(t).(string)))
			})

			s.Then(`value will be unset after Cleanup`, func(t *testcase.T) {
				subject(t)
				tbCleanupNow(t)

				_, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(!ok)
			})
		})

		s.When(`environment variable already had a value`, func(s *testcase.Spec) {
			originalValue := s.LetValue(`original value`, fixtures.Random.String())

			s.Before(func(t *testcase.T) {
				t.Must.Nil(os.Setenv(key.Get(t).(string), originalValue.Get(t).(string)))
			})

			s.Then(`os env value will be unset`, func(t *testcase.T) {
				subject(t)

				_, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(!ok)
			})

			s.Then(`old value will be restored after the Cleanup`, func(t *testcase.T) {
				subject(t)
				tbCleanupNow(t)

				v, ok := os.LookupEnv(key.Get(t).(string))
				t.Must.True(ok)
				t.Must.Equal(v, originalValue.Get(t))
			})
		})
	})
}
