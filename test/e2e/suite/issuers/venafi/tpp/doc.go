// Package tpp implements tests for the Venafi TPP issuer
package tpp

import (
	"github.com/jetstack/cert-manager/test/e2e/framework"
)

func TPPDescribe(name string, body func()) bool {
	return framework.CertManagerDescribe("[Venafi] [TPP] "+name, body)
}
