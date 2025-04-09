package function

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestHandlerI interface {
	GetAsserts() []Asserts
	GetBenchmarkRequest() Asserts
}

func NewAssert(f FunctionAssert) TestHandlerI {
	return f
}

func TestHandle(t *testing.T) {
	var Request = NewRequestBody{
		Data: Data{
			AppId: "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS",
			ObjectData: map[string]interface{}{
				"cleints_id": "9775bce2-9995-454f-bfa3-08d1b75ad2ec",
				"consultation_type":"бесплатная консультация",
			}},
	}

	req, err := json.Marshal(Request)
	if err != nil {
		t.Errorf("Error on marshal request::: %s", err.Error())
	}
	got := Handle(req)
	var resp Response
	err = json.Unmarshal([]byte(got), &resp)

	if err != nil {
		t.Errorf("Error on unmarshal response::: %s", err.Error())
	}
	fmt.Println(got)
	if resp.Status != "done" && resp.Data["message"]!="free-consultation not found" {
		t.Errorf("Failed in refund order item ::: Message >>>> %s ", resp.Data["message"].(string))
	}

	if err != nil {
		t.Errorf("Error on unmarshal response::: %s", err.Error())
	}
	fmt.Println(got)
	if resp.Status != "done" && resp.Data["message"]!="free-consultation not found"{
		t.Errorf("Failed in refund order item ::: Message >>>> %s ", resp.Data["message"].(string))
	}

}

func BenchmarkHandler(b *testing.B) {
	if !IsHTTP {
		return
	}
	a := NewAssert(FunctionAssert{})
	var start time.Time

	for i := 0; i < b.N; i++ {
		reqByte, err := json.Marshal(a.GetBenchmarkRequest().Request)
		assert.Nil(b, err)

		start = time.Now()

		response := Handle(reqByte)

		resStatus, err := ConvertResponse([]byte(response))
		
		assert.Nil(b, err)
		assert.Equal(b, "done", resStatus.Status)

		if time.Since(start) > time.Millisecond*5000 {
			assert.Nil(b, fmt.Errorf("took more time than %d ms: %v", 500, time.Since(start)))
		}
	}
}
