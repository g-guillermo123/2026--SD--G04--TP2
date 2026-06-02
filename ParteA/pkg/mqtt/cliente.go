package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"
)

// Cliente encapsula la conexión MQTT del nodo.
type Cliente struct {
	config   nodo.Configuracion
	interno  mqtt.Client
	opciones *mqtt.ClientOptions
}

// TODO 1: NuevoCliente crea la configuración inicial del cliente MQTT.
// Debe:
//   1a. Construir el tópico del testamento: nodo/{id}/estado
//   1b. Configurar el mensaje del testamento como {"estado":"offline"} con QoS 1 y retained=true.
//   1c. Configurar ClientID único, timeout de conexión y reconexión automática.
//
// Sugerencia: usar mqtt.NewClientOptions().AddBroker(...).SetClientID(...).SetWill(...)
func NuevoCliente(config nodo.Configuracion) (*Cliente, error) {

	cliente := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetClientID("sensor-greenhouse-01").
		SetWill( // testamento y última voluntad
			"sensors/greenhouse/status", // tópico del testamento
			"Estado: Offline",           // payload (mensaje del testamento)
			1,                           // QoS en 1
			true,                        // retain (almacena el mensaje en el broker)
		)

	return nil, fmt.Errorf("Error al crear cliente MQTT (TODO Realizado)") //error en caso de que falle la creación del cliente
}

// TODO 2: Conectar establece la sesión con el broker.
// Tras conectar, debe publicar un mensaje retenido {"estado":"online"} en nodo/{id}/estado.
func (c *Cliente) Conectar() error {
	// publica el estado "online" con retención y QoS 1
	client.Publish("topico" // actualizar con el tópico correcto
	, 1 // QoS 1
	, true // retención
	, "Online") // el estado actual del nodo
	return fmt.Errorf("Error al conectar (TODO Realizado)") //error en caso de que falle la conexión
}

// TODO 3: PublicarLecturas envía periódicamente las lecturas del sensor.
// Debe:
//   3a. Construir el tópico: campus/{edificio}/{aula}/sensor/temperatura
//   3b. En un ticker cada config.IntervaloSegundos, llamar sim.Leer(), serializar a JSON y publicar con QoS 1.
//   El JSON debe tener: {"nodo_id": ..., "temperatura": ..., "unidad":"C", "timestamp":"..."}
func (c *Cliente) PublicarLecturas(sim *sensor.Simulador, config nodo.Configuracion) {
	ticker := time.NewTicker(config.IntervaloSegundos * time.Second) // crea un ticker que se activa cada IntervaloSegundos
	defer ticker.Stop() // asegura que el ticker se detenga al finalizar la función
	// ¿cómo construyo el tópico?
	for {
		select {
		case <-ticker.C: // cada vez que el ticker se active (3b)
			lectura := sim.Leer() // obtiene la lectura del sensor
			token := client.Publish( // publicar el mensaje
				"topico", // actualizar tópico
				1,    // QoS 1 — at least once
				true, // retain 
				payload, // variable payload que deberia ser el JSON
			)
			token.Wait() // espera a que la publicación se complete}
		}
	}
}
// TODO 4: SuscribirComandos se une al tópico de actuadores y procesa mensajes.
// Debe:
//   4a. Suscribirse al tópico campus/{edificio}/{aula}/actuador/cmd con QoS 1.
//   4b. En el callback, deserializar el JSON, imprimir el comando recibido y simular la ejecución.
//   Ejemplo de payload esperado: {"accion":"encender_alarma", "origen":"dashboard"}
func (c *Cliente) SuscribirComandos(config nodo.Configuracion) error {
	// COMPLETAR
	return fmt.Errorf("TODO: implementar SuscribirComandos")
}

// TODO 5: Desconectar cierra limpiamente la sesión MQTT.
// Sugerencia: publicar estado offline retenido antes de desconectar.
func (c *Cliente) Desconectar() {
	log.Println("Sensor desconectándose del broker MQTT...") // mensaje de desconexión limpia, no dispara LWT
	cancel()
	client.Disconnect(250)
}
