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
	// abrimos el documento 
	f, err := os.Open("student-matt.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// se crea un nuevo lecto 
	// establecemos el numero de columas del documento
	salesData := csv.NewReader(f)
	salesData.FieldsPerRecord = 2 // 2 columas G2 Y G3

	// se lee la data del docuemtno
	records, err := salesData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//modelaremos segun el valor G3
	var r regression.Regression
	r.SetObserved("G3")
	r.SetVar(0, "G2")

	// BUCLE DE LOS DATOS
	for i, record := range records {
		// excluimos el encabezado
		if i == 0 {
			continue
		}

		// seteamos el primer valor del arreglo 
		price, err := strconv.ParseFloat(records[i][0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// seteamos el segundo valor del arreglo
		grade, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		// insertamos los valores del arreglo
		r.Train(regression.DataPoint(price, []float64{grade}))
	}

	// modelo de regresion
	r.Run()
	// imprimios la formula retornada
	fmt.Printf("Regresión Formula:\n%v\n\n", r.Formula)
}
