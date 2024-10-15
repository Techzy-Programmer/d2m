package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/charmbracelet/glamour"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/urfave/cli/v2"
)

var LogsCmd = &cli.Command{
	Name:   "logs",
	Usage:  "List and check deployment logs",
	Action: handleLogs,
	Subcommands: []*cli.Command{
		removeLogsCmd,
	},
}

var removeLogsCmd = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "Removes application logs, not the deployment ones",
	Action:  handleRemoveLogs,
}

func handleRemoveLogs(c *cli.Context) error {
	paint.Error("Not implemented yet")
	return nil
}

func handleLogs(c *cli.Context) error {
	logs := db.GetAllDeployments()
	if len(logs) == 0 {
		paint.Info("No deployments triggered yet")
		return nil
	}

	var depPrompt = []string{}

	for _, log := range logs {
		commitMsg := log.CommitMsg
		commitHsh := "-------"

		if log.Status != "success" {
			commitMsg = "[No Commit Msg] Deployment is failing"
		}

		if len(log.CommitHash) > 7 {
			commitHsh = log.CommitHash[:7]
		}

		depPrompt = append(depPrompt, fmt.Sprintf("%d, [%s] %s", log.ID, commitHsh, commitMsg))
	}

	sp := selection.New("Select a deployment to view it's log", depPrompt)
	if len(depPrompt) > 10 {
		sp.PageSize = 10
	} else {
		sp.PageSize = len(depPrompt)
	}

	choice, err := sp.RunPrompt()
	if err != nil {
		paint.ErrorF("Selection error: %v", err)
		return err
	}

	idStr := strings.Split(choice, ",")[0]
	id, errParse := strconv.ParseUint(idStr, 10, 32)
	if errParse != nil {
		paint.ErrorF("Error parsing deployment id: %v", errParse)
		return errParse
	}

	dep, errDepGet := db.GetDeploymentByID(uint(id))
	if errDepGet != nil {
		paint.ErrorF("Error getting deployment from local-db: %v", errDepGet)
		return errDepGet
	}

	renderDeploymentLogs(dep)
	return nil
}

func renderDeploymentLogs(dep db.Deployment) {
	logs := db.GetAllLogsForDeployment(dep.ID)
	if len(logs) == 0 {
		paint.Info("No logs available for this deployment")
		return
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(100),
		glamour.WithAutoStyle(),
	)

	dep.Logs = logs
	markdown := GenerateMarkdown(dep)
	out, rnErr := renderer.Render(markdown)
	if rnErr != nil {
		paint.ErrorF("Error while rendering markdown: %v", rnErr)
		return
	}

	fmt.Println(out)
}

var logLevelMap = map[string]string{
	"0": "Info",
	"1": "Warn",
	"2": "Fail",
	"3": "Pass",
}

func GenerateMarkdown(deployment db.Deployment) string {
	var builder strings.Builder

	builder.WriteString("# Deployment Report\n")
	builder.WriteString(fmt.Sprintf("- **Deployment ID:** `%d`\n", deployment.ID))
	builder.WriteString(fmt.Sprintf("- **From Branch:** `%s`\n", deployment.Branch))
	builder.WriteString(fmt.Sprintf("- **Status:** `%s`\n", deployment.Status))
	builder.WriteString(fmt.Sprintf("- **Repository:** `%s`\n", deployment.Repo))
	builder.WriteString(fmt.Sprintf("- **Commit Message:** `%s`\n", strings.TrimSpace(deployment.CommitMsg)))
	builder.WriteString(fmt.Sprintf("- **Commit Hash:** `%s`\n", deployment.CommitHash))
	builder.WriteString(fmt.Sprintf("- **Start Time:** `%s`\n", deployment.StartAt.Format(time.DateTime)))
	builder.WriteString(fmt.Sprintf("- **End Time:** `%s`\n\n", deployment.EndAt.Format(time.DateTime)))

	builder.WriteString("# Deployment Logs\n")
	for _, log := range deployment.Logs {
		builder.WriteString(fmt.Sprintf("## _%s_: %s\n", logLevelMap[fmt.Sprintf("%d", log.Level)], log.Title))
		builder.WriteString(fmt.Sprintf("- **At:** %s\n", time.Unix(log.Timestamp, 0).Format(time.DateTime)))

		if len(log.Message) > 0 {
			builder.WriteString(fmt.Sprintf("- **Message:** %s\n", log.Message))
		}

		steps := [][]string{}
		json.Unmarshal([]byte(log.Steps), &steps)

		if len(steps) > 0 {
			builder.WriteString("### Steps:\n")
			for _, step := range steps {
				level, logEntry := step[0], step[1]
				builder.WriteString(fmt.Sprintf("- **%s**: %s\n", logLevelMap[level], logEntry))
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
