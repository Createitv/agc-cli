package domain

type Capability struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Command       string            `json:"command"`
	RESTPath      string            `json:"restPath"`
	Status        string            `json:"status"`
	EndpointCount int               `json:"endpointCount"`
	Links         map[string]Link   `json:"_links,omitempty"`
	Affordances   map[string]string `json:"affordances,omitempty"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

var Capabilities = []Capability{
	{ID: "publishing", Name: "Publishing API", Description: "Manage app and atomic service release metadata, submissions, review withdrawal, release timing, and version status.", Command: "agc publishing", RESTPath: "/api/v1/publishing", Status: "registered"},
	{ID: "upload", Name: "Upload Management API", Description: "Upload APP packages, icons, screenshots, videos, PDFs, and other release files.", Command: "agc upload", RESTPath: "/api/v1/upload", Status: "registered"},
	{ID: "provisioning", Name: "Provisioning API", Description: "Manage certificates, profiles, and test devices for HarmonyOS debug and release workflows.", Command: "agc provisioning", RESTPath: "/api/v1/provisioning", Status: "registered"},
	{ID: "domains", Name: "Domain Management API", Description: "Query, download, and update atomic service domain configuration.", Command: "agc domains", RESTPath: "/api/v1/domains", Status: "registered"},
	{ID: "testing", Name: "Testing API", Description: "Manage HarmonyOS app and atomic service test versions, testers, invitations, and feedback.", Command: "agc testing", RESTPath: "/api/v1/testing", Status: "registered"},
	{ID: "reports", Name: "Reports API", Description: "Request and download AppGallery Connect reports as CSV or Excel files with filters and dimensions.", Command: "agc reports", RESTPath: "/api/v1/reports", Status: "registered"},
	{ID: "projects", Name: "Project Management API", Description: "Read developer teams, projects, app summaries, and SDK configuration files.", Command: "agc projects", RESTPath: "/api/v1/projects", Status: "registered"},
	{ID: "comments", Name: "Comments API", Description: "Fetch, analyze, export, and reply to AppGallery Connect comments.", Command: "agc comments", RESTPath: "/api/v1/comments", Status: "registered"},
	{ID: "pms", Name: "PMS API", Description: "Manage products, subscriptions, prices, promotions, and localizations.", Command: "agc pms", RESTPath: "/api/v1/pms", Status: "registered"},
	{ID: "gameplay", Name: "Game Playing Service", Description: "Manage supported playing-service records and automation surfaces exposed by Connect API.", Command: "agc gameplay", RESTPath: "/api/v1/gameplay", Status: "registered"},
	{ID: "game-items", Name: "Game Item Mall", Description: "Manage game item mall products and presentation data separately from generic PMS goods.", Command: "agc game-items", RESTPath: "/api/v1/game-items", Status: "registered"},
	{ID: "resources", Name: "Resource Package Predownload", Description: "Manage resource packages, predownload tasks, upload linkage, and task status.", Command: "agc resources", RESTPath: "/api/v1/resources", Status: "registered"},
	{ID: "cicd", Name: "CI/CD Platform", Description: "Manage Huawei CI/CD pipelines and bridge local Hvigor builds into release automation.", Command: "agc cicd", RESTPath: "/api/v1/cicd", Status: "registered"},
}

func CapabilityByID(id string) (Capability, bool) {
	for _, capability := range Capabilities {
		if capability.ID == id {
			return decorateCapability(capability), true
		}
	}
	return Capability{}, false
}

func DecoratedCapabilities() []Capability {
	out := make([]Capability, len(Capabilities))
	for i, capability := range Capabilities {
		out[i] = decorateCapability(capability)
	}
	return out
}

func decorateCapability(capability Capability) Capability {
	capability.EndpointCount = len(EndpointFamilies[capability.ID])
	capability.Links = map[string]Link{
		"self":      {Href: capability.RESTPath, Method: "GET"},
		"endpoints": {Href: capability.RESTPath + "/endpoints", Method: "GET"},
		"run":       {Href: "/api/run", Method: "POST"},
	}
	capability.Affordances = map[string]string{
		"list":      capability.Command + " list",
		"endpoints": capability.Command + " endpoints",
		"docs":      "agc docs " + capability.ID,
	}
	return capability
}
