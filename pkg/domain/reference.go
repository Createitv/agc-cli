package domain

import (
	"embed"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//go:embed reference/*/*.json
var referenceDocuments embed.FS

type ReferenceSummary struct {
	Title       string `json:"title"`
	SourceURL   string `json:"sourceUrl"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	CheckedAt   string `json:"checkedAt,omitempty"`
	Status      string `json:"status"`
	ReadCommand string `json:"readCommand"`
}

type RuleCondition struct {
	Field   string `json:"field"`
	Equals  string `json:"equals"`
	Default string `json:"default,omitempty"`
}
type ReferenceRule struct {
	UnlessPresent string         `json:"unlessPresent,omitempty"`
	MinItems      int            `json:"minItems,omitempty"`
	MaxItems      int            `json:"maxItems,omitempty"`
	MinDays       int            `json:"minDays,omitempty"`
	MaxDays       int            `json:"maxDays,omitempty"`
	ID            string         `json:"id"`
	Kind          string         `json:"kind"`
	Field         string         `json:"field,omitempty"`
	Fields        []string       `json:"fields,omitempty"`
	When          *RuleCondition `json:"when,omitempty"`
	SourceURL     string         `json:"sourceUrl"`
	Message       string         `json:"message"`
	Suggestion    string         `json:"suggestion"`
}

type InputViolation struct {
	Code       string `json:"code"`
	Field      string `json:"field"`
	Message    string `json:"message"`
	SourceURL  string `json:"sourceUrl"`
	Suggestion string `json:"suggestion"`
}

func EndpointReference(family, id string) (json.RawMessage, bool) {
	data, err := referenceDocuments.ReadFile("reference/" + family + "/" + id + ".json")
	if err != nil {
		return nil, false
	}
	// Supplemental requirements and request models are separate official pages.
	// Return their full source content with the endpoint rather than bare links.
	var supplemental []string
	if (family == "publishing" && (id == "app-file-info-v3" || id == "app-file-info")) || family == "upload" {
		supplemental = append(supplemental, "application-file-requirements")
		if family == "publishing" && id == "app-file-info-v3" {
			supplemental = append(supplemental, "lang-file-info", "file-info", "package-permission-intro-video")
		}
	}
	if family == "pms" {
		switch id {
		case "batchaddproduct-android", "batchupdateproduct-android":
			supplemental = append(supplemental, "product-import-android")
		case "batchaddproduct-harmonyosnext", "batchupdateproduct-harmonyosnext":
			supplemental = append(supplemental, "product-import-harmonyosnext")
		}
	}
	if family == "pms" && len(supplemental) > 0 {
		platform := "android"
		if supplemental[0] == "product-import-harmonyosnext" {
			platform = "harmonyosnext"
		}
		for _, model := range []string{"defaultproductpriceinfo", "productlanguageinfo", "productpriceinfo", "promotioninfo"} {
			supplemental = append(supplemental, model+"-"+platform)
		}
	}
	if family == "domains" {
		switch id {
		case "domain-api-post-domain":
			supplemental = append(supplemental, "domain-info-server", "domain-info-business")
		case "domain-api-precheck":
			supplemental = append(supplemental, "domain-info-business")
		}
	}
	if family == "provisioning" {
		switch id {
		case "provision-api-add-device":
			supplemental = append(supplemental, "add-device-info")
		case "provision-api-applyacl":
			supplemental = append(supplemental, "acl-apply-permission-info")
		}
	}
	if family == "testing" {
		switch id {
		case "test-api-add-test-package":
			supplemental = append(supplemental, "testing-package-file-info")
		case "test-api-modify-test-version":
			supplemental = append(supplemental, "testing-app-language", "testing-device-material", "testing-device-type-info", "testing-open-test-info", "testing-test-task-info", "testing-group-info", "testing-package-permission-intro-video", "testing-version-ext-info", "testing-app-review-contact", "testing-version-voc-req")
		case "test-api-query-feedback":
			supplemental = append(supplemental, "testing-feedback-list-request")
		case "test-api-delete-feedback":
			supplemental = append(supplemental, "testing-delete-feedback-request", "testing-delete-feedback-info")
		}
	}
	if family == "resources" {
		supplemental = append(supplemental, "developer-project-id-guide")
	}
	if family == "publishing" {
		supplemental = append(supplemental, "publishing-error-codes")
	}
	if len(supplemental) > 0 {
		var doc map[string]json.RawMessage
		if json.Unmarshal(data, &doc) == nil {
			var related []json.RawMessage
			for _, name := range supplemental {
				if model, err := referenceDocuments.ReadFile("reference/supplemental/" + name + ".json"); err == nil {
					related = append(related, model)
				}
			}
			doc["relatedDocuments"], _ = json.Marshal(related)
			data, _ = json.Marshal(doc)
		}
	}
	if family == "publishing" {
		data = attachPublishingLinkedModels(data)
	}
	if IsServerlessFamily(family) {
		var doc map[string]json.RawMessage
		if json.Unmarshal(data, &doc) == nil {
			var related []json.RawMessage
			_ = json.Unmarshal(doc["relatedDocuments"], &related)
			entries, _ := referenceDocuments.ReadDir("reference/supplemental")
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), "serverless-") {
					if raw, err := referenceDocuments.ReadFile("reference/supplemental/" + entry.Name()); err == nil {
						related = append(related, raw)
					}
				}
			}
			doc["relatedDocuments"], _ = json.Marshal(related)
			data, _ = json.Marshal(doc)
		}
	}
	return json.RawMessage(data), true
}

// Follow captured links through reviewed model documents. This includes
// response models and transitive prerequisites without guessing from model
// names (user and privacy AgreementInfo have different source URLs).
func attachPublishingLinkedModels(data []byte) []byte {
	var doc map[string]json.RawMessage
	if json.Unmarshal(data, &doc) != nil {
		return data
	}
	var related []json.RawMessage
	_ = json.Unmarshal(doc["relatedDocuments"], &related)
	models := map[string]json.RawMessage{}
	entries, _ := referenceDocuments.ReadDir("reference/supplemental")
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		raw, err := referenceDocuments.ReadFile("reference/supplemental/" + entry.Name())
		if err != nil {
			continue
		}
		var model struct {
			SourceURL string `json:"sourceUrl"`
			Status    string `json:"status"`
		}
		if json.Unmarshal(raw, &model) == nil && model.Status == "browser-read" &&
			(strings.Contains(model.SourceURL, "publish-api-data-") || strings.Contains(model.SourceURL, "publish-api-agreementbaseinfo") || strings.Contains(model.SourceURL, "publish-api-devicemulti") || strings.Contains(model.SourceURL, "publish-api-packagepermission")) {
			models[model.SourceURL] = raw
		}
	}
	seen := map[string]bool{}
	queue := []json.RawMessage{data}
	for _, raw := range related {
		var model struct {
			SourceURL string `json:"sourceUrl"`
		}
		_ = json.Unmarshal(raw, &model)
		seen[model.SourceURL] = true
		queue = append(queue, raw)
	}
	for len(queue) > 0 {
		raw := queue[0]
		queue = queue[1:]
		var model struct {
			Links []struct {
				SourceURL string `json:"sourceUrl"`
				Title     string `json:"title"`
			}
			Tables []struct{ Rows [][]string }
		}
		if json.Unmarshal(raw, &model) != nil {
			continue
		}
		// Captures also contain previous/next navigation. Only links naming a
		// type in a parameter table are semantic model dependencies.
		types := map[string]bool{}
		for _, table := range model.Tables {
			for _, row := range table.Rows {
				if len(row) < 3 {
					continue
				}
				name := strings.TrimPrefix(row[2], "List<")
				name = strings.TrimSuffix(name, ">")
				types[name] = true
			}
		}
		for _, link := range model.Links {
			if !types[link.Title] {
				continue
			}
			url := strings.SplitN(link.SourceURL, "#", 2)[0]
			if next, ok := models[url]; ok && !seen[url] {
				seen[url] = true
				related = append(related, next)
				queue = append(queue, next)
			}
		}
	}
	if len(related) == 0 {
		return data
	}
	doc["relatedDocuments"], _ = json.Marshal(related)
	result, err := json.Marshal(doc)
	if err != nil {
		return data
	}
	return result
}

type referenceMetadata struct {
	summary     ReferenceSummary
	description string
	rules       []ReferenceRule
}

var referenceMetadataCache sync.Map

func attachReference(e Endpoint) Endpoint {
	// Full examples and public vendor source are read only by --docs. Repeated
	// command discovery needs immutable metadata, not another full-document parse.
	if e.FamilyID == "functions" && e.Documentation != nil {
		summary := *e.Documentation
		e.Documentation = &summary
		return e
	}
	key := e.FamilyID + "/" + e.ID
	cached, ok := referenceMetadataCache.Load(key)
	if !ok {
		data, err := referenceDocuments.ReadFile("reference/" + key + ".json")
		if err != nil {
			return e
		}
		var meta referenceMetadata
		if json.Unmarshal(data, &meta.summary) != nil {
			return e
		}
		if e.FamilyID == "publishing" {
			var content struct {
				Notes []struct{ Section, Text string }
			}
			if json.Unmarshal(data, &content) == nil {
				for _, note := range content.Notes {
					if note.Section == "功能介绍" && note.Text != "" {
						meta.description = note.Text
						break
					}
				}
			}
		}
		meta.rules = referenceRules(e.FamilyID, e.ID)
		cached, _ = referenceMetadataCache.LoadOrStore(key, meta)
	}
	meta := cached.(referenceMetadata)
	summary := meta.summary
	summary.ReadCommand = e.Command + " --docs --pretty"
	e.Documentation = &summary
	if meta.description != "" {
		e.Description = meta.description
	}
	e.Rules = make([]ReferenceRule, len(meta.rules))
	for i, rule := range meta.rules {
		e.Rules[i] = rule
		e.Rules[i].Fields = append([]string(nil), rule.Fields...)
		if rule.When != nil {
			condition := *rule.When
			e.Rules[i].When = &condition
		}
	}
	return e
}

func referenceRules(family, id string) []ReferenceRule {
	data, err := referenceDocuments.ReadFile("reference/" + family + "/" + id + ".rules.json")
	if err != nil {
		return nil
	}
	var rules []ReferenceRule
	if json.Unmarshal(data, &rules) != nil {
		return nil
	}
	return rules
}

// Enrich only fields whose actual official table has been browser-read. Integer
// widths such as Integer(32) are not interpreted as value ranges.
func applyReferenceParameters(e Endpoint) Endpoint {
	// Parameter parsing uses the endpoint's own tables; related model contracts
	// are applied explicitly below, not inferred from response documentation.
	data, err := referenceDocuments.ReadFile("reference/" + e.FamilyID + "/" + e.ID + ".json")
	if err != nil {
		return e
	}
	var doc struct {
		Tables []struct {
			Section   string
			Location  string
			SourceURL string `json:"sourceUrl"`
			Rows      [][]string
		}
	}
	if json.Unmarshal(data, &doc) != nil {
		return e
	}
	if parameterData, err := referenceDocuments.ReadFile("reference/" + e.FamilyID + "/" + e.ID + ".parameters.json"); err == nil {
		var reviewed []Parameter
		if json.Unmarshal(parameterData, &reviewed) == nil {
			// Retain independently verified bounds while replacing legacy inferred
			// locations, types and requiredness with the actual parameter table.
			for i := range reviewed {
				for _, previous := range e.Parameters {
					if previous.Name == reviewed[i].Name && previous.In == reviewed[i].In {
						if reviewed[i].MinItems == 0 {
							reviewed[i].MinItems = previous.MinItems
						}
						if reviewed[i].MaxItems == 0 {
							reviewed[i].MaxItems = previous.MaxItems
						}
					}
				}
			}
			for _, previous := range e.Parameters {
				if previous.In == "path" || previous.In == "file" {
					found := false
					for _, p := range reviewed {
						if p.Name == previous.Name && p.In == previous.In {
							found = true
							break
						}
					}
					if !found {
						reviewed = append([]Parameter{previous}, reviewed...)
					}
				}
			}
			e.Parameters = reviewed
		}
	}
	for i := range e.Parameters {
		p := &e.Parameters[i]
		// Reviewed model contracts already identify the correct nesting and paragraph.
		if p.SourceURL != "" {
			continue
		}
	matchedParameter:
		for _, table := range doc.Tables {
			if table.Section != "请求参数" || table.Location != p.In {
				continue
			}
			for _, row := range table.Rows {
				if len(row) < 4 || row[0] != p.Name {
					continue
				}
				if match := regexp.MustCompile(`^List<String\((\d+)\)>$`).FindStringSubmatch(row[2]); len(match) > 1 {
					p.ItemsMaxLength, _ = strconv.Atoi(match[1])
				}
				p.Description = row[3]
				p.SourceURL = table.SourceURL
				break matchedParameter
			}
		}
	}
	extra, err := referenceDocuments.ReadFile("reference/" + e.FamilyID + "/" + e.ID + ".constraints.json")
	if err == nil {
		var constraints map[string]Parameter
		if json.Unmarshal(extra, &constraints) == nil {
			for i := range e.Parameters {
				p := &e.Parameters[i]
				c, ok := constraints[p.In+"."+p.Name]
				if !ok {
					continue
				}
				p.MinLength, p.MaxLength = c.MinLength, c.MaxLength
				p.MinItems, p.MaxItems = c.MinItems, c.MaxItems
				p.Minimum, p.Maximum, p.Enum = c.Minimum, c.Maximum, c.Enum
			}
		}
	}
	// The table marks appId mandatory but the packageName remark explicitly
	// permits either identity; the rule below enforces the documented alternative.
	if e.FamilyID == "publishing" && e.ID == "app-version-list" {
		for i := range e.Parameters {
			if e.Parameters[i].Name == "appId" {
				e.Parameters[i].Required = false
			}
		}
	}
	if e.FamilyID == "publishing" && e.ID == "app-file-info-v3" {
		langSource := "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-langfileinfo-0000002271000653"
		fileSource := "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-fileinfo-0000002236041466"
		fileFields := []Parameter{
			{Name: "deviceType", In: "body", Type: "integer", Required: true, Enum: []string{"4", "5", "8", "12", "14", "19"}, SourceURL: fileSource},
			{Name: "objectIdList", In: "body", Type: "array", Required: true, ItemsType: "string", ItemsMaxLength: 512, MaxItems: 10, SourceURL: fileSource},
			{Name: "showType", In: "body", Type: "integer", Enum: []string{"0", "1"}, SourceURL: fileSource},
		}
		for i := range e.Parameters {
			p := &e.Parameters[i]
			switch p.Name {
			case "packagePermissionIntroVideoList":
				source := "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/publish-api-packagepermissionintrovedio-0000002236201278"
				p.ItemsProperties = []Parameter{
					{Name: "lang", In: "body", Type: "string", Required: true, MaxLength: 8, SourceURL: source},
					{Name: "permissionName", In: "body", Type: "string", Required: true, MaxLength: 512, SourceURL: source},
					{Name: "deviceType", In: "body", Type: "integer", Required: true, Enum: []string{"4", "5", "8", "12", "14", "19"}, SourceURL: source},
					{Name: "objectId", In: "body", Type: "string", Required: true, MaxLength: 512, SourceURL: source},
				}
			case "appIconList", "screenShotList", "introVideoList", "rcmdVideoList", "rcmdPicList", "updateScreenShotList", "updateVideoList":
				p.ItemsProperties = []Parameter{
					{Name: "lang", In: "body", Type: "string", Required: true, MaxLength: 8, SourceURL: langSource},
					{Name: "fileInfoList", In: "body", Type: "array", Required: true, ItemsType: "object", MaxItems: 20, ItemsProperties: fileFields, SourceURL: langSource},
				}
			}
		}
	}
	if e.FamilyID == "pms" {
		var platform string
		switch e.ID {
		case "batchaddproduct-android", "batchupdateproduct-android":
			platform = "android"
		case "batchaddproduct-harmonyosnext", "batchupdateproduct-harmonyosnext":
			platform = "harmonyosnext"
		}
		if platform != "" {
			data, err := referenceDocuments.ReadFile("reference/supplemental/product-import-" + platform + ".contract.json")
			var model struct {
				Parameters []Parameter
				Rules      []ReferenceRule
			}
			if err == nil && json.Unmarshal(data, &model) == nil {
				for i := range e.Parameters {
					if e.Parameters[i].Name == "products" {
						e.Parameters[i].ItemsProperties = model.Parameters
						e.Parameters[i].ItemsRules = model.Rules
					}
				}
			}
		}
	}
	if e.FamilyID == "upload" && (e.ID == "obbfile-uploadurl" || e.ID == "obbfile-compose") {
		var bodyTables []struct {
			SourceURL string
			Rows      [][]string
		}
		for _, table := range doc.Tables {
			if table.Section == "请求参数" && table.Location == "body" {
				bodyTables = append(bodyTables, struct {
					SourceURL string
					Rows      [][]string
				}{table.SourceURL, table.Rows})
			}
		}
		if len(bodyTables) == 2 {
			model := Parameter{Type: "object", Required: true, SourceURL: bodyTables[0].SourceURL}
			for _, row := range bodyTables[1].Rows[1:] {
				p := Parameter{Name: row[0], In: "body", Required: row[1] == "M", Description: row[3], SourceURL: bodyTables[1].SourceURL, Type: "string"}
				if row[2] == "Integer(64)" {
					p.Type = "integer"
					p.Format = "int64"
				} else if row[2] == "String(64)" {
					p.MaxLength = 64
				}
				model.Properties = append(model.Properties, p)
			}
			e.DynamicBody = &model
			e.DynamicBodyKeyPattern = "^additionalProp[1-9][0-9]*$"
		}
	}
	e.RequiredParams = nil
	for _, p := range e.Parameters {
		if p.Required {
			e.RequiredParams = append(e.RequiredParams, p.Name)
		}
	}
	return e
}
