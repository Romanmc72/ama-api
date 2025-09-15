//go:build integration
// +build integration

package integration_test

import "testing"

func TestIntegrationTestSuite(t *testing.T) {
	t.Run("Question Creation Test Suite", QuestionSuite)
	t.Run("User Setup Test Suite", UserSetupSuite)
	t.Run("List Test Suite", ListSuite)
	// t.Run("User Teardown Test Suite", UserTearDownSuite)
}
