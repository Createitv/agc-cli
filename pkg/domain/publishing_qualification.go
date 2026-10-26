package domain

func applyPublishingQualificationModels(e Endpoint) Endpoint {
	if e.ID != "qualification-review-create" && e.ID != "qualification-review-update" {
		return e
	}
	const root = "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/"
	tagSource := root + "agc-help-publish-api-data-tagqualinfo-0000002236041478"
	qualSource := root + "agc-help-publish-api-data-qualificationinfo-0000002236201286"
	fileSource := root + "agc-help-publish-api-data-qualificationfileinfo-0000002271160609"
	files := []Parameter{
		{Name: "objectId", In: "body", Type: "string", Required: true, MaxLength: 128, SourceURL: fileSource},
		{Name: "fileName", In: "body", Type: "string", Required: true, MaxLength: 128, SourceURL: fileSource},
		{Name: "fileSize", In: "body", Type: "integer", Required: true, SourceURL: fileSource},
		{Name: "sha256", In: "body", Type: "string", Required: true, MaxLength: 128, SourceURL: fileSource},
	}
	qualifications := []Parameter{
		{Name: "qualificationCode", In: "body", Type: "string", Required: true, MaxLength: 32, SourceURL: qualSource},
		{Name: "qualificationFileInfos", In: "body", Type: "array", Required: true, ItemsType: "object", ItemsProperties: files, MaxItems: 10, SourceURL: qualSource},
		{Name: "validityType", In: "body", Type: "integer", Required: true, Enum: []string{"0", "1"}, SourceURL: qualSource},
		{Name: "startTime", In: "body", Type: "string", Format: "date", MaxLength: 20, SourceURL: qualSource},
		{Name: "expireTime", In: "body", Type: "string", Format: "date", MaxLength: 20, SourceURL: qualSource},
	}
	for i := range e.Parameters {
		p := &e.Parameters[i]
		if p.In != "body" || p.Name != "tagQualInfos" {
			continue
		}
		p.ItemsProperties = []Parameter{
			{Name: "tagId", In: "body", Type: "string", Required: true, MaxLength: 64, SourceURL: tagSource},
			{Name: "taskId", In: "body", Type: "string", Required: e.ID == "qualification-review-update", MaxLength: 32, SourceURL: tagSource},
			{Name: "qualificationInfos", In: "body", Type: "array", ItemsType: "object", MaxItems: 10, ItemsProperties: qualifications, SourceURL: tagSource, ItemsRules: []ReferenceRule{
				{ID: "qualification_remaining_validity", Kind: "minimum_remaining_days", Field: "body.expireTime", MinDays: 30, When: &RuleCondition{Field: "body.validityType", Equals: "0", Default: "0"}, SourceURL: qualSource, Message: "Custom qualification expiry must have more than 30 days remaining.", Suggestion: "Use a real qualification expiry date in yyyy-MM-dd format with more than 30 days remaining; do not change the source certificate's validity."},
			}},
		}
	}
	return e
}
