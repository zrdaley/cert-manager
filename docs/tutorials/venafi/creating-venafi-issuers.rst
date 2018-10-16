Venafi issuer for Jetstack cert-manager
=======================================

Venafi issuer is a cert-manager
(https://github.com/jetstack/cert-manager) extension which support
certificate management from Venafi Cloud and Venafi Venafi Platform. Also
it have fakeissuer interface for testing purpose.

Creating Venafi Cloud issuer
---------------------------------------------

Register your account at https://api.venafi.cloud/login and get API key there.

Create a secret for the issuer (in this example the issuer will be
Venafi Cloud and we'll use the default namespace)

::

    kubectl create secret generic cloudsecret --from-literal=apikey='YOUR_CLOUD_API_KEY_HERE'

Create the issuer

.. code-block:: yaml

    apiVersion: certmanager.k8s.io/v1alpha1
    kind: Issuer
    metadata:
            name: cloud-venafi-issuer
    spec:
            venafi:
                    cloudsecret: cloudsecret
                    zone: "DevOps"

You can create multiple issuers pointing to different Venafi Cloud zones, or
even have 1 issuer pointing to Venafi Platform and another pointing to Venafi Cloud.

Here's an example certificate resource file using the new issuer:

.. code-block:: yaml

    apiVersion: certmanager.k8s.io/v1alpha1
    kind: Certificate
    metadata:
            name: cert4-venafi-localhost
    spec:
            secretName: cert4-venafi-localhost
            issuerRef:
                    name: cloud-venafi-issuer
            commonName: cert4.venafi.localhost



Creating Venafi Platform issuer
--------------------------

By default one Venafi Platform issuer is alredy created when you run "make install",
it called tppvenafiissuer. You can create more issuers for different Venafi Platform
server or policies.

**Requirements for Venafi Platform policy**


1. Policy should have default template configured

2. Currently vcert (which is used in Venafi issuers) supports only user
   provided CSR. So it is must be set in the policy.

3. MSCA configuration should have http URI set before the ldap URI in
   X509 extensions, otherwise NGINX ingress controller couldn't get
   certificate chain from URL and OSCP will not work. Example:

::

    X509v3 extensions:
        X509v3 Subject Alternative Name:
        DNS:test-cert-manager1.venqa.venafi.com}}
        X509v3 Subject Key Identifier: }}
        61:5B:4D:40:F2:CF:87:D5:75:5E:58:55:EF:E8:9E:02:9D:E1:81:8E}}
        X509v3 Authority Key Identifier: }}
        keyid:3C:AC:9C:A6:0D:A1:30:D4:56:A7:3D:78:BC:23:1B:EC:B4:7B:4D:75}}X509v3 CRL Distribution Points:Full Name:
        URI:http://qavenafica.venqa.venafi.com/CertEnroll/QA%20Venafi%20CA.crl}}
        URI:ldap:///CN=QA%20Venafi%20CA,CN=qavenafica,CN=CDP,CN=Public%20Key%20Services,CN=Services,CN=Configuration,DC=venqa,DC=venafi,DC=com?certificateRevocationList?base?objectClass=cRLDistributionPoint}}{{Authority Information Access: }}
        CA Issuers - URI:http://qavenafica.venqa.venafi.com/CertEnroll/qavenafica.venqa.venafi.com_QA%20Venafi%20CA.crt}}
        CA Issuers - URI:ldap:///CN=QA%20Venafi%20CA,CN=AIA,CN=Public%20Key%20Services,CN=Services,CN=Configuration,DC=venqa,DC=venafi,DC=com?cACertificate?base?objectClass=certificationAuthority}}

4. Option in Venafi Platform CA configuration template "Automatically include CN as
   DNS SAN" should be set to true.

**Create a secret with Venafi Platform credentials:**

::

    kubectl create secret generic tppsecret --from-literal=user=admin --from-literal=password=tpppassword --namespace cert-manager-example

Create Venafi Platform issuer

.. code-block:: yaml

    apiVersion: certmanager.k8s.io/v1alpha1
    kind: Issuer
    metadata:
      name: tpp-venafi-issuer
    spec:
      venafi:
        tppsecret: tppsecret
        tppurl: https://tpp.venafi.example/vedsdk
        zone: devops\cert-manager
    status:
      conditions:
      - lastTransitionTime: 2018-08-03T12:26:58Z
        message: Venafi issuer started
        reason: Venafi issuer started
        status: "True"
        type: Ready




Create a certificate

cert.yaml:

.. code-block:: yaml

    apiVersion: certmanager.k8s.io/v1alpha1
    kind: Certificate
    metadata:
            name: hellodemo-venafi-localhost
            namespace: cert-manager-example
    spec:
            secretName: hellodemo-venafi-localhost
            issuerRef:
                    name: tppvenafiissuer
            commonName: hellodemo.venafi.localhost



