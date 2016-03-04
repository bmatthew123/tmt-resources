package accessors

import (
	"github.com/DATA-DOG/go-sqlmock"
	testhelpers "github.com/byu-oit-ssengineering/tmt-test-helpers"
	"testing"
)

func TestGetResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceAccessor(db)

	expected := Resource{"11111111-2222-3333-4444-555555555555", "test", "this is a test", "tmt.byu.edu/resources", nil}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resources WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources"))
	resource, err := ra.Get("11111111-2222-3333-4444-555555555555")
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

func TestGetAllResources(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceAccessor(db)

	expected := []Resource{
		Resource{"11111111-2222-3333-4444-555555555555", "test", "this is a test", "tmt.byu.edu/resources", nil},
		Resource{"00000000-9999-8888-7777-666666666666", "testing", "for testing purposes", "tmt.byu.edu/resources", nil},
	}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resources").
		WithArgs().
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources\n00000000-9999-8888-7777-666666666666,testing,for testing purposes.tmt.byu.edu/resources"))
	resources, err := ra.GetAll()
	if err != nil {
		t.Error("An unexpected error occurred while getting a resource %v", err)
	}

	for i := 0; i < len(resources); i++ {
		if !expected[i].Equals(resources[i]) {
			t.Errorf("Expected %v but got %v", expected, resources)
		}
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestInsertResource(t *testing.T) {
	NewGuid = func() string {
		return "123def"
	}
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("INSERT INTO resources .+ VALUES .+").
		WithArgs("123def", "test", "This is a test", "tmt.byu.edu/resources").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = ra.Insert(Resource{"123def", "test", "This is a test", "tmt.byu.edu/resources", nil})
	if err != nil {
		t.Error("An unexpected error occurred while getting a resource:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestUpdateResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("UPDATE resources SET name=(.), description=(.), apiEndpoint=(.) WHERE guid=(.)").
		WithArgs("test", "This is a test", "tmt.byu.edu/resources", "11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = ra.Update(Resource{"11111111-2222-3333-4444-555555555555", "test", "This is a test", "tmt.byu.edu/resources", nil})
	if err != nil {
		t.Error("An unexpected error occurred while getting a resource:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}

func TestDeleteResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred creating the mock database")
		return
	}

	ra := NewResourceAccessor(db)

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("DELETE FROM resources WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = ra.Delete("11111111-2222-3333-4444-555555555555")
	if err != nil {
		t.Error("An unexpected error occurred while getting a resource:\n %s", err.Error())
	}

	if err := ra.DB.Close(); err != nil {
		t.Errorf("An error occurred: %v", err)
	}
}
