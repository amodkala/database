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

func (tree RBTree) Search(key int) (data []byte, ok bool)

func (tree RBTree) Insert(key int, data []byte) error { return nil }

func (tree RBTree) Delete(key int) error

func (tree RBTree) leftRotate(node *RBTreeNode) {
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

func (tree RBTree) rightRotate(node *RBTreeNode) {
    x := node.left
    node.left = x.right
    if x.right != nil {
        x.right.parent = node
    }
    x.parent = node.parent
    switch node {
    case tree.root:
        tree.root = x
    case node.parent.left:
        node.parent.left = x
    case node.parent.right:
        node.parent.right = x
    }
    x.right = node
    node.parent = x
}
