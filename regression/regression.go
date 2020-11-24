package regression

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

var (
	// ErrNotEnoughData indica que no había suficientes puntos de datos para entrenar el modelo.
	ErrNotEnoughData = errors.New("not enough data points")
	// ErrTooManyVars indica que hay demasiadas variables para el número de observaciones que se realizan.
	ErrTooManyVars = errors.New("not enough observations to to support this many variables")
	// ErrRegressionRun indica que ya se ha llamado al método Run en el conjunto de datos entrenado.
	ErrRegressionRun = errors.New("regression has already been run")
)

// La regresión es la estructura de datos expuesta para interactuar con la API.
type Regression struct {
	names             describe
	data              []*dataPoint
	coeff             map[int]float64
	R2                float64
	Varianceobserved  float64
	VariancePrediccion float64
	initialised       bool
	Formula           string
	crosses           []featureCross
	hasRun            bool
}

type dataPoint struct {
	Observed  float64
	Variables []float64
	Prediccion float64
	Error     float64
}

type describe struct {
	obs  string
	vars map[int]string
}

// construcción más sencilla de puntos de datos de entrenamient
type DataPoints []*dataPoint

// crea un punto de datos bien formado * que se utiliza para el entrenamiento.
func DataPoint(obs float64, vars []float64) *dataPoint {
	return &dataPoint{Observed: obs, Variables: vars}
}

// Predice actualiza el valor para las características ingresadas.
func (r *Regression) Predict(vars []float64) (float64, error) {
	if !r.initialised {
		return 0, ErrNotEnoughData
	}

	// aplicar cualquier característica cruzada a vars
	for _, cross := range r.crosses {
		vars = append(vars, cross.Calculate(vars)...)
	}

	p := r.Coeff(0)
	for j := 1; j < len(r.data[0].Variables)+1; j++ {
		p += r.Coeff(j) * vars[j-1]
	}
	return p, nil
}

// SetObserved establece el nombre del valor observado.
func (r *Regression) SetObserved(name string) {
	r.names.obs = name //OBTENMOOS EL VALOR DEL DOCUMENTO
}

// GetObserved obtiene el nombre del valor observado.
func (r *Regression) GetObserved() string {
	return r.names.obs //RETORNAMOS EL VALOOR
}

//  establece el nombre de la variable i.
func (r *Regression) SetVar(i int, name string) {
	if len(r.names.vars) == 0 {
		r.names.vars = make(map[int]string, 5)
	}
	r.names.vars[i] = name //OBTENMOS LA VARIABLE DEL RETORNADO GetObserved
}

// GetVar obtiene el nombre de la variable i
func (r *Regression) GetVar(i int) string {
	x := r.names.vars[i]
	if x == "" {
		s := []string{"X", strconv.Itoa(i)}
		return strings.Join(s, "")
	}
	return x //RETORNAMOS LA VARIABLE DEL RETORNADO GetObserved
}
// Entrene la regresión con algunos puntos de datos.
func (r *Regression) Train(d ...*dataPoint) {
	r.data = append(r.data, d...)
	if len(r.data) > 2 {
		r.initialised = true // CON LOS DATOS QUE SE VAN OBTENIENDO SE ARMANDO LA REGRESION
	}
}

//  SE VAN ACTUALIZANDO LOS VALORES Y SOLO SE LLAMA UNA VEZ EEN LA FUNCION RUN
func (r *Regression) applyCrosses() {
	unusedVariableIndexCursor := len(r.data[0].Variables)
	for _, point := range r.data {
		for _, cross := range r.crosses {
			point.Variables = append(point.Variables, cross.Calculate(point.Variables)...)
		}
	}

	if len(r.names.vars) == 0 {
		r.names.vars = make(map[int]string, 5)
	}
	for _, cross := range r.crosses {
		unusedVariableIndexCursor += cross.ExtendNames(r.names.vars, unusedVariableIndexCursor)
	}
}

// SE COMPRUEBA QUE HAYAN PASADO LAS VALIADCNS INICIALES
func (r *Regression) Run() error {
	if !r.initialised {
		return ErrNotEnoughData
	}
	if r.hasRun {
		return ErrRegressionRun
	}

	//apply any features crosses
	r.applyCrosses()
	r.hasRun = true

	observations := len(r.data)
	numOfvars := len(r.data[0].Variables)

	if observations < (numOfvars + 1) {
		return ErrTooManyVars
	}

	// Create some blank variable space
	observed := mat.NewDense(observations, 1, nil)
	variables := mat.NewDense(observations, numOfvars+1, nil)

	for i := 0; i < observations; i++ {
		observed.Set(i, 0, r.data[i].Observed)
		for j := 0; j < numOfvars+1; j++ {
			if j == 0 {
				variables.Set(i, 0, 1)
			} else {
				variables.Set(i, j, r.data[i].Variables[j-1])
			}
		}
	}

	// ejecuta
	_, n := variables.Dims() // cols
	qr := new(mat.QR)
	qr.Factorize(variables)
	q := new(mat.Dense)
	reg := new(mat.Dense)
	qr.QTo(q)
	qr.RTo(reg)

	qtr := q.T()
	qty := new(mat.Dense)
	qty.Mul(qtr, observed)

	c := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		c[i] = qty.At(i, 0)
		for j := i + 1; j < n; j++ {
			c[i] -= c[j] * reg.At(i, j)
		}
		c[i] /= reg.At(i, i)
	}

	// EXPONE LOS RESULTADOS DE LA REGRESION
	r.coeff = make(map[int]float64, numOfvars)
	for i, val := range c {
		r.coeff[i] = val
		if i == 0 {
			r.Formula = fmt.Sprintf("Prediccion = %.4f", val)
		} else {
			r.Formula += fmt.Sprintf(" + %v*%.4f", r.GetVar(i-1), val)
		}
	}

	r.calcPrediccion()
	r.calculaVarianza()
	r.calcR2()
	return nil
}

// devuelve el coeficiente calculado para la variable i.
func (r *Regression) Coeff(i int) float64 {
	if len(r.coeff) == 0 {
		return 0
	}
	return r.coeff[i]
}

// GetCoeffs devuelve los coeficientes calculados. El elemento en el índice 0 es el desplazamiento.
func (r *Regression) GetCoeffs() []float64 {
	if len(r.coeff) == 0 {
		return nil
	}
	coeffs := make([]float64, len(r.coeff))
	for i := range coeffs {
		coeffs[i] = r.coeff[i]
	}
	return coeffs
}

func (r *Regression) calcPrediccion() string {
	observations := len(r.data)
	var Prediccion float64
	var output string
	for i := 0; i < observations; i++ {
		r.data[i].Prediccion, _ = r.Predict(r.data[i].Variables)
		r.data[i].Error = r.data[i].Prediccion - r.data[i].Observed

		output += fmt.Sprintf("%v. observed = %v, Prediccion = %v, Error = %v", i, r.data[i].Observed, Prediccion, r.data[i].Error)
	}
	return output
}

func (r *Regression) calculaVarianza() string {
	observations := len(r.data)
	var obtotal, prtotal, obvar, prvar float64
	for i := 0; i < observations; i++ {
		obtotal += r.data[i].Observed
		prtotal += r.data[i].Prediccion
	}
	obaverage := obtotal / float64(observations)
	praverage := prtotal / float64(observations)

	for i := 0; i < observations; i++ {
		obvar += math.Pow(r.data[i].Observed-obaverage, 2)
		prvar += math.Pow(r.data[i].Prediccion-praverage, 2)
	}
	r.Varianceobserved = obvar / float64(observations)
	r.VariancePrediccion = prvar / float64(observations)
	return fmt.Sprintf("N = %v\nVariance observed = %v\nVariance Prediccion = %v\n", observations, r.Varianceobserved, r.VariancePrediccion)
}

func (r *Regression) calcR2() string {
	r.R2 = r.VariancePrediccion / r.Varianceobserved
	return fmt.Sprintf("R2 = %.2f", r.R2)
}

func (r *Regression) calcResiduals() string {
	str := fmt.Sprintf("Residuals:\nobserved|\tPrediccion|\tResidual\n")
	for _, d := range r.data {
		str += fmt.Sprintf("%.2f|\t%.2f|\t%.2f\n", d.Observed, d.Prediccion, d.Observed-d.Prediccion)
	}
	str += "\n"
	return str
}

