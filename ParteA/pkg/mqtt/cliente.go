package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
// 	 Sugerencia: usar mqtt.NewClientOptions().AddBroker(...).SetClientID(...).SetWill(...)

func NuevoCliente(config nodo.Configuracion) (*Cliente, error) {
	topicoTestamento := fmt.Sprintf("nodo/%s/estado", config.ID) // topic for the last will, built with the node ID
	payloadTestamento := `{"estado":"offline"}`                 // serializado

	opts := mqtt.NewClientOptions().
		AddBroker(config.BrokerMQTT).      //llamada al conifgfile
		SetClientID(config.ID).            // llamada al configfile
		SetAutoReconnect(true).            // habilita reconexión automática
		SetConnectTimeout(10*time.Second). // intervalo de reconexión automática
		SetWill(                           // testamento y última voluntad
			topicoTestamento,  // tópico del testamento
			payloadTestamento, // payload (mensaje del testamento)
			1,                 // QoS en 1
			true,              // retain (almacena el mensaje en el broker)
		)

	intlClient := mqtt.NewClient(opts)

	// retorno del cliente MQTT configurado, sin conexión aún
	return &Cliente{
		config:   config,
		interno:  intlClient,
		opciones: opts,
	}, nil
}

// TODO 2: Conectar establece la sesión con el broker.
// Tras conectar, debe publicar un mensaje retenido {"estado":"online"} en nodo/{id}/estado.
func (c *Cliente) Conectar() error {
	// conexión al broker
	token := c.interno.Connect()              // inicia la conexión
	if token.Wait() && token.Error() != nil { // espera a que la conexión se establezca y verifica errores
		return fmt.Errorf("Error al conectar al broker MQTT: %v", token.Error()) // manejo de errores
	}
	log.Println("Conexión MQTT establecida exitosamente.") // aviso conexión

	// publica el estado "online" con retención y QoS 1
	topico := fmt.Sprintf("nodo/%s/estado", c.config.ID) // tópico para publicar el estado online
	payload := `{"estado":"online"}`                     // mensaje JSON indicando que el nodo está online

	pubToken := c.interno.Publish(topico, 1, true, payload) // publica el mensaje con QoS 1 y retain=true
	pubToken.WaitTimeout(5 * time.Second)                   // espera a que la publicación se complete usando WaitTimeout
	return nil                                              // retorna nil si todo fue exitoso
}

// TODO 3: PublicarLecturas envía periódicamente las lecturas del sensor.
// Debe:
//
//	3a. Construir el tópico: campus/{edificio}/{aula}/sensor/temperatura
//	3b. En un ticker cada config.IntervaloSegundos, llamar sim.Leer(), serializar a JSON y publicar con QoS 1.
//	El JSON debe tener: {"nodo_id": ..., "temperatura": ..., "unidad":"C", "timestamp":"..."}
func (c *Cliente) PublicarLecturas(sim *sensor.Simulador, config nodo.Configuracion) {

	ticker := time.NewTicker(config.IntervaloSegundos)                                     // IntervaloSegundos is already a time.Duration; no conversion needed
	defer ticker.Stop()                                                                    // asegura que el ticker se detenga al finalizar la función
	topico := fmt.Sprintf("campus/%s/%s/sensor/temperatura", config.Edificio, config.Aula) // tópico para publicar las lecturas

	for {
		select {
		case <-ticker.C: // cada vez que el ticker se active (3b)
			lectura := sim.Leer() // obtiene la lectura del sensor
			// construye la estructura del payload con los datos requeridos (para serializar)
			payloadEstructura := struct {
				NodoID      string  `json:"nodo_id"`
				Temperatura float64 `json:"temperatura"`
				Unidad      string  `json:"unidad"`
				Timestamp   string  `json:"timestamp"`
			}{
				NodoID:      config.ID,
				Temperatura: lectura,
				Unidad:      "C",
				Timestamp:   time.Now().Format(time.RFC3339), // Formato estándar ISO
			}
			payloadBytes, err := json.Marshal(payloadEstructura)
			if err != nil {
				log.Printf("Error al serializar la lectura a JSON: %v", err) // manejo de errores de serialización
				continue                                                     // salta esta iteración si hay un error
			}

			token := c.interno.Publish( // publicar el mensaje
				topico,       // tópico construido
				1,            // QoS 1 — at least once
				true,         // retain
				payloadBytes, // variable payload serializado
			)
			token.WaitTimeout(5 * time.Second) // espera a que la publicación se complete usando WaitTimeout
		}
	}
}

// TODO 4: SuscribirComandos se une al tópico de actuadores y procesa mensajes.
// Debe:
//
//	4a. Suscribirse al tópico campus/{edificio}/{aula}/actuador/cmd con QoS 1.
//	4b. En el callback, deserializar el JSON, imprimir el comando recibido y simular la ejecución.
//	Ejemplo de payload esperado: {"accion":"encender_alarma", "origen":"dashboard"}
func (c *Cliente) SuscribirComandos(config nodo.Configuracion) error {
	topico := fmt.Sprintf("campus/%s/%s/actuador/cmd", config.Edificio, config.Aula) // tópico para suscribirse a comandos

	var callback mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensaje recibido del tópico %s: %s", msg.Topic(), string(msg.Payload())) // imprime el mensaje recibido

		var comando struct {
			Accion string `json:"accion"`
			Origen string `json:"origen"`
		}

		err := json.Unmarshal(msg.Payload(), &comando) // deserializa el payload del mensaje
		if err != nil {
			log.Printf("Error al deserializar el comando recibido: %v", err) // manejo de errores de deserialización
			return                                                           // sale del callback si hay un error
		}

		// simular ejecución
		log.Printf("[ACTUADOR] Ejecutando acción: '%s' (Solicitado por: %s)", comando.Accion, comando.Origen)
		if comando.Accion == "encender_alarma" {
			fmt.Println("¡Alarma encendida en Aula!)") // simulación de acción específica
		}
	}

	// suscribirse al tópico con el callback definido
	token := c.interno.Subscribe(topico, 1, callback) // qos 1
	if token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return fmt.Errorf("error al suscribirse al tópico %s: %v", topico, token.Error())
	}

	log.Printf("Suscribiéndose al tópico de comandos: %s", topico) // mensaje de suscripción
	return nil
}

// TODO 5: Desconectar cierra limpiamente la sesión MQTT.
// Sugerencia: publicar estado offline retenido antes de desconectar.
func (c *Cliente) Desconectar() {
	log.Println("Sensor desconectándose del broker MQTT...") // mensaje de desconexión limpia, no dispara LWT
	// Publicamos manualmente el estado offline (ya que el LWT no se dispara en desconexiones ordenadas)
	topico := fmt.Sprintf("nodo/%s/estado", c.config.ID) // tópico del estado del nodo
	pubToken := c.interno.Publish(topico, 1, true, `{"estado":"offline"}`)
	pubToken.WaitTimeout(5 * time.Second)

	// desconectar dando 250 milisegundos para terminar posibles tareas activas
	c.interno.Disconnect(250)
	log.Println("Cliente MQTT desconectado con éxito.")
}
