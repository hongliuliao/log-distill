package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type DistillLog struct {
	distillValues []string
	splitSign     string
}

func NewDistillLog() *DistillLog {
	return &DistillLog{
		splitSign: ";",
	}
}

func (distillLog *DistillLog) toLine() string {
	var tmpLine string
	if len(distillLog.distillValues) == 0 {
		return tmpLine
	}
	for _, val := range distillLog.distillValues {
		tmpLine += val + distillLog.splitSign
	}
	tmpLine += "\n"
	return tmpLine
}

func (distillLog *DistillLog) fromLine(line string) error {
	if len(line) == 0 {
		return fmt.Errorf("input line is empty")
	}
	distillLog.distillValues = strings.Split(line, distillLog.splitSign)
	return nil
}

func (distillLog *DistillLog) find(key string, idx int) (bool, error) {
	if idx >= len(distillLog.distillValues) {
		return false, fmt.Errorf("idx output range, idx:%d, range:%d",
			idx, len(distillLog.distillValues))
	}
	if distillLog.distillValues[idx] != key {
		return false, fmt.Errorf("not found")
	}
	return true, nil
}

type DistillWriter struct {
	logPath string
	file    *os.File
}

func NewDistillWriter(logPath string) (*DistillWriter, error) {
	writer := &DistillWriter{
		logPath: logPath,
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("打开文件失败:", err)
		return nil, err
	}
	writer.file = file
	return writer, nil
}

func (writer *DistillWriter) Write(distill_logs []*DistillLog) error {
	start := time.Now()
	for _, dlog := range distill_logs {
		line := dlog.toLine()
		_, err := writer.file.WriteString(line)
		if err != nil {
			log.Printf("write file fail, line:%s", line)
			return err
		}
	}
	cost_ms := time.Since(start).Microseconds()
	log.Printf("write dlogs size:%d, cost_ms:%d", len(distill_logs), cost_ms)
	return nil
}

type DistillReader struct {
	logPath   string
	file      *os.File
	splitSign string
}

func NewDistillReader(logPath string) (*DistillReader, error) {
	reader := &DistillReader{
		logPath: logPath,
	}
	file, err := os.Open(logPath)
	if err != nil {
		log.Println("打开文件失败:", err)
		return nil, err
	}
	reader.file = file
	reader.splitSign = ";"
	return reader, nil
}

func (reader *DistillReader) Search(key string, idx int, limit int) []string {
	scanner := bufio.NewScanner(reader.file)
	result := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		dlog := NewDistillLog()
		dlog.fromLine(line)
		is_found, _ := dlog.find(key, idx)
		log.Printf("find for key:%s, idx:%d, is_found:%t, line:%s",
			key, idx, is_found, line)
		if is_found {
			result = append(result, line)
			if len(result) == limit {
				return result
			}
		}
	}
	return result
}
