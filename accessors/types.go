package accessors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type ResourceTypeAccessor struct {
	DB *sql.DB // Database connection
}

// Returns a new resource accessor.
func NewResourceTypeAccessor(db *sql.DB) *ResourceTypeAccessor {
	return &ResourceTypeAccessor{db}
}

// Returns the resource type information for the given resource.
//   For example, if the guid of a whiteboard is given, it will return
//   information about the whiteboard resource type.
func (ra *ResourceTypeAccessor) GetType(guid string) (Resource, error) {
	r := Resource{}
	stmt, err := ra.DB.Prepare("SELECT guid, name, description, apiEndpoint FROM resources JOIN resourceTypes ON resources.guid=resourceTypes.type WHERE resourceTypes.resourceGUID=?")
	if err != nil {
		return r, err
	}

	row := stmt.QueryRow(guid)
	err = row.Scan(&r.Guid, &r.Name, &r.Description, &r.APIEndpoint)
	return r, err
}

// Create a new resourceType.
func (ra *ResourceTypeAccessor) Insert(r, t string) error {
	stmt, err := ra.DB.Prepare("INSERT INTO resourceTypes (guid, resourceGUID, type) VALUES (?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(NewGuid(), r, t)
	return err
}
