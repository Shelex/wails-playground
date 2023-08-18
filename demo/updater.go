package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const AppPath = "Shelex/scaling-octo-garbanzo"

func (a *App) getUpdater() *selfupdate.Updater {
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})

	if err != nil {
		runtime.LogPrintf(a.ctx, "failed to initialize updater: %s", err)
	}

	return updater
}

func (a *App) doSelfUpdate() bool {
	v := semver.MustParse(a.version)

	updater := a.getUpdater()

	if _, err := updater.UpdateSelf(v, AppPath); err != nil {
		runtime.LogPrintf(a.ctx, "binary update failed: %s", err)
		return false
	}

	return true
}

func (a *App) doSelfUpdateMac() bool {
	updater := a.getUpdater()

	latest, found, _ := updater.DetectLatest(AppPath)
	if !found {
		return false
	}

	homeDir, _ := os.UserHomeDir()
	downloadPath := filepath.Join(homeDir, "Downloads", "Demo-UI.zip")
	if err := exec.Command("curl", "-L", latest.AssetURL, "-o", downloadPath).Run(); err != nil {
		runtime.LogPrintf(a.ctx, "curl error: %s", err)
		return false
	}
	var appPath string
	cmdPath, err := os.Executable()
	appPath = strings.TrimSuffix(cmdPath, "Demo-UI.app/Contents/MacOS/Demo-UI")
	if err != nil {
		appPath = "/Applications/"
	}
	if err := exec.Command("ditto", "-xk", downloadPath, appPath).Run(); err != nil {
		runtime.LogPrintf(a.ctx, "ditto error: %s", err)
		return false
	}
	if err := exec.Command("rm", downloadPath).Run(); err != nil {
		runtime.LogPrintf(a.ctx, "removing error: %s", err)
		return false
	}
	return true
}

func (a *App) checkForUpdate() (bool, string) {
	updater := a.getUpdater()
	latest, found, err := updater.DetectLatest(AppPath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "error occurred while detecting version: %s", err)
		return false, ""
	}

	if !found {
		runtime.LogPrint(a.ctx, "could not fetch latest version")
		return false, ""
	}

	v := semver.MustParse(a.version)
	if latest.Version.LTE(v) {
		runtime.LogPrintf(a.ctx, "current version %s is the latest", latest.Version)
		return false, ""
	}

	return true, latest.Version.String()
}

func (a *App) updateCheckDialog() {
	shouldUpdate, latestVersion := a.checkForUpdate()

	if !shouldUpdate {
		return
	}

	updateMessage := fmt.Sprintf("New version %s is available, would you like to update?", latestVersion)
	yes, no := "Yes", "No"
	buttons := []string{yes, no}
	dialogOpts := runtime.MessageDialogOptions{Title: "Update Available", Message: updateMessage, Type: runtime.QuestionDialog, Buttons: buttons, DefaultButton: yes, CancelButton: no}
	action, err := runtime.MessageDialog(a.ctx, dialogOpts)
	if err != nil {
		runtime.LogError(a.ctx, "error in update dialog")
	}
	runtime.LogInfo(a.ctx, action)

	if action == "No" {
		return
	}

	runtime.LogInfo(a.ctx, "update clicked")
	var updated bool
	env := runtime.Environment(a.ctx)
	if env.Platform == "darwin" {
		updated = a.doSelfUpdateMac()
	} else {
		updated = a.doSelfUpdate()
	}

	restartNow, later := "Close Now", "Meh, will do later"
	afterUpdateButtons := []string{restartNow, later}
	afterUpdateDialog := runtime.MessageDialogOptions{
		Type: runtime.InfoDialog, Buttons: afterUpdateButtons, DefaultButton: restartNow,
	}

	if updated {
		afterUpdateDialog.Title = "Update Succeeded"
		afterUpdateDialog.Message = "Update Successful. App should be restarted to take effect."
		a.version = latestVersion
	} else {
		afterUpdateDialog.Title = "Update Error"
		afterUpdateDialog.Message = "Update failed, please manually update from GitHub Releases."
		afterUpdateDialog.Buttons = []string{"That is sad :c"}
	}

	afterUpdateAction, err := runtime.MessageDialog(a.ctx, afterUpdateDialog)

	if err != nil {
		runtime.LogError(a.ctx, "error in after update dialog")
	}

	if afterUpdateAction == restartNow {
		runtime.Quit(a.ctx)
	}
}
