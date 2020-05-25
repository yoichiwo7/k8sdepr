package mix

import (
	"k8s.io/api/extensions/v1beta1"
)

type mystruct struct {
	ingress v1beta1.Ingress // want `extensions/v1beta1:Ingress is deprecated. Migrate to networking.k8s.io/v1beta1:Ingress. {deprecated=v1.14.0, removed=v1.22.0}`
}
