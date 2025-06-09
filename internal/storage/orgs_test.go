package storage

import (
	"lowerthirdsapi/internal/testutil"
	"testing"
)

func TestCreateAndGetOrganization(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, org, _ := testutil.CreateTestData(t, service)

	// Test GetOrg
	retrievedOrg, err := service.GetOrg(testutil.TestCtx, org.OrgID)
	if err != nil {
		t.Fatalf("GetOrg failed: %v", err)
	}

	// Verify organization data
	if retrievedOrg.OrgID != org.OrgID {
		t.Errorf("Expected OrgID %v, got %v", org.OrgID, retrievedOrg.OrgID)
	}
	if retrievedOrg.Name != org.Name {
		t.Errorf("Expected Name %v, got %v", org.Name, retrievedOrg.Name)
	}
}
