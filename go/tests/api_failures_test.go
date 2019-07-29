package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func assertFailResponse(t *testing.T, response *httpexpect.Response) {
	if response.Raw().StatusCode == http.StatusOK {
		assert.Fail(t, "Response Payload: ", response.Body().Raw())
	}
}

func testWrongFormat(t *testing.T, docType string) {
	c := utils.GetInsecureClient(t, utils.NODE1)

	path := fmt.Sprintf("/v1/%s", docType)
	method := "POST"

	// nil create payload
	resp := getResponse(method, path, c, utils.Nodes[utils.NODE1].ID, nil)
	assertFailResponse(t, resp)

	// wrong api format
	p := map[string]interface{}{
		"document": map[string]interface{}{"data": map[string]interface{}{"currency": "EUR"}},
	}
	resp = getResponse(method, path, c, utils.Nodes[utils.NODE1].ID, p)
	assertFailResponse(t, resp)

	path = fmt.Sprintf("/v1/%s/%s", docType, "")
	method = "PUT"

	// nil update
	resp = getResponse(method, path, c, utils.Nodes[utils.NODE1].ID, nil)
	assertFailResponse(t, resp)

	// wrong format
	resp = getResponse(method, path, c, utils.Nodes[utils.NODE1].ID, p)
	assertFailResponse(t, resp)
}

func TestEmptyAndWrongInvoiceFormat(t *testing.T) {
	testWrongFormat(t, utils.INVOICE)
}
