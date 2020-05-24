package ok

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	ext "k8s.io/api/extensions/v1beta1"
)

type myinterface interface {
	myfunc0(ext.DaemonSet) error
}
type mystruct struct {
	dsOld ext.DaemonSet
	dsNew v1.DaemonSet
	aaa   v1.DaemonSet
}

type mytype ext.DaemonSet

func myfunc1() {
	dsOld := ext.DaemonSet{}
	dsNew := v1.DaemonSet{}
	fmt.Println("ds", dsOld)
	fmt.Println("ds", dsNew)
}

func myfunc2(ds ext.DaemonSet) ext.DaemonSet {
	return ds
}
