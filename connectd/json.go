package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net"
	"net/url"

	"github.com/SUSE/connect-ng/internal/connect"
)

func certToPEM(cert *x509.Certificate) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
}

func errorToJSON(err error) string {
	var s struct {
		ErrType string `json:"err_type"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		// [optional] auxiliary error data
		Data string `json:"data,omitempty"`
	}

	// map Go x509 errors to OpenSSL verify return values
	// see: https://www.openssl.org/docs/man1.0.2/man1/verify.html
	sslErrorMap := map[int]int{
		int(x509.Expired): 10, // X509_V_ERR_CERT_HAS_EXPIRED
		// TODO: add other values as needed
	}

	if ae, ok := err.(connect.APIError); ok {
		s.ErrType = "APIError"
		s.Code = ae.Code
		s.Message = ae.Message
	} else if uerr, ok := err.(*url.Error); ok {
		ierr := errors.Unwrap(err)
		if uerr.Timeout() {
			s.ErrType = "Timeout"
			s.Message = ierr.Error()
		} else if ce, ok := ierr.(x509.CertificateInvalidError); ok {
			s.ErrType = "SSLError"
			s.Message = ierr.Error()
			s.Data = certToPEM(ce.Cert)
			s.Code = sslErrorMap[int(ce.Reason)]
		} else if ce, ok := ierr.(x509.UnknownAuthorityError); ok {
			s.ErrType = "SSLError"
			s.Message = ierr.Error()
			s.Data = certToPEM(ce.Cert)
			// this could be:
			// 18 (X509_V_ERR_DEPTH_ZERO_SELF_SIGNED_CERT),
			// 19 (X509_V_ERR_SELF_SIGNED_CERT_IN_CHAIN) or
			// 20 (X509_V_ERR_UNABLE_TO_GET_ISSUER_CERT_LOCALLY)
			s.Code = 19 // this seems to best match original behavior
		} else if ce, ok := ierr.(x509.HostnameError); ok {
			s.ErrType = "SSLError"
			s.Message = ierr.Error()
			// ruby version doesn't have this but it might be useful
			s.Data = certToPEM(ce.Certificate)
		} else if _, ok := ierr.(*net.OpError); ok {
			s.ErrType = "NetError"
			s.Message = ierr.Error()
		} else {
			connect.Debug.Printf("url.Error: %T: %v", ierr, err)
			s.Message = err.Error()
		}
	} else if je, ok := err.(connect.JSONError); ok {
		s.ErrType = "JSONError"
		s.Message = errors.Unwrap(je).Error()
	} else {
		switch err {
		case connect.ErrMalformedSccCredFile:
			s.ErrType = "MalformedSccCredentialsFile"
		case connect.ErrMissingCredentialsFile:
			s.ErrType = "MissingCredentialsFile"
		}
		connect.Debug.Printf("Error: %T: %v", err, err)
		s.Message = err.Error()
	}

	jsn, _ := json.Marshal(&s)
	return string(jsn)
}