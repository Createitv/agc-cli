package agcapi

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
)

func invokePrefetchDevice(ctx context.Context, req InvokeRequest) (InvokeResponse, error) {
	out := InvokeResponse{Method: req.Endpoint.Method, URL: req.Endpoint.Path, DryRun: req.DryRun}
	var p struct {
		BundleName string `json:"bundleName"`
		DeviceID   string `json:"deviceId"`
		HDCPath    string `json:"hdcPath"`
	}
	if json.Unmarshal(req.Body, &p) != nil || p.BundleName == "" || p.DeviceID == "" {
		return out, uploadIssue(req.Endpoint, "body.bundleName", "Provide bundleName and an explicitly selected deviceId for cf_prefetch getcache.")
	}
	if req.DryRun {
		return out, nil
	}
	if p.HDCPath == "" {
		p.HDCPath = "hdc"
	}
	command := exec.CommandContext(ctx, p.HDCPath, "-t", p.DeviceID, "shell", "cf_prefetch", "getcache", "-m", p.BundleName)
	data, err := command.CombinedOutput()
	if len(data) > 8<<20 {
		return out, uploadIssue(req.Endpoint, "response.body", "Device debug output exceeds the 8 MiB local limit.")
	}
	for code, message := range map[string]string{"1008200005": "Periodic prefetch failed; check device connectivity, prefetch configuration and mainland China function region.", "1008200009": "Device prefetch tool internal error; contact Huawei support.", "1008240000": "Install a debug HAP signed with a debugging certificate.", "1008240007": "Check the application's registered bundle name and HAP information."} {
		if strings.Contains(string(data), code) {
			b, _ := json.Marshal(map[string]string{"code": code, "message": message})
			return out, ParseAGCError(200, b)
		}
	}
	if err != nil {
		return out, uploadIssue(req.Endpoint, "body.hdcPath", "hdc did not complete; verify the executable path, explicit deviceId, device connection and debug-HAP prerequisites.")
	}
	out.StatusCode = 200
	out.ContentType = "application/json"
	out.Body, _ = json.Marshal(map[string]any{"commandCompleted": true, "dataFreshnessVerified": false, "output": string(data), "note": "Compare fetch data timestamp with the current run; an old timestamp means the latest fetch failed."})
	out.RawBody = out.Body
	return out, nil
}
