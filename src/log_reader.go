package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const READER_BATCH_SIZE int = 10

type LogReader struct {
	filePath  string
	regex     string
	scanner   *bufio.Scanner
	log_regex *regexp.Regexp
}

func NewLogReader(filePath string, regex string) (*LogReader, error) {
	reader := &LogReader{
		filePath: filePath,
		regex:    regex,
	}
	err := reader.init()
	return reader, err
}

func (reader *LogReader) init() error {
	file, err := os.Open(reader.filePath)
	if err != nil {
		log.Fatalf("openn file fail, path:%s", reader.filePath)
		return err
	}
	reader.scanner = bufio.NewScanner(file)
	reader.log_regex, err = regexp.Compile(reader.regex)
	if err != nil {
		log.Fatalf("regex compile err, regex:%s, err:%s", reader.regex, err)
		return err
	}
	return nil
}

func (reader *LogReader) GetDistillLogs() ([]*DistillLog, error) {
	start := time.Now()
	distill_logs := make([]*DistillLog, 0)

	lineNum := 0
	for reader.scanner.Scan() {
		line := reader.scanner.Text()
		lineNum++
		if lineNum > READER_BATCH_SIZE {
			break
		}
		matches := reader.log_regex.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Printf("not find match in line:%s", line)
			continue
		}
		log.Printf("match line:%s, match size:%d", line, len(matches))
		distillLog := NewDistillLog()
		for i, match := range matches {
			if i == 0 {
				continue
			}
			log.Printf("match in line:%s", match)
			distillLog.distillValues = append(distillLog.distillValues, match)
		}

		distill_logs = append(distill_logs, distillLog)
	}
	cost_ms := time.Since(start).Milliseconds()
	log.Printf("distill logs size:%s, cost_ms:%d",
		strconv.Itoa(len(distill_logs)), cost_ms)
	return distill_logs, nil
}
