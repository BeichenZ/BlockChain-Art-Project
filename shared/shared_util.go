package shared

import (
	"encoding/json"
	"log"
)

func PrettyPrint_ListMap(data []Map) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(string(b))
}

func PrettyPrint_Map(data Map) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(string(b))
}


func PrettyPrint_Path(data []PointStruct) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(string(b))
}