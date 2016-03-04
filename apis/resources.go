package apis

import (
	eden "github.com/byu-oit-ssengineering/tmt-eden"
	accessors "github.com/byu-oit-ssengineering/tmt-resources/accessors"
)

// Get all the resources.
// GET /resources
func (a *Api) GetAllResources(c *eden.Context) {
	ra := accessors.NewResourceAccessor(a.DB)
	va := accessors.NewResourceVerbAccessor(a.DB)

	resources, err := ra.GetAll()
	if err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error occurred while retrieving resources"})
		return
	}

	for i := 0; i < len(resources); i++ {
		resources[i].Verbs, _ = va.GetByResource(resources[i].Guid)
	}

	// Respond
	c.Respond(200, eden.Response{"OK", resources})
}

// Gets a resource by guid.
// GET /resources/:guid
func (a *Api) GetResource(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceAccessor(a.DB)
	va := accessors.NewResourceVerbAccessor(a.DB)

	// Parse the resource guid
	guid := c.Params[0].Value

	// Get the resource
	resource, err := ra.Get(guid)
	if err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error occurred while retrieving resource information"})
		return
	}

	resource.Verbs, _ = va.GetByResource(resource.Guid)

	// Respond
	c.Respond(200, eden.Response{"OK", resource})
}

// Create a resource.
// POST /resources name=:name, description=:description, api=:apiEndpoint
func (a *Api) InsertResource(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceAccessor(a.DB)

	// Parse name and description from POST data.
	c.Request.ParseForm()
	name, nameOk := c.Request.Form["name"]
	description, descriptionOk := c.Request.Form["description"]
	api, apiOk := c.Request.Form["api"]
	if !descriptionOk || !nameOk || !apiOk {
		c.Respond(400, eden.Response{"ERROR", "Unable to process request"})
		return
	}

	resource := accessors.Resource{Name: name[0], Description: description[0], APIEndpoint: api[0]}

	// Insert the resource and test for errors
	if err := ra.Insert(resource); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}

// Update a resource's name and/or description.
// PUT /resources/:guid name=:newName, description=:newDescription, api=:newApiEndpoint
func (a *Api) UpdateResource(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceAccessor(a.DB)

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
	if name, nameOk := c.Request.Form["name"]; nameOk {
		resource.Name = name[0]
	}
	if description, descriptionOk := c.Request.Form["description"]; descriptionOk {
		resource.Description = description[0]
	}
	if api, apiOk := c.Request.Form["api"]; apiOk {
		resource.APIEndpoint = api[0]
	}

	// Save
	if err := ra.Update(resource); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}

// Delete a resource.
// DELETE /resources/:guid
func (a *Api) DeleteResource(c *eden.Context) {
	// Create new resource accessor
	ra := accessors.NewResourceAccessor(a.DB)

	// Parse resource id
	guid := c.Params[0].Value

	// Delete the resource
	if err := ra.Delete(guid); err != nil {
		c.Respond(500, eden.Response{"ERROR", "An error has occurred"})
		return
	}

	// Respond
	c.Respond(200, eden.Response{"OK", "success"})
}
