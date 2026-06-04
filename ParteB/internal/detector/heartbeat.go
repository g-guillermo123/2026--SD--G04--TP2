package detector

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"sd-comunicacion/pkg/protocolo"
)

// Enviador envía heartbeats UDP periódicamente (Servidor)
type Enviador struct {
	destino   string
	intervalo time.Duration
	nodoID    string
	contador  int
}

func NuevoEnviador(destino string, intervalo time.Duration, nodoID string) *Enviador {
	return &Enviador{
		destino:   destino,
		intervalo: intervalo,
		nodoID:    nodoID,
	}
}

func (e *Enviador) Iniciar() {
	for {
		// Marcamos la conexión y envío (Fire and forget, es UDP) [cite: 397]
		conn, err := net.Dial("udp", e.destino)
		if err == nil {
			hb := protocolo.Heartbeat{
				NodoID:    e.nodoID,
				Timestamp: time.Now().Unix(),
				Contador:  e.contador,
			}
			data, _ := json.Marshal(hb)
			conn.Write(data)
			conn.Close()
			e.contador++
		}
		time.Sleep(e.intervalo)
	}
}

// Receptor escucha heartbeats y cambia estados (Cliente)
type Receptor struct {
	puerto  string
	timeout time.Duration
	ultimo  time.Time
	estado  string // "alive", "suspect", "dead"
}

func NuevoReceptor(puerto string, timeout time.Duration) *Receptor {
	return &Receptor{
		puerto:  puerto,
		timeout: timeout,
		estado:  "dead", // Asume muerto hasta que llega el primer heartbeat
	}
}

func (r *Receptor) Escuchar() {
	addr, err := net.ResolveUDPAddr("udp", r.puerto)
	if err != nil {
		fmt.Println("Error resolviendo puerto UDP:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error escuchando UDP:", err)
		return
	}
	defer conn.Close()

	// Goroutine que revisa el timeout periódicamente para cambiar el estado (alive -> suspect -> dead) [cite: 425]
	go func() {
		for {
			time.Sleep(r.timeout / 2)
			if r.ultimo.IsZero() {
				continue
			}

			inactividad := time.Since(r.ultimo)
			nuevoEstado := r.estado

			if inactividad > 2*r.timeout {
				nuevoEstado = "dead"
			} else if inactividad > r.timeout {
				nuevoEstado = "suspect"
			} else {
				nuevoEstado = "alive"
			}

			if nuevoEstado != r.estado {
				r.estado = nuevoEstado
				fmt.Printf("[Detector] Servidor pasó a estado: %s\n", r.estado)
			}
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf) // Escuchamos heartbeats [cite: 410]
		if err == nil {
			var hb protocolo.Heartbeat
			if err := json.Unmarshal(buf[:n], &hb); err == nil {
				r.ultimo = time.Now()
				if r.estado != "alive" {
					r.estado = "alive"
					fmt.Printf("[Detector] Servidor en estado: alive (Recibido de %s)\n", hb.NodoID)
				}
			}
		}
	}
}
