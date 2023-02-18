package main

// CPIXPipeline moves the data through the pipeline
type CPIXPipeline interface {
	Handle() error
	Next() CPIXPipeline
}

type X509CertFetcher interface {
	FetchCertificate() (X509Certificate, error)
}

type HTTPClient interface {
	CreateRequest() ([]byte, error)
	ExecuteRequestResponse(request []byte) (*CPIXDocument, error)
}

type CPIXProcessor interface {
	UnmarshalDocument() (*CPIXDocument, error)
	ValidateDeliveryKey() bool
	DecryptDocumentKey() (DocumentKey, error)
	DecryptContentKey() (ContentKey, error)
}

type Signer interface {
	ValidateSignature() bool
	SignDocument([]byte) []byte
}
