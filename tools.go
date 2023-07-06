// fwuploader-plugin-helper
// Copyright (c) 2023 Arduino LLC.  All right reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package helper

import (
	"fmt"
	"os"

	"github.com/arduino/go-paths-helper"
	semver "go.bug.st/relaxed-semver"
)

// FindToolPath retrieve the path to the given tool, if available, otherwise returns an error
func FindToolPath(toolName string, toolVersion *semver.Version) (*paths.Path, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	toolsDir := paths.New(executablePath).Parent().Parent().Parent()
	toolPath := toolsDir.Join(toolName, toolVersion.String())
	if !toolPath.IsDir() {
		return nil, fmt.Errorf("tool not found: %s@%s", toolName, toolVersion)
	}
	return toolPath, nil
}
