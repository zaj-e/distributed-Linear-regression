package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"project/regression" 
)



func main() {

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
	var r regression.Regression
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
