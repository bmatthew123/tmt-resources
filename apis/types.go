package apis

import (
	eden "github.com/byu-oit-ssengineering/tmt-eden"
	accessors "github.com/byu-oit-ssengineering/tmt-resources/accessors"
)

// Get the type of the given resource guid.
// GET /type/:guid
func (a *Api) GetResourceType(c *eden.Context) {
	// Create new resourceType accessor
	ra := accessors.NewResourceTypeAccessor(a.DB)

	// Parse the resourceType guid
	guid := c.Params[0].Value

	// Get the resourceType
	resourceType, err := ra.GetType(guid)
	if err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error occurred while retrieving resourceType information"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", resourceType})
}

// Store the type of a resource.
// POST /type resource=:resourceGUID, type=:resourceTypeGUID
func (a *Api) InsertResourceType(c *eden.Context) {
	// Create new resourceType accessor
	ra := accessors.NewResourceTypeAccessor(a.DB)

	// Parse name and description from POST data.
	c.Request.ParseForm()
	r, resourceOk := c.Request.Form["resource"]
	t, typeOk := c.Request.Form["type"]
	if !typeOk || !resourceOk {
		c.Respond(400, eden.Response{"ERROR", "Unable to process request"})
		return
	}

	// Insert the resourceType and test for errors
	if err := ra.Insert(r[0], t[0]); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}
