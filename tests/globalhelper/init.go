package globalhelper

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
)

var (
	Configuration *config.Config
)

func init() {
	var err error

	Configuration, err = config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load configuration"))
	}
}
