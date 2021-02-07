package testcase

import (
	"fmt"
	"testing"
)

// Contract meant to represent a Role Interface Contract.
// A role interface express required behavior from a consumer point of view
// and a role interface contract describes all the assumption about the behavior of supplier
// that the consumer actively uses to simply the code.
type Contract interface {
	// Test is the function that assert expected behavioral requirements from a supplier implementation.
	// These behavioral assumptions made by the Consumer in order to simplify and stabilise its own code complexity.
	// Every time a Consumer makes an assumption about the behavior of the role interface supplier,
	// it should be clearly defined it with tests under this functionality.
	Test(*testing.T)
	// Benchmark will help with what to measure.
	// When you define a role interface contract, you should clearly know what performance aspects important for your Consumer.
	// Those aspects should be expressed in a form of Benchmark,
	// so different supplier implementations can be easily A/B tested from this aspect as well.
	Benchmark(*testing.B)
}

func RunContracts(tb interface{}, contracts ...Contract) {
	for _, c := range contracts {
		switch tb := tb.(type) {
		case *testing.T:
			c.Test(tb)

		case *testing.B:
			c.Benchmark(tb)

		case *T:
			RunContracts(tb.TB, c)

		case *Spec:
			c := c // copy to avoid reference overrides from "for"
			tb.Test(fullyQualifiedName(c), func(t *T) { RunContracts(t, c) })

		default:
			panic(fmt.Errorf(`unknown test runner type: %T`, tb))
		}
	}
}
