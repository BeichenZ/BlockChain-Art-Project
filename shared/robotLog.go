package shared

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

func (r *RobotStruct) EncodeRobotLogInfo(robotLog RobotLog) string {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(robotLog)
	if err != nil {
		panic(err)
	}
	output := string(buf.Bytes())
	return output
	// fmt.Println(buf.Bytes())
}

func (r *RobotStruct) ReadFromLog() {
	robotLogContent, _ := ioutil.ReadFile("./" + r.Logname)
	buf := bytes.NewBuffer(robotLogContent)

	var decodedRobotLog RobotLog

	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&decodedRobotLog)
	if err != nil {
		panic(err)
	}

	r.RMap = decodedRobotLog.RMap
	r.CurLocation = decodedRobotLog.CurLocation
	r.CurrTask = decodedRobotLog.CurrTask
	r.CurPath = decodedRobotLog.CurPath

	fmt.Printf("%+v", decodedRobotLog.RMap)
	fmt.Println(decodedRobotLog.CurLocation)
	fmt.Println(decodedRobotLog.CurrTask)
	fmt.Println(decodedRobotLog.CurPath)
	fmt.Println("finshed loading from log")
}

func (r *RobotStruct) CreateLog() (*os.File, error) {
	file, err := os.Create("./" + r.Logname)
	if err != nil {
		fmt.Println("error creating robot log")
	}
	return file, err
}

func (r *RobotStruct) ProduceLogInfo() RobotLog {
	robotLog := RobotLog{
		CurrTask:    r.CurrTask,
		CurPath:     r.CurPath,
		RMap:        r.RMap,
		CurLocation: r.CurLocation,
	}
	return robotLog
}

func (r *RobotStruct) LocateLog() (*os.File, error) {
	file, err := os.Open(r.Logname)
	return file, err
}

func (r *RobotStruct) ClearLog() {

}

func (r *RobotStruct) WriteToLog() {
	os.Remove(r.Logname)
	newLog, _ := r.CreateLog()
	robotLogInfo := r.ProduceLogInfo()
	encodedRobotLog := r.EncodeRobotLogInfo(robotLogInfo)
	newLog.WriteString(encodedRobotLog)
	newLog.Close()
	fmt.Println("Successfully wrote to robot's log")
}
