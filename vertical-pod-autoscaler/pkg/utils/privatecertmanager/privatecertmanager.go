package privatecertmanager

import (
	"fmt"

	"github.com/open-policy-agent/cert-controller/pkg/rotator"
	"k8s.io/apimachinery/pkg/types"
	controller "sigs.k8s.io/controller-runtime"
)

const (
	serviceName     = "vpa-webhook"
	vwhName         = "vpa-webhook-webhook-configuration"
	caName          = "vpa-ca"
	caOrganization  = "vpa"
	secretNamespace = "kube-system"
	secretName      = "vpa-webhook-server-cert"
	certDir         = "/etc/tls-certs/"
)

var dnsName = fmt.Sprintf("%s.%s.svc", serviceName, secretNamespace)

// CreatePrivateCert creates all certs for webhooks. This function is called from main.go.
func CreatePrivateCert(mgr controller.Manager, certRequired bool) (chan struct{}, error) {
	setupFinished := make(chan struct{})
	if !certRequired {
		close(setupFinished)
		return setupFinished, nil
	}

	return setupFinished, rotator.AddRotator(mgr, &rotator.CertRotator{
		SecretKey: types.NamespacedName{
			Namespace: secretNamespace,
			Name:      secretName,
		},
		CertDir:        certDir,
		CAName:         caName,
		CAOrganization: caOrganization,
		DNSName:        dnsName,
		IsReady:        setupFinished,
		Webhooks: []rotator.WebhookInfo{{
			Type: rotator.Validating,
			Name: vwhName,
		}},
	})
}
