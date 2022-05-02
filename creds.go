// Copyright (C) 2022 Alexander Sowitzki
//
// This program is free software: you can redistribute it and/or modify it under the terms of the
// GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License along with this program.
// If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const credDirEnvName = "CREDENTIALS_DIRECTORY"

// CredsDir returns the path to the directory where systemd credentials reside.
// An error is returned if systemd did not st the env var CREDENTIALS_DIRECTORY.
func CredsDir() (string, error) {
	credsDir, ok := os.LookupEnv(credDirEnvName)
	if !ok {
		return "", fmt.Errorf("env CREDENTIALS_DIRECTORY not set")
	}

	return credsDir, nil
}

// UnmarshalYAMLCreds unmarshals the YAML credential file called name into dst.
func (s Service) UnmarshalYAMLCreds(dst interface{}, name string) error {
	dir, err := CredsDir()
	if err != nil {
		return err
	}

	credFile, err := os.Open(path.Join(dir, name))
	if err != nil {
		return fmt.Errorf("open cred file: %w", err)
	}

	err = yaml.NewDecoder(credFile).Decode(dst)

	closeErr := credFile.Close()

	switch {
	case err != nil && closeErr != nil:
		return fmt.Errorf("decode cred: %w; close cred file: %v", err, closeErr)
	case err != nil:
		return fmt.Errorf("decode cred: %w", err)
	case closeErr != nil:
		return fmt.Errorf("close cred file: %w", closeErr)
	default:
		return nil
	}
}
