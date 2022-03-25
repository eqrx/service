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
	"errors"
	"fmt"
	"net"
	"os"
)

// notifySocketEnvName is the env name that contains the path to the systemd notify socket.
const notifySocketEnvName = "NOTIFY_SOCKET"

// errNotifyNotSet indicates that we are not running as a systemd notify service.
var errNotifyNotSet = errors.New("env NOTIFY_SOCKET not set")

// newNotifySocket creates a new systemd notify socket.
func newNotifySocket() (*net.UnixConn, error) {
	notifySocket, hasNotifySocket := os.LookupEnv(notifySocketEnvName)

	if !hasNotifySocket {
		return nil, errNotifyNotSet
	}

	if err := os.Unsetenv(notifySocketEnvName); err != nil {
		return nil, fmt.Errorf("unset systemd notify socket env: %w", err)
	}

	socketAddr := &net.UnixAddr{Name: notifySocket, Net: "unixgram"}

	notify, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return nil, fmt.Errorf("open systemd notify socket: %w", err)
	}

	return notify, nil
}

// MarkReady tells systemd that this service is ready and running.
func (s Service) MarkReady() error {
	if _, err := s.notify.Write([]byte("READY=1")); err != nil {
		return fmt.Errorf("write to systemd notify: %w", err)
	}

	return nil
}

// MarkStopping tells systemd that this service is about to stop.
func (s Service) MarkStopping() error {
	if _, err := s.notify.Write([]byte("STOPPING=1")); err != nil {
		return fmt.Errorf("write to systemd notify: %w", err)
	}

	return nil
}

// MarkStatus tells systemd that this service has the given status.
func (s Service) MarkStatus(status string) error {
	if _, err := s.notify.Write([]byte("STATUS=" + status)); err != nil {
		return fmt.Errorf("write to systemd notify: %w", err)
	}

	return nil
}
