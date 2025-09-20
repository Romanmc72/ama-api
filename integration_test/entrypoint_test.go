//go:build integration
// +build integration

package integration_test

import (
	"flag"
	"testing"
)

const (
	flagHelp          = "Which test suite to run, can be 'all' | 'setup' | 'teardown' or for short 'a' | 's' | 't'. Defaults to 'all', fails on bad flag being passed in."
	SuitesRunAll      = "all"
	SuitesRunSetUp    = "setup"
	SuitesRunTearDown = "teardown"
)

var suite = flag.String("suite", "", flagHelp)
var suiteMap = map[string]string{
	SuitesRunAll:      SuitesRunAll,
	"a":               SuitesRunAll,
	SuitesRunSetUp:    SuitesRunSetUp,
	"s":               SuitesRunSetUp,
	SuitesRunTearDown: SuitesRunTearDown,
	"t":               SuitesRunTearDown,
}

func TestMain(m *testing.M) {
	flag.Parse()
	run, ok := suiteMap[*suite]
	if ok {
		*suite = run
	}
	m.Run()
}

func runSetUp(t *testing.T) {
	t.Run("Question Creation Test Suite", QuestionSetUpSuite)
	t.Run("User Setup Test Suite", UserSetUpSuite)
	t.Run("List Test Suite", ListSuite)
}

func runTearDown(t *testing.T) {
	t.Run("User Teardown Test Suite", UserTearDownSuite)
	t.Run("Question Teardown Test Suite", QuestionTearDownSuite)
}

func TestSetUpSuite(t *testing.T) {
	if *suite != SuitesRunSetUp {
		t.Skip("Skipping the setup suite")
	}
	runSetUp(t)
}

func TestTearDownSuite(t *testing.T) {
	if *suite != SuitesRunTearDown {
		t.Skip("Skipping the teardown suite")
	}
	runTearDown(t)
}

func TestAllSuite(t *testing.T) {
	if *suite != SuitesRunAll {
		t.Skip("Skipping the all suite")
	}
	runSetUp(t)
	runTearDown(t)
}

func TestOtherSuite(t *testing.T) {
	if *suite != SuitesRunAll && *suite != SuitesRunSetUp && *suite != SuitesRunTearDown {
		t.Fatalf("invalid run flag passed in, '%s'. %s", *suite, flagHelp)
	}
}
