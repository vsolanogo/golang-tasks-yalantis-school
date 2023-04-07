package main

import (
	"fmt"
	"net" 
	"task3/myprotodata"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

type sorterServer struct{}

func (k *sorterServer) GetRequest(ctx context.Context, request *myprotodata.Request) (*myprotodata.Response, error) {
	if(request.Action == "ADD") {
		response := &myprotodata.Response{}
		response.Error = ""
		response.Payload = []string{""}

		for i := range request.Payload {
			fmt.Println(request.Payload[i])
			key := strings.Split(request.Payload[i], ",")[0]
			insert(root, keyValueStruct{key, request.Payload[i]})
		}

		return response, nil
	}

	if(request.Action == "START") {
		inOrder(root)

		response := &myprotodata.Response{}
		response.Error = ""
		response.Payload = []string{"Performed sort"}

		return response, nil
	}
	
	if(request.Action == "GET") {
		fmt.Println("Send results")
		var sortedArrToSend []string;

		for i := range sortedArr {
			sortedArrToSend = append(sortedArrToSend, sortedArr[i].value)
		}

		response := &myprotodata.Response{}
		response.Error = ""
		response.Payload = sortedArrToSend

		return response, nil
	}	

	if(request.Action == "STOP") {
		fmt.Println("Discard tree")
	
		root = nil
		sortedArr = nil

		response := &myprotodata.Response{}
		response.Error = ""
		response.Payload = []string{""}

		return response, nil
	}

	response := &myprotodata.Response{}
	response.Error = ""
	response.Payload = []string{""}

	return response, nil
}

func main() {
	l, err := net.Listen("tcp", ":4567")

	if err != nil {
		fmt.Println(err)
	}

	grpcServer := grpc.NewServer()
	myprotodata.RegisterSorterMessageServiceServer(grpcServer, &sorterServer{})
	grpcServer.Serve(l)
}