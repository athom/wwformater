package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"sync"

	"github.com/athom/goset"
)

var (
	lineLength = 18
	lineHigh   = 8
)

var (
	stopCharacters = []string{
		`,`,
		`，`,
		`、`,

		`.`,
		`。`,

		`？`,
		`?`,

		`:`,

		`!`,

		`"`,
		`”`,
		`“`,
	}
)

func makeCenter(line string) (r string) {
	r = line
	r = strings.TrimLeft(r, " ")
	r = strings.TrimRight(r, " ")
	x := zhLen(r)
	n := (lineLength - x) / 2
	for n > 0 {
		r = " " + r
		n--
	}
	return
}

func makeCenterVitial(lines []string) []string {
	n := lineHigh / 2
	for n > 0 {
		lines = append(lines, "")
		n--
	}
	return lines
}

func zhLen(line string) int {
	x := bytes.Runes([]byte(line))
	return len(x)
}

func runesToString(s []rune) (r string) {
	for _, c := range s {
		r += runeToString(c)
	}
	return r
}

func runeToString(c rune) (r string) {
	return fmt.Sprintf("%c", c)
}

func breakLongLine(line string) (r []string) {
	n := zhLen(line)
	if n <= lineLength {
		r = append(r, line)
		return
	}

	runes := bytes.Runes([]byte(line))
	startIndex := 0
	for index, _ := range runes {
		if index-startIndex == lineLength {
			//if index > 0 && index%lineLength == 0 {
			if goset.IsIncluded(stopCharacters, runeToString(runes[index])) {
				index += 1
			}
			zhLine := runes[startIndex:index]
			r = append(r, runesToString(zhLine))
			startIndex = index
		}
	}
	zhLine := runes[startIndex:n]
	line = runesToString(zhLine)
	if len(line) > 0 {
		r = append(r, line)
	}
	return
}

func parse(fileName string) (r []string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	text := string(b)
	lines := strings.Split(text, "\n")
	meetTitle := false
	for _, line := range lines {
		if line == "\n" {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		if !meetTitle {
			line = makeCenter(line)
			r = append(r, line)
			r = makeCenterVitial(r)
			meetTitle = true
			continue
		}

		line = strings.TrimLeft(line, " ")
		if zhLen(line) <= lineLength {
			r = append(r, line)
			continue
		}

		newLines := breakLongLine(line)
		r = append(r, newLines...)
	}
	return
}

func work(fileName string) {
	// output file
	outputFile := `formatted_` + fileName
	f, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines := parse(fileName)
	for _, line := range lines {
		f.WriteString(line + "\n")
	}

	if printStdout {
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

func batchWork(fileNames []string) {
	wg := sync.WaitGroup{}
	for _, fn := range fileNames {
		wg.Add(1)

		go func(fn string) {
			work(fn)
			wg.Done()
		}(fn)
	}
	wg.Wait()
}

var inputFileName string
var printStdout bool

func main() {
	flag.StringVar(&inputFileName, "f", "notset", "file name")
	flag.BoolVar(&printStdout, "d", false, "show result in stdout")
	flag.Parse()

	if inputFileName == `notset` {
		if len(os.Args) > 1 {
			fns := os.Args[1:]
			batchWork(fns)
		}
		return
	}

	// output file
	work(inputFileName)
}

func TestbreakLongLine() {
	lines := breakLongLine(`對了，你媽媽不是副校長嗎，你這次被抓緊政教處，你回家不就死定了，聽說這個滅絕師太不像老姑婆，雖然老姑婆兇是兇，還打人，但並不會輕易通知家長。可這滅絕師太據說極爲殘忍沒人性，無論大小錯都通知家長`)
	if lines[0] != `對了，你媽媽不是副校長嗎，你這次被抓` {
		panic(lines[0])
	}
	if lines[1] != `緊政教處，你回家不就死定了，聽說這個` {
		panic(lines[1])
	}
	if lines[2] != `滅絕師太不像老姑婆，雖然老姑婆兇是兇，` {
		panic(lines[2])
	}
	if lines[3] != `還打人，但並不會輕易通知家長。可這滅` {
		panic(lines[3])
	}
	if lines[4] != `絕師太據說極爲殘忍沒人性，無論大小錯` {
		panic(lines[4])
	}
	if lines[5] != `都通知家長` {
		panic(lines[5])
	}
	return
}
