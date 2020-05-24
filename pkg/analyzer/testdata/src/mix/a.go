package mix

import (
	"k8s.io/api/extensions/v1beta1"
)

type mystruct struct {
	ingress v1beta1.Ingress // want `extensions/v1beta1:Ingress is deprecated in v1.14.0. Migrate to networking.k8s.io/v1beta1:Ingress.`
}
