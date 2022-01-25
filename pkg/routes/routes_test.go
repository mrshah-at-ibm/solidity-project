package routes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/zap"

	. "github.com/onsi/gomega"

	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer/executerfakes"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/routes"
)

var _ = Describe("Routes", func() {

	var exec executerfakes.FakeExecuterInterface
	var rts *routes.Routes

	var req *http.Request
	var rr *httptest.ResponseRecorder
	var err error
	var handler http.HandlerFunc

	BeforeEach(func() {
		logger, err := zap.NewDevelopment()
		Expect(err).To(BeNil())

		defer logger.Sync()

		exec = executerfakes.FakeExecuterInterface{}
		rts, err = routes.New(logger, &exec)
		Expect(err).To(BeNil())
	})

	Describe("mintTransaction", func() {

		BeforeEach(func() {
			req, err = http.NewRequest("GET", "/mint/0x01", nil)
			Expect(err).To(BeNil())

			rr = httptest.NewRecorder()
			handler = http.HandlerFunc(rts.MintTransactionRoute)
			exec.MintTokenReturns(types.NewReceipt([]byte{}, false, 1), nil)
		})
		It("should return 200 on successful call", func() {
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
		It("should return 500 if error in mint", func() {
			exec.MintTokenReturns(nil, errors.New("Fake error"))
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("Fake error"))

		})
	})
	Describe("BurnToken", func() {
		BeforeEach(func() {
			req, err = http.NewRequest("GET", "/token/0x01/burn", nil)
			Expect(err).To(BeNil())

			rr = httptest.NewRecorder()
			handler = http.HandlerFunc(rts.BurnTokenRoute)

		})
		It("should return 200 on successful call", func() {
			exec.BurnTokenReturns(types.NewReceipt([]byte{}, false, 1), nil)
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
		It("should return 500 if error in mint", func() {
			exec.BurnTokenReturns(nil, errors.New("Fake error"))
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("Fake error"))

		})
	})

	Describe("TransferToken", func() {
		BeforeEach(func() {

			bodybytes := strings.NewReader("{\"from\": \"0x1\", \"to\": \"0x2\"}")
			req, err = http.NewRequest("POST", "/token/0x01/transfer", bodybytes)
			Expect(err).To(BeNil())

			rr = httptest.NewRecorder()
			handler = http.HandlerFunc(rts.TransferTokenRoute)

		})
		It("should return 200 on successful call", func() {
			exec.TransferTokenReturns(types.NewReceipt([]byte{}, false, 1), nil)
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
		It("should return 500 if error in mint", func() {
			exec.TransferTokenReturns(nil, errors.New("Fake error"))
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			Expect(rr.Body.String()).To(ContainSubstring("Fake error"))
		})
	})

})
