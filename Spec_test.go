package testcase_test

import (
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamluzsi/testcase/assert"
	"github.com/adamluzsi/testcase/fixtures"
	"github.com/adamluzsi/testcase/internal"
	"github.com/adamluzsi/testcase/internal/spechelper"

	"github.com/adamluzsi/testcase"
)

func TestSpec_DSL(t *testing.T) {
	var sideEffect []string
	var actualSE []string

	valueName := strconv.Itoa(rand.Int())
	nest1Value := rand.Int()
	nest2Value := rand.Int()
	nest3Value := rand.Int()

	s := testcase.NewSpec(t)

	s.Before(func(t *testcase.T) {
		actualSE = make([]string, 0)
	})

	s.Describe(`nest-lvl-1`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) int { return t.I(valueName).(int) }
		s.Let(valueName, func(t *testcase.T) interface{} { return nest1Value })

		s.When(`nest-lvl-2`, func(s *testcase.Spec) {
			s.Let(valueName, func(t *testcase.T) interface{} { return nest2Value })

			s.After(func(t *testcase.T) {
				sideEffect = append(sideEffect, "after1")
			})

			s.Before(func(t *testcase.T) {
				actualSE = append(actualSE, `before1`)
				sideEffect = append(sideEffect, `before1`)
			})

			s.Around(func(t *testcase.T) func() {
				actualSE = append(actualSE, `around1-begin`)
				sideEffect = append(sideEffect, `around1-begin`)
				return func() {
					sideEffect = append(sideEffect, `around1-end`)
				}
			})

			s.And(`nest-lvl-3`, func(s *testcase.Spec) {
				s.Let(valueName, func(t *testcase.T) interface{} { return nest3Value })

				s.After(func(t *testcase.T) {
					sideEffect = append(sideEffect, "after2")
				})

				s.Before(func(t *testcase.T) {
					actualSE = append(actualSE, `before2`)
					sideEffect = append(sideEffect, `before2`)
				})

				s.Around(func(t *testcase.T) func() {
					actualSE = append(actualSE, `around2-begin`)
					sideEffect = append(sideEffect, `around2-begin`)
					return func() {
						sideEffect = append(sideEffect, `around2-end`)
					}
				})

				s.Then(`lvl-3`, func(t *testcase.T) {
					expectedSE := []string{`before1`, `around1-begin`, `before2`, `around2-begin`}
					assert.Must(t).Equal(expectedSE, actualSE)
					assert.Must(t).Equal(nest3Value, t.I(valueName))
					assert.Must(t).Equal(nest3Value, subject(t))
				})
			})

			s.Then(`lvl-2`, func(t *testcase.T) {
				expectedSE := []string{`before1`, `around1-begin`}
				assert.Must(t).Equal(expectedSE, actualSE)
				assert.Must(t).Equal(nest2Value, t.I(valueName))
				assert.Must(t).Equal(nest2Value, subject(t))
			})
		})

		s.Then(`lvl-1`, func(t *testcase.T) {
			expectedSE := []string{}
			assert.Must(t).Equal(expectedSE, actualSE)
			t.Must.Equal(nest1Value, t.I(valueName))
			t.Must.Equal(nest1Value, subject(t))
		})
	})
}

func TestSpec_subSpecIsExecuted(t *testing.T) {
	var ran bool
	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)
		s.Describe(`nest-lvl-1`, func(s *testcase.Spec) {
			s.When(`nest-lvl-2`, func(s *testcase.Spec) {
				s.And(`nest-lvl-3`, func(s *testcase.Spec) {
					s.Then(`lvl-3`, func(t *testcase.T) {
						ran = true
					})
				})
			})
		})
	})
	assert.Must(t).True(ran)
}

func TestSpec_Context(t *testing.T) {

	var allSideEffect [][]string
	var sideEffect []string

	valueName := strconv.Itoa(rand.Int())
	nest1Value := rand.Int()
	nest2Value := rand.Int()
	nest3Value := rand.Int()

	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)

		s.Around(func(t *testcase.T) func() {
			sideEffect = make([]string, 0)
			return func() { allSideEffect = append(allSideEffect, sideEffect) }
		})

		s.Context(`nest-lvl-1`, func(s *testcase.Spec) {
			subject := func(t *testcase.T) int { return t.I(valueName).(int) }
			s.Let(valueName, func(t *testcase.T) interface{} { return nest1Value })

			s.Context(`nest-lvl-2`, func(s *testcase.Spec) {
				s.Let(valueName, func(t *testcase.T) interface{} { return nest2Value })

				s.After(func(t *testcase.T) {
					sideEffect = append(sideEffect, "after1")
				})

				s.Before(func(t *testcase.T) {
					sideEffect = append(sideEffect, `before1`)
				})

				s.Around(func(t *testcase.T) func() {
					sideEffect = append(sideEffect, `around1-begin`)
					return func() { sideEffect = append(sideEffect, `around1-end`) }
				})

				s.Context(`nest-lvl-3`, func(s *testcase.Spec) {
					s.Let(valueName, func(t *testcase.T) interface{} { return nest3Value })

					s.After(func(t *testcase.T) {
						sideEffect = append(sideEffect, "after2")
					})

					s.Before(func(t *testcase.T) {
						sideEffect = append(sideEffect, `before2`)
					})

					s.Around(func(t *testcase.T) func() {
						sideEffect = append(sideEffect, `around2-begin`)
						return func() { sideEffect = append(sideEffect, `around2-end`) }
					})

					s.Test(`lvl-3`, func(t *testcase.T) {
						t.Must.Equal([]string{`before1`, `around1-begin`, `before2`, `around2-begin`}, sideEffect)
						t.Must.Equal(nest3Value, t.I(valueName))
						t.Must.Equal(nest3Value, subject(t))
					})
				})

				s.Test(`lvl-2`, func(t *testcase.T) {
					t.Must.Equal([]string{`before1`, `around1-begin`}, sideEffect)
					t.Must.Equal(nest2Value, t.I(valueName))
					t.Must.Equal(nest2Value, subject(t))
				})
			})

			s.Test(`lvl-1`, func(t *testcase.T) {
				t.Must.Equal([]string{}, sideEffect)
				t.Must.Equal(nest1Value, t.I(valueName))
				t.Must.Equal(nest1Value, subject(t))
			})
		})
	})

	//t.Logf(`%#v`, allSideEffect)
	assert.Must(t).ContainExactly([][]string{
		{},
		{"before1", "around1-begin", "around1-end", "after1"},
		{"before1", "around1-begin", "before2", "around2-begin", "around2-end", "after2", "around1-end", "after1"},
	}, allSideEffect)
}

func TestSpec_Describe_executedAsAGroupInTheEndOfThe(t *testing.T) {
	var (
		ran  bool
		once int
	)
	s := testcase.NewSpec(t)
	s.Describe(`executed in the end of the Describe block`, func(s *testcase.Spec) {
		s.Test(``, func(t *testcase.T) {
			ran = true
			once++
		})
	})
	assert.Must(t).True(ran)
	assert.Must(t).Equal(1, once)
}

func TestSpec_ParallelSafeVariableSupport(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()

	valueName := strconv.Itoa(rand.Int())
	nest1Value := rand.Int()
	nest2Value := rand.Int()
	nest3Value := rand.Int()

	s.Describe(`nest-lvl-1`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) int { return t.I(valueName).(int) }
		s.Let(valueName, func(t *testcase.T) interface{} { return nest1Value })

		s.When(`nest-lvl-2`, func(s *testcase.Spec) {
			s.Let(valueName, func(t *testcase.T) interface{} { return nest2Value })

			s.And(`nest-lvl-3`, func(s *testcase.Spec) {
				s.Let(valueName, func(t *testcase.T) interface{} { return nest3Value })

				s.Test(`lvl-3`, func(t *testcase.T) {
					t.Must.Equal(nest3Value, t.I(valueName))
					t.Must.Equal(nest3Value, subject(t))
				})
			})

			s.Test(`lvl-2`, func(t *testcase.T) {
				t.Must.Equal(nest2Value, t.I(valueName))
				t.Must.Equal(nest2Value, subject(t))
			})
		})

		s.Test(`lvl-1`, func(t *testcase.T) {
			t.Must.Equal(nest1Value, t.I(valueName))
			t.Must.Equal(nest1Value, subject(t))
		})
	})
}

func TestSpec_InvalidUsages(t *testing.T) {
	stub := &internal.StubTB{}
	s := testcase.NewSpec(stub)
	valueName := strconv.Itoa(rand.Int())
	nest1Value := rand.Int()
	nest2Value := rand.Int()
	nest3Value := rand.Int()

	willFatal := isFatalFn(stub)

	failNowSpecs := func(t *testing.T, s *testcase.Spec, expectedToFailNow bool) {
		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.Before(func(t *testcase.T) {})
		}))

		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.After(func(t *testcase.T) {})
		}))

		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.Around(func(t *testcase.T) func() { return func() {} })
		}))

		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.Let(strconv.Itoa(rand.Int()), func(t *testcase.T) interface{} { return nil })
		}))

		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.Parallel()
		}))

		assert.Must(t).Equal(expectedToFailNow, willFatal(func() {
			s.LetValue(`value`, rand.Int())
		}))
	}

	shouldFailNowForHooking := func(t *testing.T, s *testcase.Spec) { failNowSpecs(t, s, true) }
	shouldNotFailForHooking := func(t *testing.T, s *testcase.Spec) { failNowSpecs(t, s, false) }

	topSpec := s

	s.Describe(`nest-lvl-1`, func(s *testcase.Spec) {

		shouldFailNowForHooking(t, topSpec)

		shouldNotFailForHooking(t, s)
		s.Let(valueName, func(t *testcase.T) interface{} { return nest1Value })

		shouldNotFailForHooking(t, s)
		s.Test(`lvl-1-first`, func(t *testcase.T) {})

		shouldFailNowForHooking(t, s)

		s.When(`nest-lvl-2`, func(s *testcase.Spec) {
			shouldNotFailForHooking(t, s)
			s.Let(valueName, func(t *testcase.T) interface{} { return nest2Value })

			shouldNotFailForHooking(t, s)
			s.And(`nest-lvl-3`, func(s *testcase.Spec) {
				s.Let(valueName, func(t *testcase.T) interface{} { return nest3Value })

				s.Test(`lvl-3`, func(t *testcase.T) {})

				shouldFailNowForHooking(t, s)
			})

			shouldFailNowForHooking(t, s)

			s.Test(`lvl-2`, func(t *testcase.T) {})

			shouldFailNowForHooking(t, s)

			s.And(`nest-lvl-2-2`, func(s *testcase.Spec) {
				shouldNotFailForHooking(t, s)
				s.Test(`nest-lvl-2-2-then`, func(t *testcase.T) {})
				shouldFailNowForHooking(t, s)
			})

			shouldFailNowForHooking(t, s)

		})

		shouldFailNowForHooking(t, s)

		s.Test(`lvl-1-last`, func(t *testcase.T) {})

		shouldFailNowForHooking(t, s)

	})
}

func TestSpec_FriendlyVarNotDefined(t *testing.T) {
	stub := &internal.StubTB{}
	s := testcase.NewSpec(stub)
	willFatalWithMessage := willFatalWithMessageFn(stub)

	s.Let(`var1`, func(t *testcase.T) interface{} { return `hello-world` })
	s.Let(`var2`, func(t *testcase.T) interface{} { return `hello-world` })
	tct := testcase.NewT(stub, s)

	t.Run(`var1 var found`, func(t *testing.T) {
		assert.Must(t).Equal(`hello-world`, tct.I(`var1`).(string))
	})

	t.Run(`not existing var will panic with friendly msg`, func(t *testing.T) {
		panicMSG := willFatalWithMessage(t, func() { tct.I(`not-exist`) })
		msg := strings.Join(panicMSG, " ")
		assert.Must(t).Contain(msg, `Variable "not-exist" is not found`)
		assert.Must(t).Contain(msg, `Did you mean?`)
		assert.Must(t).Contain(msg, `var1`)
		assert.Must(t).Contain(msg, `var2`)
	})
}

func TestSpec_Let_valuesAreDeterministicallyCached(t *testing.T) {
	s := testcase.NewSpec(t)

	var testCase1Value int
	var testCase2Value int

	type TestStruct struct {
		Value string
	}

	s.Describe(`Let`, func(s *testcase.Spec) {
		s.Let(`int`, func(t *testcase.T) interface{} { return rand.Int() })

		s.Then(`regardless of multiple call, let value remain the same for each`, func(t *testcase.T) {
			value := t.I(`int`).(int)
			testCase1Value = value
			t.Must.Equal(value, t.I(`int`).(int))
		})

		s.Then(`for every then block then block value is reevaluated`, func(t *testcase.T) {
			value := t.I(`int`).(int)
			testCase2Value = value
			t.Must.Equal(value, t.I(`int`).(int))
		})

		s.And(`the value is accessible from the hooks as well`, func(s *testcase.Spec) {
			var value int

			s.Before(func(t *testcase.T) {
				value = t.I(`int`).(int)
			})

			s.Then(`it will remain the same value in the test case as well compared to the before block`, func(t *testcase.T) {
				t.Must.NotEqual(0, value)
				t.Must.Equal(value, t.I(`int`).(int))
			})
		})

		s.And(`struct value can be modified by hooks for preparation purposes like setting up mocks expectations`, func(s *testcase.Spec) {
			s.Let(`struct`, func(t *testcase.T) interface{} {
				return &TestStruct{}
			})

			s.Before(func(t *testcase.T) {
				value := t.I(`struct`).(*TestStruct)
				value.Value = "testing"
			})

			s.Then(`the value can be seen from the test case scope`, func(t *testcase.T) {
				t.Must.Equal(`testing`, t.I(`struct`).(*TestStruct).Value)
			})
		})
	})

	assert.Must(t).NotEqual(testCase1Value, testCase2Value)
}

func TestSpec_Let_valueScopesAppliedOnHooks(t *testing.T) {
	s := testcase.NewSpec(t)

	var leaker int
	s.Context(`1`, func(s *testcase.Spec) {
		s.Before(func(t *testcase.T) {
			leaker = t.I(`value`).(int)
		})

		s.Let(`value`, func(t *testcase.T) interface{} {
			return 24
		})

		s.Context(`2`, func(s *testcase.Spec) {
			s.Let(`value`, func(t *testcase.T) interface{} {
				return 42
			})

			s.Test(`testCase`, func(t *testcase.T) {
				t.Must.Equal(42, leaker)
			})
		})
	})

}

func hackCallParallel(tb testing.TB) {
	switch tb := tb.(type) {
	case *testing.T:
		tb.Parallel()
	case *internal.RecorderTB:
		hackCallParallel(tb.TB)
	case *testcase.T:
		hackCallParallel(tb.TB)
	default:
		tb.Fatalf(`%T don't implement #Parallel`, tb)
	}
}

func TestSpec_Parallel(t *testing.T) {
	s := testcase.NewSpec(t)

	isPanic := func(block func()) (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		block()
		return false
	}

	s.Describe(`Parallel`, func(s *testcase.Spec) {

		s.When(`no parallel set on top level nesting`, func(s *testcase.Spec) {
			s.And(`on each sub level`, func(s *testcase.Spec) {
				s.Then(`it will accept T#Parallel call`, func(t *testcase.T) {
					assert.Must(t).True(!isPanic(func() { hackCallParallel(t.TB) }))
				})
			})
			s.Then(`it will accept T#Parallel call`, func(t *testcase.T) {
				assert.Must(t).True(!isPanic(func() { hackCallParallel(t.TB) }))
			})
		})

		s.When(`on the first level there is no parallel configured`, func(s *testcase.Spec) {
			s.And(`on the second one, yes`, func(s *testcase.Spec) {
				s.Parallel()

				s.And(`parallel will be "inherited" for each nested spec`, func(s *testcase.Spec) {
					s.Then(`it will panic on T#Parallel call`, func(t *testcase.T) {
						assert.Must(t).True(isPanic(func() { hackCallParallel(t.TB) }))
					})
				})

				s.Then(`it panic on T#Parallel call`, func(t *testcase.T) {
					assert.Must(t).True(isPanic(func() { hackCallParallel(t.TB) }))
				})
			})

			s.Then(`it will accept T#Parallel call`, func(t *testcase.T) {
				assert.Must(t).True(!isPanic(func() { hackCallParallel(t.TB) }))
			})

		})

	})
}

func TestSpec_testsWithName_shouldRun(t *testing.T) {
	var a, b bool
	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)

		s.Test(``, func(t *testcase.T) {
			a = true
		})

		s.Test(``, func(t *testcase.T) {
			b = true
		})
	})

	assert.Must(t).True(a)
	assert.Must(t).True(b)
}

func TestSpec_NoSideEffect(t *testing.T) {
	s := testcase.NewSpec(t)

	isPanic := func(block func()) (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		block()
		return false
	}

	s.Describe(`NoSideEffect`, func(s *testcase.Spec) {
		s.When(`on the first level there is no parallel configured`, func(s *testcase.Spec) {
			s.And(`on the second one, yes`, func(s *testcase.Spec) {
				s.NoSideEffect()

				s.And(`parallel will be "inherited" for each nested spec`, func(s *testcase.Spec) {
					s.Then(`it will panic on T#parallel call`, func(t *testcase.T) {
						assert.Must(t).True(isPanic(func() { hackCallParallel(t.TB) }))
					})
				})

				s.Then(`it panic on T#parallel call`, func(t *testcase.T) {
					assert.Must(t).True(isPanic(func() { hackCallParallel(t.TB) }))
				})
			})

			s.Then(`it will accept T#parallel call`, func(t *testcase.T) {
				assert.Must(t).True(!isPanic(func() { hackCallParallel(t.TB) }))
			})

		})
	})
}

func TestSpec_Let_FallibleValue(t *testing.T) {
	s := testcase.NewSpec(t)

	s.Let(`fallible`, func(t *testcase.T) interface{} {
		return t.TB
	})

	s.Then(`fallible receive the same testing object as this spec`, func(t *testcase.T) {
		t.Must.Equal(t.TB, t.I(`fallible`))
	})
}

func TestSpec_LetValue_ValueDefinedAtDeclarationWithoutTheNeedOfFunctionCallback(t *testing.T) {
	s := testcase.NewSpec(t)

	s.LetValue(`value`, 42)

	s.Then(`the testCase variable will be accessible`, func(t *testcase.T) {
		t.Must.Equal(42, t.I(`value`))
	})

	for kind, example := range map[reflect.Kind]interface{}{
		reflect.String:     "hello world",
		reflect.Bool:       true,
		reflect.Int:        int(42),
		reflect.Int8:       int8(42),
		reflect.Int16:      int16(42),
		reflect.Int32:      int32(42),
		reflect.Int64:      int64(42),
		reflect.Uint:       uint(42),
		reflect.Uint8:      uint8(42),
		reflect.Uint16:     uint16(42),
		reflect.Uint32:     uint32(42),
		reflect.Uint64:     uint64(42),
		reflect.Float32:    float32(42),
		reflect.Float64:    float64(42),
		reflect.Complex64:  complex64(42),
		reflect.Complex128: complex128(42),
	} {
		kind := kind
		example := example

		s.Context(kind.String(), func(s *testcase.Spec) {
			s.LetValue(kind.String(), example)

			s.Then(`it will return the value`, func(t *testcase.T) {
				t.Must.Equal(example, t.I(kind.String()))
			})
		})
	}
}

func TestSpec_LetValue_mutableValuesAreNotAllowed(t *testing.T) {
	stub := &internal.StubTB{}
	s := testcase.NewSpec(stub)

	var finished bool
	internal.RecoverExceptGoexit(func() {
		type SomeStruct struct {
			Text string
		}
		s.LetValue(`mutable values are not allowed`, &SomeStruct{Text: `hello world`})
		finished = true
	})
	assert.Must(t).True(!finished)
	assert.Must(t).True(stub.IsFailed)
}

func TestSpec_Cleanup_inACleanupWithinACleanup(t *testing.T) {
	t.Run(`spike`, func(t *testing.T) {
		var ran bool
		t.Run(``, func(t *testing.T) {
			t.Cleanup(func() {
				t.Cleanup(func() {
					t.Cleanup(func() {
						ran = true
					})
				})
			})
		})
		assert.Must(t).True(ran)
	})

	var ran bool
	s := testcase.NewSpec(t)
	s.After(func(t *testcase.T) {
		t.Cleanup(func() {
			t.Cleanup(func() {
				ran = true
			})
		})
	})
	s.Test(``, func(t *testcase.T) {})

	s.Finish()

	assert.Must(t).True(ran)
}

func BenchmarkTest_Spec(b *testing.B) {
	b.Log(`this is actually a testCase`)
	b.Log(`it will run a bench *testing.B.N times`)

	var total int
	b.Run(``, func(b *testing.B) {
		s := testcase.NewSpec(b)
		s.Test(``, func(t *testcase.T) {
			time.Sleep(time.Millisecond)
			total++
		})
	})

	assert.Must(b).True(1 < total)
}

func BenchmarkTest_Spec_eachBenchmarkingRunsWithFreshState(b *testing.B) {
	b.Log(`this is actually a testCase`)
	b.Log(`it will run a bench *testing.B.N times but in parallel with b.RunParallel`)

	s := testcase.NewSpec(b)

	type mutable struct{ used bool }
	s.Let(`mutable`, func(t *testcase.T) interface{} {
		return &mutable{used: false}
	})

	s.Before(func(t *testcase.T) {
		assert.Must(t).True(!t.I(`mutable`).(*mutable).used)
	})

	b.Log(`each benchmarking runs with fresh state to avoid side effects between bench mark iterations`)
	s.Test(``, func(t *testcase.T) {
		// A bit sleeping here makes measuring the average runtime speed really really really easy and much faster in general.
		// else the value would be so small, that it becomes difficult for the testing package benchmark suite to measure it with small number of samplings.
		time.Sleep(time.Millisecond)

		m := t.I(`mutable`).(*mutable)
		assert.Must(t).True(!m.used)
		m.used = true
	})
}

type UnknownTestingTB struct {
	testing.TB
	logs [][]interface{}
}

func (tb *UnknownTestingTB) Log(args ...interface{}) {
	tb.logs = append(tb.logs, args)
}

func TestSpec_Test_withSomethingThatImplementsTestcaseTB(t *testing.T) {
	rtb := &internal.RecorderTB{TB: &internal.StubTB{}}

	var tb testcase.TBRunner = rtb // implements check
	s := testcase.NewSpec(tb)

	s.Test(`passthrough`, func(t *testcase.T) {
		t.FailNow()
	})

	rtb.CleanupNow()
	assert.Must(t).True(rtb.IsFailed)
}

func TestSpec_Sequential(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	spechelper.OrderAsDefined(t)
	s := testcase.NewSpec(t)

	// --- somewhere in a spec helper setup function --- //

	// Sequential is used to ensure testing executed sequential,
	// because here we define a variable that meant to be used sequentially only.
	s.Sequential()

	// --- in a spec file --- //

	// somewhere else in a spec where the code itself has no side effect
	// we use parallel to allow the possibility of run testCase concurrently.
	s.Parallel()
	var bTestRan bool

	s.Test(`A`, func(t *testcase.T) {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
		assert.Must(t).True(!bTestRan,
			`testCase A ran in parallel with testCase B, but this was not expected after using testcase.Spec#Sequence`)
	})

	s.Test(`B`, func(t *testcase.T) {
		bTestRan = true
	})
}

func TestSpec_HasSideEffect(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	spechelper.OrderAsDefined(t)

	s := testcase.NewSpec(t)
	s.HasSideEffect()
	s.NoSideEffect()
	var bTestRan bool

	s.Test(`A`, func(t *testcase.T) {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
		assert.Must(t).True(!bTestRan,
			`testCase A ran in parallel with testCase B, but this was not expected after using testcase.Spec#Sequence`)
	})

	s.Test(`B`, func(t *testcase.T) {
		bTestRan = true
	})
}

func TestSpec_Sequential_scoped(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	spechelper.OrderAsDefined(t)

	s := testcase.NewSpec(t)
	s.Parallel()

	s.Context(`scope A`, func(s *testcase.Spec) {
		s.Sequential()

		s.Context(`sub-scope of A`, func(s *testcase.Spec) {
			s.Parallel()
			var bTestRan bool

			s.Test(`A`, func(t *testcase.T) {
				runtime.Gosched()
				time.Sleep(time.Millisecond)
				assert.Must(t).True(!bTestRan,
					`testCase A ran in parallel with testCase B, but this was not expected after using testcase.Spec#Sequence`)
			})

			s.Test(`B`, func(t *testcase.T) {
				bTestRan = true
			})
		})
	})

	s.Context(`scope B`, func(s *testcase.Spec) {
		var wg = &sync.WaitGroup{}
		wg.Add(1)

		s.Test(`A`, func(t *testcase.T) {
			runtime.Gosched()
			c := make(chan struct{})
			go func() {
				defer close(c)
				wg.Wait() // if wait is done, we are in the happy group
			}()

			select {
			case <-c: // happy case
			case <-time.After(time.Second):
				t.Fatal(`B testCase probably not running concurrently`) // timed orderingOutput
			}
		})

		s.Test(`B`, func(t *testcase.T) {
			defer wg.Done()
		})
	})
}

func TestSpec_Sequential_callingItAfterContextDeclarationYieldPanic(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Context(``, func(s *testcase.Spec) {})
	assert.Must(t).Panic(func() { s.HasSideEffect() })
}

func TestSpec_Skip(t *testing.T) {
	var out []int

	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)
		s.Sequential()

		s.Context(`skipped ones`, func(s *testcase.Spec) {
			s.Skip(`WIP or something like that`)

			s.Test(`will be skipped`, func(t *testcase.T) {
				out = append(out, 0)
			})

			s.Test(`will be skipped as well`, func(t *testcase.T) {
				out = append(out, 1)
			})

			s.Context(`skipped as well just like parent tests`, func(s *testcase.Spec) {

				s.Test(`will be skipped`, func(t *testcase.T) {
					out = append(out, 0)
				})

			})
		})

		s.Test(`will run`, func(t *testcase.T) {
			out = append(out, 42)
		})
	})

	assert.Must(t).Equal([]int{42}, out)
}

func TestSpec_panicDoNotLeakOutFromTestingScope(t *testing.T) {
	var noPanic bool
	func() {
		defer recover()
		rtb := &internal.RecorderTB{TB: &internal.StubTB{}}
		defer rtb.CleanupNow()
		s := testcase.NewSpec(rtb)
		s.Test(``, func(t *testcase.T) { panic(`die`) })
		s.Test(``, func(t *testcase.T) { noPanic = true })
	}()
	assert.Must(t).True(noPanic)
}

func TestSpec_panicDoNotLeakOutFromTestingScope_poc(t *testing.T) {
	t.Skip(`POC for manual testing`)
	s := testcase.NewSpec(t)
	s.Test(``, func(t *testcase.T) { panic(`die`) })
	s.Test(``, func(t *testcase.T) { t.Log(`OK`) })
}

func BenchmarkTest_Spec_hooksInBenchmarkCalledInEachRun(b *testing.B) {

	var (
		beforeTimes int
		deferTimes  int
		afterTimes  int
		testTimes   int
	)

	b.Run(``, func(b *testing.B) {
		s := testcase.NewSpec(b)
		s.Sequential()

		s.Before(func(t *testcase.T) { t.Defer(func() { deferTimes++ }) })
		s.Before(func(t *testcase.T) { beforeTimes++ })
		s.After(func(t *testcase.T) { afterTimes++ })

		var flag bool
		s.Around(func(t *testcase.T) func() {
			return func() { flag = false }
		})

		s.Test(``, func(t *testcase.T) {
			assert.Must(t).True(!flag)
			flag = true // mutate so we expect After to restore the flag state to "false"

			testTimes++
			time.Sleep(time.Millisecond) // A bit sleep here helps the benchmark to make a faster conclusion.
		})
	})

	assert.Must(b).NotEqual(0, testTimes)
	assert.Must(b).Equal(testTimes, beforeTimes)
	assert.Must(b).Equal(testTimes, afterTimes)
	assert.Must(b).Equal(testTimes, deferTimes)
}

func TestSpec_hooksAlignWithCleanup(t *testing.T) {
	var afters []string
	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)

		s.After(func(t *testcase.T) {
			afters = append(afters, `First After`)
		})

		s.Before(func(t *testcase.T) {
			t.Defer(func() { afters = append(afters, `Defer`) })
		})

		s.Before(func(t *testcase.T) {
			t.Cleanup(func() { afters = append(afters, `Cleanup`) })
		})

		s.After(func(t *testcase.T) {
			afters = append(afters, `Last After`)
		})

		s.Test(``, func(t *testcase.T) {})
	})
	assert.Must(t).Equal([]string{`Last After`, `Cleanup`, `Defer`, `First After`}, afters)
}

func BenchmarkTest_Spec_SkipBenchmark1(b *testing.B) {
	var (
		allowedTestRan   bool
		forbiddenTestRan bool
	)
	b.Run(``, func(b *testing.B) {
		s := testcase.NewSpec(b)
		s.Test(``, func(t *testcase.T) {
			allowedTestRan = true
			t.SkipNow()
		})

		s.Test(``, func(t *testcase.T) {
			forbiddenTestRan = true
			t.SkipNow()
		}, testcase.SkipBenchmark())

		s.Then(``, func(t *testcase.T) {
			forbiddenTestRan = true
			t.SkipNow()
		}, testcase.SkipBenchmark())
	})

	assert.Must(b).True(allowedTestRan)
	assert.Must(b).True(!forbiddenTestRan)
}

func BenchmarkTest_Spec_SkipBenchmark2(b *testing.B) {
	var (
		allowedTestRan   bool
		forbiddenTestRan bool
	)

	b.Run(``, func(b *testing.B) {
		s := testcase.NewSpec(b)
		s.Context(``, func(s *testcase.Spec) {
			s.Test(``, func(t *testcase.T) {
				allowedTestRan = true
				t.SkipNow()
			})
		})

		s.Context(``, func(s *testcase.Spec) {
			s.SkipBenchmark()

			s.Test(``, func(t *testcase.T) {
				forbiddenTestRan = true
				t.SkipNow()
			})
		})
	})

	assert.Must(b).True(allowedTestRan)
	assert.Must(b).True(!forbiddenTestRan)
}

type stubB struct {
	*internal.StubTB
	TestingB *testing.B
}

func (b *stubB) Run(name string, fn func(b *testing.B)) bool {
	return b.TestingB.Run(name, fn)
}

func BenchmarkTest_Spec_SkipBenchmark_invalidUse(b *testing.B) {
	stub := &internal.StubTB{}
	stb := &stubB{
		StubTB:   stub,
		TestingB: b,
	}
	s := testcase.NewSpec(stb)

	s.Test(``, func(t *testcase.T) { t.SkipNow() })

	var finished bool
	internal.RecoverExceptGoexit(func() {
		s.SkipBenchmark()
		finished = false
	})
	assert.Must(b).True(!finished)
	assert.Must(b).True(stub.IsFailed)
	assert.Must(b).Contain(stub.Logs, "you can't use .SkipBenchmark after you already used when/and/then")
}

func BenchmarkTest_Spec_Test_flaky(b *testing.B) {
	s := testcase.NewSpec(b)
	var hasRun bool
	s.Test(``, func(t *testcase.T) {
		hasRun = true
		t.SkipNow()
	}, testcase.Flaky(time.Second))
	assert.Must(b).True(!hasRun)
}

func TestSpec_Test_FailNowWithCustom(t *testing.T) {
	rtb := &internal.RecorderTB{TB: &internal.StubTB{}}
	s := testcase.NewSpec(rtb)

	var failCount int
	s.Test(``, func(t *testcase.T) {
		failCount++
		t.FailNow()
	})

	rtb.CleanupNow()
	assert.Must(t).Equal(1, failCount)
	assert.Must(t).True(rtb.IsFailed)
}

func TestSpec_Test_flaky_withoutFlakyFlag_willFailAndNeverRunAgain(t *testing.T) {
	stub := &internal.StubTB{}
	s := testcase.NewSpec(stub)
	var total int
	s.Test(``, func(t *testcase.T) { total++; t.FailNow() })
	stub.Finish()
	assert.Must(t).Equal(1, total)
}

func TestSpec_Test_flakyByTimeout_willRunAgainWithinTheTimeoutDurationUntilItPasses(t *testing.T) {
	s := testcase.NewSpec(t)

	var failedOnce bool
	s.Test(``, func(t *testcase.T) {
		if failedOnce {
			return
		}

		failedOnce = true
		t.FailNow()
	}, testcase.Flaky(2))
}

func TestSpec_Test_flakyByRetry_willRunAgainWithTheProvidedRetry(t *testing.T) {
	s := testcase.NewSpec(t)

	var retryUsed bool
	retry := testcase.Retry{
		Strategy: testcase.RetryStrategyFunc(func(condition func() bool) {
			retryUsed = true
			for condition() {
			}
		}),
	}

	var failedOnce bool
	s.Test(``, func(t *testcase.T) {
		if failedOnce {
			return
		}
		failedOnce = true
		t.FailNow()
	}, testcase.Flaky(retry))
	s.Finish()
	assert.Must(t).True(retryUsed)
}

func TestSpec_Test_flakyByStrategy_willRunAgainBasedOnTheStrategy(t *testing.T) {
	s := testcase.NewSpec(t)

	var strategyCallCount, testCount int
	strategy := testcase.RetryStrategyFunc(func(condition func() bool) {
		for condition() {
			strategyCallCount++
		}
	})

	s.Test(``, func(t *testcase.T) {
		testCount++
		if fixtures.Random.Bool() {
			t.FailNow()
		}
	}, testcase.Flaky(strategy))

	assert.Must(t).Equal(strategyCallCount, testCount)
}

func TestSpec_Test_flakyFlagWithInvalidValue_willPanics(t *testing.T) {
	s := testcase.NewSpec(t)
	assert.Must(t).Panic(func() { testcase.Flaky("foo") })
	assert.Must(t).Panic(func() { s.Test(``, func(t *testcase.T) {}, testcase.Flaky("foo")) })
}

func TestSpec_Test_flakyByRetryCount_willRunAgainWithinTheAcceptedRetryCount(t *testing.T) {
	s := testcase.NewSpec(t)

	var failedOnce bool
	s.Test(``, func(t *testcase.T) {
		if failedOnce {
			return
		}

		failedOnce = true
		t.FailNow()
	}, testcase.Flaky(42))
}

// This testCase will artificially create a scenario where one of the before block will be held up,
// and the other testCase is expected to finishNow ahead of time.
// If the preparation is not done concurrently as well,
// then the testCase will panic with the reason for failure.
// I know, panic not an ideal way to represent failed testCase, but this approach is deterministic.
func TestSpec_Parallel_testPrepareActionsExecutedInParallel(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()

	s.Around(func(t *testcase.T) func() {
		timer := time.NewTimer(time.Second)
		go func() {
			if _, ok := <-timer.C; ok {
				panic(`it was expected that #Before run parallel as well in case of Spec#Parallel is used`)
			}
		}()
		return func() { timer.Stop() }
	})

	total := 2
	var wg sync.WaitGroup
	wg.Add(total)
	s.Before(func(t *testcase.T) {
		wg.Done() // check in that we race tests ready for the Finish
		wg.Wait() // wait for the Finish OR stuck on DEADLOCK if execution is not parallel
	})
	if runtime.NumCPU() < total {
		t.Skip(`test require at least 2 CPU core to able to run concurrently`)
	}
	for i := 0; i < total; i++ {
		s.Test(``, func(t *testcase.T) {})
	}
}

func TestSpec_executionOrder(t *testing.T) {
	t.Skip(`WIP`)

	t.Run(`Non parallel testCase will run in randomized order`, func(t *testing.T) {
		testcase.Retry{Strategy: testcase.Waiter{WaitDuration: time.Second}}.Assert(t, func(tb testing.TB) {
			var m sync.Mutex
			total := fixtures.Random.IntBetween(32, 128)
			out := make([]int, 0, total)
			s := testcase.NewSpec(tb)

			s.Describe(``, func(s *testcase.Spec) {
				// No Parallel flag
				for j := 0; j < total; j++ {
					v := j // pass by value
					s.Test(``, func(t *testcase.T) {
						m.Lock()
						defer m.Unlock()
						out = append(out, v)
					})
				}
			})

			assert.Must(tb).True(!sort.IsSorted(sort.IntSlice(out)))
		})
	})
}

func TestSpec_Finish(t *testing.T) {
	s := testcase.NewSpec(t)
	var ran bool
	s.Test(``, func(t *testcase.T) { ran = true })
	assert.Must(t).True(!ran)
	s.Finish()
	assert.Must(t).True(ran)
}

func TestSpec_Finish_finishedSpecIsImmutable(t *testing.T) {
	stub := &internal.StubTB{}
	s := testcase.NewSpec(stub)
	s.Before(func(t *testcase.T) {})
	s.Finish()

	var finished bool
	internal.RecoverExceptGoexit(func() {
		s.Before(func(t *testcase.T) {})
		finished = true
	})
	assert.Must(t).True(!finished, `it was expected that FailNow prevents finishing of the process`)
	assert.Must(t).True(stub.IsFailed, `it was expected that the test fails`)
}

func TestSpec_Finish_runOnlyOnce(t *testing.T) {
	s := testcase.NewSpec(t)
	var count int
	s.Test(``, func(t *testcase.T) { count++ })

	assert.Must(t).Equal(0, count)
	s.Finish()
	assert.Must(t).Equal(1, count)
	s.Finish()
	assert.Must(t).Equal(1, count, `should not repeat the test execution`)
}

func TestSpec_eachContextRunsOnce(t *testing.T) {
	var (
		testInDescribe int
		testInContext  int
		testOnTopLevel int
	)

	t.Run(``, func(t *testing.T) {
		s := testcase.NewSpec(t)
		s.Describe(``, func(s *testcase.Spec) {
			s.Test(``, func(t *testcase.T) { testInDescribe++ })
		})
		s.Context(``, func(s *testcase.Spec) {
			s.Test(``, func(t *testcase.T) { testInContext++ })
		})
		s.Test(``, func(t *testcase.T) { testOnTopLevel++ })
	})

	assert.Must(t).Equal(1, testInDescribe)
	assert.Must(t).Equal(1, testInContext)
	assert.Must(t).Equal(1, testOnTopLevel)
}

func TestSpec_Finish_describeBlocksRunWhenTheyCloseAndNotAfter(t *testing.T) {
	var (
		testInDescribe int
		testOnTopLevel int
	)

	s := testcase.NewSpec(t)
	s.Describe(``, func(s *testcase.Spec) {
		s.Test(``, func(t *testcase.T) { testInDescribe++ })
	})
	s.Test(``, func(t *testcase.T) { testOnTopLevel++ })

	assert.Must(t).Equal(1, testInDescribe)
	assert.Must(t).Equal(0, testOnTopLevel)
	s.Finish()
	assert.Must(t).Equal(1, testInDescribe)
	assert.Must(t).Equal(1, testOnTopLevel)
}

func TestSpec_Describe_withCustomTB(t *testing.T) {
	var ran bool
	s := testcase.NewSpec(&CustomTB{TB: t})
	s.Describe(`subcontext`, func(s *testcase.Spec) {
		s.Test(``, func(t *testcase.T) { ran = true })
	})
	assert.Must(t).True(ran)
}

func BenchmarkTest_Spec_Describe(b *testing.B) {
	b.Run(`*testing.B`, func(b *testing.B) {
		var ran bool
		s := testcase.NewSpec(b)
		s.Describe(``, func(s *testcase.Spec) {
			s.Test(``, func(t *testcase.T) {
				time.Sleep(time.Millisecond)
				ran = true
			})
		})
		assert.Must(b).True(ran)
	})
	b.Run(`withCustomTB`, func(b *testing.B) {
		var ran bool
		s := testcase.NewSpec(&CustomTB{TB: b})
		s.Describe(``, func(s *testcase.Spec) {
			s.Test(``, func(t *testcase.T) {
				time.Sleep(time.Millisecond)
				ran = true
			})
		})
		assert.Must(b).True(ran)
	})
}

func BenchmarkTest_Spec_Describe_withCustomTB(b *testing.B) {
	var ran bool
	s := testcase.NewSpec(&CustomTB{TB: b})
	s.Describe(`subcontext`, func(s *testcase.Spec) {
		s.Test(``, func(t *testcase.T) {
			ran = true
			t.SkipNow()
		})
	})
	assert.Must(b).True(ran)
}

func TestSpec_Describe_withSomeTestRunner(t *testing.T) {
	type SomeTestTB struct{ testing.TB }
	var ran bool
	s := testcase.NewSpec(SomeTestTB{TB: t})
	s.Describe(`subcontext`, func(s *testcase.Spec) {
		s.Test(``, func(t *testcase.T) { ran = true })
	})
	s.Finish()
	assert.Must(t).True(ran)
}

func TestNewSpec_withTestingT_optionsPassed(t *testing.T) {
	s := testcase.NewSpec(t, testcase.Flaky(time.Second))
	s.Test(``, func(t *testcase.T) {
		if fixtures.Random.Bool() {
			t.FailNow()
		}
	})
	s.Finish()
}

func TestNewSpec_withTestcaseT_optionsPassed(t *testing.T) {
	testcase.NewSpec(t).Test(``, func(t *testcase.T) {
		stub := &internal.StubTB{}
		s := testcase.NewSpec(t, testcase.Flaky(time.Second))
		s.Test(``, func(t *testcase.T) {
			if fixtures.Random.Bool() {
				t.FailNow()
			}
		})
		s.Finish()
		assert.Must(t).True(!stub.IsFailed)
	})
}

func TestNewSpec_withTestcaseT_InheritContext(t *testing.T) {
	s := testcase.NewSpec(t)

	n := testcase.Var{Name: "n"} // intentionally without Init
	nGet := func(t *testcase.T) int { return n.Get(t).(int) }
	n.Let(s, func(t *testcase.T) interface{} { return t.Random.Int() })

	s.Test(``, func(t *testcase.T) {
		sub := testcase.NewSpec(t)
		defer sub.Finish()
		sub.Test(``, func(t *testcase.T) { t.Log(nGet(t)) })
	})
}

func TestSpec_Test_whenTestingTBIsGivenThatDoesNotSupportTBRunner_executesOnFinish(t *testing.T) {
	s := testcase.NewSpec(t)
	var ran bool
	s.Test(``, func(t *testcase.T) { ran = true })
	assert.Must(t).True(!ran)
	s.Finish()
	assert.Must(t).True(ran)
}

func TestSpec_Test_testingTBNoTBRunner_ordered(t *testing.T) {
	testcase.SetEnv(t, testcase.EnvKeySeed, "42")
	testcase.SetEnv(t, testcase.EnvKeyOrdering, string(testcase.OrderingAsRandom))
	s := testcase.NewSpec(t)

	var out []int
	s.Test(``, func(t *testcase.T) { out = append(out, 1) })
	s.Test(``, func(t *testcase.T) { out = append(out, 2) })
	s.Test(``, func(t *testcase.T) { out = append(out, 3) })
	s.Test(``, func(t *testcase.T) { out = append(out, 4) })
	s.Test(``, func(t *testcase.T) { out = append(out, 5) })
	s.Finish()

	assert.Must(t).Equal([]int{3, 4, 5, 1, 2}, out)
}
