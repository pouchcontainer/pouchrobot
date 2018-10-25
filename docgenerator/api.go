// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docgenerator

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// generateAPIDoc will generate api doc via Swagger2Markup.
// First, execute command `java -jar swagger2markup-cli-{release-version}.jar convert -i...`

// Swagger2MarkupReleaseVersion represents version of swagger2markup.
var swagger2markupReleaseVersion = "1.3.1"
var swagger2markupJar = fmt.Sprintf("swagger2markup-cli-%s.jar", swagger2markupReleaseVersion)
var swagger2markupPath = fmt.Sprintf("/root/%s", swagger2markupJar)

var swagger2markupConfig = "/go/src/github.com/pouchcontainer/pouchrobot/config.properties"

//var swaggerYML = "/go/src/github.com/alibaba/pouch/apis/swagger.yml"
// var targetAPIFile = "/go/src/github.com/alibaba/pouch/docs/api/HTTP_API"

func (g *Generator) generateAPIDoc() error {
	if g.RootDir == "" {
		return fmt.Errorf("API doc generation fails with no root dir set down")
	}

	swaggerYML := filepath.Join(g.RootDir, g.SwaggerPath)
	targetAPIFile := filepath.Join(g.RootDir, g.APIDocPath)
	args := []string{"-jar",
		swagger2markupPath,
		"convert",
		"-i",
		swaggerYML,
		"-f",
		targetAPIFile,
		"-c",
		swagger2markupConfig,
	}
	cmd := exec.Command("java", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to use swagger2markdown to generate API docs: %v", err)
	}
	return nil
}
