package apis

import (
	eden "github.com/byu-oit-ssengineering/tmt-eden"
	accessors "github.com/byu-oit-ssengineering/tmt-resources/accessors"
)

// Get a list of the verbs associated to a resource.
// GET /verbs
func (a *Api) GetResourceVerbs(c *eden.Context) {
	ra := accessors.NewResourceVerbAccessor(a.DB)

	guid := c.Params[0].Value

	verbs, err := ra.GetByResource(guid)
	if err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error occurred while retrieving resources"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", verbs})
}

// Associate a verb to a resource.
// POST /verbs resource=:resourceGUID, verb=:verb, description=:description
func (a *Api) AddVerb(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceVerbAccessor(a.DB)

	// Parse name and description from POST data.
	c.Request.ParseForm()
	resourceGUID, resourceGUIDOk := c.Request.Form["resourceGUID"]
	verb, verbOk := c.Request.Form["verb"]
	description, descriptionOk := c.Request.Form["description"]
	if !descriptionOk || !verbOk || !resourceGUIDOk {
		c.Respond(400, eden.Response{"ERROR", "Not enough information given"})
		return
	}

	resource := accessors.ResourceVerb{ResourceGUID: resourceGUID[0], Verb: verb[0], Description: description[0]}

	// Insert the resource and test for errors
	if err := ra.Add(resource); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}

// Update a verb's description.
// PUT /verbs/:guid description=:newDescription
func (a *Api) UpdateVerb(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceVerbAccessor(a.DB)

	// Parse resource guid
	guid := c.Params[0].Value

	// Get the resource
	resource, err := ra.Get(guid)
	if err != nil {
		c.Respond(500, eden.Response{"ERROR", "An unexpected error occurred"})
		return
	}

	// Update given fields
	c.Request.ParseForm()
	if description, descriptionOk := c.Request.Form["description"]; descriptionOk {
		resource.Description = description[0]
	} else {
		c.Respond(400, eden.Response{"ERROR", "No new description has been specified"})
		return
	}

	// Save
	if err := ra.Update(guid, resource.Description); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}

// Delete a verb-resource type association.
// DELETE /verbs/:guid
func (a *Api) RemoveVerb(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceVerbAccessor(a.DB)

	// Parse resource id
	guid := c.Params[0].Value

	// Delete the resource
	if err := ra.Remove(guid); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}
