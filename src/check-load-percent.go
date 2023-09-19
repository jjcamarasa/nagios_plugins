/*
Basado en el check-load de https://github.com/atc0005/go-check-plugins.git
DESCRIPCION:	
Este paquete calcula el uso de cpu en porcentajes en base al numero de CPUs del sistema. Solo funciona en linux. Saca perfdata para visualizar graficas. Probado en debian 10/11 y RedHat

VERSION:
		1.0	2023-06-27 - Version inicial
		1.1	2023-09-13 - Afegits flags per a msg, interval i disable perfdata.
USO:	Sin argumentos, usa por defecto el warning al 85% y el critico al 95%.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func NagiosExit(ret int, msg string, perf string) {
	OutStates := [4]string{"OK", "WARNING", "CRITICAL", "UNKNOWN"}
	if len(perf) > 0 {
		fmt.Printf("%s - %s | %s\n", OutStates[ret], msg, perf)
	} else {
		fmt.Printf("%s - %s\n", OutStates[ret], msg)
	}
	os.Exit(ret)
}

func GetLoadAvg() (loadavgs [3]float64, _ error) {
	contentbytes, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return loadavgs, fmt.Errorf("Failed to load /proc/loadavg: %s", err)
	}
	content := string(contentbytes)
	cols := strings.Split(content, " ")
	for i := 0; i < 3; i++ {
		f, err := strconv.ParseFloat(cols[i], 64)
		if err != nil {
			return loadavgs, fmt.Errorf("Failed to parse loadavg metrics: %s", err)
		}
		loadavgs[i] = f
	}
	return loadavgs, nil
}

func main() {
	var warning int
	var critical int
	var noperfdata bool
	var time int
	var interval int = -1
	intervalText := [3]string{"minute", "5 minutes", "15 minutes"}
	var outmsg string
	NagiosRetCode := 3
	NagiosMsg := "Failed to execute check"
	var PerfData string

	flag.IntVar(&warning, "warning", 85, "Percent of cpu usage that is consedered as Warning")
	flag.IntVar(&critical, "critical", 95, "Percent of cpu usage that is consedered as Critical")
	flag.BoolVar(&noperfdata, "noperfdata", false, "Print performance data for charts")
	flag.StringVar(&outmsg, "msg", "cpu average load", "Message to display in output")
	flag.IntVar(&time, "time", 0, "Interval time to show, accepted input 1,5,15 minutes. Default show all.")

	flag.Parse()

	switch time {
	case 0:
		interval = 0
	case 1:
		interval = 0
	case 5:
		interval = 1
	case 15:
		interval = 2
	default:
		NagiosRetCode = 2
		NagiosMsg = "Bad options: time Interval must be 1, 5 or 15"
		NagiosExit(NagiosRetCode, NagiosMsg, PerfData)
	}

	loadavgs, err := GetLoadAvg()
	if err != nil {
		NagiosRetCode = 3
		NagiosMsg = "Failed to acces /proc/loadavg"
	} else {
		numCPU := runtime.NumCPU()
		for i, load := range loadavgs {
			loadavgs[i] = load / float64(numCPU)
		}
		var PercentLoad [3]int
		var StatusLoad [3]int
		for i := range loadavgs {
			PercentLoad[i] = int(loadavgs[i] * 100)
			if PercentLoad[i] > 100 {
			 PercentLoad[i] = 100
			}
			if PercentLoad[i] > critical {
				StatusLoad[i] = 2
				NagiosRetCode = 2
				break
			}
			if PercentLoad[i] > warning {
				StatusLoad[i] = 1
				NagiosRetCode = 1
			} else {
				StatusLoad[i] = 0
				NagiosRetCode = 0
			}
		}
		var MsgOutput string = fmt.Sprintf("%s 1m:%d%%, 5m:%d%%, 15m:%d%%", outmsg, PercentLoad[0], PercentLoad[1], PercentLoad[2])
		var FullPerfData = fmt.Sprintf("'total 1m'=%d%%;%d;%d 'total 5m'=%d%%;%d;%d 'total 15m'=%d%%;%d;%d", PercentLoad[0], warning, critical, PercentLoad[1], warning, critical, PercentLoad[2], warning, critical)
		if time != 0 {
			MsgOutput = fmt.Sprintf("%s in last %s: %d%%", outmsg, intervalText[interval], PercentLoad[interval])
			FullPerfData = fmt.Sprintf("'%dm'=%d%%;%d;%d", time, PercentLoad[interval], warning, critical)
		} 

		if NagiosRetCode != 3 {
			NagiosMsg = MsgOutput
			PerfData = FullPerfData
		}
	}

	if noperfdata == true {
	 PerfData = ""
	}
	 NagiosExit(NagiosRetCode, NagiosMsg, PerfData)
}
