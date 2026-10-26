package domain

// The metadata update accepts model objects directly. Their names in the
// reference are type names, not extra JSON wrapper properties.
func applyPublishingMetadataModels(e Endpoint) Endpoint {
	if e.ID != "app-info-update-v3" {
		return e
	}
	const root = "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/"
	customerSource := root + "agc-help-publish-api-data-custinfo-0000002504067925"
	deviceSource := root + "agc-help-publish-api-data-devicetypeinfo-0000002236201270"
	contactSource := root + "agc-help-publish-api-data-appreviewcontact-0000002505438989"
	vocSource := root + "agc-help-publish-api-data-versionvoc-0000002495156924"
	for i := range e.Parameters {
		p := &e.Parameters[i]
		if p.In != "body" {
			continue
		}
		switch p.Name {
		case "customInfo":
			p.Properties = []Parameter{
				{Name: "phone", In: "body", Type: "string", MaxLength: 50, SourceURL: customerSource},
				{Name: "email", In: "body", Type: "string", MaxLength: 50, SourceURL: customerSource},
				{Name: "qqNum", In: "body", Type: "string", MaxLength: 50, SourceURL: customerSource},
				{Name: "internalPath", In: "body", Type: "string", MaxLength: 100, SourceURL: customerSource},
			}
		case "deviceTypes":
			p.ItemsProperties = []Parameter{
				{Name: "deviceType", In: "body", Type: "integer", Required: true, Enum: []string{"4", "5", "8", "12", "14", "19"}, SourceURL: deviceSource},
				{Name: "appAdapters", In: "body", Type: "string", MaxLength: 32, SourceURL: deviceSource},
				{Name: "remoteControl", In: "body", Type: "string", MaxLength: 16, Enum: []string{"1"}, SourceURL: deviceSource},
			}
		case "appReviewPhoneInfo", "appReviewEmailInfo":
			p.Properties = []Parameter{{Name: "account", In: "body", Type: "string", MaxLength: 500, SourceURL: contactSource}}
		case "versionVocs":
			p.ItemsProperties = []Parameter{
				{Name: "versionId", In: "body", Type: "string", Required: true, MaxLength: 32, SourceURL: vocSource},
				{Name: "packageId", In: "body", Type: "string", Required: true, MaxLength: 32, SourceURL: vocSource},
				{Name: "vocId", In: "body", Type: "string", Required: true, MaxLength: 32, SourceURL: vocSource},
			}
		}
	}
	return e
}
