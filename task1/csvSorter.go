package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type keyValueStruct struct {
	key string
	value string
}

var sortedArr []keyValueStruct

type Node struct {
	 data keyValueStruct
	 left *Node
	 right *Node
 }

var root *Node = nil

func createNewNode(newData keyValueStruct) *Node {
	// type can be omitted
	 var newNode *Node = new(Node)
	 newNode.data = newData
	 newNode.left = nil
	 newNode.right = nil
	 return newNode
 }

func inOrder(root *Node) {
	if root == nil{
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

	 // check node, can be nil
	 /**
	 	if node == nil {

	 	}
	  */


	 var compareValue = strings.Compare(newData.key, node.data.key)

	 // redundant parentheses
	 if(compareValue == 0 ) {
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


// return error to top level
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
		// it looks better then if else
		/**
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		 */
	}

	return inputLines
}
 	

func main() {
	minusR := flag.Bool("r", false, "-r Sort input lines in reverse order")
	minusF := flag.Int("f", 1, "-f N Sort input lines by value number N (starts from 1)")
	minusH := flag.Bool("h", false, "-h	The first line is a header that must be ignored during sorting but included in the output")
	minusI := flag.String("i", "", "Use a file with the name file-name as an input.")
	minusO := flag.String("o", "", "Use a file with the name file-name as an output.")

	// ignore errors, you will receive a panic in err case
	/**
	CommandLine.Parse(os.Args[1:])
	 */
	flag.Parse()

	valueReversed := *minusR
	columnToSort := *minusF - 1
	firstLineIsHeader := *minusH
	fileName := *minusI
	fileNameOutput := *minusO

	// if code useless, remove it, or leave the comment
	// fmt.Println("-r ", valueReversed)
	// fmt.Println("-f ", columnToSort)
	// fmt.Println("-h ", firstLineIsHeader)
	// fmt.Println("-i ", fileName)
	// fmt.Println("-o ", fileNameOutput)

	inputLines := make([]string, 0)

	// redundant parentheses
	if(len(fileName) == 0) {
		fmt.Println("input your comma-separated values:")

		reader := bufio.NewReader(os.Stdin)

		for {
			// do not ignore error
			text, _ := reader.ReadString('\n')
			/**
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("could not read value from cmd, err: ", err)
					os.Exit(1)
				}
			 */


			if text=="\n" {
				// redundant semicolon
				break;
			}

			inputLines = append(inputLines, text)
		}

	} else {
		inputLines = lineByLine(fileName)
	}

	columnsCount := len(strings.Split(inputLines[0], ","))

	fmt.Println("columns count we'll work with:", columnsCount)

	normalizedLines := make([]keyValueStruct, 0)

	var startNormalizingFromLine int

	// redundant parentheses
	// for bool variable use template that start with is: is1stLineHeader
	if(firstLineIsHeader) {
		startNormalizingFromLine = 1
	} else {
		startNormalizingFromLine = 0
	}

	// the better way to avoid if with else statement
	/**
	startNormalizingFromLine := 0
	if firstLineIsHeader {
		startNormalizingFromLine
	}
	 */

	for i:= startNormalizingFromLine; i< len(inputLines); i++ {
		// declare variable inputLines[i], it will look more pretty
		if len(strings.Split(inputLines[i], ",")) == columnsCount {
			key:= strings.Split(inputLines[i], ",")[columnToSort]

			normalizedLines = append(normalizedLines,keyValueStruct{key, inputLines[i]} )
		}
	}

	for i:= range normalizedLines{ 
		insert(root, normalizedLines[i])
	}

	inOrder(root)


	// redundant parentheses
	if(valueReversed) {
		sortedArr = reverseArray(sortedArr)
	}

	// try to avoid big if else statements
	// redundant parentheses
	if(len(fileNameOutput) == 0) {
		// redundant parentheses
		if(firstLineIsHeader) {
			fmt.Println(inputLines[0])
		}  
		
		for i:= range sortedArr{ 
			fmt.Printf(sortedArr[i].value)
		}
	} else {
		f1, err := os.Create(fileNameOutput)
		
		if err != nil {
			fmt.Println("Cannot create file", err)
			// exit with error
			return
		}

		defer f1.Close()

		// redundant parentheses
		if(firstLineIsHeader) {
			// unhandled error
			fmt.Fprintf(f1, string([]byte(inputLines[0])))
		}

		for i:= range sortedArr{
			// unhandled error
			fmt.Fprintf(f1, string([]byte(sortedArr[i].value)))
		}

		fmt.Println("See results in " + fileNameOutput + " file")
	}

	fmt.Printf("broken csv lines: ")
	fmt.Println(len(inputLines) - len(normalizedLines))
}

