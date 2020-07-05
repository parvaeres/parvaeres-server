package parvaeres

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetInClusterClient(t *testing.T) {
	Convey("When trying to get a client outside the cluster", t, func() {
		clientset, err := getInClusterClient()
		Convey("Then it fails", func() {
			So(err, ShouldNotBeNil)
			So(clientset, ShouldBeNil)
			So(err.Error(), ShouldStartWith, "unable to load in-cluster configuration")

		})
	})
}
