package domain

// verifiedParameters overrides the method-based legacy inference only for
// contracts verified against the linked Huawei reference or upstream responses.
// Optional entries expose known locations without asserting requiredness that
// has not been verified. This is intentionally not a complete API schema.
func verifiedParameters(family, id string, inferred []Parameter) []Parameter {
	p := func(name, location string, required bool, description string) Parameter {
		return Parameter{Name: name, In: location, Required: required, Description: description}
	}
	appHeader := p("appId", "header", true, "Huawei application ID")
	switch family + "/" + id {
	case "provisioning/provision-api-get-fingerprints", "provisioning/provision-api-getacl", "testing/test-api-get-test-grouplist":
		return []Parameter{appHeader}
	case "domains/domain-api-get-domain":
		return []Parameter{appHeader, p("category", "query", true, "Domain category")}
	case "testing/test-api-query-test-user":
		return []Parameter{appHeader, p("groupId", "header", true, "Test group ID")}
	case "upload/upload-url-new":
		return []Parameter{p("contentLength", "query", true, "Upload file length in bytes"), p("appId", "query", false, "Huawei application ID"), p("suffix", "query", false, "Upload file suffix"), p("fileName", "query", false, "Upload file name")}
	case "publishing/app-info-query":
		return []Parameter{p("appId", "query", true, "Huawei application ID"), p("lang", "query", false, "Application information language")}
	case "projects/queryservice":
		return []Parameter{p("projectId", "query", true, "Huawei project ID"), p("appID", "query", true, "Huawei application ID; parameter name is case sensitive")}
	case "comments/comapi-getreviews-harmonyos", "comments/com-rating-harmonyos":
		return harmonyOSCommentsParameters(id)
	case "reports/appdownloadexport":
		return append(inferred, p("language", "query", false, "Report language"), p("startTime", "query", false, "Start date in YYYYMMDD format"), p("endTime", "query", false, "End date in YYYYMMDD format"))
	case "reports/iapexport":
		return append(inferred, p("currency", "query", true, "Report currency"))
	case "reports/orderanalysisexport":
		return append(inferred, p("teamId", "header", true, "Huawei team ID"))
	default:
		return inferred
	}
}

// Huawei query tables for agcapi-comapi-getreviews-harmonyos-0000002470893976
// and agcapi-com-rating-harmonyos-0000002503933921, updated 2026-04-29.
func harmonyOSCommentsParameters(id string) []Parameter {
	p := func(name, typ, format string, required bool, description string) Parameter {
		return Parameter{Name: name, In: "query", Type: typ, Format: format, Required: required, Description: description}
	}
	result := []Parameter{
		p("appId", "string", "", true, "Huawei application ID"),
		p("beginTime", "integer", "int64", true, "Start timestamp in milliseconds since Unix epoch; query interval must not exceed six months"),
		p("endTime", "integer", "int64", true, "End timestamp in milliseconds since Unix epoch; query interval must not exceed six months"),
		p("countries", "string", "", true, "Comma-separated country codes, such as CN; all countries must belong to the same regional site"),
		p("ratings", "string", "", false, "Comma-separated ratings from 1 to 5"),
		p("appVersions", "string", "", false, "Comma-separated application versions"),
		p("sort", "integer", "int32", false, "0: descending time (default); 1: descending score; 2: ascending score"),
		p("page", "integer", "int32", false, "Page number; defaults to 1"),
	}
	if id == "comapi-getreviews-harmonyos" {
		result = append(result, p("limit", "integer", "int32", false, "Page size; defaults to 20, maximum 100"),
			p("content", "string", "", false, "Comment content keyword"),
			p("devReplyStates", "string", "", false, "Comma-separated developer reply states: 0, 1, 6, 3"),
			p("langs", "string", "", false, "Comma-separated ISO 639-1 language codes"))
	} else {
		result = append(result, p("limit", "integer", "int32", false, "Page size; defaults to 50, maximum 100"))
	}
	return result
}
