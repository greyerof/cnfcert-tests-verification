package globalhelper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
)

// LaunchTests stats tests based on given parameters.
func LaunchTests(testCaseName string, tcNameForReport string) error {
	containerEngine, err := container.SelectEngine()
	if err != nil {
		return fmt.Errorf("failed to select engine: %w", err)
	}

	err = os.Setenv("TNF_CONTAINER_CLIENT", containerEngine)
	if err != nil {
		return fmt.Errorf("failed to set TNF_CONTAINER_CLIENT: %w", err)
	}

	glog.V(5).Info(fmt.Sprintf("container engine set to %s", containerEngine))
	testArgs := []string{
		"-k", os.Getenv("KUBECONFIG"),
		"-c", Configuration.General.DockerConfigDir + "/config",
		"-t", Configuration.General.TnfConfigDir,
		"-o", Configuration.General.TnfReportDir,
		"-i", fmt.Sprintf("%s:%s", Configuration.General.TnfImage, Configuration.General.TnfImageTag),
		"-l", testCaseName,
	}

	cmd := exec.Command(fmt.Sprintf("./%s", Configuration.General.TnfEntryPointScript))
	cmd.Args = append(cmd.Args, testArgs...)
	cmd.Dir = Configuration.General.TnfRepoPath

	debugTnf, err := Configuration.DebugTnf()
	if err != nil {
		return fmt.Errorf("failed to set env var TNF_LOG_LEVEL: %w", err)
	}

	if debugTnf {
		outfile := Configuration.CreateLogFile(getTestSuiteName(testCaseName), tcNameForReport)

		defer outfile.Close()

		_, err = outfile.WriteString(fmt.Sprintf("Running test: %s\n", tcNameForReport))
		if err != nil {
			return fmt.Errorf("failed to write to debug file: %w", err)
		}

		cmd.Stdout = outfile
		cmd.Stderr = outfile
	}

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run tc %s: %w", testCaseName, err)
	}

	CopyClaimFileToTcFolder(testCaseName, tcNameForReport)

	return err
}

func getTestSuiteName(testCaseName string) string {
	if strings.Contains(testCaseName, globalparameters.NetworkSuiteName) {
		return globalparameters.NetworkSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.AffiliatedCertificationSuiteName) {
		return globalparameters.AffiliatedCertificationSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.LifecycleSuiteName) {
		return globalparameters.LifecycleSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.PlatformAlterationSuiteName) {
		return globalparameters.PlatformAlterationSuiteName
	}

	if strings.Contains(testCaseName, globalparameters.ObservabilitySuiteName) {
		return globalparameters.ObservabilitySuiteName
	}

	if strings.Contains(testCaseName, globalparameters.AccessControlSuiteName) {
		return globalparameters.AccessControlSuiteName
	}

	panic(fmt.Sprintf("can't retrieve test suite name from test case name %s", testCaseName))
}
