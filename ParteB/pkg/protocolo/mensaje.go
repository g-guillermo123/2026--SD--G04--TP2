package protocolo

// Heartbeat es el mensaje enviado periódicamente por UDP
type Heartbeat struct {
	NodoID    string `json:"nodo_id"`
	Timestamp int64  `json:"timestamp"`
	Contador  int    `json:"contador"`
}

// Lectura para el RPC RegistrarLectura
type Lectura struct {
	SensorID    string  `json:"sensor_id"`
	Temperatura float64 `json:"temperatura"`
	Timestamp   int64   `json:"timestamp"`
}

// RespuestaLectura para la confirmación de RegistrarLectura
type RespuestaLectura struct {
	ID      int    `json:"id"`
	Mensaje string `json:"mensaje"`
}

// ConsultaUltimaLectura para el RPC ObtenerUltimaLectura
type ConsultaUltimaLectura struct {
	SensorID string `json:"sensor_id"`
}
