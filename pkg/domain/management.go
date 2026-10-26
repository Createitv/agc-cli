package domain

// TemplateCodeUploadEndpoint is an internal stage of the reviewed SAM deploy
// workflow, not an application gallery or a public storage-management API.
func TemplateCodeUploadEndpoint() Endpoint {
	return Endpoint{ID: "template-code-upload", FamilyID: "templates", Name: "SAM code artifact upload", Method: "PUT", Path: "/v0/{bucket}/{object_name}", Authentication: "agcli-management-oauth", Direction: "developer-to-huawei", SourceURL: "https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0", RequestContentType: "application/x-www-form-urlencoded"}
}
