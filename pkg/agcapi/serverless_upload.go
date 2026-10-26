package agcapi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func uploadHeader(headers map[string]string, name string) string {
	for k, v := range headers {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}
func uploadIssue(e domain.Endpoint, field, message string) error {
	return &RequestValidationError{Violations: []domain.InputViolation{{Code: "upload_prerequisite", Field: field, Message: message, SourceURL: e.SourceURL, Suggestion: "Use the returned signed headers or documented resumable size, offset and SHA256; regenerate upload authorization if file bytes changed. Read " + e.Command + " --docs --pretty."}}}
}
func validateServerlessUploadFile(e domain.Endpoint, path string, headers map[string]string) error {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return uploadIssue(e, "filePath", "A readable regular local file is required.")
	}
	f, err := os.Open(path)
	if err != nil {
		return uploadIssue(e, "filePath", "Cannot read the selected local file.")
	}
	defer f.Close()
	if e.FamilyID == "hosting" {
		digest := uploadHeader(headers, "x-amz-content-sha256")
		if headers == nil {
			return nil
		} // Legacy caller only checks file readability; invoke always supplies headers.
		if !regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(digest) {
			return uploadIssue(e, "header.x-amz-content-sha256", "Supply the signed file's 64-character SHA256.")
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return uploadIssue(e, "filePath", "Cannot hash the complete upload file.")
		}
		if !strings.EqualFold(digest, hex.EncodeToString(h.Sum(nil))) {
			return uploadIssue(e, "header.x-amz-content-sha256", "Signed SHA256 does not match the selected file; do not send mismatching signed bytes.")
		}
		if length := uploadHeader(headers, "Content-Length"); length != "" && length != strconv.FormatInt(info.Size(), 10) {
			return uploadIssue(e, "header.Content-Length", "Signed content length differs from local file length.")
		}
		return nil
	}
	offsetText, totalText := uploadHeader(headers, "X-Agc-File-Offset"), uploadHeader(headers, "X-Agc-File-Size")
	var offset int64
	if offsetText != "" {
		offset, err = strconv.ParseInt(offsetText, 10, 64)
		if err != nil || offset < 0 {
			return uploadIssue(e, "header.X-Agc-File-Offset", "Use a nonnegative integer byte offset.")
		}
	}
	total := info.Size()
	if totalText != "" {
		total, err = strconv.ParseInt(totalText, 10, 64)
		if err != nil || total < 0 {
			return uploadIssue(e, "header.X-Agc-File-Size", "Use a nonnegative integer total file size.")
		}
		if offset > total || info.Size() > total-offset {
			return uploadIssue(e, "header.X-Agc-File-Size", "Offset and uploaded chunk exceed the documented total file length.")
		}
	}
	if digest := uploadHeader(headers, "X-Agc-Sha256"); digest != "" {
		if !regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(digest) {
			return uploadIssue(e, "header.X-Agc-Sha256", "Use the complete source file's SHA256 in hexadecimal.")
		}
		if offset == 0 && total == info.Size() {
			h := sha256.New()
			if _, err = io.Copy(h, f); err != nil {
				return uploadIssue(e, "filePath", "Cannot hash the selected file.")
			}
			if !strings.EqualFold(digest, hex.EncodeToString(h.Sum(nil))) {
				return uploadIssue(e, "header.X-Agc-Sha256", "SHA256 differs from the selected complete file.")
			}
		}
	}
	return nil
}
func resumableViolations(e domain.Endpoint, headers map[string]string) []domain.InputViolation {
	var issues []domain.InputViolation
	if e.FamilyID != "storage" {
		return issues
	}
	add := func(field, msg string) {
		issues = append(issues, domain.InputViolation{Code: "resumable_prerequisite", Field: "header." + field, Message: msg, SourceURL: e.SourceURL, Suggestion: "Provide the complete file size and SHA256 when supplying a resumable file offset; see " + e.Command + " --docs --pretty."})
	}
	if len(headers) > 40 {
		add("headers", fmt.Sprintf("Explicit header count %d exceeds the documented limit of 40.", len(headers)))
	}
	if e.ID == "upload-file" && uploadHeader(headers, "X-Agc-File-Offset") != "" {
		for _, n := range []string{"X-Agc-File-Size", "X-Agc-Sha256"} {
			if uploadHeader(headers, n) == "" {
				add(n, n+" is required for resumable uploads.")
			}
		}
	}
	return issues
}
