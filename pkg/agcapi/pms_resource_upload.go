package agcapi

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"sort"
)

// Spool multipart data to a temporary file so large uploads do not accumulate
// in memory. The caller owns closing and removing the returned file.
func preparePMSResourceUpload(path string, body []byte) (result *os.File, size int64, contentType string, err error) {
	input, err := os.Open(path)
	if err != nil {
		return nil, 0, "", fmt.Errorf("open upload file: %w", err)
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return nil, 0, "", err
	}
	if !info.Mode().IsRegular() {
		return nil, 0, "", fmt.Errorf("upload file must be a regular file")
	}
	var fields map[string]string
	if err = json.Unmarshal(body, &fields); err != nil || fields == nil {
		return nil, 0, "", fmt.Errorf("multipart resource fields must be a JSON object of strings")
	}
	result, err = os.CreateTemp("", "agc-resource-upload-*")
	if err != nil {
		return nil, 0, "", err
	}
	defer func() {
		if err != nil {
			result.Close()
			os.Remove(result.Name())
		}
	}()
	writer := multipart.NewWriter(result)
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err = writer.WriteField(key, fields[key]); err != nil {
			return
		}
	}
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{"name": "file", "filename": fields["fileName"]}))
	fileType := mime.TypeByExtension(filepath.Ext(path))
	if fileType == "" {
		fileType = "application/octet-stream"
	}
	header.Set("Content-Type", fileType)
	var part io.Writer
	if part, err = writer.CreatePart(header); err != nil {
		return
	}
	if _, err = io.Copy(part, input); err != nil {
		return
	}
	if err = writer.Close(); err != nil {
		return
	}
	contentType = writer.FormDataContentType()
	if size, err = result.Seek(0, io.SeekEnd); err != nil {
		return
	}
	_, err = result.Seek(0, io.SeekStart)
	return
}
