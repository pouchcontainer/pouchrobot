package doc

import (
	"fmt"
	"os/exec"
)

// generateAPIDoc will generate api doc via Swagger2Markup.
// First, execute command `java -jar swagger2markup-cli-{release-version}.jar convert -i...`

// Swagger2MarkupReleaseVersion represents version of swagger2markup.
var swagger2markupReleaseVersion = "1.3.1"
var swagger2markupJar = fmt.Sprintf("swagger2markup-cli-%s.jar", swagger2markupReleaseVersion)
var swagger2markupPath = fmt.Sprintf("/root/%s", swagger2markupJar)

var swaggerYML = "/go/src/github.com/alibaba/pouch/apis/swagger.yml"
var swagger2markupConfig = "/go/src/github.com/allencloud/automan/config.properties"
var targetAPIFile = "/go/src/github.com/alibaba/pouch/docs/api/HTTP_API"

func (g *Generator) generateAPIDoc() error {
	args := []string{"-jar", swagger2markupPath, "convert", "-i", swaggerYML, "-f", targetAPIFile, "-c", swagger2markupConfig}
	cmd := exec.Command("java", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to use swagger2markdown to generate API docs: %v", err)
	}
	return nil
}
