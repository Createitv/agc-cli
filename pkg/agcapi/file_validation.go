package agcapi

import (
	"os"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// File prerequisites use the same structured feedback in CLI, REST and API.
func ValidateEndpointFileInput(e domain.Endpoint, path string, headerSets ...map[string]string) error {
	if (e.FamilyID == "storage" && e.ID == "upload-file") || (e.FamilyID == "hosting" && e.ID == "upload-version") {
		var headers map[string]string
		if len(headerSets) > 0 {
			headers = headerSets[0]
		}
		return validateServerlessUploadFile(e, path, headers)
	}
	if e.FamilyID != "pms" || e.ID != "resourcemanage-harmonyosnext" {
		return nil
	}
	source := e.SourceURL
	for _, p := range e.Parameters {
		if p.In == "file" && p.Name == "file" && p.SourceURL != "" {
			source = p.SourceURL
		}
	}
	code, message := "file_required", "Multipart resource upload requires a local file."
	if path != "" {
		file, err := os.Open(path)
		if err == nil {
			info, statErr := file.Stat()
			file.Close()
			if statErr == nil && info.Mode().IsRegular() {
				return nil
			}
		}
		code, message = "invalid_file", "The selected upload file must be readable and regular."
	}
	return &RequestValidationError{Violations: []domain.InputViolation{{Code: code, Field: "filePath", Message: message, SourceURL: source, Suggestion: "Select a readable local image with --file PATH in CLI or filePath in REST; supply requestId, fileType=1, zoneId=1 and fileName separately."}}}
}
