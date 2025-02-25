# GoDS

This repository contains a collection of **thread-safe**, **generic** data structures for Go. Each data structure is implemented in a separate package and includes an ELI5 (Explain Like I'm 5) section and example usage.

## Data Structures

### Fundamental Data

- [Stack](./pkg/stack/) (LIFO stack)
- [Queue](./pkg/queue/) (FIFO queue)
- [Deque](./pkg/deque/) (Double-ended queue)
- [Priority Queue](./pkg/priorityqueue/) (heap-based)

### Linked

- Singly Linked List
- Doubly Linked List
- Circular Linked List

### Set-Based

- Hash Set (unordered, unique elements)
- Ordered Set (tree-based, sorted unique elements)

### Map-Based

- Hash Map (key-value storage)
- Ordered Map (sorted key-value storage)

### Tree

- Binary Search Tree (BST)
- AVL Tree (self-balancing BST)
- Red-Black Tree (used in many real-world DBs)
- Trie (Prefix Tree) (useful for autocompletion/search)

### Graph

- Adjacency List Graph
- Adjacency Matrix Graph
- Directed & Undirected Graphs

### Specialized

- Bloom Filter (probabilistic set membership test)
- LRU Cache (least recently used cache)
- Skip List (alternative to BSTs)
