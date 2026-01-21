// Package updater provides self-updating functionality for the application.
package updater

import (
	"fmt"
	"log/slog"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// Update performs the self-update process.
func Update(currentVersion string, repo string) error {
	v, err := semver.Parse(currentVersion)
	if err != nil {
		return fmt.Errorf("could not parse current version: %w", err)
	}

	slog.Info("checking for updates...", "current_version", currentVersion)

	latest, err := selfupdate.UpdateSelf(v, repo)
	if err != nil {
		return fmt.Errorf("binary update failed: %w", err)
	}

	if latest.Version.Equals(v) {
		slog.Info("current binary is the latest version", "version", currentVersion)
	} else {
		slog.Info("successfully updated to new version",
			"new_version", latest.Version,
			"release_notes", latest.ReleaseNotes,
		)
		slog.Info("please restart the application")
	}

	return nil
}

// Check checks for updates without applying them.
func Check(currentVersion string, repo string) (*selfupdate.Release, bool, error) {
	v, err := semver.Parse(currentVersion)
	if err != nil {
		return nil, false, fmt.Errorf("error parsing version: %w", err)
	}

	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		return nil, false, fmt.Errorf("error checking for update: %w", err)
	}

	if !found || latest.Version.LE(v) {
		return nil, false, nil
	}

	return latest, true, nil
}
