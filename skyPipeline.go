package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func NewSkyCPIXPipeline() CPIXPipeline {
	return SkyCertFetcher{
		Namespace:    "Sky",
		CertLocation: "Sky cert location", // get from a config store
	}
}

// ---------------Concrete implementation of X509CertFetcher interface-------------------
type SkyCertFetcher struct {
	Namespace    string
	CertLocation string
	Certificate  X509Certificate
}

func (sf SkyCertFetcher) FetchCert() (X509Certificate, error) {
	var err error
	fmt.Println("Access Secret Store and retrieve the Sky certificate")
	if err != nil {
		return "", err
	}
	return "Sky-X509Certtificate", nil
}

// CPIXPipeline interface implementation
func (sf SkyCertFetcher) Handle() error {
	cert, err := sf.FetchCert()
	if err != nil {
		return err
	}
	sf.Certificate = cert
	return nil
}

func (sf SkyCertFetcher) Next() CPIXPipeline {
	return SkyHttpClient{
		DRMServerUrl: "https://sky.drmserver.com", // get from a config store
		Certificate:  sf.Certificate,
	}
}

// ---------------Concrete implementation of HTTPClient interface-------------------
type SkyHttpClient struct {
	DRMServerUrl string
	Certificate  X509Certificate
	Payload      []byte
}

func (shc SkyHttpClient) CreateRequest() ([]byte, error) {
	var err error
	fmt.Println("Create GET HTTP request")
	fmt.Println("Construct mTLS connection, use X509Certificate")
	fmt.Println("Set mTLS header")
	if err != nil {
		return nil, err
	}

	return []byte("httpRequest"), nil
}

func (shc SkyHttpClient) ExecuteRequestResponse(request []byte) ([]byte, error) {
	var err error

	fmt.Println("Send GET HTTP request to Sky DRM Server")
	fmt.Println("Receive data from Sky DRM Server")
	status := 200
	if !(status == http.StatusOK) {
		return nil, errors.Errorf("error from DRM server %v", status)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error from DRM server")
	}
	payload := []byte("CPIXDocument")
	return payload, nil
}

// CPIXPipeline interface implementation
func (shc SkyHttpClient) Handle() error {
	req, err := shc.CreateRequest()
	if err != nil {
		return err
	}
	bytes, err := shc.ExecuteRequestResponse(req)
	if err != nil {
		return err
	}
	shc.Payload = bytes
	return nil
}
func (shc SkyHttpClient) Next() CPIXPipeline {
	return SkyCPIXProcessor{
		CpixBytes: shc.Payload,
	}
}

// ---------------Concrete implementation of CPIXProcessor interface-------------------
type SkyCPIXProcessor struct {
	CpixBytes []byte
	CpixDoc   CPIXDocument
	Key       ContentKey
}

func (scp SkyCPIXProcessor) UnarshalDocument() (*CPIXDocument, error) {
	fmt.Println("Unmarshalling Payload into CPIXDocument")
	var err error
	if err != nil {
		return nil, err
	}
	doc := CPIXDocument{}
	return &doc, nil
}

func (scp SkyCPIXProcessor) ValidateDeliveryKey() bool {
	fmt.Println("Noop. No delivery key in the document")
	return true
}

func (scp SkyCPIXProcessor) DecryptDocumentKey() (DocumentKey, error) {
	fmt.Println("Noop. No document key in the document")
	return "", nil
}

func (scp SkyCPIXProcessor) DecryptContentKey() (ContentKey, error) {
	fmt.Println("Content key is in clear, returning it")
	var contentKey ContentKey
	return contentKey, nil
}

// CPIXPipeline interface implementation
func (scp SkyCPIXProcessor) Handle() error {
	doc, err := scp.UnarshalDocument()
	if err != nil {
		return err
	}
	scp.CpixDoc = *doc
	key, err := scp.DecryptContentKey()
	if err != nil {
		return err
	}
	scp.Key = key
	return nil
}

func (scp SkyCPIXProcessor) Next() CPIXPipeline {
	return SkyHandOff{
		Key: scp.Key,
	}
}

// ---------------Concrete implementation of HandOff interface-------------------
type SkyHandOff struct {
	Key           ContentKey
	KafkaEndpoint string
	KafkaTopic    string
}

func (sho SkyHandOff) PushDownstream() error {
	var err error
	sho.KafkaEndpoint = "kafkaEndpoint" // read from a config store
	sho.KafkaTopic = "Sky"
	fmt.Println("Put the content key on kafka queue, topic 'Sky'")
	if err != nil {
		return errors.Wrap(err, "could not push data downstream")
	}
	return nil
}

// CPIXPipeline interface implementation
func (sho SkyHandOff) Handle() error {
	return sho.PushDownstream()
}

func (sho SkyHandOff) Next() CPIXPipeline {
	return nil
}
