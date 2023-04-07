package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type keyValueStruct struct {
	key   string
	value string
}

var sortedArr []keyValueStruct

type Node struct {
	data  keyValueStruct
	left  *Node
	right *Node
}

var root *Node = nil

func createNewNode(newData keyValueStruct) *Node {
	var newNode *Node = new(Node)
	newNode.data = newData
	newNode.left = nil
	newNode.right = nil
	return newNode
}

func inOrder(root *Node) {
	if root == nil {
		return
	}

	inOrder(root.left)

	sortedArr = append(sortedArr, root.data)
	inOrder(root.right)
}

func insert(node *Node, newData keyValueStruct) {
	if root == nil {
		root = &Node{data: newData, left: nil, right: nil}
		return
	}

	var compareValue = strings.Compare(newData.key, node.data.key)

	if compareValue == 0 {
		compareValue = strings.Compare(newData.key+newData.value, node.data.key)
	}

	if compareValue < 0 {
		if node.left == nil {
			node.left = createNewNode(newData)
		} else {
			insert(node.left, newData)
		}
	} else if compareValue > 0 {
		if node.right == nil {
			node.right = createNewNode(newData)
		} else {
			insert(node.right, newData)
		}
	}
}

func reverseArray(myArr []keyValueStruct) []keyValueStruct {
	for i := 0; i < len(myArr)/2; i++ {
		j := len(myArr) - i - 1
		myArr[i], myArr[j] = myArr[j], myArr[i]
	}
	return myArr
}

func lineByLine(file string) []string {
	var err error
	inputLines := make([]string, 0)

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("erorr reading file")
	}
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		inputLines = append(inputLines, line)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
	}

	return inputLines
}

func readFromStdin() []string {
	fmt.Println("input your comma-separated values:")
	inputLines := make([]string, 0)

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		if text == "\n" {
			break
		}

		inputLines = append(inputLines, text)
	}

	return inputLines
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

	flag.Parse()

	valueReversed := *minusR
	columnToSort := *minusF - 1
	firstLineIsHeader := *minusH
	fileName := *minusI
	fileNameOutput := *minusO
	directoryName := *minusD

	fmt.Println("-r ", valueReversed)
	fmt.Println("-f ", columnToSort)
	fmt.Println("-h ", firstLineIsHeader)
	fmt.Println("-i ", fileName)
	fmt.Println("-o ", fileNameOutput)
	fmt.Println("-d ", directoryName)

	var startNormalizingFromLine int

	if firstLineIsHeader {
		startNormalizingFromLine = 1
	} else {
		startNormalizingFromLine = 0
	}

	toPrint := make([]string, 0)

	if len(directoryName) == 0 {
		inputLines := make([]string, 0)

		if len(fileName) == 0 {
			inputLines = readFromStdin()
		} else {
			inputLines = lineByLine(fileName)
		}

		columnsCount := len(strings.Split(inputLines[0], ","))

		fmt.Println("columns count we'll work with:", columnsCount)

		normalizedLines := make([]keyValueStruct, 0)

		for i := startNormalizingFromLine; i < len(inputLines); i++ {
			if len(strings.Split(inputLines[i], ",")) == columnsCount {
				key := strings.Split(inputLines[i], ",")[columnToSort]

				normalizedLines = append(normalizedLines, keyValueStruct{key, inputLines[i]})
			}
		}

		for i := range normalizedLines {
			insert(root, normalizedLines[i])
		}

		inOrder(root)

		if valueReversed {
			sortedArr = reverseArray(sortedArr)
		}

		if len(fileNameOutput) == 0 {
			if firstLineIsHeader {
				fmt.Println(inputLines[0])
			}

			for i := range sortedArr {
				fmt.Printf(sortedArr[i].value)
			}
		} else {

			if firstLineIsHeader {
				toPrint = append(toPrint, inputLines[0])
			}

			for i := range sortedArr {
				toPrint = append(toPrint, sortedArr[i].value)
			}

			writeResultToFile(fileNameOutput, toPrint)
		}

		fmt.Printf("broken csv lines: ")
		fmt.Println(len(inputLines) - len(normalizedLines))

	} else {
		var files []string

		folder := "./" + directoryName

		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})

		if err != nil {
			panic(err)
		}

		generator := func(done <-chan interface{}, integers ...string) <-chan string {
			stringStream := make(chan string)
			go func() {
				defer close(stringStream)
				for _, i := range integers {
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

					for {
						line, err := r.ReadString('\n')

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

		add := func(
			done <-chan interface{},
			intStream <-chan string,
		) <-chan string {
			addedStream := make(chan string)
			go func() {
				defer close(addedStream)

				for i := range intStream {
					insert(root, keyValueStruct{strings.Split(i, ",")[columnToSort], i})
				}

			}()
			return addedStream
		}

		done := make(chan interface{})
		defer close(done)

		strsStream := generator(done, files[1:]...)
		pipeline := add(done, read(done, strsStream))

		for range pipeline {
		}

		inOrder(root)
		
		if valueReversed {
			sortedArr = reverseArray(sortedArr)
		}

		if len(fileNameOutput) == 0 {
			for i := range sortedArr {
				fmt.Printf(sortedArr[i].value)
			}
		} else {
			for i := range sortedArr {
				toPrint = append(toPrint, sortedArr[i].value)
			}

			writeResultToFile(fileNameOutput, toPrint)
		}

	}

}
