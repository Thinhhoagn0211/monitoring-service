package util

import (
	"fmt"
	"os"
)

// GetAvailableDrives returns a list of available drives on the system
func GetAvailableDrives() []string {
	var drives []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := fmt.Sprintf("%c:\\", drive)
		if _, err := os.Stat(drivePath); !os.IsNotExist(err) {
			drives = append(drives, drivePath)
		}
	}
	return drives
}
