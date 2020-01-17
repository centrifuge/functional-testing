package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/centrifuge/functional-testing/go/utils"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndUpdateDocumentFromOrigin(t *testing.T) {
	// nodes
	e := utils.GetInsecureClient(t, utils.NODE1)
	e1 := utils.GetInsecureClient(t, utils.NODE2)

	// create invoice
	payload := map[string]interface{}{
		"scheme":       "generic",
		"data":         map[string]interface{}{},
		"write_access": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj := CreateDocument(t, e, utils.Nodes[utils.NODE1].ID, payload, http.StatusAccepted)

	docIdentifier := obj.Value("header").Path("$.document_id").String().NotEmpty().Raw()

	// update invoice
	payload = map[string]interface{}{
		"scheme":       "generic",
		"data":         map[string]interface{}{},
		"write_access": []string{utils.Nodes[utils.NODE2].ID},
	}

	obj = UpdateDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier, payload, http.StatusAccepted)
	GetDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier)

	// Receiver has document
	GetDocument(t, e1, utils.Nodes[utils.NODE2].ID, docIdentifier)
}

func TestCreateAndUpdateDocumentFromCollaborator(t *testing.T) {
	// nodes
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

	// update document
	payload = map[string]interface{}{
		"scheme":       "generic",
		"data":         map[string]interface{}{},
		"write_access": []string{utils.Nodes[utils.NODE1].ID},
	}

	obj = UpdateDocument(t, e1, utils.Nodes[utils.NODE2].ID, docIdentifier, payload, http.StatusAccepted)
	GetDocument(t, e1, utils.Nodes[utils.NODE2].ID, docIdentifier)

	// Receiver has document
	GetDocument(t, e, utils.Nodes[utils.NODE1].ID, docIdentifier)
}

func GetDocument(t *testing.T, e *httpexpect.Expect, auth string, docIdentifier string) *httpexpect.Value {
	objGet := utils.AddCommonHeaders(e.GET(fmt.Sprintf("/v1/documents/%s", docIdentifier)), auth).
		Expect().Status(http.StatusOK)
	assertOkResponse(t, objGet, http.StatusOK)
	objGet.JSON().Path("$.header.document_id").String().Equal(docIdentifier)
	return objGet.JSON()
}

func CreateDocument(t *testing.T, e *httpexpect.Expect, auth string, payload map[string]interface{}, status int) *httpexpect.Object {
	path := fmt.Sprintf("/v1/%s", "documents")
	method := "POST"
	resp := getResponse(method, path, e, auth, payload).Status(status)
	assertOkResponse(t, resp, status)
	obj := resp.JSON().Object()
	txID := getTransactionID(t, obj)
	waitTillSuccess(t, e, auth, txID)
	return obj
}

func UpdateDocument(t *testing.T, e *httpexpect.Expect, auth string, documentID string, payload map[string]interface{}, status int) *httpexpect.Object {
	path := fmt.Sprintf("/v1/documents/%s", documentID)
	method := "PUT"
	resp := getResponse(method, path, e, auth, payload).Status(status)
	assertOkResponse(t, resp, status)
	obj := resp.JSON().Object()
	txID := getTransactionID(t, obj)
	waitTillSuccess(t, e, auth, txID)
	return obj
}

func getResponse(method, path string, e *httpexpect.Expect, auth string, payload map[string]interface{}) *httpexpect.Response {
	return utils.AddCommonHeaders(e.Request(method, path), auth).
		WithJSON(payload).
		Expect()
}

func assertOkResponse(t *testing.T, response *httpexpect.Response, status int) {
	if response.Raw().StatusCode != status {
		assert.Fail(t, "Response Payload: ", response.Body().Raw())
	}
}

func getTransactionID(t *testing.T, resp *httpexpect.Object) string {
	txID := resp.Value("header").Path("$.job_id").String().Raw()
	if txID == "" {
		t.Error("transaction ID empty")
	}

	return txID
}

func waitTillSuccess(t *testing.T, e *httpexpect.Expect, auth string, txID string) {
	for {
		resp := utils.AddCommonHeaders(e.GET("/v1/jobs/"+txID), auth).Expect().Status(200).JSON().Object()
		status := resp.Path("$.status").String().Raw()
		if status == "pending" {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if status == "failed" {
			t.Error(resp.Path("$.message").String().Raw())
		}

		break
	}
}
