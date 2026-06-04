package telemetria

import (
	"fmt"
	"sync"
	"time"

	"sd-comunicacion/pkg/protocolo"
)

// Telemetria es el servicio RPC.
type Telemetria struct {
	mu         sync.Mutex // Mutex para proteger el mapa de lecturas concurrentes
	lecturas   map[string]protocolo.Lectura
	contadorID int
}

// NuevaTelemetria inicializa el servicio
func NuevaTelemetria() *Telemetria {
	return &Telemetria{
		lecturas: make(map[string]protocolo.Lectura),
	}
}

// RegistrarLectura recibe datos del cliente y los guarda en el mapa.
// Requiere 2 argumentos (entrada y salida como puntero) y retorna un error[cite: 279].
func (t *Telemetria) RegistrarLectura(args protocolo.Lectura, resp *protocolo.RespuestaLectura) error {
	t.mu.Lock()         // Bloqueamos para escritura segura
	defer t.mu.Unlock() // Desbloqueamos al salir de la función

	// Guardamos la lectura en el mapa
	t.lecturas[args.SensorID] = args
	t.contadorID++

	// Preparamos la respuesta
	resp.ID = t.contadorID
	resp.Mensaje = "Lectura registrada con éxito"

	fmt.Printf("[%s] Lectura de %s: %.2f°C\n", time.Now().Format(time.Kitchen), args.SensorID, args.Temperatura)
	return nil
}

// ObtenerUltimaLectura devuelve la última lectura guardada de un sensor.
func (t *Telemetria) ObtenerUltimaLectura(args protocolo.ConsultaUltimaLectura, resp *protocolo.Lectura) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	lectura, existe := t.lecturas[args.SensorID]
	if !existe {
		return fmt.Errorf("no hay lecturas para el sensor %s", args.SensorID)
	}

	*resp = lectura
	return nil
}
