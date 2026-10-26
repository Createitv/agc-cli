package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Createitv/agc-cli/pkg/cloudfunctions"
	"github.com/Createitv/agc-cli/pkg/domain"
	"github.com/Createitv/agc-cli/pkg/output"
	"github.com/Createitv/agc-cli/pkg/project"
	"github.com/spf13/cobra"
)

func functionEndpointCommand(opts *options, e domain.Endpoint) *cobra.Command {
	r := domain.FunctionRequest{DryRun: true}
	var invoke, docs bool
	var payload, payloadFile, environment, environmentFile string
	var memory, timeout, concurrency int
	var handler string
	var layers []string
	cmd := &cobra.Command{Use: e.ID, Short: e.Name, RunE: func(cmd *cobra.Command, args []string) (runErr error) {
		defer func() {
			if runErr != nil && opts.output == "json" {
				cmd.Root().SilenceErrors = true
				cmd.Root().SilenceUsage = true
				_ = output.Write(cmd.OutOrStdout(), domain.ErrorEnvelope{Error: cloudfunctions.Feedback(e.ID, runErr)}, output.JSON, opts.pretty)
			}
		}()
		if docs {
			raw, ok := domain.EndpointReference(e.FamilyID, e.ID)
			if !ok {
				return fmt.Errorf("source document missing")
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[json.RawMessage]{Data: raw}, output.Format(opts.output), opts.pretty)
		}
		if !invoke {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[domain.Endpoint]{Data: e}, output.Format(opts.output), opts.pretty)
		}
		if payload != "" && payloadFile != "" {
			return fmt.Errorf("payload and payload-file are mutually exclusive")
		}
		if environment != "" && environmentFile != "" {
			return fmt.Errorf("environment and environment-file are mutually exclusive")
		}
		if payloadFile != "" {
			b, err := os.ReadFile(payloadFile)
			if err != nil {
				return fmt.Errorf("cannot read payload-file")
			}
			r.Payload = b
		} else if payload != "" {
			r.Payload = json.RawMessage(payload)
		}
		if environmentFile != "" {
			b, err := os.ReadFile(environmentFile)
			if err != nil {
				return fmt.Errorf("cannot read environment-file")
			}
			environment = string(b)
		}
		if environment != "" {
			var env map[string]string
			if err := json.Unmarshal([]byte(environment), &env); err != nil || env == nil {
				return fmt.Errorf("environment must be a JSON object of string values")
			}
			r.Environment = &env
		}
		if cmd.Flags().Changed("memory") {
			r.Memory = &memory
		}
		if cmd.Flags().Changed("timeout") {
			r.Timeout = &timeout
		}
		if cmd.Flags().Changed("instance-concurrency") {
			r.InstanceConcurrency = &concurrency
		}
		if cmd.Flags().Changed("handler") {
			r.Handler = &handler
		}
		if cmd.Flags().Changed("layers") {
			r.Layers = &layers
		}
		if r.ProjectID == "" {
			config, err := project.Load(opts.project)
			if err == nil {
				r.ProjectID = config.ProjectID
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("cannot load project context")
			}
		}
		result, err := (&cloudfunctions.Client{}).Execute(cmd.Context(), e.ID, r)
		if err != nil {
			return err
		}
		return output.Write(cmd.OutOrStdout(), domain.Envelope[cloudfunctions.Response]{Data: result}, output.Format(opts.output), opts.pretty)
	}}
	f := cmd.Flags()
	f.BoolVar(&invoke, "invoke", false, "Run this native workflow")
	f.BoolVar(&docs, "docs", false, "Read official vendor source evidence and limitations")
	f.BoolVar(&r.DryRun, "dry-run", true, "Validate locally without reading authorization or contacting Huawei")
	f.StringVar(&r.Name, "name", "", "Function name")
	f.StringVar(&r.SessionFile, "session-file", "", "Explicit user OAuth session JSON (never publishing credentials)")
	f.StringVar(&r.BaseURL, "base-url", "", "Project service base URL, or session projectBaseUrl")
	f.StringVar(&r.ProjectID, "project-id", "", "Selected project; must match session productId")
	f.StringVar(&r.CodePath, "code-path", "", "Local ZIP deployment package")
	f.StringVar(&r.Runtime, "runtime", "", "Runtime, defaults to reviewed AGCLI nodejs14.x contract")
	f.StringVar(&handler, "handler", "", "Function entry handler")
	f.StringVar(&r.Description, "description", "", "Function description")
	f.StringVar(&r.VersionDescription, "version-description", "", "Required publish description")
	f.StringVar(&r.Version, "version", "", "Function version or $latest")
	f.StringVar(&r.AliasName, "alias-name", "", "Invoke alias; mutually exclusive with version")
	f.IntVar(&memory, "memory", 0, "Memory in MiB")
	f.IntVar(&timeout, "timeout", 0, "Function timeout in seconds")
	f.IntVar(&concurrency, "instance-concurrency", 0, "Concurrent instances")
	f.StringVar(&environment, "environment", "", "Environment JSON object; values are never echoed")
	f.StringVar(&environmentFile, "environment-file", "", "Environment JSON file")
	f.StringSliceVar(&layers, "layers", nil, "Layer names or URNs with versions, max five")
	f.StringVar(&payload, "payload", "", "Arbitrary valid JSON invocation event")
	f.StringVar(&payloadFile, "payload-file", "", "JSON invocation event file")
	f.StringVar(&r.FuzzyName, "fuzzy-name", "", "Function search filter")
	f.IntVar(&r.PageSize, "page-size", 0, "Page size, defaults 10")
	f.IntVar(&r.PageIndex, "page-index", 0, "Page index, defaults 1")
	return cmd
}
