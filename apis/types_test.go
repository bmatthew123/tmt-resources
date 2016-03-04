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

func TestGetResourceType(t *testing.T) {
	db, err := testhelpers.GetMockDB()
	if err != nil {
		t.Error("An unexpected error occurred instantiating accessor")
		return
	}
	api := &Api{db}

	expected := accessors.Resource{"11111111-2222-3333-2222-111111111111", "test", "this is a test", "tmt.byu.edu/resourceTypes", nil}
	columns := []string{"guid", "name", "description", "apiEndpoint"}
	sqlmock.ExpectPrepare()
	sqlmock.ExpectQuery("SELECT guid, name, description, apiEndpoint FROM resources JOIN resourceTypes ON resources.guid=resourceTypes.type WHERE resourceTypes.resourceGUID=(.)").
		WithArgs("11111111-2222-3333-4444-555555555555").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("11111111-2222-3333-2222-111111111111,test,this is a test,tmt.byu.edu/resourceTypes"))

	// Create context, call API
	var result []byte
	var output testResponseResource
	c := testhelpers.NewTestingContext("", httprouter.Params{httprouter.Param{Key: "guid", Value: "11111111-2222-3333-4444-555555555555"}}, api.GetResourceType)
	testhelpers.CallAPI(api.GetResourceType, c, &result)

	// Parse output
	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Compare output to expected output
	if !expected.Equals(output.Data) {
		t.Errorf("Expected: %v, but got %v", expected, output)
	}
}

func TestInsertResourceType(t *testing.T) {
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
	sqlmock.ExpectExec("INSERT INTO resourceTypes .+ VALUES .+").
		WithArgs("123def", "11111111-2222-3333-4444-555555555555", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create context and call API
	var result []byte
	var output eden.Response
	c := testhelpers.NewTestingContext("resource=11111111-2222-3333-4444-555555555555&type=test", nil, api.InsertResourceType)
	testhelpers.CallAPI(api.InsertResourceType, c, &result)

	err = json.Unmarshal(result, &output)
	if err != nil {
		t.Error(err.Error())
	}

	// Ensure correct output
	if output.Data != "success" {
		t.Errorf("expected to get 'success' but got %v instead", output)
	}
}
