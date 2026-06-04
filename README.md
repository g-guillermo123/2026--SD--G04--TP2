# 2026 SD - G04 - Práctica Guiada 02

Repositorio con la resolución de las Partes A y B de la Práctica Guiada 2 de la cátedra Sistemas Distribuidos.

## Integrantes

- Gallo Guillermo Ariel
- Pedernera Theisen Nahuel Thomas

---

## Estructura del Proyecto

El trabajo práctico se divide en dos secciones principales, cada una resolviendo distintos desafíos de comunicación en sistemas distribuidos. Ambas partes cuentan con su propio entorno dockerizado y automatizado con `Makefile`.

* **Parte A** Comunicación IoT utilizando MQTT y CoAP (Nodo Smart Campus).
* **Parte B** Servicio de Telemetría con Detección de Fallos (RPC + Heartbeat UDP).

---

## Desarrollo de la Parte A: Comunicación IoT mediante MQTT y servidores CoAP

En la Parte A, se implementa la arquitectura base para una red IoT distribuida (un entorno de Smart Campus), donde se combinan dos protocolos: MQTT (para el envío masivo y asíncrono de eventos) y CoAP (para consultas directas tipo REST bajo demanda).
El broker MQTT(NanoMQ) se encarga de ser el intermediario entre el cliente publicador (sensor simulado) y los nodos consumidores subscriptos al tópico dónde el sensor publica.
Los Nodos (dispositivos IoT) corren de forma independiente mediante gorutinas:
* Publicador MQTT (Publisher): lee periódicamente un componente simulador que genera mediciones de temperatura y las envía al Broker en un formato JSON serializado.
* Suscriptor MQTT (Subscriber): se mantiene escuchando un tópico específico, si el broker le transmite una orden, el nodo reacciona de inmediato ejecutando la acción.
* Servidor CoAP Local: levanta en paralelo un servidor web. Esto permite que cualquier cliente externo en la red pueda hacerle una consulta directa (un GET /temperatura o GET/PUT /config) directamente al nodo sin necesidad de pasar por el broker MQTT.

### Ejecución (Vía Docker Compose)

Para levantar el entorno completo, navegar al directorio `ParteB` y utilizar tres terminales:

```bash
# Terminal 1: Levantar broker y ver registro de logs
make run
make docker compose logs -f

# Terminal 2: Levantar el primer cliente interactivo
make docker-nodo1

# Terminal 3: Levantar el segundo cliente interactivo
make docker-nodo2

# Terminal 4 (opcional): Envío de comando GET para simular lectura
coap-client -m get coap://127.0.0.1:5683/temperatura
```

---

## Desarrollo de la Parte B: Telemetría y Detección de Fallos

En la Parte B implementamos un sistema distribuido enfocado en la comunicación confiable y la tolerancia a fallos. Simulamos nodos sensores (clientes) que reportan datos a un servidor central. 

Las características técnicas principales incluyen:
- **Comunicación RPC con JSON (`net/rpc/jsonrpc`):** Envío periódico de lecturas de temperatura simuladas.
- **Tolerancia a fallos de red:** Los clientes capturan los errores de conexión RPC y soportan caídas del servidor sin interrumpir su ciclo de vida.
- **Detector de Fallos Pasivo (Heartbeat UDP):** El servidor emite latidos periódicos a los nodos. Los clientes corren una rutina concurrente (goroutine) que monitorea los tiempos de llegada para determinar la salud del servidor, transitando por los estados `alive` -> `suspect` -> `dead`.

### Ejecución (Vía Docker Compose)

Para levantar el entorno completo, navegar al directorio `ParteB` y utilizar tres terminales:

```bash
# Terminal 1: Levantar servidor y ver logs
make docker-up
make docker-logs

# Terminal 2: Levantar el primer cliente interactivo
make docker-cliente1

# Terminal 3: Levantar el segundo cliente interactivo
make docker-cliente2
```
