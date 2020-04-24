package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

type AttributeMapRequest map[string]AttributeRequest
type AttributeRequest struct {
	Type          string         `json:"type" enums:"integer,decimal,string,bytes,timestamp,monetary"`
	Value         string         `json:"value"`
	MonetaryValue interface{}    `json:"monetary_value,omitempty"`
}

func TestCreateNFT(t *testing.T) {

	e := utils.GetInsecureClient(t, utils.NODE1)
	e1 := utils.GetInsecureClient(t, utils.NODE2)

	// create document
	payload := map[string]interface{}{
		"scheme":       "generic",
		"data":         map[string]interface{}{},
		"attributes": AttributeMapRequest{
			"Originator": AttributeRequest{
				Type: "bytes",
				Value: "0xdF4c909513fc38e4565593887423704b28CC5a82",
			},
			"AssetValue": AttributeRequest{
				Type: "decimal",
				Value: "100",
			},
			"AssetIdentifier": AttributeRequest{
				Type: "bytes",
				Value: "0x6f39076b2df9d504098928306ceb57df02997645fc51b61ec1bcfb1c498b64f9",
			},
			"MaturityDate": AttributeRequest{
				Type: "timestamp",
				Value: "2020-12-02T15:04:05.999999999Z",
			},
		},
		"write_access": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj := CreateDocument(t, e, utils.Nodes[utils.NODE1].ID, payload, http.StatusAccepted)
	docIdentifier := obj.Value("header").Path("$.document_id").String().NotEmpty().Raw()
	GetDocument(t, e1, utils.Nodes[utils.NODE2].ID, docIdentifier)

	payload = map[string]interface{}{
		"document_id":           docIdentifier,
		"registry_address":      utils.RegistryAddress,
		"deposit_address":       "0x44a0579754d6c94e7bb2c26bfa7394311cc50ccb", // Centrifuge address
		"proof_fields":          defaultProofFields(),
		"asset_manager_address": utils.AssetAddress,
	}

	obj = MintInvoiceUnpaidNFT(t, e, utils.Nodes[utils.NODE1].ID, utils.RegistryAddress, payload)
	doc := GetDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
	assert.True(t, len(doc.Path("$.header.nfts[0].token_id").String().Raw()) > 0, "successful tokenId should have length 77")
	assert.True(t, len(doc.Path("$.header.nfts[0].token_index").String().Raw()) > 0, "successful tokenIndex should have a value")
}

func defaultProofFields() []string {
	originator := "cd_tree.attributes[0xe24e7917d4fcaf79095539ac23af9f6d5c80ea8b0d95c9cd860152bff8fdab17].byte_val"
	assetVal := "cd_tree.attributes[0xcd35852d8705a28d4f83ba46f02ebdf46daf03638b40da74b9371d715976e6dd].byte_val"
	assetID := "cd_tree.attributes[0xbbaa573c53fa357a3b53624eb6deab5f4c758f299cffc2b0b6162400e3ec13ee].byte_val"
	matDate := "cd_tree.attributes[0xe5588a8a267ed4c32962568afe216d4ba70ae60576a611e3ca557b84f1724e29].byte_val"
	signature := "signatures_tree.signatures[0xdf4c909513fc38e4565593887423704b28cc5a82000000000000000000000000f4808bbc8975d7cc026c0da09c6da7928c530e88]"

	return []string{originator, assetVal, assetID, matDate, signature}
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
