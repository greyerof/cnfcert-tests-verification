package tests

import (
	"context"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
)

const (
	WaitingTime    = 5 * time.Minute
	istioNamespace = "istio-system"
)

var _ = Describe("platform-alteration-service-mesh-usage-installed", func() {

	execute.BeforeAll(func() {
		By("Install istio")
		cmd := exec.Command("/bin/bash", "-c", "curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.17.1 | sh - "+
			"&& istio-1.17.1/bin/istioctl install --set profile=demo -y")
		err := cmd.Run()
		Expect(err).ToNot(HaveOccurred(), "Error installing istio")
	})

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56594
	It("istio is installed", func() {
		By("Define a test pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56596
	It("istio is installed but proxy containers does not exist [negative]", func() {
		By("Define a test pod without istio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56597
	It("istio is installed but proxy container exist on one pod only [negative]", func() {
		By("Define first pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		putb := pod.DefinePod("lifecycle-putb", tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = Describe("platform-alteration-service-mesh-usage-uninstalled", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56595
	It("istio is not installed", func() {
		By("Check if Istio resource exists")
		gvr := schema.GroupVersionResource{Group: "install.istio.io", Version: "v1alpha1", Resource: "istiooperators"}

		_, err := globalhelper.APIClient.DynamicClient.Resource(gvr).Namespace(istioNamespace).Get(context.TODO(),
			"installed-state", metav1.GetOptions{})

		if err == nil {
			By("Uninstall istio")
			cmd := exec.Command("/bin/bash", "-c", "curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.17.1 | sh - "+
				"&& istio-1.17.1/bin/istioctl uninstall --purge")
			err := cmd.Run()
			Expect(err).ToNot(HaveOccurred(), "Error uninstalling istio")
		}

		By("Define a test pod with istio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseSkipped)
		Expect(err).ToNot(HaveOccurred())
	})
})
