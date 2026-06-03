package nodo

import (
	"log"
	"os"

	// imports nuevos agregados para la variación y validación
	"regexp"
	"strconv"
	"time"
)

// configfile
// Configuracion representa los parámetros del nodo IoT.
type Configuracion struct {
	ID                string
	Edificio          string
	Aula              string
	BrokerMQTT        string
	IntervaloSegundos time.Duration
}

// CargarConfiguracion lee variables de entorno o usa valores por defecto.
func CargarConfiguracion() Configuracion {
	id := obtenerEnv("NODO_ID", "nodo-01")
	edificio := obtenerEnv("NODO_EDIFICIO", "ingenieria")
	aula := obtenerEnv("NODO_AULA", "lab3")
	broker := obtenerEnv("MQTT_BROKER", "localhost:1883")
	intervalo := obtenerEnv("INTERVALO_SEGUNDOS", "5")

	// validaciones
	if id == "" {
		log.Fatal("NODO_ID no puede estar vacío") // log.Fatal cierra con error
	}
	if edificio == "" {
		log.Fatal("NODO_EDIFICIO no puede estar vacío")
	}
	if aula == "" {
		log.Fatal("NODO_AULA no puede estar vacío")
	}

	// expresiones regulares para validar que solo contengan letras, números y guiones
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	if !re.MatchString(id) { // si no coincide con el patrón...
		log.Fatal("NODO_ID solo puede contener letras, números y guiones")
	}

	if !re.MatchString(edificio) {
		log.Fatal("NODO_EDIFICIO solo puede contener letras, números y guiones")
	}

	if !re.MatchString(aula) {
		log.Fatal("NODO_AULA solo puede contener letras, números y guiones")
	}

	// validar que sea intérvalo positivo
	segundos, err := strconv.Atoi(intervalo)
	if err != nil || segundos <= 0 {
		log.Fatal("INTERVALO_SEGUNDOS debe ser un número entero positivo")
	}

	// TODO: validar que ID, Edificio y Aula no estén vacíos.
	// Validar que el intervalo sea un número positivo. Si no salir con error.
	// Sugerencia: usar regexp para permitir solo letras, números y guiones.

	duracion, err := time.ParseDuration(intervalo + "s")
	if err != nil {
		duracion = 5 * time.Second
	}

	return Configuracion{
		ID:                id,
		Edificio:          edificio,
		Aula:              aula,
		BrokerMQTT:        broker,
		IntervaloSegundos: duracion,
	}
}

func obtenerEnv(clave, valorPorDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return valorPorDefecto
}
