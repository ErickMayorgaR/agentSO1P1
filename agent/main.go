package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/shirou/gopsutil/cpu"
)

type RAMInfo struct {
	TotalRAM      uint64  `json:"totalRAM"`
	RAMEnUso      uint64  `json:"RAMEnUso"`
	RAMLibre      uint64  `json:"RAMLibre"`
	PorcentajeUso float64 `json:"porcentajeUso"`
}

type PIDCPU struct {
	PID    uint64 `json:"PID"`
	Name   string `json:"Nombre"`
	Status uint64 `json:"Status"`
	Size   uint64 `json:"Size"`
	UID    uint64 `json:"UID"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			sendRAMInfo()
			sendUsageCPU()
			<-time.After(1 * time.Second)
		}

	}()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.Handle("/getPIDCPU", corsOptions.Handler(http.HandlerFunc(getPIDCPU)))
	http.Handle("/killProcess", corsOptions.Handler(http.HandlerFunc(killProcess)))

	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("Server started on :8080")

}

func sendRAMInfo() {
	data, err := os.ReadFile("/proc/ram_201901758")
	if err != nil {
		fmt.Println(err)
		return
	}

	ramInfo := parseRAMInfo(string(data))
	ramInfo.PorcentajeUso = calculateUsagePercent(ramInfo.RAMEnUso, ramInfo.TotalRAM)

	jsonData, err := json.Marshal(ramInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendInfo([]byte(jsonData), "/insertRAMInformation")
}

func parseRAMInfo(data string) RAMInfo {
	var ramInfo RAMInfo
	fmt.Sscanf(data, "Total_RAM: %d\nRAM_en_Uso: %d\nRAM_libre: %d\nPorcentaje_en_uso: %f\n",
		&ramInfo.TotalRAM, &ramInfo.RAMEnUso, &ramInfo.RAMLibre, &ramInfo.PorcentajeUso)
	return ramInfo
}

func calculateUsagePercent(usedMemory, totalMemory uint64) float64 {
	return (float64(usedMemory) / float64(totalMemory) * 100)
}

func sendUsageCPU() {
	// Obtiene el uso de la CPU.
	// 0 significa todos los nucleos y false significa uso promedio del CPU de los procesos.
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Convierte la rebanada de valores flotantes a un único valor flotante.
	usoCPU := cpuPercent[0]

	// Crea un objeto JSON con la estructura UsoCPU: Porcentaje.
	usoCPUStruct := struct {
		UsoCPU float64 `json:"UsoCPU"`
	}{
		UsoCPU: usoCPU,
	}

	// Codifica el objeto JSON.
	usoCPUJSON, err := json.Marshal(usoCPUStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cpuPercent)
	sendInfo(usoCPUJSON, "/insertCPUInformation")

}

func sendInfo(info []byte, url string) {
	direccionServer := os.Getenv("direccionServer")

	req, err := http.NewRequest("POST", "http://"+direccionServer+":5000"+url, bytes.NewBuffer(info))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Realiza la petición
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Status)
}

func getPIDCPU(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("/proc/cpu_201901758")
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}
	pidCpuArray := []PIDCPU{}

	lines := strings.Split(string(data), "\n")
	pidCpu := PIDCPU{}
	for _, line := range lines {
		if line != "" {
			err2 := json.Unmarshal([]byte(line), &pidCpu)
			if err2 != nil {
				fmt.Printf("Error decoding input string: %v\n", err2)
				http.Error(w, "Error decoding input string", http.StatusBadRequest)
				return
			}

			pidCpuArray = append(pidCpuArray, pidCpu)
		}

	}

	// Codifica el arreglo de structs PIDCPU a un JSON
	jsonData, err := json.Marshal(pidCpuArray)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Escribe el JSON en la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func killProcess(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PID int `json:"pid"`
	}
	err1 := json.NewDecoder(r.Body).Decode(&body)
	if err1 != nil {
		// Enviar una respuesta de error personalizada
		data := struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{
			Status: "error",
			Error:  err1.Error(),
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	// Ejecutar el comando kill -9 PID
	err := exec.Command("kill", "-9", strconv.Itoa(body.PID)).Run()
	if err != nil {
		// Enviar una respuesta de error personalizada
		data := struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{
			Status: "error",
			Error:  "Error al eliminar PID: " + err.Error(),
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	// Enviar una respuesta de éxito
	data := struct {
		Status  string `json:"status"`
		Mensaje string `json:"mensaje"`
	}{
		Status:  "ok",
		Mensaje: "Proceso eliminado exitosamente",
	}
	json.NewEncoder(w).Encode(data)

}
