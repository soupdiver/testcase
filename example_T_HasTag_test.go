package testcase_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/adamluzsi/testcase"
)

func ExampleT_HasTag() {
	var t *testing.T
	var s = testcase.NewSpec(t)

	s.Let(`db`, func(t *testcase.T) interface{} {
		db, err := sql.Open(`driverName`, `dataSourceName`)
		t.Must.Nil(err)

		if t.HasTag(`black box`) {
			// tests with black box  use http testCase server or similar things and high level tx management not maintainable.
			t.Defer(db.Close)
			return db
		}

		tx, err := db.BeginTx(context.Background(), nil)
		t.Must.Nil(err)
		t.Defer(tx.Rollback)
		return tx
	})
}
