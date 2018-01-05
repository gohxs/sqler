package main

import (
	"log"
	"strings"
)

type Node struct {
	children map[string]*Node
}

func NewNode() *Node {
	return &Node{map[string]*Node{}}
}

func (n *Node) Add(s string) {
	if s == "" {
		return
	}
	var first, rest string
	res := strings.SplitN(s, " ", 2)
	log.Printf("Res: %#v", res)
	first = res[0]
	if len(res) > 1 {
		rest = res[1]
	}

	var childNode *Node
	var ok bool
	childNode, ok = n.children[first]
	if !ok {
		log.Println("Creating node")
		childNode = NewNode()
	}
	log.Println("Adding rest")
	childNode.Add(rest)
	n.children[first] = childNode

}

func main() {

	root := NewNode()
	phrases := []string{
		"select 123",
		"select 123 456",
		"select 456",
	}

	for _, v := range phrases {
		root.Add(v)
	}
	log.Printf("Res: %#v", root)

}
