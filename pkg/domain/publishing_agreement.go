package domain

// AgreementInfo has separate official models for user and privacy agreements.
// Keep nested rules in the domain so CLI, REST and OpenAPI share the same
// contract. Optional baseInfo becomes subject to its required name only when
// the caller supplies that object.
func applyPublishingAgreementModels(e Endpoint) Endpoint {
	privacy := e.ID == "privacy-agreement-update"
	if !privacy && e.ID != "user-agreement-create" && e.ID != "user-agreement-update" {
		return e
	}
	const root = "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/"
	infoSource := root + "agc-help-publish-api-data-agreementinfo-user-0000002236041482"
	baseSource := root + "agc-help-publish-api-data-agreementbaseinfo-user-0000002236201290"
	base := []Parameter{
		{Name: "name", In: "body", Type: "string", Required: true, MaxLength: 100, SourceURL: baseSource},
		{Name: "agreementTitle", In: "body", Type: "string", MaxLength: 100, SourceURL: baseSource},
		{Name: "userAgreementUrl", In: "body", Type: "string", MaxLength: 512, SourceURL: baseSource},
	}
	if privacy {
		infoSource = root + "agc-help-publish-api-data-agreementinfo-privacy-0000002294845630"
		baseSource = root + "publish-api-agreementbaseinfo-privacy-0000002328924929"
		base = []Parameter{
			{Name: "name", In: "body", Type: "string", Required: true, MaxLength: 100, SourceURL: baseSource},
			{Name: "appName", In: "body", Type: "string", MaxLength: 255, SourceURL: baseSource},
			{Name: "companyName", In: "body", Type: "string", MaxLength: 1024, SourceURL: baseSource},
			{Name: "briefDesc", In: "body", Type: "string", MaxLength: 500, SourceURL: baseSource},
			{Name: "effectiveTime", In: "body", Type: "integer", SourceURL: baseSource},
		}
	}
	for i := range e.Parameters {
		p := &e.Parameters[i]
		if p.In != "body" || p.Name != "agreementInfos" {
			continue
		}
		p.MaxItems = 100
		p.ItemsProperties = []Parameter{
			{Name: "lang", In: "body", Type: "string", Required: true, MaxLength: 20, Enum: []string{"zh", "en", "ug", "bo", "zh-Hant-HK", "zh-Hant-TW"}, SourceURL: infoSource},
			{Name: "baseInfo", In: "body", Type: "object", Properties: base, SourceURL: infoSource},
		}
		if privacy {
			p.ItemsProperties = append(p.ItemsProperties,
				Parameter{Name: "personalInfo", In: "body", Type: "object", Properties: publishingPrivacyFields("PersonalInfo"), SourceURL: infoSource},
				Parameter{Name: "devicePermissions", In: "body", Type: "array", ItemsType: "object", ItemsProperties: publishingPrivacyFields("DevicePermission"), MaxItems: 200, SourceURL: infoSource},
				Parameter{Name: "thirdSdkInfo", In: "body", Type: "object", Properties: publishingPrivacyFields("ThirdSdkInfo"), SourceURL: infoSource},
				Parameter{Name: "thirdSharedInfos", In: "body", Type: "array", ItemsType: "object", ItemsProperties: publishingPrivacyFields("ThirdSharedInfo"), MaxItems: 200, SourceURL: infoSource},
				Parameter{Name: "extInfo", In: "body", Type: "object", Properties: publishingPrivacyFields("ExtendInfo"), SourceURL: infoSource},
			)
		}
	}
	return e
}
