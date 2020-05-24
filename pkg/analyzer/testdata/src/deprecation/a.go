package deprecation

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	ext "k8s.io/api/extensions/v1beta1"
)

type myinterface interface {
	myfunc0(ext.DaemonSet) error // want `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.`
}
type mystruct struct {
	dsOld ext.DaemonSet // want `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.`
	dsNew v1.DaemonSet
	aaa   v1.DaemonSet
}

type mytype ext.DaemonSet // want `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.`

func myfunc1() {
	dsOld := ext.DaemonSet{} // want `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.`
	dsNew := v1.DaemonSet{}
	fmt.Println("ds", dsOld)
	fmt.Println("ds", dsNew)
}

func myfunc2(ds ext.DaemonSet) ext.DaemonSet { // want `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.` `extensions/v1beta1:DaemonSet is deprecated in v1.9.0. Migrate to apps/v1:DaemonSet.`
	return ds
}
