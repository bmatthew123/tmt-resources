package apis

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	eden "github.com/byu-oit-ssengineering/tmt-eden"
	accessors "github.com/byu-oit-ssengineering/tmt-resources/accessors"
	testhelpers "github.com/byu-oit-ssengineering/tmt-test-helpers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"testing"
)

type testResponseResource struct {
	Status string
	Data   accessors.Resource
}

type testResponseResourceArray struct {
	Status string
	Data   []accessors.Resource
}

func TestGetResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	expected := accessors.Resource{"11111111-2222-3333-4444-555555555555", "test", "this is a test", "tmt.byu.edu/resources",
		[]accessors.ResourceVerb{accessors.ResourceVerb{"22222222-2222-2222-2222-222222222222", "11111111-2222-3333-4444-555555555555", "edit", "can edit"}}}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resources WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources"))

	columns = []string{"guid", "resourceGUID", "name", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE resourceGUID=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("22222222-2222-2222-2222-222222222222,11111111-2222-3333-4444-555555555555,edit,can edit"))

	// Create context, call API
	var result []byte
	var output testResponseResource
	c := testhelpers.NewTestingContext("", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.GetResource)
	testhelpers.CallAPI(api.GetResource, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Compare output to expected output
	if !expected.Equals(output.Data) {
		t.Errorf("Expected: %v, but got %v", expected, output.Data)
	}
}

func TestGetAllResources(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	expected := []accessors.Resource{
		accessors.Resource{"11111111-2222-3333-4444-555555555555", "test", "this is a test", "tmt.byu.edu/resources",
			[]accessors.ResourceVerb{accessors.ResourceVerb{"22222222-2222-2222-2222-222222222222", "11111111-2222-3333-4444-555555555555", "edit", "can edit"}}},
		accessors.Resource{"00000000-9999-8888-7777-666666666666", "testing", "for testing purposes", "tmt.byu.edu/resources",
			[]accessors.ResourceVerb{accessors.ResourceVerb{"33333333-3333-3333-3333-333333333333", "00000000-9999-8888-7777-666666666666", "edit", "can edit"}}},
	}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resources").
		WithArgs().
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources\n00000000-9999-8888-7777-666666666666,testing,for testing purposes,tmt.byu.edu/resources"))

	columns = []string{"guid", "resourceGUID", "name", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE resourceGUID=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("22222222-2222-2222-2222-222222222222,11111111-2222-3333-4444-555555555555,edit,can edit"))

	columns = []string{"guid", "resourceGUID", "name", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE resourceGUID=(.)").
		WithArgs("00000000-9999-8888-7777-666666666666").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("33333333-3333-3333-3333-333333333333,00000000-9999-8888-7777-666666666666,edit,can edit"))

	// Create context, call API
	var result []byte
	var output testResponseResourceArray
	c := testhelpers.NewTestingContext("", nil, api.GetAllResources)
	testhelpers.CallAPI(api.GetAllResources, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Compare output to expected output
	for i := 0; i < len(expected); i++ {
		if !expected[i].Equals(output.Data[i]) {
			t.Errorf("Expected: %v, but got %v instead", expected, output.Data)
		}
	}
}

func TestInsertResource(t *testing.T) {
	accessors.NewGuid = func() string {
		return "123def"
	}
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("INSERT INTO resources .+ VALUES .+").
		WithArgs("123def", "test", "This is a test", "tmt.byu.edu/resources").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("name=test&description=This%20is%20a%20test&api=tmt.byu.edu/resources", nil, api.InsertResource)
	testhelpers.CallAPI(api.InsertResource, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output.Data)
	}
}

func TestUpdateResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resources WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,test,this is a test,tmt.byu.edu/resources"))
	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("UPDATE resources SET name=(.), description=(.), apiEndpoint=(.) WHERE guid=(.)").
		WithArgs("changed", "testing", "tmt.byu.edu/resources", "11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("name=changed&description=testing&api=tmt.byu.edu/resources", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.UpdateResource)
	testhelpers.CallAPI(api.UpdateResource, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output.Data)
	}
}

func TestDeleteResource(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("DELETE FROM resources WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.DeleteResource)
	testhelpers.CallAPI(api.DeleteResource, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output.Data)
	}
}
