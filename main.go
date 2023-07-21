package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/vmunzenmayer/p44CheckContainer/models"
	"github.com/xuri/excelize/v2"
)

var result = []ContainerExcel{}
var resultP44 = []ContainerExcel{}

type ContainerExcel struct {
	NumeroOrdenVenta string
	Contenedor       string
	Naviera          string
	OrigenPuerto     string
	DestinoPuerto    string
	DatosEnP44       bool
}

func fetchContainersFromExcel(filename string) []ContainerExcel {
	//result := []ContainerExcel{}
	containers, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := containers.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	rows, err := containers.GetRows("Hoja1")
	if err != nil {
		log.Fatalln(err)
	}

	for _, row := range rows[1:] {
		var container ContainerExcel
		container.NumeroOrdenVenta = row[0]
		container.Contenedor = row[1]
		container.Naviera = row[2]
		container.OrigenPuerto = row[3]
		container.DestinoPuerto = row[4]
		container.DatosEnP44 = false

		result = append(result, container)
	}

	return result
}

func processContainer(container ContainerExcel) {
	url := os.Getenv("P44_URL")
	method := "GET"
	token := os.Getenv("P44_TOKEN")

	urlContainer := url + container.Contenedor
	client := &http.Client{}
	req, err := http.NewRequest(method, urlContainer, nil)

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseObject models.Response
	json.Unmarshal(body, &responseObject)

	if len(responseObject.Data) > 0 {
		container.DatosEnP44 = true
	}

	resultP44 = append(resultP44, container)
}

func worker(jobs <-chan ContainerExcel, wg *sync.WaitGroup) {
	defer wg.Done()

	for container := range jobs {
		processContainer(container)
	}
}

func main() {
	start := time.Now()
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))

	if err != nil {
		log.Fatalln("The NUM_WORKERS variable is required in the .env file")
	}

	if numWorkers < 1 {
		log.Fatalln("The NUM_WORKERS variable must be greater than 0")
	}

	containers := fetchContainersFromExcel("containers.xlsx")

	if len(containers) < 1 {
		log.Fatalln("There are no containers in Excel file")
	}

	jobs := make(chan ContainerExcel, len(containers))
	var wg sync.WaitGroup

	// Inicialización de las goroutines (workers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// Enviar trabajos a las goroutines
	for _, container := range containers {
		jobs <- container
	}
	close(jobs)

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	excelResult := excelize.NewFile()
	sheetName := "Sheet1"

	defer func() {
		if err := excelResult.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rowNum := 1
	excelResult.SetCellValue(sheetName, fmt.Sprint("A", rowNum), "NUMERO ORDEN")
	excelResult.SetCellValue(sheetName, fmt.Sprint("B", rowNum), "CONTENEDOR")
	excelResult.SetCellValue(sheetName, fmt.Sprint("C", rowNum), "NAVIERA")
	excelResult.SetCellValue(sheetName, fmt.Sprint("D", rowNum), "PUERTO ORIGEN")
	excelResult.SetCellValue(sheetName, fmt.Sprint("E", rowNum), "PUERTO DESTINO")
	excelResult.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "¿DATOS EN P44?")

	for _, row := range resultP44 {

		rowNum++
		excelResult.SetCellValue(sheetName, fmt.Sprint("A", rowNum), row.NumeroOrdenVenta)
		excelResult.SetCellValue(sheetName, fmt.Sprint("B", rowNum), row.Contenedor)
		excelResult.SetCellValue(sheetName, fmt.Sprint("C", rowNum), row.Naviera)
		excelResult.SetCellValue(sheetName, fmt.Sprint("D", rowNum), row.OrigenPuerto)
		excelResult.SetCellValue(sheetName, fmt.Sprint("E", rowNum), row.DestinoPuerto)

		if row.DatosEnP44 {
			excelResult.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "SI")
		} else {
			excelResult.SetCellValue(sheetName, fmt.Sprint("F", rowNum), "NO")
		}
	}

	if err := excelResult.SaveAs("result.xlsx"); err != nil {
		fmt.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("P44 Check Containers took %s. %d containers has been processed", elapsed, len(containers))
}
