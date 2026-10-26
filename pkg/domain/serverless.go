package domain

import "encoding/json"

// Serverless contracts are generated from browser-read official parameter tables,
// then source-reviewed overrides. They do not share the publishing credential type.
func init() {
	raw, err := referenceDocuments.ReadFile("reference/serverless/endpoints.json")
	if err != nil {
		panic(err)
	}
	var endpoints []Endpoint
	if err = json.Unmarshal(raw, &endpoints); err != nil {
		panic(err)
	}
	for _, e := range endpoints {
		e.Body = inferBody(e.Method, e.Parameters)
		for _, p := range e.Parameters {
			if p.Required {
				e.RequiredParams = append(e.RequiredParams, p.Name)
			}
		}
		EndpointFamilies[e.FamilyID] = append(EndpointFamilies[e.FamilyID], e)
	}
	raw, err = referenceDocuments.ReadFile("reference/management/endpoints.json")
	if err != nil {
		panic(err)
	}
	// Unmarshal reuses non-nil slice elements and retains omitted struct fields.
	// Management parameters must not inherit enums/bounds from the REST catalog.
	endpoints = nil
	if err = json.Unmarshal(raw, &endpoints); err != nil {
		panic(err)
	}
	for _, e := range endpoints {
		e.Body = inferBody(e.Method, e.Parameters)
		for _, p := range e.Parameters {
			if p.Required {
				e.RequiredParams = append(e.RequiredParams, p.Name)
			}
		}
		EndpointFamilies[e.FamilyID] = append(EndpointFamilies[e.FamilyID], e)
	}
	for _, c := range []struct{ id, name, description string }{
		{"cloud-auth", "Cloud Authentication", "Client sign-in and server user management; separate from developer authorization."},
		{"clouddb", "Cloud DB for Object", "Client REST data operations plus source-reviewed schema, zone and async-process management with explicit user sessions."},
		{"storage", "Cloud Storage", "File bytes, metadata and listing through server REST."},
		{"hosting", "Cloud Hosting", "Website sites, environments, signed ZIP uploads and version lifecycle."},
		{"templates", "Serverless Templates", "Authentication callbacks and source-reviewed SAM application/plan/deployment workflows; console gallery templates are distinct."},
		{"cache", "Cloud Cache", "Official AGCLI source management and Redis-command transport using an explicit user session."},
		{"prefetch", "Prefetch", "Browser-reviewed developer-server protocol simulation and device debugging; not a cloud configuration API."},
	} {
		Capabilities = append(Capabilities, Capability{ID: c.id, Name: c.name, Description: c.description, Command: "agc " + c.id, RESTPath: "/api/v1/" + c.id, Status: "registered"})
	}
}

func IsServerlessFamily(family string) bool {
	switch family {
	case "cloud-auth", "clouddb", "storage", "hosting", "templates", "functions", "cache", "prefetch":
		return true
	}
	return false
}
