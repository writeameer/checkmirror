package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func PrettyPrintJson(v interface{}) string {
	outBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("ERROR: PrettyPrintJson error '%s' for v='%+v'", err.Error(), v)
		elog.Error(102, fmt.Sprintf("ERROR: PrettyPrintJson error '%s' for v='%+v'", err.Error(), v))
	}
	outBytes = append(outBytes, "\n"...)
	return fmt.Sprintf("%-512s", string(outBytes)) // Pad right to avoid "Friendly HTTP error messages" issue in IE & Chrome
}

func JsonPrintIt(resp *JsonResponse) string {
	resp.Version = VERSION
	resp.Time = time.Now().UTC().Format(time.RFC3339)
	return PrettyPrintJson(resp)
}

func JsonErrorIt(msg string) string {
	jsonErr := JsonResponse{
		Error:   msg,
		Version: VERSION,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}
	outBytes, err := json.MarshalIndent(jsonErr, "", "  ")
	if err != nil {
		log.Printf("ERROR: JsonErrorIt marshal error '%s'", err.Error())
		elog.Error(103, fmt.Sprintf("ERROR: JsonErrorIt marshal error '%s'", err.Error()))
	}
	return fmt.Sprintf("%-512s", outBytes) // Pad right to avoid "Friendly HTTP error messages" issue in IE & Chrome
}
