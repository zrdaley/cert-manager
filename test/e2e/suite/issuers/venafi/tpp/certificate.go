package tpp

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/jetstack/cert-manager/test/e2e/framework"
	vaddon "github.com/jetstack/cert-manager/test/e2e/suite/issuers/venafi/addon"
	"github.com/jetstack/cert-manager/test/util"
)

var _ = TPPDescribe("with a properly configured Issuer", func() {
	f := framework.NewDefaultFramework("venafi-tpp-certificate")

	var (
		issuer                *cmapi.Issuer
		tppAddon              = &vaddon.VenafiTPP{}
		certificateName       = "test-venafi-cert"
		certificateSecretName = "test-venafi-cert-tls"
	)

	BeforeEach(func() {
		tppAddon.Namespace = f.Namespace.Name
	})

	f.RequireAddon(tppAddon)

	// Create the Issuer resource
	BeforeEach(func() {
		var err error

		By("Creating a Venafi Issuer resource")
		issuer = tppAddon.Details().BuildIssuer()
		issuer, err = f.CertManagerClientSet.CertmanagerV1alpha1().Issuers(f.Namespace.Name).Create(issuer)
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for Issuer to become Ready")
		err = util.WaitForIssuerCondition(f.CertManagerClientSet.CertmanagerV1alpha1().Issuers(f.Namespace.Name),
			issuer.Name,
			cmapi.IssuerCondition{
				Type:   cmapi.IssuerConditionReady,
				Status: cmapi.ConditionTrue,
			})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		By("Cleaning up")
		f.CertManagerClientSet.CertmanagerV1alpha1().Issuers(f.Namespace.Name).Delete(issuer.Name, nil)
	})

	It("should obtain a signed certificate for a single domain", func() {
		certClient := f.CertManagerClientSet.CertmanagerV1alpha1().Certificates(f.Namespace.Name)
		secretClient := f.KubeClientSet.CoreV1().Secrets(f.Namespace.Name)

		By("Creating a Certificate")
		_, err := certClient.Create(
			util.NewCertManagerBasicCertificate(certificateName, certificateSecretName, issuer.Name, cmapi.IssuerKind),
		)
		Expect(err).NotTo(HaveOccurred())
		By("Verifying the Certificate is valid")
		err = util.WaitCertificateIssuedValid(certClient, secretClient, certificateName, time.Second*30)
		Expect(err).NotTo(HaveOccurred())
	})
})
