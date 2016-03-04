package main

import (
	"fmt"
	eden "github.com/byu-oit-ssengineering/tmt-eden"
	apis "github.com/byu-oit-ssengineering/tmt-resources/apis"
)

// Responds with the allowed HTTP methods for this microservice.
func Options(c *eden.Context) {
	c.Response.Header().Add("Access-Control-Allow-Origin", "*")
	c.Response.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.Response.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response.WriteHeader(200)
}

func main() {
	r := eden.New()
	r.Use(eden.Authorize)

	a, err := apis.New()
	if err != nil {
		panic(err)
	}
	// Register api paths

	// Resources
	r.GET("/resources", a.GetAllResources)
	r.GET("/resources/:guid", a.GetResource)
	r.POST("/resources", a.InsertResource)
	r.PUT("/resources/:guid", a.UpdateResource)
	r.DELETE("/resources/:guid", a.DeleteResource)

	// Resource Verbs
	r.GET("/verbs/:guid", a.GetResourceVerbs)
	r.POST("/verbs", a.AddVerb)
	r.PUT("/verbs/:guid", a.UpdateVerb)
	r.DELETE("/verbs/:guid", a.RemoveVerb)

	// Resource Types
	r.GET("/type/:guid", a.GetResourceType)
	r.POST("/type", a.InsertResourceType)

	// For HTTP requests that aren't GET or POST, due to this being a cross-domain micro-service
	//   from the TMT, the browser is required to send an OPTIONS preflight request
	//   to determine which HTTP methods are allowed. This will inform the TMT
	//   that GET, POST, PUT, and DELETE methods are available. See the Options function for how this is done.
	// Also, this is a generalized version. If it is needed to be more specific, it could register an OPTIONS
	//   url for each possible path to more specifically lock down cross-origin requests.
	r.Register("OPTIONS", "/*path", Options)

	// Run the server
	if err := r.Run(":9000"); err != nil {
		fmt.Println(err)
	}
}
