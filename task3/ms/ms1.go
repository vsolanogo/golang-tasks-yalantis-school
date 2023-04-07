package main

import (
	"log"
	"task3/myprotodata"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bufio"
	"flag"
	"io"
	"os"
	"path/filepath"
)

func lineByLine(file string, client myprotodata.SorterMessageServiceClient, firstLineIsHeader bool) {
	var err error

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("erorr reading file")
	}
	defer f.Close()

	r := bufio.NewReader(f)

	var skippedFirstLine = !firstLineIsHeader

	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}

		if !skippedFirstLine {
			skippedFirstLine = !skippedFirstLine

			continue
		}

		response, requesterr := client.GetRequest(context.Background(), &myprotodata.Request{Action: "ADD", Payload: []string{line}})

		if requesterr != nil {
			log.Fatal("error service", requesterr)
		}
	
		fmt.Println(response.Error)
		fmt.Println(response.Payload)
	}

	return 
}

func readFromStdin(client myprotodata.SorterMessageServiceClient, firstLineIsHeader bool) {
	fmt.Println("input your comma-separated values:")

	reader := bufio.NewReader(os.Stdin)

	var skippedFirstLine = !firstLineIsHeader

	for {
		text, _ := reader.ReadString('\n')

		if !skippedFirstLine {
			skippedFirstLine = !skippedFirstLine

			continue
		}

		if text == "\n" {
			break
		}

		response, err := client.GetRequest(context.Background(), &myprotodata.Request{Action: "ADD", Payload: []string{text}})

		if err != nil {
			log.Fatal("error service", err)
		}
	
		fmt.Println(response.Error)
		fmt.Println(response.Payload)
	}

	return 
}

func readFromFolder(dirname string, client myprotodata.SorterMessageServiceClient, firstLineIsHeader bool) {
	var files []string

	folder := "./" + dirname

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	generator := func(done <-chan interface{}, strs ...string) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for _, i := range strs {
				select {
				case <-done:
					return
				case stringStream <- i:
				}
			}
		}()
		return stringStream
	}

	read := func(
		done <-chan interface{},
		intStream <-chan string,
	) <-chan string {
		multipliedStream := make(chan string)

		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				f, err := os.Open(i)
				if err != nil {
					fmt.Println("erorr reading file")
				}
				defer f.Close()

				r := bufio.NewReader(f)

				var skippedFirstLine = !firstLineIsHeader

				for {
					line, err := r.ReadString('\n')

					if !skippedFirstLine {
						skippedFirstLine = !skippedFirstLine

						continue
					}

					select {
						case <-done:
							return
						case multipliedStream <- line:
					}

					if err == io.EOF {
						break
					} else if err != nil {
						fmt.Printf("error reading file %s", err)
						break
					}
				}

			}
		}()
		return multipliedStream
	}

	send := func(
		done <-chan interface{},
		strStream <-chan string,
	) <-chan string {
		addedStream := make(chan string)
		go func() {
			defer close(addedStream)

			for i := range strStream {
				response, err := client.GetRequest(context.Background(), &myprotodata.Request{Action: "ADD", Payload: []string{i}})

				if err != nil {
					log.Fatal("error service", err)
				}
			
				fmt.Println(response.Error)
				fmt.Println(response.Payload)
			}

		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	strsStream := generator(done, files[1:]...)
	pipeline :=  send(done, read(done, strsStream))

	for v := range pipeline {
		fmt.Println(v)
	}

	return 
}

func writeResultToFile(fileNameOutput string, arr []string) {
	f1, err := os.Create(fileNameOutput)

	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}

	defer f1.Close()

	for i := range arr {
		fmt.Fprintf(f1, string([]byte(arr[i])))
	}

	fmt.Println("See results in " + fileNameOutput + " file")
}

func main() {
	minusR := flag.Bool("r", false, "-r Sort input lines in reverse order")
	minusF := flag.Int("f", 1, "-f N Sort input lines by value number N (starts from 1)")
	minusH := flag.Bool("h", false, "-h	The first line is a header that must be ignored during sorting but included in the output")
	minusI := flag.String("i", "", "Use a file with the name file-name as an input.")
	minusO := flag.String("o", "", "Use a file with the name file-name as an output.")
	minusD := flag.String("d", "", "dirname")
	minusS := flag.Bool("s", false, "Start building tree in MS2")
	minusG := flag.Bool("g", false, "Get sorted results")
	minusDel := flag.Bool("del", false, "Discard tree")
	
	flag.Parse()

	valueReversed := *minusR
	columnToSort := *minusF - 1
	firstLineIsHeader := *minusH
	fileName := *minusI
	fileNameOutput := *minusO
	directoryName := *minusD
	startBuilding := *minusS
	getSortedResults := *minusG
	discardResults := *minusDel

	fmt.Println("-r ", valueReversed)
	fmt.Println("-f ", columnToSort)
	fmt.Println("-h ", firstLineIsHeader)
	fmt.Println("-i ", fileName)
	fmt.Println("-o ", fileNameOutput)
	fmt.Println("-d ", directoryName)
	fmt.Println("-s ", startBuilding)
	fmt.Println("-g ", getSortedResults)
	fmt.Println("-del ", discardResults)

	c, err := grpc.Dial("127.0.0.1:4567", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	client := myprotodata.NewSorterMessageServiceClient(c)

	if(discardResults) {
		response, err := client.GetRequest(context.Background(), &myprotodata.Request{Action: "STOP", Payload: nil})

		if err != nil {
			log.Fatal("error service", err)
		}
	
		fmt.Println(response.Error)
		fmt.Println(response.Payload)
	} else if(getSortedResults) {
		response, err := client.GetRequest(context.Background(), &myprotodata.Request{Action: "GET", Payload: nil})

		if err != nil {
			log.Fatal("error service", err)
		}
	
		if len(fileNameOutput) > 0 {
			toPrint := make([]string, 0)

			for i := range response.Payload {
				toPrint = append(toPrint, response.Payload[i])
			}

			writeResultToFile(fileNameOutput, toPrint)
			fmt.Println("find results in output file")
		} else {
			fmt.Println(response.Error)
			fmt.Println(response.Payload)
		}
		
	} else if(startBuilding) {
		
		response, err := client.GetRequest(context.Background(), &myprotodata.Request{Action: "START", Payload: nil})

		if err != nil {
			log.Fatal("error service", err)
		}
	
		fmt.Println(response.Error)
		fmt.Println(response.Payload)
	} else if len(directoryName) == 0 {
		if len(fileName) == 0 {
			readFromStdin(client, firstLineIsHeader)
		} else {
			lineByLine(fileName, client, firstLineIsHeader)
		}
	} else if len(directoryName) > 0 {
		readFromFolder(directoryName, client, firstLineIsHeader)
	}

}
 