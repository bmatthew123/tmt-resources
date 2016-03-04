package accessors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Resource struct that reflects the resources table.
type ResourceVerb struct {
	Guid         string `json:"guid"`
	ResourceGUID string `json:"resourceGUID"`
	Verb         string `json:"verb"`
	Description  string `json:"description"`
}

type ResourceVerbAccessor struct {
	DB *sql.DB // Database connection
}

// Returns a new resource verb accessor.
func NewResourceVerbAccessor(db *sql.DB) *ResourceVerbAccessor {
	return &ResourceVerbAccessor{db}
}

func (ra *ResourceVerbAccessor) Get(guid string) (ResourceVerb, error) {
	var r ResourceVerb
	stmt, err := ra.DB.Prepare("SELECT * FROM resourceVerbs WHERE guid=?")
	if err != nil {
		return r, err
	}

	row := stmt.QueryRow(guid)
	err = row.Scan(&r.Guid, &r.ResourceGUID, &r.Verb, &r.Description)
	return r, err
}

// Gets all verbs associated to a given resource by that resource's guid.
func (ra *ResourceVerbAccessor) GetByResource(resource string) ([]ResourceVerb, error) {
	verbs := make([]ResourceVerb, 0)
	stmt, err := ra.DB.Prepare("SELECT * FROM resourceVerbs WHERE resourceGUID=?")
	if err != nil {
		return verbs, err
	}

	rows, err := stmt.Query(resource)
	if err != nil {
		return verbs, err
	}
	defer rows.Close()

	for rows.Next() {
		var r ResourceVerb
		rows.Scan(&r.Guid, &r.ResourceGUID, &r.Verb, &r.Description)
		verbs = append(verbs, r)
	}

	return verbs, nil
}

// Associate a new verb to a resource.
func (ra *ResourceVerbAccessor) Add(r ResourceVerb) error {
	stmt, err := ra.DB.Prepare("INSERT INTO resourceVerbs (guid, resourceGUID, name, description) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(NewGuid(), r.ResourceGUID, r.Verb, r.Description)
	return err
}

// Update the description for a verb on a resource type. The guid passed in
//   is the guid of the resource/verb association.
func (ra *ResourceVerbAccessor) Update(guid, description string) error {
	stmt, err := ra.DB.Prepare("UPDATE resourceVerbs SET description=? WHERE guid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(description, guid)
	return err
}

// Disassociate a verb from a resource type.
func (ra *ResourceVerbAccessor) Remove(guid string) error {
	stmt, err := ra.DB.Prepare("DELETE FROM resourceVerbs WHERE guid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(guid)
	return err
}
