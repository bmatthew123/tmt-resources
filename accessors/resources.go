package accessors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Resource struct that reflects the resources table.
type Resource struct {
	Guid        string         `json:"guid"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	APIEndpoint string         `json:"apiEndpoint"`
	Verbs       []ResourceVerb `json:"verbs"`
}

type ResourceAccessor struct {
	DB *sql.DB // Database connection
}

// Returns a new resource accessor.
func NewResourceAccessor(db *sql.DB) *ResourceAccessor {
	return &ResourceAccessor{db}
}

// Gets the resource with the given id.
func (ra *ResourceAccessor) Get(guid string) (Resource, error) {
	r := Resource{}
	stmt, err := ra.DB.Prepare("SELECT * FROM resources WHERE guid=?")
	if err != nil {
		return r, err
	}

	row := stmt.QueryRow(guid)
	err = row.Scan(&r.Guid, &r.Name, &r.Description, &r.APIEndpoint)
	return r, err
}

// Gets the resource with the given id.
func (ra *ResourceAccessor) GetAll() ([]Resource, error) {
	resources := make([]Resource, 0)
	stmt, err := ra.DB.Prepare("SELECT * FROM resources")
	if err != nil {
		return resources, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return resources, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Resource
		rows.Scan(&r.Guid, &r.Name, &r.Description, &r.APIEndpoint)
		resources = append(resources, r)
	}

	return resources, nil
}

// Create a new resource.
func (ra *ResourceAccessor) Insert(r Resource) error {
	stmt, err := ra.DB.Prepare("INSERT INTO resources (guid, name, description, apiEndpoint) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(NewGuid(), r.Name, r.Description, r.APIEndpoint)
	return err
}

// Renames a resource with the given id to have the provided name.
func (ra *ResourceAccessor) Update(r Resource) error {
	stmt, err := ra.DB.Prepare("UPDATE resources SET name=?, description=?, apiEndpoint=? WHERE guid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.Name, r.Description, r.APIEndpoint, r.Guid)
	return err
}

// Delete a resource.
func (ra *ResourceAccessor) Delete(guid string) error {
	stmt, err := ra.DB.Prepare("DELETE FROM resources WHERE guid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(guid)
	return err
}

// Helper function

// Determines whether two Resources are equal.
func (r *Resource) Equals(other Resource) bool {
	if len(r.Verbs) != len(other.Verbs) {
		return false
	}
	for i := 0; i < len(r.Verbs); i++ {
		if r.Verbs[i] != other.Verbs[i] {
			return false
		}
	}

	if r.Guid != other.Guid || r.Name != other.Name || r.Description != other.Description || r.APIEndpoint != other.APIEndpoint {
		return false
	}
	return true
}
