package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cast"
)

const (
	botToken        = ""
	chatID          = ""
	baseUrl         = "https://api.admin.u-code.io"
	logFunctionName = "ucode-template"
	IsHTTP          = true // if this is true banchmark test works.
)

/*
Answer below questions before starting the function.

When the function invoked?
 - cleints -> AFTER -> CREATE
What does it do?
- Explain the purpose of the function.(O'zbekcha yozilsa ham bo'ladi.)
bemorlar doktorlarni oldiga бесплатная консультация olish uchun borganida,
 adminkada client ning бесплатная консультация si sonini kamaytirish uchun(-1)
*/

// Request structures
type (
	// Handle request body
	NewRequestBody struct {
		RequestData HttpRequest `json:"request_data"`
		Auth        AuthData    `json:"auth"`
		Data        Data        `json:"data"`
	}

	HttpRequest struct {
		Method  string      `json:"method"`
		Path    string      `json:"path"`
		Headers http.Header `json:"headers"`
		Params  url.Values  `json:"params"`
		Body    []byte      `json:"body"`
	}

	AuthData struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	// Function request body >>>>> GET_LIST, GET_LIST_SLIM, CREATE, UPDATE
	Request struct {
		Data map[string]interface{} `json:"data"`
	}

	// most common request structure -> UPDATE, MULTIPLE_UPDATE, CREATE, DELETE
	Data struct {
		AppId      string                 `json:"app_id"`
		Method     string                 `json:"method"`
		ObjectData map[string]interface{} `json:"object_data"`
		ObjectIds  []string               `json:"object_ids"`
		TableSlug  string                 `json:"table_slug"`
		UserId     string                 `json:"user_id"`
	}

	FunctionRequest struct {
		BaseUrl     string  `json:"base_url"`
		TableSlug   string  `json:"table_slug"`
		AppId       string  `json:"app_id"`
		Request     Request `json:"request"`
		DisableFaas bool    `json:"disable_faas"`
	}
	GetListFunctionRequest struct {
		BaseUrl     string                 `json:"base_url"`
		TableSlug   string                 `json:"table_slug"`
		AppId       string                 `json:"app_id"`
		Request     map[string]interface{} `json:"request"`
		DisableFaas bool                   `json:"disable_faas"`
	}
)

// Response structures
type (
	// Create function response body >>>>> CREATE
	Datas struct {
		Data struct {
			Data struct {
				Data map[string]interface{} `json:"data"`
			} `json:"data"`
		} `json:"data"`
	}

	// ClientApiResponse This is get single api response >>>>> GET_SINGLE_BY_ID, GET_SLIM_BY_ID
	ClientApiResponse struct {
		Data ClientApiData `json:"data"`
	}

	ClientApiData struct {
		Data ClientApiResp `json:"data"`
	}

	ClientApiResp struct {
		Response map[string]interface{} `json:"response"`
	}

	Response struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	// GetListClientApiResponse This is get list api response >>>>> GET_LIST, GET_LIST_SLIM
	GetListClientApiResponse struct {
		Data GetListClientApiData `json:"data"`
	}

	GetListClientApiData struct {
		Data GetListClientApiResp `json:"data"`
	}

	GetListClientApiResp struct {
		Response []map[string]interface{} `json:"response"`
	}

	// ClientApiUpdateResponse This is single update api response >>>>> UPDATE
	ClientApiUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			TableSlug string                 `json:"table_slug"`
			Data      map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	// ClientApiMultipleUpdateResponse This is multiple update api response >>>>> MULTIPLE_UPDATE
	ClientApiMultipleUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Data struct {
				Objects []map[string]interface{} `json:"objects"`
			} `json:"data"`
		} `json:"data"`
	}

	ResponseStatus struct {
		Status string `json:"status"`
	}
)

// Testing types
type (
	Asserts struct {
		Request  NewRequestBody
		Response Response
	}

	FunctionAssert struct{}
)

func (f FunctionAssert) GetAsserts() []Asserts {
	// var appId = os.Getenv("APP_ID")
	return []Asserts{
		{
			Request: NewRequestBody{
				Data: Data{
					AppId: "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS",
					ObjectData: map[string]interface{}{
						"cleints_id":        "accd5dcf-c6a9-4b49-a5f8-7474f01dcae6",
						"consultation_type": "платная консультация",
					},
				},
			},
			Response: Response{
				Status: "done",
			},
		},
	}
}

func (f FunctionAssert) GetBenchmarkRequest() Asserts {
	// var appId = os.Getenv("APP_ID")
	return Asserts{
		Request: NewRequestBody{
			Data: Data{
				AppId: "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS",
				ObjectData: map[string]interface{}{
					"cleints_id":        "accd5dcf-c6a9-4b49-a5f8-7474f01dcae6",
					"consultation_type": "платная консультация",
				}},
		},
		Response: Response{
			Status: "done",
		},
	}
}

// func main() {
// 	body := `
// 	{
// 		"data":{
// 			"app_id":"P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS",
// 			"object_data":{
// 				"cleints_id":"93a01c6d-7f29-41b8-b706-5f3a59af84fa"
// 			}
// 		}
// 	}
// 	`
// 	fmt.Println(Handle([]byte(body)))
// }

// Handle a serverless request
func Handle(req []byte) string {
	var (
		response Response
		request  NewRequestBody
	)

	// defer func() {
	// 	responseByte, _ := json.Marshal(response)
	// 	Send(string(responseByte))
	// }()

	err := json.Unmarshal(req, &request)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling request", "error": err.Error()}
		response.Status = "error"
		responseByte, _ := json.Marshal(response)
		return string(responseByte)
	}

	if request.Data.AppId == "" {
		response.Data = map[string]interface{}{"message": "App id required"}
		response.Status = "error"
		responseByte, _ := json.Marshal(response)
		return string(responseByte)
	}

	if request.Data.ObjectData["consultation_type"] != "бесплатная консультация" {
		response.Data = map[string]interface{}{"message": "paid-consultation"}
		response.Data = map[string]interface{}{}
		response.Status = "done" //if all will be ok else "error"
		responseByte, _ := json.Marshal(response)

		return string(responseByte)
	}

	getObjectRequest := map[string]interface{}{
		"cleints_id": cast.ToString(request.Data.ObjectData["cleints_id"]),
		"status":     []string{"активный"},
	}

	res, response, err := GetListSlimObject(GetListFunctionRequest{
		BaseUrl:     baseUrl,
		TableSlug:   "subscription_report",
		AppId:       request.Data.AppId,
		Request:     getObjectRequest,
		DisableFaas: true,
	})
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while GetListSlimObject", "error": err.Error()}
		response.Status = "error"
		responseByte, _ := json.Marshal(response)
		return string(responseByte)
	}

	if len(res.Data.Data.Response) > 0 {
		var (
			oldestTime time.Time
			oldestData map[string]interface{}
		)
		for _, v := range res.Data.Data.Response {
			if consultationCount, ok := v["consultation_count"].(float64); ok && consultationCount > 0 {
				if endDateStr, ok := v["end_date"].(string); ok {
					endDate, err := time.Parse(time.RFC3339, endDateStr)
					if err != nil {
						continue
					}
					if oldestTime.IsZero() || endDate.Before(oldestTime) {
						oldestTime = endDate
						oldestData = v
					}
				}
			}
		}
		if oldestData == nil {
			response.Data = map[string]interface{}{"message": "Consultation Not Found", "error": err}
			response.Status = "error"
			responseByte, _ := json.Marshal(response)
			return string(responseByte)
		}
		updateRequest := Request{
			Data: map[string]interface{}{
				"guid":               oldestData["guid"],
				"consultation_count": cast.ToInt(oldestData["consultation_count"]) - 1,
			},
		}
		res2, response2, err2 := UpdateObject(FunctionRequest{
			BaseUrl:     baseUrl,
			TableSlug:   "subscription_report",
			AppId:       request.Data.AppId,
			Request:     updateRequest,
			DisableFaas: true,
		})
		if err2 != nil {
			response2.Data = map[string]interface{}{"message": "Error while UpdateObject", "error": err.Error()}
			response2.Status = "error"
			responseByte, _ := json.Marshal(response2)
			return string(responseByte)
		}
		response2.Data = res2.Data.Data
		response2.Data = map[string]interface{}{}
		response2.Status = "done" //if all will be ok else "error"
		responseByte, _ := json.Marshal(response2)

		return string(responseByte)

	} else {
		response.Data = map[string]interface{}{"message": "free-consultation not found"}
		response.Status = "error" //if all will be ok else "error"
		responseByte, _ := json.Marshal(response)

		return string(responseByte)
	}
}
