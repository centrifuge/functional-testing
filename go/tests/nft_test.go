package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func TestCreateNFT(t *testing.T) {

	e := utils.GetInsecureClient(t, utils.NODE1)
	e1 := utils.GetInsecureClient(t, utils.NODE2)

	// create document
	payload := map[string]interface{}{
		"scheme":       "generic",
		"data":         map[string]interface{}{},
		"write_access": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj := CreateDocument(t, e, utils.Nodes[utils.NODE1].ID, payload, http.StatusAccepted)
	docIdentifier := obj.Value("header").Path("$.document_id").String().NotEmpty().Raw()
	GetDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
	GetDocument(t, e1, utils.Nodes[utils.NODE2].ID, docIdentifier)

	payload = map[string]interface{}{
		"document_id":           docIdentifier,
		"deposit_address":       "0x44a0579754d6c94e7bb2c26bfa7394311cc50ccb", // Centrifuge address
		"proof_fields":          []string{"cd_tree" + ".next_version"},
		"asset_manager_address": utils.AssetAddress,
	}

	obj = MintInvoiceUnpaidNFT(t, e, utils.Nodes[utils.NODE1].ID, utils.RegistryAddress, payload)
	doc := GetDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
	assert.True(t, len(doc.Path("$.header.nfts[0].token_id").String().Raw()) > 0, "successful tokenId should have length 77")
	assert.True(t, len(doc.Path("$.header.nfts[0].token_index").String().Raw()) > 0, "successful tokenIndex should have a value")
}

func MintInvoiceUnpaidNFT(t *testing.T, e *httpexpect.Expect, auth, registry string, payload map[string]interface{}) *httpexpect.Object {
	path := fmt.Sprintf("/v1/nfts/registries/%s/mint", registry)
	method := "POST"
	resp := getResponse(method, path, e, auth, payload).Status(http.StatusAccepted)
	assertOkResponse(t, resp, http.StatusAccepted)
	obj := resp.JSON().Object()
	txID := getTransactionID(t, obj)
	waitTillSuccess(t, e, auth, txID)
	return obj
}
