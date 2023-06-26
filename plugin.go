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
	"io"

	"github.com/arduino/go-paths-helper"
	semver "go.bug.st/relaxed-semver"
)

// Plugin is the interface that the implementations of the Firmware Uploader plugins
// must follow
type Plugin interface {
	// GetPluginInfo returns information about the plugin
	GetPluginInfo() *PluginInfo

	// UploadFirmware performs a firmware upload on the board
	UploadFirmware(portAddress string, firmwarePath *paths.Path, feedback *PluginFeedback) error

	// UploadCertificate performs a certificate upload on the board
	UploadCertificate(portAddress string, certificatePath *paths.Path, feedback *PluginFeedback) error

	// GetFirmwareVersion retrieve the firmware version installed on the board
	GetFirmwareVersion(portAddress string, feedback *PluginFeedback) (*semver.RelaxedVersion, error)
}

// PluginInfo is a set of information describing the plugin
type PluginInfo struct {
	// Name of the plugin
	Name string

	// Version of the plugin
	Version *semver.Version
}

// PluginFeedback is a struct that provides ways for the plugin to give feedback to
// the user.
type PluginFeedback struct {
	stdOut, stdErr io.Writer
}

// Out returns an output stream for console-like feedback
func (f *PluginFeedback) Out() io.Writer {
	if f.stdOut == nil {
		return io.Discard
	}
	return f.stdOut
}

// Err returns an error stream for console-like feedback
func (f *PluginFeedback) Err() io.Writer {
	if f.stdErr == nil {
		return io.Discard
	}
	return f.stdErr
}
