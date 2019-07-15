package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvoiceUnpaidNFT(t *testing.T) {
	// nodes
	e := utils.GetInsecureClient(t, utils.NODE1)

	// create invoice
	currency := "USD"
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"number":        "12324",
			"status":        "unpaid",
			"sender":        utils.Nodes[utils.NODE1].ID,
			"document_type": "invoice",
			"date_due":      "2018-09-26T23:12:37.902198664Z",
			"gross_amount":  "40",
			"currency":      currency,
			"net_amount":    "40",
		},
	}

	obj := CreateDocument(t, utils.INVOICE, e, utils.Nodes[utils.NODE1].ID, payload, http.StatusAccepted)

	docIdentifier := obj.Value("header").Path("$.document_id").String().NotEmpty().Raw()

	doc := GetDocument(t, utils.INVOICE, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
	doc.Path("$.data.currency").String().Equal(currency)

	// mint invoice unpaid NFT
	nftPayload := map[string]interface{}{
		"document_id":     docIdentifier,
		"deposit_address": "0x44a0579754d6c94e7bb2c26bfa7394311cc50ccb", // Centrifuge address
	}
	obj = MintInvoiceUnpaidNFT(t, e, utils.Nodes[utils.NODE1].ID, nftPayload)
	doc = GetDocument(t, utils.INVOICE, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
	assert.True(t, len(doc.Path("$.header.nfts[0].token_id").String().Raw()) > 0, "successful tokenId should have length 77")
	assert.True(t, len(doc.Path("$.header.nfts[0].token_index").String().Raw()) > 0, "successful tokenIndex should have a value")
}

func MintInvoiceUnpaidNFT(t *testing.T, e *httpexpect.Expect, auth string, payload map[string]interface{}) *httpexpect.Object {
	path := fmt.Sprintf("/v1/invoices/%s/mint/unpaid", payload["document_id"])
	method := "POST"
	resp := getResponse(method, path, e, auth, payload).Status(http.StatusAccepted)
	assertOkResponse(t, resp, http.StatusAccepted)
	obj := resp.JSON().Object()
	txID := getTransactionID(t, obj)
	waitTillSuccess(t, e, auth, txID)
	return obj
}
