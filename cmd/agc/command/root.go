package command

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	var clientID string
	var clientKey string
	var name string
	var credentialsPath string
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
			if err := agcapi.ValidateCredential(credential); err != nil {
				return err
			}
			path, err := credentialPath(credentialsPath)
			if err != nil {
				return err
			}
			if err := agcapi.SaveCredential(path, credential); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.Credential]{Data: credential}, output.Format(opts.output), opts.pretty)
		},
	}
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
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]agcapi.Credential]{Data: store.Accounts}, output.Format(opts.output), opts.pretty)
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
			return output.Write(cmd.OutOrStdout(), domain.Envelope[agcapi.Credential]{Data: credential}, output.Format(opts.output), opts.pretty)
		},
	}
	check.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")
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
			srv := &http.Server{Addr: addr, Handler: server.Handler(), ReadHeaderTimeout: 5 * time.Second}
			cmd.Printf("agc web-server listening on http://localhost%s\n", addr)
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
	cmd.Flags().StringVar(&addr, "addr", ":8421", "Listen address")
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
					if headerMap[parameter.Name] == "" {
						return fmt.Errorf("missing --header %s=value", parameter.Name)
					}
				case "body":
					if bodyFile == "" && fieldMap[parameter.Name] == "" {
						return fmt.Errorf("missing --field %s=value or --body", parameter.Name)
					}
				case "file":
					if paramMap[parameter.Name] == "" && fieldMap[parameter.Name] == "" && bodyFile == "" {
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
					body, err = json.Marshal(fieldMap)
					if err != nil {
						return err
					}
				}
			}
			if token == "" {
				token = os.Getenv("AGC_ACCESS_TOKEN")
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()
			if token == "" && !dryRun {
				token, err = endpointAccessToken(ctx, opts, credentialsPath, baseURL)
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
	cmd.Flags().StringVar(&baseURL, "base-url", "https://connect-api.cloud.huawei.com", "Connect API base URL")
	cmd.Flags().StringVar(&token, "token", "", "Bearer token; defaults to AGC_ACCESS_TOKEN")
	cmd.Flags().StringVar(&credentialsPath, "credentials-path", "", "Override credentials file path")
	cmd.Flags().StringVar(&outFile, "out", "", "Write raw response body to a file")
	cmd.Flags().BoolVar(&invoke, "invoke", false, "Invoke the registered endpoint")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Build the request without sending it")
	return cmd
}

func endpointAccessToken(ctx context.Context, opts *options, credentialsPathOverride, baseURL string) (string, error) {
	path, err := credentialPath(credentialsPathOverride)
	if err != nil {
		return "", err
	}
	store, err := agcapi.LoadCredentials(path)
	if err != nil {
		return "", err
	}
	credential, ok, err := resolveCredential(opts, store)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("no active credential profile; pass --token, set AGC_ACCESS_TOKEN, or run agc auth login")
	}
	token, err := agcapi.AccessToken(ctx, nil, baseURL, credential)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
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
		credential, ok := agcapi.ActiveCredential(store)
		return credential, ok, nil
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
