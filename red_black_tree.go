package rbtree

import (
    "sync"
)

type Colour int

const (
    red = iota
    black
)

type RBTree struct {
    mu *sync.Mutex
    root *RBTreeNode
}

type RBTreeNode struct {
    key int
    data []byte
    colour Colour
    left *RBTreeNode
    right *RBTreeNode
    parent *RBTreeNode
}

func Search(tree RBTree, key int) (data []byte, ok bool)
func Insert(tree RBTree, key int, data []byte) error
func Delete(tree RBTree, key int) error

func leftRotate(tree RBTree, node *RBTreeNode) {
    tree.mu.Lock()
    defer tree.mu.Unlock()

    y := node.right
    node.right = y.left
    if y.left != nil {
        y.left.parent = node
    }
    y.parent = node.parent
    switch node {
    case tree.root:
        tree.root = y
    case node.parent.left:
        node.parent.left = y
    case node.parent.right:
        node.parent.right = y
    }
    y.left = node
    node.parent = y
}

func rightRotate(tree RBTree, node *RBTreeNode) {
    tree.mu.Lock()
    defer tree.mu.Unlock()
}
