package main

import "log"

func main() {
	reader, err := NewLogReader("./test/test.log", `log_id:(\d+),\s+uid:(\d+)`)
	if err != nil {
		log.Printf("init log reader err:%s", err)
		return
	}
	dLogs, err := reader.GetDistillLogs()
	if err != nil {
		log.Printf("read distill logs fail, err:%s", err)
		return
	}

	writer, err := NewDistillWriter("./output/dlog.log")
	if err != nil {
		log.Printf("distill log writer init fail, err:%s", err)
		return
	}
	writer.Write(dLogs)
	log.Println("hello world")
}
