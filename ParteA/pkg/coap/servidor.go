package coap

import (
	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"
	"sync"

	"bytes"
	"encoding/json"
	"log"
	"time"

	gocoap "github.com/plgd-dev/go-coap/v3" // alias para evitar confusión (mismo nombre)
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
)

// ServidorCoAP expone recursos REST sobre UDP.
type ServidorCoAP struct {
	sim    *sensor.Simulador
	config nodo.Configuracion
	mu     sync.RWMutex
	modo   string
}

// NuevoServidor crea la instancia del servidor CoAP.
func NuevoServidor(sim *sensor.Simulador, config nodo.Configuracion) *ServidorCoAP {
	return &ServidorCoAP{
		sim:    sim,
		config: config,
		modo:   "automatico",
	}
}

// TODO 6: Iniciar arranca el servidor UDP en el puerto 5683.
// Debe:
//
//	6a. Crear router con mux.NewRouter().
//	6b. Registrar handler GET /temperatura que devuelva JSON con la última lectura.
//	    El JSON debe incluir: nodo_id, temperatura, unidad, timestamp.
//	6c. Registrar handler PUT /config que actualice s.modo y otros parámetros desde el body JSON.
//	6d. Registrar handler GET /config que devuelva la configuración actual en JSON.
//	6e. Llamar coap.ListenAndServe("udp", ":5683", router).
func (s *ServidorCoAP) Iniciar() {
	r := mux.NewRouter()                                              // creación del router
	r.Handle("/temperatura", mux.HandlerFunc(s.handleGetTemperatura)) // registro del handler para GET /temperatura
	r.Handle("/config", mux.HandlerFunc(s.handleConfig))              // registro del handler para PUT y GET /config

	log.Println("Servidor CoAP iniciado en udp://localhost:5683")

	err := gocoap.ListenAndServe("udp", ":5683", r) // arranque del servidor CoAP
	if err != nil {
		log.Fatalf("Error crítico en el servidor CoAP: %v", err)
	}
}

// funciones handlers
// GET /temperatura
func (s *ServidorCoAP) handleGetTemperatura(w mux.ResponseWriter, r *mux.Message) {
	// Verificar que sea una petición GET
	code := r.Code()
	if code != codes.GET {
		w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Método no permitido")))
		return
	}

	// obtener lectura del simulador
	lectura := s.sim.ObtenerUltima()

	s.mu.Lock()
	nodoID := s.config.ID
	s.mu.Unlock()

	respuesta := struct {
		NodoID      string  `json:"nodo_id"`
		Temperatura float64 `json:"temperatura"`
		Unidad      string  `json:"unidad"`
		Timestamp   string  `json:"timestamp"`
	}{
		NodoID:      nodoID,
		Temperatura: lectura,
		Unidad:      "C",
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	jsonBytes, err := json.Marshal(respuesta)
	if err != nil {
		w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("Error interno serializando JSON")))
		return
	}

	// responder exitosamente (Content 2.05 en CoAP)
	w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(jsonBytes))
}

// manejador de /config para GET y PUT
func (s *ServidorCoAP) handleConfig(w mux.ResponseWriter, r *mux.Message) {
	code := r.Code()

	switch code {
	case codes.GET:
		s.handleGetConfig(w, r)
	case codes.PUT:
		s.handlePutConfig(w, r)
	default:
		w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Usa GET o PUT")))
	}
}

// GET /config
func (s *ServidorCoAP) handleGetConfig(w mux.ResponseWriter, r *mux.Message) {
	s.mu.RLock()
	configActual := s.config
	s.mu.RUnlock()

	jsonBytes, err := json.Marshal(configActual)
	if err != nil {
		w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("Error al serializar config")))
		return
	}

	w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(jsonBytes))
}

// PUT /config
func (s *ServidorCoAP) handlePutConfig(w mux.ResponseWriter, r *mux.Message) {
	body, err := r.ReadBody()
	if err != nil {
		w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("No se pudo leer el cuerpo")))
		return
	}

	var nuevaConfig struct {
		Modo              string `json:"modo"`
		IntervaloSegundos int    `json:"intervalo_segundos"`
	}

	if err := json.Unmarshal(body, &nuevaConfig); err != nil {
		w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("JSON inválido")))
		return
	}

	s.mu.Lock()
	if nuevaConfig.Modo != "" {
		s.modo = nuevaConfig.Modo
	}
	if nuevaConfig.IntervaloSegundos > 0 {
		s.config.IntervaloSegundos = time.Duration(nuevaConfig.IntervaloSegundos) * time.Second
	}
	log.Printf("[CoAP PUT] Configuración actualizada: Modo=%s, Intervalo=%v", s.modo, s.config.IntervaloSegundos)
	s.mu.Unlock()

	w.SetResponse(codes.Changed, message.TextPlain, bytes.NewReader([]byte("Configuración actualizada con éxito")))
}
