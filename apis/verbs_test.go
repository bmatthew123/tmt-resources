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

type testResponseVerbArray struct {
	Status string
	Data   []accessors.ResourceVerb
}

type testResponseVerb struct {
	Status string
	Data   accessors.ResourceVerb
}

func TestGetResourceVerbs(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	expected := []accessors.ResourceVerb{
		accessors.ResourceVerb{"11111111-2222-3333-4444-555555555555", "11111111-2222-3333-2222-111111111111", "test", "allows testing"},
		accessors.ResourceVerb{"11111111-2222-3333-4444-666666666666", "11111111-2222-3333-2222-111111111111", "create", "can create tests"},
	}
	columns := []string{"guid", "resourceGUID", "verb", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE resourceGUID=(.)").
		WithArgs("11111111-2222-3333-2222-111111111111").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,11111111-2222-3333-2222-111111111111,test,allows testing\n11111111-2222-3333-4444-666666666666,11111111-2222-3333-2222-111111111111,create,can create tests"))

	// Create context, call API
	var result []byte
	var output testResponseVerbArray
	c := testhelpers.NewTestingContext("", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-2222-111111111111"}}, api.GetResourceVerbs)
	testhelpers.CallAPI(api.GetResourceVerbs, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Compare output to expected output
	for i := 0; i < len(output.Data); i++ {
		if output.Data[i] != expected[i] {
			t.Errorf("Expected: %v, but got %v instead", expected, output)
		}
	}
}

func TestInsertResourceVerb(t *testing.T) {
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
	sqlmock.ExpectExec("INSERT INTO resourceVerbs .+ VALUES .+").
		WithArgs("123def", "11111111-2222-3333-4444-555555555555", "test", "allows testing").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("resourceGUID=11111111-2222-3333-4444-555555555555&verb=test&description=allows%20testing", nil, AddVerb)
	testhelpers.CallAPI(AddVerb, c, &result)

	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Error(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output)
	}
}

func TestUpdateVerb(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	columns := []string{"guid", "resourceGUID", "verb", "description"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT (.) FROM resourceVerbs WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-4444-555555555555,00000000-9999-8888-7777-666666666666,test,this is a test"))
	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("UPDATE resourceVerbs SET description=(.) WHERE guid=(.)").
		WithArgs("testing", "11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("description=testing", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.UpdateVerb)
	testhelpers.CallAPI(api.UpdateVerb, c, &result)

	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Error(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output)
	}
}

func TestDeleteResourceVerb(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	sqlmock.ExpectPrepare()
	sqlmock.ExpectExec("DELETE FROM resourceVerbs WHERE guid=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.RemoveVerb)
	testhelpers.CallAPI(api.RemoveVerb, c, &result)

	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Error(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output)
	}
}
