package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vmunzenmayer/p44CheckContainer/models"
	"github.com/xuri/excelize/v2"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("P44_URL")
	method := "GET"
	token := os.Getenv("P44_TOKEN")

	containers, err := excelize.OpenFile("containers.xlsx")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
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

	rows, err := containers.GetRows("Hoja1")
	if err != nil {
		fmt.Println(err)
		return
	}

	rowNum := 1
	result.SetCellValue(sheetName, fmt.Sprint("A", rowNum), "NUMERO ORDEN")
	result.SetCellValue(sheetName, fmt.Sprint("B", rowNum), "CONTENEDOR")
	result.SetCellValue(sheetName, fmt.Sprint("C", rowNum), "NAVIERA")
	result.SetCellValue(sheetName, fmt.Sprint("D", rowNum), "PUERTO ORIGEN")
	result.SetCellValue(sheetName, fmt.Sprint("E", rowNum), "PUERTO DESTINO")
	result.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "¿DATOS EN P44?")

	for _, row := range rows[1:] {
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

		rowNum++
		result.SetCellValue(sheetName, fmt.Sprint("A", rowNum), row[0])
		result.SetCellValue(sheetName, fmt.Sprint("B", rowNum), row[1])
		result.SetCellValue(sheetName, fmt.Sprint("C", rowNum), row[2])
		result.SetCellValue(sheetName, fmt.Sprint("D", rowNum), row[3])
		result.SetCellValue(sheetName, fmt.Sprint("E", rowNum), row[4])

		if len(responseObject.Data) == 0 {
			result.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "NO")
		} else {
			result.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "SI")
		}

		if err := result.SaveAs("result.xlsx"); err != nil {
			fmt.Println(err)
		}
	}
}
