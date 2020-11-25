package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"project/regression"
	"strconv"
)

var r regression.Regression

type PartialData struct {
	G1 string `json:"g1"`
	G2 string `json:"g2"`
}

type FullData struct {
	G1 string `json:"g1"`
	G2 string `json:"g2"`
	G3 string `json:"g3"`
}

func main () {
	router := mux.NewRouter()

	router.HandleFunc("/predictGrade", predictGrade).Methods("POST")
	//feed
	//formula

	http.ListenAndServe(":3000", router)
}

func predictGrade(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var partialData PartialData
	_ = json.NewDecoder(request.Body).Decode(&partialData)
	var fullData FullData
	parsedG1, _ := strconv.ParseFloat(partialData.G1, 64)
	parsedG2, _ := strconv.ParseFloat(partialData.G2, 64)
	fullData.G1 = partialData.G1
	fullData.G2 = partialData.G2
	outputG3String := fmt.Sprintf("%f", r.TwoVariableGradePrediction(parsedG1, parsedG2)+10.0)
	fullData.G3 = outputG3String
	fmt.Println(fullData)
	json.NewEncoder(writer).Encode(fullData)
}

func init() {

	nodehosts := []string {
		"192.168.0.12:8000",
		//"192.168.0.4:8001",
		//"192.168.0.4:8002",
	}


	// abrimos el documento 
	f, err := os.Open("students-mat.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// se crea un nuevo lecto 
	// establecemos el numero de columas del documento
	salesData := csv.NewReader(f)
	salesData.FieldsPerRecord = 3 // 2 columas G2 Y G3

	// se lee la data del docuemtno
	records, err := salesData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//modelaremos segun el valor G3

	r.NodesDir = nodehosts
	r.SetObserved("G3")
	r.SetVar(0, "G1")
	r.SetVar(1, "G2")

	// BUCLE DE LOS DATOS
	for i, record := range records {
		// excluimos el encabezado
		if i == 0 {
			continue
		}

		// seteamos el primer valor del arreglo 
		G3, err := strconv.ParseFloat(records[i][0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// seteamos el segundo valor del arreglo
		G1, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		G2, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// insertamos los valores del arreglo
		go r.Train(regression.DataPoint(G3, []float64{G2, G1}))
	}

	// modelo de regresion
	r.Run()
	// imprimios la formula retornada
	fmt.Printf("Regresión Formula:\n%v\n\n", r.Formula)
	//fmt.Printf("Prediction Formula:\n%v\n\n", r.CalcPrediccion())
}
