package isolated

import (
	"fmt"
	"os"
	"testing"
	"time"

	"code.cloudfoundry.org/cli/integration/helpers/fakeservicebroker"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccversion"
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	CFEventuallyTimeout   = 300 * time.Second
	CFConsistentlyTimeout = 500 * time.Millisecond
	RealIsolationSegment  = "persistent_isolation_segment"
)

var (
	// Suite Level
	apiURL            string
	skipSSLValidation bool
	ReadOnlyOrg       string
	ReadOnlySpace     string

	// Per Test Level
	homeDir string
)

func TestIsolated(t *testing.T) {
	RegisterFailHandler(Fail)
	reporters := []Reporter{}

	prBuilderReporter := helpers.GetPRBuilderReporter()
	if prBuilderReporter != nil {
		reporters = append(reporters, prBuilderReporter)
	}

	RunSpecsWithDefaultAndCustomReporters(t, "Isolated Integration Suite", reporters)
}

var _ = SynchronizedBeforeSuite(func() []byte {
	GinkgoWriter.Write([]byte("==============================Global FIRST Node Synchronized Before Each=============================="))
	SetDefaultEventuallyTimeout(CFEventuallyTimeout)
	SetDefaultConsistentlyDuration(CFConsistentlyTimeout)

	helpers.SetupSynchronizedSuite(func() {
		helpers.EnableFeatureFlag("diego_docker")

		if helpers.IsVersionMet(ccversion.MinVersionShareServiceV3) {
			helpers.EnableFeatureFlag("service_instance_sharing")
		}
	})

	fakeservicebroker.Setup()

	GinkgoWriter.Write([]byte("==============================End of Global FIRST Node Synchronized Before Each=============================="))

	return nil
}, func(_ []byte) {
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================Global Node %d Synchronized Before Each==============================", GinkgoParallelNode())))
	// Ginkgo Globals
	SetDefaultEventuallyTimeout(CFEventuallyTimeout)
	SetDefaultConsistentlyDuration(CFConsistentlyTimeout)

	// Setup common environment variables
	helpers.TurnOffColors()

	ReadOnlyOrg, ReadOnlySpace = helpers.SetupReadOnlyOrgAndSpace()
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================End of Global Node %d Synchronized Before Each==============================", GinkgoParallelNode())))
})

var _ = SynchronizedAfterSuite(func() {
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================Global Node %d Synchronized After Each==============================", GinkgoParallelNode())))
	homeDir = helpers.SetHomeDir()
	helpers.SetAPI()
	helpers.LoginCF()
	fakeservicebroker.Cleanup()
	helpers.QuickDeleteOrg(ReadOnlyOrg)
	helpers.DestroyHomeDir(homeDir)
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================End of Global Node %d Synchronized After Each==============================", GinkgoParallelNode())))
}, func() {
	outputRoot := os.Getenv(helpers.PRBuilderOutputEnvVar)
	if outputRoot != "" {
		helpers.WriteFailureSummary(outputRoot, "summary_isi.txt")
	}
})

var _ = BeforeEach(func() {
	GinkgoWriter.Write([]byte("==============================Global Before Each=============================="))
	homeDir = helpers.SetHomeDir()
	apiURL, skipSSLValidation = helpers.SetAPI()
	GinkgoWriter.Write([]byte("==============================End of Global Before Each=============================="))
})

var _ = AfterEach(func() {
	GinkgoWriter.Write([]byte("==============================Global After Each==============================\n"))
	helpers.DestroyHomeDir(homeDir)
	GinkgoWriter.Write([]byte("==============================End of Global After Each=============================="))
})
