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
    mu *sync.RWMutex
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

func New() RBTree {
    return RBTree{
        root: nil,
    }
}

func (tree RBTree) Search(key int) (data []byte, ok bool) {
    if tree.root == nil {
        return []byte{}, false
    }

    tree.mu.RLock()
    defer tree.mu.RUnlock()

    curr := tree.root
    for curr != nil && curr.key != key {
        if key < curr.key {
            curr = curr.left
        } else {
            curr = curr.right
        }
    }

    if curr == nil {
        return nil, false
    }
    return curr.data, true
}

func (tree RBTree) Insert(key int, data []byte) { 


    if tree.root == nil {
        tree.root = &RBTreeNode{
            key: key,
            data: data,
            colour: red,
            left: nil,
            right: nil,
            parent: nil,
        }
        return
    }

    tree.mu.Lock()
    defer tree.mu.Unlock()

    var temp *RBTreeNode
    var child string

    for curr := tree.root; curr != nil; {

        if curr.key == key {
            curr.data = data
            return
        }

        temp = curr
        if key < curr.key {
            curr = curr.left
            child = "left"
        } else {
            curr = curr.right
            child = "right"
        }
    }

    node := &RBTreeNode{
        key: key,
        data: data,
        colour: red,
        left: nil,
        right: nil,
        parent: nil,
    }

    node.parent = temp
    switch child {
    case "left":
        temp.left = node
    case "right":
        temp.right = node
    }

    tree.rebalance(node)
}

func (tree RBTree) Delete(key int) error

func (tree RBTree) rebalance(node *RBTreeNode) {

    for colour := node.parent.colour; colour == red; {
        if node.parent == node.parent.parent.left {
        }
    }

    tree.root.colour = black
}

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

func (tree RBTree) isBalanced() bool {
    if tree.root == nil {
        return true
    }
    return true
}

