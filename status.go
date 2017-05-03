package gohttp

import (
	"net/http"
)

const (
	StatusContinue           = http.StatusContinue           // 100  RFC 7231, 6.2.1
	StatusSwitchingProtocols = http.StatusSwitchingProtocols // 101  RFC 7231, 6.2.2
	StatusProcessing         = http.StatusProcessing         // 102  RFC 2518, 10.1

	StatusOK                   = http.StatusOK                   // 200  RFC 7231, 6.3.1
	StatusCreated              = http.StatusCreated              // 201  RFC 7231, 6.3.2
	StatusAccepted             = http.StatusAccepted             // 202  RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = http.StatusNonAuthoritativeInfo // 203  RFC 7231, 6.3.4
	StatusNoContent            = http.StatusNoContent            // 204  RFC 7231, 6.3.5
	StatusResetContent         = http.StatusResetContent         // 205  RFC 7231, 6.3.6
	StatusPartialContent       = http.StatusPartialContent       // 206  RFC 7233, 4.1
	StatusMultiStatus          = http.StatusMultiStatus          // 207  RFC 4918, 11.1
	StatusAlreadyReported      = http.StatusAlreadyReported      // 208  RFC 5842, 7.1
	StatusIMUsed               = http.StatusIMUsed               // 226  RFC 3229, 10.4.1

	StatusMultipleChoices  = http.StatusMultipleChoices  // 300  RFC 7231, 6.4.1
	StatusMovedPermanently = http.StatusMovedPermanently // 301  RFC 7231, 6.4.2
	StatusFound            = http.StatusFound            // 302  RFC 7231, 6.4.3
	StatusSeeOther         = http.StatusSeeOther         // 303  RFC 7231, 6.4.4
	StatusNotModified      = http.StatusNotModified      // 304  RFC 7232, 4.1
	StatusUseProxy         = http.StatusUseProxy         // 305  RFC 7231, 6.4.5

	StatusTemporaryRedirect = http.StatusTemporaryRedirect // 307  RFC 7231, 6.4.7
	StatusPermanentRedirect = http.StatusPermanentRedirect // 308  RFC 7538, 3

	StatusBadRequest                   = http.StatusBadRequest                   // 400  RFC 7231, 6.5.1
	StatusUnauthorized                 = http.StatusUnauthorized                 // 401  RFC 7235, 3.1
	StatusPaymentRequired              = http.StatusPaymentRequired              // 402  RFC 7231, 6.5.2
	StatusForbidden                    = http.StatusForbidden                    // 403  RFC 7231, 6.5.3
	StatusNotFound                     = http.StatusNotFound                     // 404  RFC 7231, 6.5.4
	StatusMethodNotAllowed             = http.StatusMethodNotAllowed             // 405  RFC 7231, 6.5.5
	StatusNotAcceptable                = http.StatusNotAcceptable                // 406  RFC 7231, 6.5.6
	StatusProxyAuthRequired            = http.StatusProxyAuthRequired            // 407  RFC 7235, 3.2
	StatusRequestTimeout               = http.StatusRequestTimeout               // 408  RFC 7231, 6.5.7
	StatusConflict                     = http.StatusConflict                     // 409  RFC 7231, 6.5.8
	StatusGone                         = http.StatusGone                         // 410  RFC 7231, 6.5.9
	StatusLengthRequired               = http.StatusLengthRequired               // 411  RFC 7231, 6.5.10
	StatusPreconditionFailed           = http.StatusPreconditionFailed           // 412  RFC 7232, 4.2
	StatusRequestEntityTooLarge        = http.StatusRequestEntityTooLarge        // 413  RFC 7231, 6.5.11
	StatusRequestURITooLong            = http.StatusRequestURITooLong            // 414  RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = http.StatusUnsupportedMediaType         // 415  RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = http.StatusRequestedRangeNotSatisfiable // 416  RFC 7233, 4.4
	StatusExpectationFailed            = http.StatusExpectationFailed            // 417  RFC 7231, 6.5.14
	StatusTeapot                       = http.StatusTeapot                       // 418  RFC 7168, 2.3.3
	StatusUnprocessableEntity          = http.StatusUnprocessableEntity          // 422  RFC 4918, 11.2
	StatusLocked                       = http.StatusLocked                       // 423  RFC 4918, 11.3
	StatusFailedDependency             = http.StatusFailedDependency             // 424  RFC 4918, 11.4
	StatusUpgradeRequired              = http.StatusUpgradeRequired              // 426  RFC 7231, 6.5.15
	StatusPreconditionRequired         = http.StatusPreconditionRequired         // 428  RFC 6585, 3
	StatusTooManyRequests              = http.StatusTooManyRequests              // 429  RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = http.StatusRequestHeaderFieldsTooLarge  // 431  RFC 6585, 5
	StatusUnavailableForLegalReasons   = http.StatusUnavailableForLegalReasons   // 451  RFC 7725, 3

	StatusInternalServerError           = http.StatusInternalServerError           // 500  RFC 7231, 6.6.1
	StatusNotImplemented                = http.StatusNotImplemented                // 501  RFC 7231, 6.6.2
	StatusBadGateway                    = http.StatusBadGateway                    // 502  RFC 7231, 6.6.3
	StatusServiceUnavailable            = http.StatusServiceUnavailable            // 503  RFC 7231, 6.6.4
	StatusGatewayTimeout                = http.StatusGatewayTimeout                // 504  RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = http.StatusHTTPVersionNotSupported       // 505  RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = http.StatusVariantAlsoNegotiates         // 506  RFC 2295, 8.1
	StatusInsufficientStorage           = http.StatusInsufficientStorage           // 507  RFC 4918, 11.5
	StatusLoopDetected                  = http.StatusLoopDetected                  // 508  RFC 5842, 7.2
	StatusNotExtended                   = http.StatusNotExtended                   // 510  RFC 2774, 7
	StatusNetworkAuthenticationRequired = http.StatusNetworkAuthenticationRequired // 511 RFC 6585, 6

    StatusPrepareError  = 600 //define by me



)
