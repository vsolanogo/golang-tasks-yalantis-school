package main

import (
	"testing"
)

func TestTree(t *testing.T) {
	shouldBeLike := []keyValueStruct{
		keyValueStruct{key: "aaaaaaa", value: "value1"},
		keyValueStruct{key: "bbb", value: "value1"},
		keyValueStruct{key: "ccc", value: "value1"},
		keyValueStruct{key: "gdgggg", value: "value1"},
		keyValueStruct{key: "z", value: "value1"},
		keyValueStruct{key: "zzzzzzzzz", value: "value1"},
	}

	insert(root, shouldBeLike[4])
	insert(root, shouldBeLike[3])
	insert(root, shouldBeLike[5])
	insert(root, shouldBeLike[2])
	insert(root, shouldBeLike[0])
	insert(root, shouldBeLike[1])
	inOrder(root)

	if len(sortedArr) != len(shouldBeLike) {
		t.Error(`len(sortedArr) != length`)
	}

	for i := 0; i < len(sortedArr); i++ {
		if sortedArr[i].key != shouldBeLike[i].key {
			t.Error(`incorrect sort`)
		}
	}
}
