package es

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type IRoundTripChecker interface {
	Before(req *http.Request)
	After(req *http.Request, resp *http.Response, reqErr error) // should not close body
}

type TransportCheckRoundTripper struct {
	Transport      http.RoundTripper
	NewCheckerFunc func() IRoundTripChecker
}

func NewTransportCheckRoundTripper(
	transport http.RoundTripper,
	newCheckerFunc func() IRoundTripChecker,
) *TransportCheckRoundTripper {
	return &TransportCheckRoundTripper{
		Transport:      transport,
		NewCheckerFunc: newCheckerFunc,
	}
}

// will not close body
func (m *TransportCheckRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if m.NewCheckerFunc != nil {
		var checker IRoundTripChecker
		func() {
			defer func() {
				if iRec := recover(); iRec != nil {
					logrus.Errorf("check before round trip panic: %s.", err)
				}
			}()
			checker = m.NewCheckerFunc()
			if checker == nil {
				return
			}

			checker.Before(req)
		}()

		defer func() {
			defer func() {
				if iRec := recover(); iRec != nil {
					logrus.Errorf("check after round trip panic: %s", iRec)
				}
			}()

			if checker != nil {
				checker.After(req, resp, err)
			}
		}()

	}

	resp, err = m.Transport.RoundTrip(req)
	if nil != err {
		logrus.Errorf("transport check round trip failed. error: %s.", err)
		return
	}

	return
}
