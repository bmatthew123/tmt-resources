package accessors

import (
	"github.com/DATA-DOG/go-sqlmock"
	testhelpers "github.com/byu-oit-ssengineering/tmt-test-helpers"
	"testing"
)

func TestGetResourceVerb(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceVerbAccessor(db)

	expected := ResourceVerb{"11111111-2222-3333-4444-555555555555", "11111111-1111-1111-1111-111111111111", "test", "allows testing"}
	columns := []string{"guid", "resourceGUID", "verb", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE guid=(.)").
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,11111111-1111-1111-1111-111111111111,test,allows testing"))
	resourceVerb, err := ra.Get("11111111-1111-1111-1111-111111111111")
	if err != nil {
		t.Errorf("An unexpected error occurred while getting a resourceVerb %v", err)
	}

	if resourceVerb != expected {
		t.Errorf("Expected %v but got %v", expected, resourceVerb)
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestGetResourceVerbsByResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceVerbAccessor(db)

	expected := []ResourceVerb{
		ResourceVerb{"11111111-2222-3333-4444-555555555555", "11111111-1111-1111-1111-111111111111", "test", "allows testing"},
		ResourceVerb{"00000000-9999-8888-7777-666666666666", "11111111-1111-1111-1111-111111111111", "create", "allows creation of tests"},
	}
	columns := []string{"guid", "resourceGUID", "verb", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE resourceGUID=(.)").
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,11111111-1111-1111-1111-111111111111,test,allows testing\n00000000-9999-8888-7777-666666666666,11111111-1111-1111-1111-111111111111,create,allows creation of tests"))
	resourceVerbs, err := ra.GetByResource("11111111-1111-1111-1111-111111111111")
	if err != nil {
		t.Errorf("An unexpected error occurred while getting a resourceVerb %v", err)
	}

	for i := 0; i < len(resourceVerbs); i++ {
		if resourceVerbs[i] != expected[i] {
			t.Errorf("Expected %v but got %v", expected, resourceVerbs)
		}
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestInsertResourceVerb(t *testing.T) {
	NewGuid = func() string {
		return "123def"
	}
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceVerbAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("INSERT INTO resourceVerbs .+ VALUES .+").
		WithArgs("123def", "11111111-1111-1111-1111-111111111111", "test", "allows testing").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = ra.Add(ResourceVerb{ResourceGUID: "11111111-1111-1111-1111-111111111111", Verb: "test", Description: "allows testing"})
	if err != nil {
		t.Error("An unexpected error occurred while getting a resourceVerb:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestUpdateResourceVerb(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceVerbAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("UPDATE resourceVerbs SET description=(.) WHERE guid=(.)").
		WithArgs("This is a test", "11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = ra.Update("11111111-2222-3333-4444-555555555555", "This is a test")
	if err != nil {
		t.Error("An unexpected error occurred while getting a resourceVerb:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestDeleteResourceVerb(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceVerbAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("DELETE FROM resourceVerbs WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = ra.Remove("11111111-2222-3333-4444-555555555555")
	if err != nil {
		t.Error("An unexpected error occurred while getting a resourceVerb:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}
