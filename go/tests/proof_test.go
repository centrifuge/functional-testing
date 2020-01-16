package tests

import (
	"net/http"
	"testing"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
)

func TestProofGenerationWithMultipleFields(t *testing.T) {
	e := utils.GetInsecureClient(t, utils.NODE1)

	payload := map[string]interface{}{
		"scheme": "generic",
		"data":   map[string]interface{}{},
	}

	obj := CreateDocument(t, e, utils.Nodes[utils.NODE1].ID, payload, http.StatusAccepted)

	docIdentifier := obj.Value("header").Path("$.document_id").String().NotEmpty().Raw()

	proofPayload := map[string]interface{}{
		"fields": []string{"generic.scheme", "cd_tree.document_identifier"},
	}

	objProof := GetProof(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier, proofPayload)
	objProof.Path("$.header.document_id").String().Equal(docIdentifier)
	objProof.Path("$.field_proofs[0].property").String().Equal("0x0005000000000001") // generic.scheme
	objProof.Path("$.field_proofs[0].sorted_hashes").NotNull()
	objProof.Path("$.field_proofs[1].property").String().Equal("0x0100000000000009") // cd_tree.document_identifier
	objProof.Path("$.field_proofs[1].sorted_hashes").NotNull()
}

func GetProof(t *testing.T, e *httpexpect.Expect, auth string, documentID string, payload map[string]interface{}) *httpexpect.Object {
	obj := utils.AddCommonHeaders(e.POST("/v1/documents/"+documentID+"/proofs"), auth).
		WithJSON(payload).
		Expect().Status(http.StatusOK)
	assertOkResponse(t, obj, http.StatusOK)
	return obj.JSON().Object()
}
