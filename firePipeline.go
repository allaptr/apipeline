package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type FireSigner struct {
	Certificate X509Certificate
	CpixDoc     CPIXDocument
}

func NewFireCPIXPipeline() CPIXPipeline {
	return FireCertFetcher{
		Namespace:    "Fire",
		CertLocation: "Fire cert location", // get from a config store
	}
}

// ---------------Concrete implementation of X509CertFetcher interface-------------------
type FireCertFetcher struct {
	Namespace    string
	CertLocation string
	Certificate  X509Certificate
}

func (ff FireCertFetcher) FetchCert() (X509Certificate, error) {
	var err error
	fmt.Println("Access Secret Store and retrieve the Fire certificate")
	if err != nil {
		return "", err
	}
	return "Fire-X509Certtificate", nil
}

// CPIXPipeline interface implementation
func (ff FireCertFetcher) Handle() error {
	cert, err := ff.FetchCert()
	if err != nil {
		return err
	}
	ff.Certificate = cert
	return nil
}

func (ff FireCertFetcher) Next() CPIXPipeline {
	return FireHttpClient{
		DRMServerUrl: "https://fire.drmserver.com",
		Certificate:  ff.Certificate,
	}
}

func (fs FireSigner) ValidateSignature() bool {
	fmt.Println("Validating Signature in CPIXDocument")
	return true
}

func (fs FireSigner) SignDocument(input []byte) []byte {
	fmt.Println("Signing POST Request")
	return []byte("Signed POST Request")
}

// ---------------Concrete implementation of HTTPClient interface-------------------
type FireHttpClient struct {
	DRMServerUrl string
	Certificate  X509Certificate
	CpixBytes    []byte
}

func (fhc FireHttpClient) CreateRequest() ([]byte, error) {
	var err error
	fmt.Println("Create POST HTTP request")
	reqData := []byte("POST Request Data")
	fmt.Println("Sign the request using X509Certificate")
	fs := FireSigner{
		Certificate: fhc.Certificate,
	}
	signedReqData := fs.SignDocument(reqData)
	if err != nil {
		return nil, err
	}
	return signedReqData, nil
}

func (fhc FireHttpClient) ExecuteRequestResponse(request []byte) ([]byte, error) {
	var err error
	fmt.Println("Construct TLS connection")
	fmt.Println("Send POST HTTP request to Fire DRM Server")
	fmt.Println("Receive data from Fire DRM Server")
	status := 200
	if !(status == http.StatusOK) {
		return nil, errors.Errorf("error from Fire DRM server %v", status)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error from DRM server")
	}
	payload := []byte("cpixdocument")
	return payload, nil
}

// CPIXPipeline interface implementation
func (fhc FireHttpClient) Handle() error {
	req, err := fhc.CreateRequest()
	if err != nil {
		return err
	}
	payload, err := fhc.ExecuteRequestResponse(req)
	if err != nil {
		return err
	}
	fhc.CpixBytes = payload
	return nil
}
func (fhc FireHttpClient) Next() CPIXPipeline {
	return FireCPIXProcessor{
		Certificate: fhc.Certificate,
		CpixBytes:   fhc.CpixBytes,
	}
}

// ---------------Concrete implementation of CPIXProcessor interface-------------------
type FireCPIXProcessor struct {
	Certificate X509Certificate
	CpixBytes   []byte
	CpixDoc     CPIXDocument
	DocKey      DocumentKey
	Key         ContentKey
}

func (scp FireCPIXProcessor) UnarshalDocument() (*CPIXDocument, error) {
	fmt.Println("Unmarshalling Payload into CPIXDocument")
	var err error
	if err != nil {
		return nil, err
	}
	doc := CPIXDocument{}
	fs := FireSigner{
		Certificate: scp.Certificate,
		CpixDoc:     scp.CpixDoc,
	}
	if !fs.ValidateSignature() {
		return nil, errors.New("Document Signature cannot be validated")
	}

	return &doc, nil
}

func (fcp FireCPIXProcessor) ValidateDeliveryKey() bool {
	fmt.Println("Validated Delivery Key")
	return true
}

func (fcp FireCPIXProcessor) DecryptDocumentKey() (DocumentKey, error) {
	fmt.Println("Decrypted Document Key using X509Certificate")
	var documentKey DocumentKey
	documentKey = "Document-Key"
	return documentKey, nil
}

func (fcp FireCPIXProcessor) DecryptContentKey() (ContentKey, error) {
	fmt.Println("Decrypted Content Key using Document Key")
	fmt.Println("Return Content Key in clear")
	var contentKey ContentKey
	contentKey = "Content-Key"
	return contentKey, nil
}

// CPIXPipeline interface implementation
func (fcp FireCPIXProcessor) Handle() error {
	doc, err := fcp.UnarshalDocument()
	if err != nil {
		return err
	}
	fcp.CpixDoc = *doc
	if !fcp.ValidateDeliveryKey() {
		return errors.New("Error validating delivery key")
	}
	dkey, err := fcp.DecryptDocumentKey()
	if err != nil {
		return err
	}
	fcp.DocKey = dkey
	key, err := fcp.DecryptContentKey()
	if err != nil {
		return err
	}
	fcp.Key = key
	return nil
}

func (fcp FireCPIXProcessor) Next() CPIXPipeline {
	return FireHandOff{
		Key: fcp.Key,
	}
}

// ---------------Concrete implementation of HandOff interface-------------------
type FireHandOff struct {
	Key           ContentKey
	KafkaEndpoint string
	KafkaTopic    string
}

func (fho FireHandOff) PushDownstream() error {
	var err error
	fho.KafkaEndpoint = "kafkaEndpoint" // read from a config store
	fho.KafkaTopic = "Fire"
	fmt.Println("Put the content key on kafka queue, topic 'Fire'")
	if err != nil {
		return errors.Wrap(err, "could not push data downstream")
	}
	return nil
}

// CPIXPipeline interface implementation
func (fho FireHandOff) Handle() error {
	return fho.PushDownstream()
}

func (fho FireHandOff) Next() CPIXPipeline {
	return nil
}
