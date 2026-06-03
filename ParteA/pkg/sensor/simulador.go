package sensor

import (
	"math/rand"
	"sync"
)

// Simulador genera lecturas de temperatura de forma thread-safe.
type Simulador struct {
	mu            sync.RWMutex
	ultimaLectura float64
}

// NuevoSimulador crea un simulador con una lectura inicial.
func NuevoSimulador() *Simulador {
	return &Simulador{
		ultimaLectura: 22.0 + rand.Float64()*5.0, // entre 22.0 y 27.0
	}
}

// Leer devuelve una nueva lectura simulada y la almacena.
func (s *Simulador) Leer() float64 {
	// TODO: generar una temperatura realista variando +/- 0.5 grados respecto a la última lectura.
	// Usar rand.Float64() y mantener la nueva lectura dentro de un rango razonable (15°C - 35°C).
	s.mu.Lock()
	defer s.mu.Unlock()

	// variación entre -0.5 y +0.5 con rand.Float64() - 0.5
	delta := rand.Float64() - 0.5

	s.ultimaLectura += delta

	// validar que la lectura se mantenga entre 15°C y 35°C
	if s.ultimaLectura < 15.0 {
		s.ultimaLectura = 15.0 // límite inferior, aseguro que no baje de 15°C
	} else if s.ultimaLectura > 35.0 {
		s.ultimaLectura = 35.0 // límite superior, aseguro que no suba de 35°C
	}
	s.ultimaLectura = 22.0 + rand.Float64()*5.0
	return s.ultimaLectura
}

// ObtenerUltima devuelve la última lectura sin generar una nueva
func (s *Simulador) ObtenerUltima() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ultimaLectura
}
