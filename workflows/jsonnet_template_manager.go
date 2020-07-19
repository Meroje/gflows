package workflows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/jbrunton/gflows/adapters"
	"github.com/jbrunton/gflows/config"
	"github.com/spf13/afero"
)

type JsonnetTemplateManager struct {
	fs      *afero.Afero
	logger  *adapters.Logger
	context *config.GFlowsContext
}

func NewJsonnetTemplateManager(fs *afero.Afero, logger *adapters.Logger, context *config.GFlowsContext) *JsonnetTemplateManager {
	return &JsonnetTemplateManager{
		fs:      fs,
		logger:  logger,
		context: context,
	}
}

func (manager *JsonnetTemplateManager) GetWorkflowSources() []string {
	files := []string{}
	err := manager.fs.Walk(manager.context.WorkflowsDir, func(path string, f os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext == ".jsonnet" || ext == ".libsonnet" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func (manager *JsonnetTemplateManager) GetWorkflowTemplates() []string {
	sources := manager.GetWorkflowSources()
	var templates []string
	for _, source := range sources {
		if filepath.Ext(source) == ".jsonnet" {
			templates = append(templates, source)
		}
	}
	return templates
}

// GetWorkflowDefinitions - get workflow definitions for the given context
func (manager *JsonnetTemplateManager) GetWorkflowDefinitions() ([]*WorkflowDefinition, error) {
	templates := manager.GetWorkflowTemplates()
	definitions := []*WorkflowDefinition{}
	for _, templatePath := range templates {
		vm := createVM(manager.context)
		workflowName := manager.getWorkflowName(manager.context.WorkflowsDir, templatePath)
		input, err := manager.fs.ReadFile(templatePath)
		if err != nil {
			return []*WorkflowDefinition{}, err
		}

		destinationPath := filepath.Join(manager.context.GitHubDir, "workflows/", workflowName+".yml")
		definition := &WorkflowDefinition{
			Name:        workflowName,
			Source:      templatePath,
			Destination: destinationPath,
			Status:      ValidationResult{Valid: true},
		}

		workflow, err := vm.EvaluateSnippet(templatePath, string(input))
		if err != nil {
			definition.Status.Valid = false
			definition.Status.Errors = []string{strings.Trim(err.Error(), " \n\r")}
		} else {
			meta := strings.Join([]string{
				"# File generated by gflows, do not modify",
				fmt.Sprintf("# Source: %s", templatePath),
			}, "\n")
			definition.Content = meta + "\n" + workflow
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func (manager *JsonnetTemplateManager) getWorkflowName(workflowsDir string, filename string) string {
	_, templateFileName := filepath.Split(filename)
	return strings.TrimSuffix(templateFileName, filepath.Ext(templateFileName))
}

func createVM(context *config.GFlowsContext) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: context.EvalJPaths(),
	})
	vm.StringOutput = true
	return vm
}
