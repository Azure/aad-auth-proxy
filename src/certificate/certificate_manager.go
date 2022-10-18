package certificate

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// This manages certificates which should have publish permissions to accounts we are trying to ingest metrics.
// This converts pfx certificate containing private key to Tls certificate.
// The certificate is read from a CSI driver mounted path.
type CertificateManager struct {
	path        string
	certificate *tls.Certificate
}

// Creates a new certificate manager.
func NewCerificateManager(path string) (certManager *CertificateManager, err error) {
	if path == "" {
		return nil, errors.New("certificate path cannot be empty")
	}

	return &CertificateManager{
		path:        path,
		certificate: nil,
	}, nil
}

func (manager *CertificateManager) GetTlsCertificate() (tlsCert *tls.Certificate, hasCertChanged bool, err error) {
	certificate, err := manager.fetchTlsCertificate()
	if err != nil {
		return nil, hasCertChanged, err
	}

	// certificate not yet read or it has changed
	if manager.certificate == nil || md5.Sum(certificate.Leaf.Raw) != md5.Sum(manager.certificate.Leaf.Raw) {
		log.WithField("certificatePath", manager.path).Info("loaded a new certificate")
		hasCertChanged = true
		manager.certificate = certificate
	}

	return manager.certificate, hasCertChanged, err
}

// Reads certifiate from a local mount path.
// CSI driver would have downloaded latest file to this path.
func (manager *CertificateManager) readCertificateFromLocal() ([]byte, error) {
	content, err := ioutil.ReadFile(manager.path)
	if err != nil {
		log.WithFields(log.Fields{
			"certificatePath": manager.path,
		}).Error("failed to read certificate: ", err)
		return nil, err
	}

	return content, nil
}

// Converts pfx certificate containing private key to Tls certificate.
func (manager *CertificateManager) fetchTlsCertificate() (tlsCert *tls.Certificate, err error) {
	content, err := manager.readCertificateFromLocal()
	if err != nil {
		return nil, err
	}

	// Decode chain returns private key, leaf cert(public cert in 0 indes), extra public certs and error if any
	privateKey, leafCert, certs, err := pkcs12.DecodeChain(content, "")
	if err != nil {
		log.WithFields(log.Fields{
			"certificatePath": manager.path,
		}).Error("failed to decode certificate: ", err)
		return nil, err
	}

	// there can be more than one public cert and we need the last one in the list
	if len(certs) > 0 {
		leafCert = certs[len(certs)-1]
	}

	tlsCert = &tls.Certificate{
		Certificate: [][]byte{leafCert.Raw},
		PrivateKey:  privateKey,
		Leaf:        leafCert,
	}

	return tlsCert, nil
}
