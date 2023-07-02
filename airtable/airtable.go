package airtable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// AirtableRecord represents a record in an Airtable table
type AirtableRecord struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// AirtableResponse represents the response from the Airtable API
type AirtableResponse struct {
	Records []AirtableRecord `json:"records"`
}

type Textfile struct {
	Name string
	Body string
}

func FetchFromAirtable(apiKey, baseID, tableName, nameID, bodyID string) ([]Textfile, error) {
	url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s", baseID, tableName)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data AirtableResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Now data.Records contains all records from the Airtable table
	var result []Textfile
	for _, record := range data.Records {
		if record.Fields[bodyID] != nil {
			result = append(result, Textfile{
				Name: fmt.Sprint(record.Fields[nameID]),
				Body: fmt.Sprint(record.Fields[bodyID]),
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}
