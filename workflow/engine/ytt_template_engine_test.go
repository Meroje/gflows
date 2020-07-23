package engine

import (
	"strings"
	"testing"

	"github.com/jbrunton/gflows/config"
	"github.com/jbrunton/gflows/content"
	"github.com/jbrunton/gflows/fixtures"
	"github.com/jbrunton/gflows/workflow"
	"github.com/stretchr/testify/assert"
)

func newYttTemplateEngine(config string) (*content.Container, *config.GFlowsContext, *YttTemplateEngine) {
	adaptersContainer, context, _ := fixtures.NewTestContext(config)
	container := content.NewContainer(adaptersContainer)
	templateEngine := NewYttTemplateEngine(container.FileSystem(), container.Logger(), context, container.ContentWriter())
	return container, context, templateEngine
}

func TestGenerateYttWorkflowDefinitions(t *testing.T) {
	container, _, templateEngine := newYttTemplateEngine("")
	fs := container.FileSystem()
	fs.WriteFile(".gflows/workflows/test/config.yml", []byte(""), 0644)

	definitions, _ := templateEngine.GetWorkflowDefinitions()

	expectedContent := "# File generated by gflows, do not modify\n# Source: .gflows/workflows/test\n"
	expectedJson, _ := workflow.YamlToJson(expectedContent)
	expectedDefinition := workflow.Definition{
		Name:        "test",
		Source:      ".gflows/workflows/test",
		Destination: ".github/workflows/test.yml",
		Content:     expectedContent,
		Status:      workflow.ValidationResult{Valid: true},
		JSON:        expectedJson,
	}
	assert.Equal(t, []*workflow.Definition{&expectedDefinition}, definitions)
}

func TestGetYttWorkflowSources(t *testing.T) {
	container, _, templateEngine := newYttTemplateEngine("")
	fs := container.FileSystem()
	fs.WriteFile(".gflows/workflows/my-workflow/config1.yml", []byte("config1"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/config2.yaml", []byte("config2"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/config3.txt", []byte("config3"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/invalid.ext", []byte("ignored"), 0644)
	fs.WriteFile(".gflows/workflows/invalid-dir.yml", []byte("ignored"), 0644)

	sources := templateEngine.GetWorkflowSources()

	assert.Equal(t, []string{".gflows/workflows/my-workflow/config1.yml", ".gflows/workflows/my-workflow/config2.yaml", ".gflows/workflows/my-workflow/config3.txt"}, sources)
}

func TestGetYttWorkflowTemplates(t *testing.T) {
	container, _, templateEngine := newYttTemplateEngine("")
	fs := container.FileSystem()
	fs.WriteFile(".gflows/workflows/my-workflow/config1.yml", []byte("config1"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/nested-dir/config2.yaml", []byte("config2"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/config3.txt", []byte("config3"), 0644)
	fs.WriteFile(".gflows/workflows/my-workflow/invalid.ext", []byte("ignored"), 0644)
	fs.WriteFile(".gflows/workflows/invalid-dir.yml", []byte("ignored"), 0644)
	fs.WriteFile(".gflows/workflows/another-workflow/config.yml", []byte("config"), 0644)
	fs.WriteFile(".gflows/workflows/jsonnet/foo.jsonnet", []byte("jsonnet"), 0644)

	templates := templateEngine.GetWorkflowTemplates()

	assert.Equal(t, []string{".gflows/workflows/another-workflow", ".gflows/workflows/my-workflow"}, templates)
}

func TestGetAllYttLibs(t *testing.T) {
	config := strings.Join([]string{
		"templates:",
		"  engine: ytt",
		"  defaults:",
		"    ytt:",
		"      libs: [common, config]",
		"  overrides:",
		"    my-workflow:",
		"      ytt:",
		"        libs: [my-lib]",
	}, "\n")
	_, _, engine := newYttTemplateEngine(config)

	assert.Equal(t, []string{".gflows/common", ".gflows/config", ".gflows/my-lib"}, engine.getAllYttLibs())
}

func TestGetYttLibs(t *testing.T) {
	config := strings.Join([]string{
		"templates:",
		"  engine: ytt",
		"  defaults:",
		"    ytt:",
		"      libs: [common, config]",
		"  overrides:",
		"    my-workflow:",
		"      ytt:",
		"        libs: [my-lib]",
	}, "\n")
	_, _, engine := newYttTemplateEngine(config)

	assert.Equal(t, []string{".gflows/common", ".gflows/config", ".gflows/my-lib"}, engine.getYttLibs("my-workflow"))
	assert.Equal(t, []string{".gflows/common", ".gflows/config"}, engine.getYttLibs("other-workflow"))
}
