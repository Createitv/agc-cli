package domain

// These are the reviewed privacy-policy request models. Each field retains
// the exact model URL; the same structure drives CLI/REST validation and
// OpenAPI. Integer widths describe types, not numeric bounds. All nested
// fields below are optional, as marked in the official model tables.
func publishingPrivacyFields(model string) []Parameter {
	const root = "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/"
	sources := map[string]string{
		"PersonalInfo":           "agc-help-publish-api-data-personalinfo-0000002328805173",
		"RequiredInfo":           "agc-help-publish-api-data-requiredinfo-0000002295005294",
		"RecommendInfo":          "agc-help-publish-api-data-recommendinfo-0000002328805177",
		"AdsInfo":                "agc-help-publish-api-data-adsinfo-0000002295005298",
		"FullServiceMode":        "agc-help-publish-api-data-fullservicemode-0000002294845638",
		"ContractObligationMode": "agc-help-publish-api-data-contractobligationmode-0000002328924941",
		"LawRegulationMode":      "agc-help-publish-api-data-lawregulationmode-0000002295005342",
		"ServiceTypeInfo":        "agc-help-publish-api-data-servicetypeinfo-0000002328805213",
		"SceneItem":              "agc-help-publish-api-data-sceneitem-0000002294845634",
		"DataItem":               "agc-help-publish-api-data-dataitem-0000002328924933",
		"DevicePermission":       "agc-help-publish-api-data-devicepermission-0000002294845674",
		"Permission":             "agc-help-publish-api-data-permission-0000002328924977",
		"ThirdSdkInfo":           "agc-help-publish-api-data-thirdsdkinfo-0000002328805217",
		"SdkInfo":                "agc-help-publish-api-data-sdkinfo-0000002295005346",
		"ThirdSharedInfo":        "agc-help-publish-api-data-thirdsharedinfo-0000002294845678",
		"ExtendInfo":             "agc-help-publish-api-data-extendinfo-0000002328924981",
		"ContactInfo":            "agc-help-publish-api-data-contactinfo-0000002328805225",
		"StorageInfo":            "agc-help-publish-api-data-storageinfo-0000002295005350",
		"VersionChangeDesc":      "agc-help-publish-api-data-versionchangedesc-0000002294845686",
		"ManagePersonalInfo":     "agc-help-publish-api-data-managepersionalinfo-0000002328924985",
		"InfoAccess":             "agc-help-publish-api-data-infoaccess-0000002328805233",
		"InfoModify":             "agc-help-publish-api-data-infomodify-0000002295005354",
		"InfoDelete":             "agc-help-publish-api-data-infodelete-0000002294845690",
		"InfoCancel":             "agc-help-publish-api-data-infocancel-0000002328924989",
		"ServiceStop":            "agc-help-publish-api-data-servicestop-0000002328805237",
		"CustomInfo":             "agc-help-publish-api-data-custominfo-0000002295005358",
		"CustomChapter":          "agc-help-publish-api-data-customchapter-0000002294845694",
		"ExtFieldCustomInfo":     "agc-help-publish-api-data-extfieldcustominfo-0000002328924997",
	}
	source := root + sources[model]
	text := func(name string, length int) Parameter {
		return Parameter{Name: name, In: "body", Type: "string", MaxLength: length, SourceURL: source}
	}
	object := func(name, child string) Parameter {
		return Parameter{Name: name, In: "body", Type: "object", Properties: publishingPrivacyFields(child), SourceURL: source}
	}
	list := func(name, child string, maximum int) Parameter {
		return Parameter{Name: name, In: "body", Type: "array", ItemsType: "object", ItemsProperties: publishingPrivacyFields(child), MaxItems: maximum, SourceURL: source}
	}
	switch model {
	case "PersonalInfo":
		return []Parameter{object("requiredInfo", "RequiredInfo"), object("recommendInfo", "RecommendInfo"), object("adsInfo", "AdsInfo"), object("fullServiceMode", "FullServiceMode"), object("contractObligationMode", "ContractObligationMode"), object("lawRegulationMode", "LawRegulationMode")}
	case "RequiredInfo":
		return []Parameter{list("sceneItems", "SceneItem", 200)}
	case "RecommendInfo":
		return []Parameter{list("sceneItems", "SceneItem", 200), text("closeSetting", 500)}
	case "AdsInfo":
		return []Parameter{list("sceneItems", "SceneItem", 200), text("pushContent", 500), text("closeSetting", 500)}
	case "FullServiceMode":
		return []Parameter{text("selectPath", 500), text("personalCollectType", 500), text("basicServiceType", 500), text("extraServiceType", 500)}
	case "ContractObligationMode":
		return []Parameter{text("serviceType", 500), text("personalCollectReason", 500), list("serviceTypeInfos", "ServiceTypeInfo", 100)}
	case "LawRegulationMode":
		return []Parameter{list("serviceTypeInfos", "ServiceTypeInfo", 100)}
	case "ServiceTypeInfo":
		return []Parameter{text("serviceType", 100), list("sceneItems", "SceneItem", 200)}
	case "SceneItem":
		return []Parameter{list("dataItems", "DataItem", 200), text("sceneId", 64), text("sceneReason", 500), text("collectModeId", 64), text("collectMode", 500), text("legalBasisId", 64), text("legalBasis", 500), text("remark", 500)}
	case "DataItem":
		return []Parameter{text("itemId", 64), text("itemDesc", 500), {Name: "sensitiveData", In: "body", Type: "integer", Enum: []string{"0", "1"}, SourceURL: source}}
	case "DevicePermission":
		permissions := list("permissions", "Permission", 200)
		permissions.MinItems = 1
		return []Parameter{permissions, text("reasonId", 64), text("reason", 500), text("groupId", 100), text("groupName", 256)}
	case "Permission":
		return []Parameter{text("permissionName", 100), text("permissionDesc", 500)}
	case "ThirdSdkInfo":
		return []Parameter{list("sdkInfos", "SdkInfo", 50)}
	case "SdkInfo":
		return []Parameter{text("name", 100), text("providerName", 100), text("version", 100), text("privacyDeclare", 512), list("privacyLabels", "SceneItem", 200), text("scene", 200), text("collectType", 100), text("purpose", 500)}
	case "ThirdSharedInfo":
		return []Parameter{text("name", 100), text("purpose", 500), text("personalType", 500), text("officialPrivacyUrl", 1024), text("serviceType", 100), text("scene", 200), text("shareType", 100), list("dataItems", "DataItem", 200)}
	case "ExtendInfo":
		return []Parameter{object("contactInfo", "ContactInfo"), object("storageInfo", "StorageInfo"), text("childPrivacyUrl", 512), list("versionChangeDescs", "VersionChangeDesc", 0), text("abstracts", 300), object("managePersonalInfo", "ManagePersonalInfo"), object("customInfo", "CustomInfo"), text("extField", 8192), object("extFieldCustomInfo", "ExtFieldCustomInfo")}
	case "ContactInfo":
		return []Parameter{
			{Name: "email", In: "body", Type: "array", ItemsType: "string", MaxItems: 10, ItemsMaxLength: 100, SourceURL: source},
			{Name: "phoneNum", In: "body", Type: "array", ItemsType: "string", MaxItems: 10, ItemsMaxLength: 31, SourceURL: source},
			text("companyAddress", 512), text("customContactInfo", 500),
		}
	case "StorageInfo":
		return []Parameter{text("storagePeriod", 512), text("country", 100), {Name: "isCrossBorder", In: "body", Type: "boolean", SourceURL: source}, text("crossBorderReason", 500), text("crossBorderReceiveName", 100), list("crossBorderDeclarations", "SceneItem", 200), text("crossBorderCountry", 100)}
	case "VersionChangeDesc":
		return []Parameter{text("chapter", 50), text("desc", 500)}
	case "ManagePersonalInfo":
		return []Parameter{object("infoAccess", "InfoAccess"), object("infoModify", "InfoModify"), object("infoDelete", "InfoDelete"), object("infoCancel", "InfoCancel"), object("serviceStop", "ServiceStop")}
	case "InfoAccess", "InfoModify", "InfoDelete":
		return []Parameter{text("personalSettingPath", 500), text("personalType", 500)}
	case "InfoCancel":
		return []Parameter{text("personalSettingPath", 500), text("serviceType", 200)}
	case "ServiceStop":
		return []Parameter{text("personalSettingPath", 500)}
	case "CustomInfo":
		return []Parameter{text("title", 100), text("customInfoDesc", 5000), list("customInfos", "CustomChapter", 5)}
	case "CustomChapter":
		return []Parameter{text("title", 100), text("customInfoDesc", 3000)}
	case "ExtFieldCustomInfo":
		return []Parameter{text("childPrivacyTitle", 512), text("childPrivacyCustomInfoDesc", 5000)}
	}
	return nil
}
