package command

import (
	_ "embed"
	"fmt"

	"github.com/Createitv/agc-cli/pkg/domain"
	"github.com/Createitv/agc-cli/pkg/output"
	"github.com/spf13/cobra"
)

//go:embed skills/agc-cli.md
var agcSkillContent string

type skillSummary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
}

func skillsCommand(opts *options) *cobra.Command {
	summary := skillSummary{
		Name:        "agc-cli",
		Description: "Discover and invoke AppGallery Connect workflows with the agc CLI.",
		Command:     "agc skills show agc-cli",
	}
	cmd := &cobra.Command{Use: "skills", Short: "Read bundled agent skills offline"}
	var remote bool
	list := &cobra.Command{
		Use: "list", Short: "List bundled agent skills", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if remote {
				return runSkillsCLI(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), "--yes", "skills", "add", "Createitv/agc-cli", "--list")
			}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[[]skillSummary]{Data: []skillSummary{summary}}, output.Format(opts.output), opts.pretty)
		},
	}
	list.Flags().BoolVar(&remote, "remote", false, "List published skills using the online installer")
	cmd.AddCommand(list)
	installer := remoteSkillsCommand()
	install, _, _ := installer.Find([]string{"install"})
	installer.RemoveCommand(install)
	cmd.AddCommand(install)
	cmd.AddCommand(&cobra.Command{
		Use: "show <name>", Short: "Read a bundled agent skill", Args: cobra.ExactArgs(1),
		ValidArgs: []string{summary.Name},
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != summary.Name {
				return fmt.Errorf("unknown skill %q; run 'agc skills list' to list bundled skills", args[0])
			}
			if output.Format(opts.output) == output.Markdown {
				_, err := fmt.Fprint(cmd.OutOrStdout(), agcSkillContent)
				return err
			}
			detail := struct {
				skillSummary
				Content string `json:"content"`
			}{summary, agcSkillContent}
			return output.Write(cmd.OutOrStdout(), domain.Envelope[any]{Data: detail}, output.Format(opts.output), opts.pretty)
		},
	})
	return cmd
}
