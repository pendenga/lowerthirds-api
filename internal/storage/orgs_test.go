package storage

import (
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/storage/testutil"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAndGetOrganization(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Create test organization
	orgID := uuid.New()
	org := &entities.Organization{
		OrgID: orgID,
		Name:  "Test Organization 2",
	}

	// Test CreateOrg
	err := service.CreateOrg(testutil.TestCtx, org)
	if err != nil {
		t.Fatalf("CreateOrg failed: %v", err)
	}

	// Test GetOrg
	retrievedOrg, err := service.GetOrg(testutil.TestCtx, orgID)
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
