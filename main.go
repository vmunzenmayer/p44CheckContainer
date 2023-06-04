package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vmunzenmayer/p44CheckContainer/models"
	"github.com/xuri/excelize/v2"
)

func main() {
	url := "https://api.clearmetal.com/v1/trips?equipment_numbers="
	method := "GET"
	token := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2ODI2NzQzNTMsIm5iZiI6MTY4MjY3NDM1MywianRpIjoiYzc1ODU3OTItZDY1NC00NzQyLWI3NWUtMDFkOTEwY2U2MDhjIiwiZXhwIjoyMjAxMDc0MzUzLCJpZGVudGl0eSI6eyJ2ZXJzaW9uIjoxLCJ1c2VyIjoiY21wY19hcGlfdXNlckBjbGVhcm1ldGFsLmNvbSIsInVzZXJfaWQiOm51bGwsInRlbmFudCI6ImNtcGMiLCJhY2NvdW50X3R5cGUiOiJhcGlfdXNlciIsInN1Yl90ZW5hbnRzIjpbXSwiY2FuX2ltcGVyc29uYXRlIjpmYWxzZSwiaWdub3JlX3VzZXJfYWNjZXNzX3Jlc3RyaWN0aW9ucyI6ZmFsc2V9LCJmcmVzaCI6ZmFsc2UsInR5cGUiOiJhY2Nlc3MiLCJjc3JmIjoiZDM1YjNmZTQtZmQ3Yi00YWZiLTliZmItZDk2YmUzNWFlMTMyIn0.6oG-0d5G9bgi82qeWSgjsS3lfJh9MOHhqpDX1fHByuI"

	containers, err := excelize.OpenFile("containers.xlsx")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		// Close the spreadsheet.
		if err := containers.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	result := excelize.NewFile()
	sheetName := "Sheet1"

	defer func() {
		if err := result.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the Hoja1.
	rows, err := containers.GetRows("Hoja1")
	if err != nil {
		fmt.Println(err)
		return
	}

	rowNum := 1
	result.SetCellValue(sheetName, fmt.Sprint("A", rowNum), "CONTENEDOR")
	result.SetCellValue(sheetName, fmt.Sprint("B", rowNum), "¿DATOS EN P44?")

	for _, row := range rows[1:] {
		//fmt.Println(row[1], "\t")
		urlContainer := url + row[1]
		client := &http.Client{}
		req, err := http.NewRequest(method, urlContainer, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Authorization", token)

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var responseObject models.Response
		json.Unmarshal(body, &responseObject)
		//fmt.Println(responseObject.Data)

		rowNum++
		// Set value of a cell.
		result.SetCellValue(sheetName, fmt.Sprint("A", rowNum), row[1])

		if len(responseObject.Data) == 0 {
			result.SetCellValue(sheetName, fmt.Sprint("B", rowNum), "NO")
		} else {
			result.SetCellValue(sheetName, fmt.Sprint("B", rowNum), "SI")
		}
	}

	// Save spreadsheet by the given path.
	if err := result.SaveAs("result.xlsx"); err != nil {
		fmt.Println(err)
	}
}
