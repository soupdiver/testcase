package random_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/adamluzsi/testcase/assert"
	"github.com/adamluzsi/testcase/random"

	"github.com/adamluzsi/testcase"
)

func TestRandomizer(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Let(`randomizer`, func(t *testcase.T) interface{} {
		return &random.Random{Source: rand.NewSource(time.Now().Unix())}
	})
	s.Let(`source`, func(t *testcase.T) interface{} {
		return rand.NewSource(time.Now().Unix())
	})
	SpecRandomizerMethods(s)
}

func SpecRandomizerMethods(s *testcase.Spec) {
	var randomizer = func(t *testcase.T) *random.Random {
		return t.I(`randomizer`).(*random.Random)
	}

	s.Describe(`Int`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) int {
			return randomizer(t).Int()
		}

		s.Then(`it returns a non-negative pseudo-random int`, func(t *testcase.T) {
			out := subject(t)
			assert.Must(t).True(0 <= out)
		})

		s.Then(`it returns distinct value on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t))
		})
	})

	s.Describe(`Float32`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) float32 {
			return randomizer(t).Float32()
		}

		s.Then(`it returns, as a float32, a pseudo-random number in [0.0,1.0).`, func(t *testcase.T) {
			assert.Must(t).True(0 <= subject(t))
		})

		s.Then(`it returns distinct value on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t))
		})
	})

	s.Describe(`Float64`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) float64 {
			return randomizer(t).Float64()
		}

		s.Then(`it returns, as a float64, a pseudo-random number in [0.0,1.0).`, func(t *testcase.T) {
			assert.Must(t).True(0 <= subject(t))
		})

		s.Then(`it returns distinct value on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t))
		})
	})

	s.Describe(`IntN`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) int {
			return randomizer(t).IntN(t.I(`n`).(int))
		}

		s.Let(`n`, func(t *testcase.T) interface{} {
			return randomizer(t).IntN(42) + 42 // ensure it is not zero for the test
		})

		s.Test(`returns with random number excluding the received`, func(t *testcase.T) {
			out := subject(t)
			assert.Must(t).True(0 <= out)
			assert.Must(t).True(out < t.I(`n`).(int))
		})
	})

	s.Describe(`IntBetween`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) int {
			return randomizer(t).IntBetween(t.I(`min`).(int), t.I(`max`).(int))
		}

		s.Let(`min`, func(t *testcase.T) interface{} {
			return randomizer(t).IntN(42)
		})

		s.Let(`max`, func(t *testcase.T) interface{} {
			// +1 in the end to ensure that `max` is bigger than `min`
			return randomizer(t).IntN(42) + t.I(`min`).(int) + 1
		})

		s.Then(`it will return a value between the range`, func(t *testcase.T) {
			out := subject(t)
			assert.Must(t).True(t.I(`min`).(int) <= out, `expected that from <= than out`)
			assert.Must(t).True(out <= t.I(`max`).(int), `expected that out is <= than max`)
		})

		s.And(`min and max is in the negative range`, func(s *testcase.Spec) {
			s.LetValue(`min`, -128)
			s.LetValue(`max`, -64)

			s.Then(`it will return a value between the range`, func(t *testcase.T) {
				out := subject(t)
				assert.Must(t).True(t.I(`min`).(int) <= out, `expected that from <= than out`)
				assert.Must(t).True(out <= t.I(`max`).(int), `expected that out is <= than max`)
			})
		})

		s.And(`min and max equal`, func(s *testcase.Spec) {
			s.Let(`max`, func(t *testcase.T) interface{} { return t.I(`min`) })

			s.Then(`it returns the min and max value since the range can only have one value`, func(t *testcase.T) {
				t.Must.Equal(t.I(`max`), subject(t))
			})
		})
	})

	s.Describe(`ElementFromSlice`, func(s *testcase.Spec) {
		s.Test(`E2E`, func(t *testcase.T) {
			pool := []int{1, 2, 3, 4, 5}
			resSet := make(map[int]struct{})
			for i := 0; i < 1024; i++ {
				res := randomizer(t).ElementFromSlice(pool).(int)
				resSet[res] = struct{}{}
				t.Must.Contain(pool, res)
			}
			assert.Must(t).True(len(resSet) > 1, fmt.Sprintf(`%#v`, resSet))
		})
	})

	s.Describe(`KeyFromMap`, func(s *testcase.Spec) {
		s.Test(`E2E`, func(t *testcase.T) {
			var keys = []int{1, 2, 3, 4, 5}
			var srcMap = make(map[int]struct{})
			for _, k := range keys {
				srcMap[k] = struct{}{}
			}
			t.Must.Contain(keys, randomizer(t).KeyFromMap(srcMap).(int))
		})

		s.Test(`randomness`, func(t *testcase.T) {
			var keys = []int{1, 2, 3, 4, 5}
			var srcMap = make(map[int]struct{})
			for _, k := range keys {
				srcMap[k] = struct{}{}
			}
			resSet := make(map[int]struct{})
			for i := 0; i < 1024; i++ {
				res := randomizer(t).KeyFromMap(srcMap).(int)
				resSet[res] = struct{}{}
				t.Must.Contain(keys, res)
			}
			assert.Must(t).True(len(resSet) > 1, fmt.Sprintf(`%#v`, resSet))
		})
	})

	s.Describe(`StringN`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) string {
			return randomizer(t).StringN(t.I(`length`).(int))
		}
		s.Let(`length`, func(t *testcase.T) interface{} {
			return randomizer(t).IntN(42) + 5
		})

		s.Then(`it create a string with a given length`, func(t *testcase.T) {
			t.Must.Equal(t.I(`length`).(int), len(subject(t)),
				`it was expected to create string with the given length`)
		})

		s.Then(`it create random strings on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t),
				`it was expected to create different strings`)
		})
	})

	s.Describe(`StringNWithCharset`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) string {
			return randomizer(t).StringNWithCharset(t.I(`length`).(int), t.I(`charset`).(string))
		}
		s.Let(`length`, func(t *testcase.T) interface{} {
			return randomizer(t).IntN(42) + 5
		})

		s.Let(`charset`, func(t *testcase.T) interface{} {
			return "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
		})

		s.Then(`it create a string with a given length`, func(t *testcase.T) {
			t.Must.Equal(t.I(`length`).(int), len(subject(t)),
				`it was expected to create string with the given length`)
		})

		s.Then(`it create random strings on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t),
				`it was expected to create different strings`)
		})

		s.Test(`charset defines what characters will be randomly used`, func(t *testcase.T) {
			for _, edge := range []struct {
				charset string
			}{
				{charset: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"},
				{charset: "0123456789"},
				{charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
				{charset: "-$!/%"},
			} {
				t.Set(`charset`, edge.charset)
				for _, char := range subject(t) {
					t.Must.Contain(edge.charset, string(char))
				}
			}
		})
	})

	s.Describe(`Bool`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) bool {
			return randomizer(t).Bool()
		}

		s.Then(`it return with random bool on each calls`, func(t *testcase.T) {
			var bools = map[bool]struct{}{}
			for i := 0; i <= 1024; i++ {
				bools[subject(t)] = struct{}{}
			}
			t.Must.Equal(2, len(bools))
		})
	})

	s.Describe(`String`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) string {
			return randomizer(t).String()
		}

		s.Then(`it create strings with different lengths`, func(t *testcase.T) {
			var lengths = make(map[int]struct{})
			for i := 0; i < 1024; i++ {
				lengths[len(subject(t))] = struct{}{}
			}
			t.Must.True(1 < len(lengths))
		})

		s.Then(`it create random strings on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t),
				`it was expected to create different strings`)
		})
	})

	s.Describe(`TimeBetween`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) time.Time {
			return randomizer(t).TimeBetween(t.I(`from`).(time.Time), t.I(`to`).(time.Time))
		}

		s.Let(`from`, func(t *testcase.T) interface{} {
			return time.Now().UTC()
		})

		s.Let(`to`, func(t *testcase.T) interface{} {
			return t.I(`from`).(time.Time).Add(24 * time.Hour)
		})

		s.Then(`it will return a date between the given time range including 'from' and excluding 'to'`, func(t *testcase.T) {
			out := subject(t)
			assert.Must(t).True(t.I(`from`).(time.Time).Unix() <= out.Unix(), `expected that from <= than out`)
			assert.Must(t).True(out.Unix() < t.I(`to`).(time.Time).Unix(), `expected that out is < than to`)
		})

		s.Then(`it will generate different time on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t))
		})

		s.And(`from is before 1970-01-01 (unix timestamp 0)`, func(s *testcase.Spec) {
			s.Let(`from`, func(t *testcase.T) interface{} {
				return time.Unix(0, 0).UTC().AddDate(0, -1, 0)
			})

			s.Let(`to`, func(t *testcase.T) interface{} {
				return t.I(`from`).(time.Time).AddDate(0, 0, 1)
			})

			s.Then(`it will generate a random time between 'from' and 'to'`, func(t *testcase.T) {
				out := subject(t)
				assert.Must(t).True(t.I(`from`).(time.Time).Unix() <= out.Unix(), `expected that from <= than out`)
				assert.Must(t).True(out.Unix() < t.I(`to`).(time.Time).Unix(), `expected that out is < than to`)
			})
		})

		s.Then(`result is safe to format into RFC3339`, func(t *testcase.T) {
			t1 := subject(t)
			t2, _ := time.Parse(time.RFC3339, t1.Format(time.RFC3339))
			t.Must.Equal(t1.UTC(), t2.UTC())
		})
	})

	s.Describe(`Time`, func(s *testcase.Spec) {
		var subject = func(t *testcase.T) time.Time {
			return randomizer(t).Time()
		}

		s.Then(`it will generate different time on each call`, func(t *testcase.T) {
			assert.Must(t).NotEqual(subject(t), subject(t))
		})

		s.Then(`result is safe to format into RFC3339`, func(t *testcase.T) {
			t1 := subject(t)
			t2, _ := time.Parse(time.RFC3339, t1.Format(time.RFC3339))
			t.Must.Equal(t1.UTC(), t2.UTC())
		})
	})

	s.Describe(`TimeN`, func(s *testcase.Spec) {
		var (
			from = s.Let(`from`, func(t *testcase.T) interface{} {
				return time.Now()
			})
			fromGet = func(t *testcase.T) time.Time { return from.Get(t).(time.Time) }
			years   = s.Let(`years`, func(t *testcase.T) interface{} {
				return t.Random.IntN(42)
			})
			months = s.Let(`months`, func(t *testcase.T) interface{} {
				return t.Random.IntN(42)
			})
			days = s.Let(`days`, func(t *testcase.T) interface{} {
				return t.Random.IntN(42)
			})
		)
		var subject = func(t *testcase.T) time.Time {
			return randomizer(t).TimeN(fromGet(t), years.Get(t).(int), months.Get(t).(int), days.Get(t).(int))
		}

		getMaxDate := func(t *testcase.T) time.Time {
			return fromGet(t).AddDate(years.Get(t).(int), months.Get(t).(int), days.Get(t).(int))
		}

		s.Then(`it will return a value greater or equal with "from"`, func(t *testcase.T) {
			t.Must.True(fromGet(t).Unix() <= subject(t).Unix())
		})

		s.Then(`it will return a value less or equal with the maximum expected date that is: "from"+years+months+days`, func(t *testcase.T) {
			t.Must.True(subject(t).Unix() <= getMaxDate(t).Unix())
		})

		s.And(`years is negative`, func(s *testcase.Spec) {
			years.Let(s, func(t *testcase.T) interface{} {
				return t.Random.IntN(42) * -1
			})
			months.Let(s, func(t *testcase.T) interface{} {
				return t.Random.IntN(12) * -1
			})
			days.Let(s, func(t *testcase.T) interface{} {
				return t.Random.IntN(29) * -1
			})

			s.Then(`time shift backwards`, func(t *testcase.T) {
				t.Must.True(subject(t).Unix() <= fromGet(t).Unix())
				t.Must.True(getMaxDate(t).Unix() <= subject(t).Unix())
			})
		})

		s.Then(`stress test`, func(t *testcase.T) {
			min := fromGet(t).Unix()
			max := getMaxDate(t).Unix()
			for i := 0; i < 42; i++ {
				sub := subject(t).Unix()
				t.Must.True(min <= sub)
				t.Must.True(sub <= max)
			}
		})

		s.Then(`result is safe to format into RFC3339`, func(t *testcase.T) {
			t1 := subject(t)
			t2, _ := time.Parse(time.RFC3339, t1.Format(time.RFC3339))
			t.Log("t1:", t1.UnixNano(), "t2:", t2.UnixNano())
			t.Must.Equal(t1.UTC(), t2.UTC())
		})

		s.Then(`using it is race safe`, func(t *testcase.T) {
			rdz := randomizer(t)
			f := fromGet(t)
			y := years.Get(t).(int)
			m := months.Get(t).(int)
			d := days.Get(t).(int)
			blk := func() { rdz.TimeN(f, y, m, d) }
			testcase.Race(blk, blk, blk)
		})
	})
}
