package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	client "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const daemonSetReadyRetryInterval = 5 * time.Second

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and wait until all  replicas are up and running.
func CreateAndWaitUntilDaemonSetIsReady(daemonSet *v1.DaemonSet, timeout time.Duration) error {
	APIClient := client.Get()

	runningDaemonSet, err := APIClient.DaemonSets(daemonSet.Namespace).Create(
		context.Background(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isDaemonSetReady(runningDaemonSet.Namespace, runningDaemonSet.Name)
		if err != nil {
			glog.Fatal(fmt.Sprintf(
				"daemonset %s is not ready, retry in 5 seconds", runningDaemonSet.Name))

			return false
		}

		return status
	}, timeout, daemonSetReadyRetryInterval).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

func isDaemonSetReady(namespace string, name string) (bool, error) {
	APIClient := client.Get()

	daemonSet, err := APIClient.DaemonSets(namespace).Get(
		context.Background(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable == 0 {
		return true, nil
	}

	return false, nil
}
