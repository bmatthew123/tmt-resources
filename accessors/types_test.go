package accessors

import (
	"github.com/DATA-DOG/go-sqlmock"
	testhelpers "github.com/byu-oit-ssengineering/tmt-test-helpers"
	"testing"
)

func TestGetType(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceTypeAccessor(db)

	expected := Resource{"11111111-2222-3333-4444-555555555555", "test", "this is a test", "tmt.byu.edu/resources", nil}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT guid, name, description, apiEndpoint FROM resources JOIN resourceTypes ON resources.guid=resourceTypes.type WHERE resourceTypes.resourceGUID=(.)").
		WithArgs("11111111-2222-3333-2222-111111111111").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources"))
	resource, err := ra.GetType("11111111-2222-3333-2222-111111111111")
	if err != nil {
		t.Error("An unexpected error occurred while getting a resource %v", err)
	}

	if !expected.Equals(resource) {
		t.Errorf("Expected %v but got %v", expected, resource)
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestInsertType(t *testing.T) {
	NewGuid = func() string {
		return "123def"
	}
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceTypeAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("INSERT INTO resourceTypes (.+) VALUES (.+)").
		WithArgs("123def", "11111111-2222-3333-2222-111111111111", "55555555-6666-7777-8888-999999999999").
		WillReturnResult(sqlmock.NewResult(1, 1))

	ra.Insert("11111111-2222-3333-2222-111111111111", "55555555-6666-7777-8888-999999999999")

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}
