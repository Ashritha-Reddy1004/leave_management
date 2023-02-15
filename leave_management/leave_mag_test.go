package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetAdminCred(t *testing.T) {
	var jsonStr = []byte(`{"adminname":"Ashritha","password":"12345678"}`)
	req, err := http.NewRequest("POST", "/setadmincredentials", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Contents", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetAdminCred)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := map[string]string{"adminname": "Ashritha", "password": "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"}

	var got map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Errorf("Cannot unmarshal resp to interafce, err=%v", err)
	}

	if strings.Compare(fmt.Sprintf("%v", got["data"].([]interface{})[0]), fmt.Sprintf("%v", expected)) != 0 {
		t.Errorf("handler returned unexpected body: got %s want %v", rr.Body.String(), expected)
	}

}
