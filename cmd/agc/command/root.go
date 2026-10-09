package command

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
	"github.com/Createitv/agc-cli/pkg/output"
	"github.com/Createitv/agc-cli/pkg/project"
	"github.com/Createitv/agc-cli/pkg/server"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

type options struct {
	output  string
	pretty  bool
	timeout time.Duration
	profile string
	project string
}

func NewRootCommand() *cobra.Command {
	opts := &options{}
	cmd := &cobra.Command{
		Use:   "agc",
		Short: "AppGallery Connect command center",
		Long:  "agc automates AppGallery Connect workflows from the terminal, CI, local REST API, and web command center.",
	}
	cmd.PersistentFlags().StringVar(&opts.output, "output", "json", "Output format: json, table, markdown")
	cmd.PersistentFlags().BoolVar(&opts.pretty, "pretty", false, "Pretty-print JSON output")
	cmd.PersistentFlags().DurationVar(&opts.timeout, "timeout", 60*time.Second, "Request timeout")
	cmd.PersistentFlags().StringVar(&opts.profile, "profile", "", "Credential profile name")
	cmd.PersistentFlags().StringVar(&opts.project, "project", ".", "Project directory")

	cmd.AddCommand(versionCommand(opts))
	cmd.AddCommand(capabilitiesCommand(opts))
	cmd.AddCommand(endpointsCommand(opts))
	cmd.AddCommand(openAPICommand(opts))
	cmd.AddCommand(authCommand(opts))
	cmd.AddCommand(initCommand(opts))
	cmd.AddCommand(webServerCommand())
	cmd.AddCommand(docsCommand(opts))
	cmd.AddCommand(skillsCommand(opts))
	for _, capability := range domain.DecoratedCapabilities() {
		cmd.AddCommand(moduleCommand(opts, capability))
	}
	return cmd
}

func versionCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), map[string]string{"version": version, "commit": commit, "buildDate": buildDate}, output.Format(opts.output), opts.pretty)
		},
	}
}

func capabilitiesCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "List AppGallery Connect API families supported by agc-cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]domain.Capability]{
				Data: domain.DecoratedCapabilities(),
				Affordances: map[string]string{
					"serve": "agc web-server",
				},
			}, output.Format(opts.output), opts.pretty)
		},
	}
}

func endpointsCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "endpoints",
		Short: "List every registered AppGallery Connect endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]domain.Endpoint]{
				Data: domain.AllEndpoints(),
				Affordances: map[string]string{
					"capabilities": "agc capabilities",
				},
			}, output.Format(opts.output), opts.pretty)
		},
	}
}

func openAPICommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Export the local REST and endpoint invocation OpenAPI contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.OpenAPISpec(), output.Format(opts.output), opts.pretty)
		},
	}
}

func authCommand(opts *options) *cobra.Command {
	var serviceAccountFile string
	var apiClientFile string
	var activate bool
	var clientID string
	var clientKey string
	var name string
	var credentialsPath string
	var remote bool
	var checkBaseURL string
	cmd := &cobra.Command{Use: "auth", Short: "Manage AppGallery Connect credentials"}
	login := &cobra.Command{
		Use:   "login",
		Short: "Save a service account or API client credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				name = "default"
			}
			mode := "service-account"
			if clientID != "" || clientKey != "" {
				mode = "api-client"
			}
			credential := agcapi.Credential{Name: name, Mode: mode, ServiceAccountFile: serviceAccountFile, ClientID: clientID, ClientKey: clientKey}
			if apiClientFile != "" {
				if cmd.Flags().Changed("client-id") || cmd.Flags().Changed("client-key") || cmd.Flags().Changed("service-account-file") {
					return fmt.Errorf("--api-client-file cannot be combined with --client-id, --client-key, or --service-account-file")
				}
				imported, err := agcapi.LoadAPIClientCredential(apiClientFile)
				if err != nil {
					return err
				}
				imported.Name = name
				credential = imported
			}
			if err := agcapi.ValidateCredential(credential); err != nil {
				return err
			}
			path, err := credentialPath(credentialsPath)
			if err != nil {
				return err
			}
			if err := agcapi.SaveCredentialWithActivation(path, credential, activate); err != nil {
				return err
			}
			saved, err := agcapi.LoadCredentials(path)
			if err != nil {
				return err
			}
			credential, _ = agcapi.CredentialByName(saved, name)
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.CredentialView]{Data: credential.View()}, output.Format(opts.output), opts.pretty)
		},
	}
	login.Flags().StringVar(&apiClientFile, "api-client-file", "", "API client JSON file with client_id/client_secret or clientId/clientKey")
	login.Flags().BoolVar(&activate, "activate", true, "Make this profile active; false preserves the current active profile")
	login.Flags().StringVar(&serviceAccountFile, "service-account-file", "", "Path to Huawei service account JSON")
	login.Flags().StringVar(&clientID, "client-id", "", "API client ID")
	login.Flags().StringVar(&clientKey, "client-key", "", "API client key")
	login.Flags().StringVar(&name, "name", "default", "Credential profile name")
	login.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")

	list := &cobra.Command{
		Use:   "list",
		Short: "List saved credential profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := credentialPath(credentialsPath)
			if err != nil {
				return err
			}
			store, err := agcapi.LoadCredentials(path)
			if err != nil {
				return err
			}
			views := make([]agcapi.CredentialView, 0, len(store.Accounts))
			for _, credential := range store.Accounts {
				views = append(views, credential.View())
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]agcapi.CredentialView]{Data: views}, output.Format(opts.output), opts.pretty)
		},
	}
	list.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")

	check := &cobra.Command{
		Use:   "check",
		Short: "Show active credential profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := credentialPath(credentialsPath)
			if err != nil {
				return err
			}
			store, err := agcapi.LoadCredentials(path)
			if err != nil {
				return err
			}
			credential, ok, err := resolveCredential(opts, store)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("no active credential profile")
			}
			if remote {
				ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
				defer cancel()
				token, err := agcapi.AccessToken(ctx, nil, checkBaseURL, credential)
				if err != nil {
					return fmt.Errorf("remote token generation: %w", err)
				}
				config, err := project.Load(opts.project)
				if err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("load project context: %w", err)
				}
				projectID := config.ProjectID
				if projectID == "" && credential.Mode == "service-account" {
					account, err := agcapi.LoadServiceAccount(credential.ServiceAccountFile)
					if err != nil {
						return err
					}
					projectID = account.ProjectID
				}
				if credential.Mode == "api-client" && config.AppID == "" {
					return output.Write(cmd.OutOrStdout(), domain.Envelope[map[string]any]{Data: map[string]any{"credential": credential.View(), "tokenVerified": true, "readVerified": false, "remoteVerified": false, "message": "Token verified; run agc init --app-id to verify app read permission."}}, output.Format(opts.output), opts.pretty)
				}
				endpointID := "queryprojectlist"
				family := "projects"
				query := map[string]string{}
				params := map[string]string{}
				if projectID != "" {
					endpointID = "queryprojectdetail"
					params["projectId"] = projectID
				}
				if credential.Mode == "api-client" {
					endpointID = "app-info-query"
					family = "publishing"
					params = map[string]string{}
					query["appId"] = config.AppID
					query["lang"] = "en-US"
				}
				var endpoint domain.Endpoint
				for _, candidate := range domain.EndpointsByFamily(family) {
					if candidate.ID == endpointID {
						endpoint = candidate
						break
					}
				}
				headers := map[string]string{}
				if credential.ClientID != "" {
					headers["client_id"] = credential.ClientID
				}
				response, err := agcapi.InvokeEndpoint(ctx, nil, agcapi.InvokeRequest{Endpoint: endpoint, BaseURL: checkBaseURL, Params: params, Query: query, Headers: headers, AccessToken: token.AccessToken})
				if err != nil {
					return fmt.Errorf("remote authorization check: %w", err)
				}
				return output.Write(cmd.OutOrStdout(), domain.Envelope[map[string]any]{Data: map[string]any{"credential": credential.View(), "remoteVerified": true, "tokenVerified": true, "readVerified": true, "endpoint": endpointID, "statusCode": response.StatusCode}}, output.Format(opts.output), opts.pretty)
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.CredentialView]{Data: credential.View()}, output.Format(opts.output), opts.pretty)
		},
	}
	check.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")
	check.Flags().BoolVar(&remote, "remote", false, "Verify token and available app or project read permission")
	check.Flags().StringVar(&checkBaseURL, "base-url", "https://connect-api.cloud.huawei.com", "Connect API base URL for remote verification")
	token := &cobra.Command{
		Use:   "token",
		Short: "Create an AppGallery Connect authorization token for the active credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := credentialPath(credentialsPath)
			if err != nil {
				return err
			}
			store, err := agcapi.LoadCredentials(path)
			if err != nil {
				return err
			}
			credential, ok, err := resolveCredential(opts, store)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("no active credential profile")
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()
			token, err := agcapi.AccessToken(ctx, nil, "", credential)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.TokenResponse]{Data: token}, output.Format(opts.output), opts.pretty)
		},
	}
	token.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")
	cmd.AddCommand(login, list, check, token)
	return cmd
}

func initCommand(opts *options) *cobra.Command {
	var appID string
	var projectID string
	var packageName string
	var profile string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Pin an AppGallery Connect app context to .agc/project.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if appID == "" {
				return fmt.Errorf("--app-id is required")
			}
			config := project.Config{AppID: appID, ProjectID: projectID, PackageName: packageName, Profile: profile}
			if err := project.Save(opts.project, config); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[project.Config]{Data: config}, output.Format(opts.output), opts.pretty)
		},
	}
	cmd.Flags().StringVar(&appID, "app-id", "", "AppGallery Connect app ID")
	cmd.Flags().StringVar(&projectID, "project-id", "", "AppGallery Connect project ID")
	cmd.Flags().StringVar(&packageName, "package-name", "", "HarmonyOS bundle/package name")
	cmd.Flags().StringVar(&profile, "default-profile", "", "Default credential profile for this project")
	return cmd
}

func webServerCommand() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:   "web-server",
		Short: "Start the local REST API for Command Center and agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverToken := os.Getenv("AGC_SERVER_TOKEN")
			if !loopbackListenAddress(addr) && serverToken == "" {
				return fmt.Errorf("non-loopback listening requires AGC_SERVER_TOKEN; set it or use --addr 127.0.0.1:8421")
			}
			srv := &http.Server{Addr: addr, Handler: server.HandlerWithToken(serverToken), ReadHeaderTimeout: 5 * time.Second}
			cmd.Printf("agc web-server listening on http://%s\n", addr)
			go func() {
				<-cmd.Context().Done()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_ = srv.Shutdown(ctx)
			}()
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&addr, "addr", "127.0.0.1:8421", "Listen address")
	return cmd
}

//go:embed assets/agc-cli/SKILL.md
var agcCLISkill []byte

func skillsCommand(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Install agent skills",
	}
	var agent string
	var force bool
	add := &cobra.Command{
		Use:   "add",
		Short: "Add the agc-cli skill to an agent",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			agentDirs := map[string]string{
				"copilot": ".github/skills",
				"claude":  ".claude/skills",
				"codex":   ".agents/skills",
			}
			agents := []string{agent}
			if agent == "all" {
				agents = []string{"copilot", "claude", "codex"}
			}
			for _, name := range agents {
				if _, ok := agentDirs[name]; !ok {
					return fmt.Errorf("unsupported agent %q; choose copilot, claude, codex, or all", name)
				}
			}

			type skillInstall struct {
				Agent string `json:"agent"`
				Path  string `json:"path"`
			}
			paths := make([]string, 0, len(agents))
			for _, name := range agents {
				path := filepath.Join(agentDirs[name], "agc-cli", "SKILL.md")
				info, err := os.Lstat(filepath.Join(opts.project, path))
				if err == nil && !force {
					return fmt.Errorf("skill already exists at %s; use --force to replace it", path)
				} else if err == nil && info.IsDir() {
					return fmt.Errorf("skill destination %s is a directory", path)
				} else if err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("check skill destination %s: %w", path, err)
				}
				paths = append(paths, path)
			}

			installed := make([]skillInstall, 0, len(agents))
			for i, path := range paths {
				destination := filepath.Join(opts.project, path)
				if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
					return fmt.Errorf("create skill directory %s: %w", filepath.Dir(destination), err)
				}
				if force {
					if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
						return fmt.Errorf("replace skill at %s: %w", path, err)
					}
				}
				file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
				if err != nil {
					return fmt.Errorf("install skill at %s: %w", path, err)
				}
				_, writeErr := file.Write(agcCLISkill)
				closeErr := file.Close()
				if writeErr != nil {
					_ = os.Remove(destination)
					return fmt.Errorf("write skill at %s: %w", path, writeErr)
				}
				if closeErr != nil {
					_ = os.Remove(destination)
					return fmt.Errorf("close skill at %s: %w", path, closeErr)
				}
				installed = append(installed, skillInstall{Agent: agents[i], Path: filepath.ToSlash(path)})
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]skillInstall]{Data: installed}, output.Format(opts.output), opts.pretty)
		},
	}
	add.Flags().StringVar(&agent, "agent", "", "Target agent: copilot, claude, codex, or all")
	add.Flags().BoolVar(&force, "force", false, "Replace an existing agc-cli skill")
	_ = add.MarkFlagRequired("agent")
	cmd.AddCommand(add)
	return cmd
}

func docsCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "docs [module]",
		Short: "Print the local documentation path for a module",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			module := "AGC_CLI_FULL_PLAN"
			if len(args) == 1 {
				module = args[0]
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[map[string]string]{Data: map[string]string{"path": "docs/features/" + module + ".md"}}, output.Format(opts.output), opts.pretty)
		},
	}
}

func moduleCommand(opts *options, capability domain.Capability) *cobra.Command {
	cmd := &cobra.Command{
		Use:   capability.ID,
		Short: capability.Name,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List scaffolded operations for " + capability.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[domain.Capability]{
				Data: capability,
				Warnings: []string{
					"Endpoint adapters for this API family must be implemented from the official Huawei Connect API reference before production use.",
				},
			}, output.Format(opts.output), opts.pretty)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "endpoints",
		Short: "List registered endpoints for " + capability.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]domain.Endpoint]{
				Data: domain.EndpointsByFamily(capability.ID),
				Affordances: map[string]string{
					"family": capability.Command + " list",
				},
			}, output.Format(opts.output), opts.pretty)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show implementation status for " + capability.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.Write(cmd.OutOrStdout(), domain.Envelope[map[string]string]{
				Data: map[string]string{"id": capability.ID, "status": capability.Status, "restPath": capability.RESTPath, "endpoints": fmt.Sprint(capability.EndpointCount)},
			}, output.Format(opts.output), opts.pretty)
		},
	})
	for _, endpoint := range domain.EndpointsByFamily(capability.ID) {
		cmd.AddCommand(endpointCommand(opts, endpoint))
	}
	return cmd
}

func endpointCommand(opts *options, endpoint domain.Endpoint) *cobra.Command {
	var params []string
	var query []string
	var headers []string
	var fields []string
	var bodyFile string
	var uploadFile string
	var baseURL string
	var invoke bool
	var dryRun bool
	var token string
	var credentialsPath string
	var outFile string
	cmd := &cobra.Command{
		Use:   endpoint.ID,
		Short: endpoint.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !invoke {
				return output.Write(cmd.OutOrStdout(), domain.Envelope[domain.Endpoint]{Data: endpoint}, output.Format(opts.output), opts.pretty)
			}
			paramMap, err := parsePairs(params)
			if err != nil {
				return err
			}
			queryMap, err := parsePairs(query)
			if err != nil {
				return err
			}
			fieldMap, err := parsePairs(fields)
			if err != nil {
				return err
			}
			headerMap, err := parsePairs(headers)
			if err != nil {
				return err
			}
			config, configErr := project.Load(opts.project)
			if configErr != nil && !os.IsNotExist(configErr) {
				return fmt.Errorf("load project context: %w", configErr)
			}
			for _, parameter := range endpoint.Parameters {
				value := ""
				switch parameter.Name {
				case "appId", "appID":
					value = config.AppID
				case "projectId":
					value = config.ProjectID
				}
				if value != "" {
					switch parameter.In {
					case "path":
						if paramMap[parameter.Name] == "" {
							paramMap[parameter.Name] = value
						}
					case "query":
						if queryMap[parameter.Name] == "" {
							queryMap[parameter.Name] = value
						}
					case "header":
						if !hasHeader(headerMap, parameter.Name) {
							headerMap[parameter.Name] = value
						}
					case "body":
						if bodyFile == "" && fieldMap[parameter.Name] == "" {
							fieldMap[parameter.Name] = value
						}
					}
				}
			}
			var bodyFields map[string]json.RawMessage
			if bodyFile != "" {
				bodyData, err := os.ReadFile(bodyFile)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(bodyData, &bodyFields); err != nil {
					return fmt.Errorf("body must be a JSON object: %w", err)
				}
				if bodyFields == nil {
					return fmt.Errorf("body must be a JSON object")
				}
			}
			for _, parameter := range endpoint.Parameters {
				if !parameter.Required {
					continue
				}
				switch parameter.In {
				case "path":
					if paramMap[parameter.Name] == "" {
						return fmt.Errorf("missing --param %s=value", parameter.Name)
					}
				case "query":
					if queryMap[parameter.Name] == "" {
						return fmt.Errorf("missing --query %s=value", parameter.Name)
					}
				case "header":
					if !hasHeaderValue(headerMap, parameter.Name) {
						return fmt.Errorf("missing --header %s=value", parameter.Name)
					}
				case "body":
					if (bodyFile == "" && fieldMap[parameter.Name] == "") || (bodyFile != "" && (len(bodyFields[parameter.Name]) == 0 || string(bodyFields[parameter.Name]) == "null" || string(bodyFields[parameter.Name]) == `""`)) {
						return fmt.Errorf("missing --field %s=value or --body", parameter.Name)
					}
				case "file":
					if paramMap[parameter.Name] == "" && fieldMap[parameter.Name] == "" && bodyFile == "" && uploadFile == "" {
						return fmt.Errorf("missing --param %s=value, --field %s=value, or --body", parameter.Name, parameter.Name)
					}
				}
			}
			var body []byte
			if bodyFile != "" {
				body, err = os.ReadFile(bodyFile)
				if err != nil {
					return err
				}
			} else {
				if len(fieldMap) > 0 {
					body, err = agcapi.MarshalEndpointFields(endpoint, fieldMap)
					if err != nil {
						return err
					}
				}
			}
			if err := agcapi.ValidateEndpointBody(endpoint, body); err != nil {
				return err
			}
			explicitAuthHeader := hasHeaderValue(headerMap, "Authorization") || hasHeaderValue(headerMap, "oauth2Token")
			if token == "" && !explicitAuthHeader {
				token = os.Getenv("AGC_ACCESS_TOKEN")
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()
			if token == "" && !explicitAuthHeader && !dryRun && endpoint.Path != "{uploadUrl}" {
				var clientID string
				token, clientID, err = endpointAuthorization(ctx, opts, credentialsPath, baseURL)
				if clientID != "" && !hasHeader(headerMap, "client_id") {
					headerMap["client_id"] = clientID
				}
				if err != nil {
					return err
				}
			}
			resp, err := agcapi.InvokeEndpoint(ctx, nil, agcapi.InvokeRequest{
				Endpoint:    endpoint,
				BaseURL:     baseURL,
				Params:      paramMap,
				Query:       queryMap,
				Headers:     headerMap,
				Body:        body,
				FilePath:    uploadFile,
				AccessToken: token,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
			}
			if outFile != "" && len(resp.RawBody) > 0 {
				if err := os.WriteFile(outFile, resp.RawBody, 0644); err != nil {
					return err
				}
				resp.Body = nil
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.InvokeResponse]{Data: resp}, output.Format(opts.output), opts.pretty)
		},
	}
	cmd.Flags().StringArrayVar(&params, "param", nil, "Path parameter as key=value; repeatable")
	cmd.Flags().StringArrayVar(&query, "query", nil, "Query parameter as key=value; repeatable")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "HTTP header as key=value; repeatable")
	cmd.Flags().StringArrayVar(&fields, "field", nil, "JSON body field as key=value; repeatable")
	cmd.Flags().StringVar(&bodyFile, "body", "", "JSON body file")
	cmd.Flags().StringVar(&uploadFile, "file", "", "Raw file for signed upload endpoints")
	cmd.Flags().StringVar(&baseURL, "base-url", "https://connect-api.cloud.huawei.com", "Connect API base URL")
	cmd.Flags().StringVar(&token, "token", "", "Bearer token; defaults to AGC_ACCESS_TOKEN")
	cmd.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")
	cmd.Flags().StringVar(&outFile, "out", "", "Write raw response body to a file")
	cmd.Flags().BoolVar(&invoke, "invoke", false, "Invoke the registered endpoint")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Build the request without sending it")
	return cmd
}

func endpointAuthorization(ctx context.Context, opts *options, credentialsPathOverride, baseURL string) (string, string, error) {
	path, err := credentialPath(credentialsPathOverride)
	if err != nil {
		return "", "", err
	}
	store, err := agcapi.LoadCredentials(path)
	if err != nil {
		return "", "", err
	}
	credential, ok, err := resolveCredential(opts, store)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", fmt.Errorf("no active credential profile; pass --token, set AGC_ACCESS_TOKEN, or run agc auth login")
	}
	token, err := agcapi.AccessToken(ctx, nil, baseURL, credential)
	if err != nil {
		return "", "", err
	}
	return token.AccessToken, credential.ClientID, nil
}

func resolveCredential(opts *options, store agcapi.CredentialStore) (agcapi.Credential, bool, error) {
	profileName := opts.profile
	if profileName == "" {
		config, err := project.Load(opts.project)
		if err == nil {
			profileName = config.Profile
		} else if !os.IsNotExist(err) {
			return agcapi.Credential{}, false, fmt.Errorf("load project profile: %w", err)
		}
	}
	if profileName == "" {
		profileName = os.Getenv("AGC_PROFILE")
	}
	if profileName == "" {
		credential, ok := agcapi.ActiveCredential(store)
		if ok {
			return credential, true, nil
		}
		return agcapi.CredentialFromEnvironment()
	}
	credential, ok := agcapi.CredentialByName(store, profileName)
	if !ok {
		return agcapi.Credential{}, false, fmt.Errorf("credential profile %q not found", profileName)
	}
	return credential, true, nil
}

func parsePairs(items []string) (map[string]string, error) {
	out := map[string]string{}
	for _, item := range items {
		key, value, ok := strings.Cut(item, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid key=value pair %q", item)
		}
		out[key] = value
	}
	return out, nil
}

func credentialPath(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	if path := os.Getenv("AGC_CREDENTIALS_PATH"); path != "" {
		return path, nil
	}
	return agcapi.CredentialsPath()
}

func hasHeader(headers map[string]string, name string) bool {
	for key := range headers {
		if strings.EqualFold(key, name) {
			return true
		}
	}
	return false
}

func hasHeaderValue(headers map[string]string, name string) bool {
	for key, value := range headers {
		if strings.EqualFold(key, name) && value != "" {
			return true
		}
	}
	return false
}

func loopbackListenAddress(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
