package main

import "fmt"

type Node struct {
	size                         int // количество занятых ключей
	key                          [3]int
	first, second, third, fourth *Node
	parent                       *Node // указатель на родителя нужен, т.к. адрес корня может меняться при удалении
}

func NewNode(k int) *Node {
	return &Node{
		size:   1,
		key:    [3]int{k, 0, 0},
		first:  nil,
		second: nil,
		third:  nil,
		fourth: nil,
		parent: nil,
	}
}

func NewNodeWithChildren(k int, firstNew, secondNew, thirdNew, fourthNew, parentNew *Node) *Node {
	return &Node{
		size:   1,
		key:    [3]int{k, 0, 0},
		first:  firstNew,
		second: secondNew,
		third:  thirdNew,
		fourth: fourthNew,
		parent: parentNew,
	}
}

// find возвращает true, если ключ k находится в вершине, иначе false
func (n *Node) find(k int) bool {
	for i := 0; i < n.size; i++ {
		if n.key[i] == k {
			return true
		}
	}
	return false
}

// swap
func (n *Node) swap(x, y *int) {
	*x, *y = *y, *x
}

// sort сортирует ключи в вершине
func (n *Node) sort() {
	if n.size == 1 {
		return
	}
	if n.size == 2 {
		n.sort2(&n.key[0], &n.key[1])
	}
	if n.size == 3 {
		n.sort3(&n.key[0], &n.key[1], &n.key[2])
	}
}

// sort2 вызывается для сортировки 2-х ключей в вершине
func (n *Node) sort2(x, y *int) {
	if *x > *y {
		n.swap(x, y)
	}
}

// sort3 вызывается для сортировки 3-х ключей в вершине
func (n *Node) sort3(x, y, z *int) {
	if *x > *y {
		n.swap(x, y)
	}
	if *x > *z {
		n.swap(x, z)
	}
	if *y > *z {
		n.swap(y, z)
	}
}

// insertToNode вставляет ключ k в вершину
func (n *Node) insertToNode(k int) {
	n.key[n.size] = k
	n.size++
	n.sort()
}

// removeFromNode удаляет ключ k из вершины
func (n *Node) removeFromNode(k int) {
	if n.size >= 1 && n.key[0] == k {
		n.key[0] = n.key[1]
		n.key[1] = n.key[2]
		n.size--
	} else if n.size == 2 && n.key[1] == k {
		n.key[1] = n.key[2]
		n.size--
	}
}

// becomeNode2 преобразовывает в вершину с двумя потомками
func (n *Node) becomeNode2(k int, firstNew, secondNew *Node) {
	n.key[0] = k
	n.first = firstNew
	n.second = secondNew
	n.third = nil
	n.fourth = nil
	n.parent = nil
	n.size = 1
}

// isLeaf проверяет является ли вершина листом
// проверка используется при вставке и удалении
func (n *Node) isLeaf() bool {
	return n.first == nil && n.second == nil && n.third == nil
}

// ------ Вставка ключа ------------------------------------

// insert вставляет ключ k в дерево с корнем p;
// всегда возвращаем корень дерева, т.к. он может меняться
func (n *Node) insert(p *Node, k int) *Node {
	// если дерево пусто, то создаем первую 2-3 вершину (корень)
	if p == nil {
		return NewNode(k)
	}

	if p.isLeaf() {
		p.insertToNode(k)
	} else if k <= p.key[0] {
		p.insert(p.first, k)
	} else if p.size == 1 || (p.size == 2 && k <= p.key[1]) {
		p.insert(p.second, k)
	} else {
		p.insert(p.third, k)
	}

	return p.split(p)
}

// --------- Разделение вершины ----------------------
// split вызывается при вставке ключа в вершину; разделение необходимо, когда количество ключей в вершине = 3
func (n *Node) split(item *Node) *Node {
	if item.size < 3 {
		return item
	}
	// создаем две новые вершины, которые имеют такого же родителя, как и разделяющий элемент
	var x *Node = NewNodeWithChildren(item.key[0], item.first, item.second, nil, nil, item.parent)
	var y *Node = NewNodeWithChildren(item.key[2], item.third, item.fourth, nil, nil, item.parent)

	// переустанавливаем "родителей" и "сыновей"
	if x.first != nil {
		x.first.parent = x
	}
	if x.second != nil {
		x.second.parent = x
	}
	if y.first != nil {
		y.first.parent = y
	}
	if y.second != nil {
		y.second.parent = y
	}

	if item.parent != nil {
		item.parent.insertToNode(item.key[1])

		if item.parent.first == item {
			item.parent.first = nil
		} else if item.parent.second == item {
			item.parent.second = nil
		} else if item.parent.third == item {
			item.parent.third = nil
		}
		// Перераспраделение ключей при разделении
		if item.parent.first == nil {
			item.parent.fourth = item.parent.third
			item.parent.third = item.parent.second
			item.parent.second = y
			item.parent.first = x
		} else if item.parent.second == nil {
			item.parent.fourth = item.parent.third
			item.parent.third = y
			item.parent.second = x
		} else {
			item.parent.fourth = y
			item.parent.third = x
		}

		var tmp *Node = item.parent
		item = nil
		return tmp
	} else {
		x.parent = item
		y.parent = item
		item.becomeNode2(item.key[1], x, y)
		return item
	}
}

// ------ Поиск ключа k в 2-3 дереве с корнем p ------------------------------------
// search ищет Node, в которой находится ключ и возвращает ее, если ключа нет - возвращает nil
// вызывается при удалении ключа
func (n *Node) search(p *Node, k int) *Node {
	if p == nil {
		return nil
	}

	if p.find(k) {
		return p
	} else if k < p.key[0] {
		return n.search(p.first, k)
	} else if p.size == 2 && k < p.key[1] || p.size == 1 {
		return n.search(p.second, k)
	} else if p.size == 2 {
		return n.search(p.third, k)
	}
	return nil
}

// ------ Поиск узла с минимальным элементом в 2-3 дереве с корнем p ------------------------------------

// searchMin вызывается при удалении
// находит эквивалентный ключ из листовой вершины для удаляемого ключа
func (n *Node) searchMin(p *Node) *Node {
	if p == nil {
		return nil
	}

	if p.first == nil {
		return p
	} else {
		return n.searchMin(p.first)
	}
}

// ------ Удаление ключа ------------------------------------

// remove удаляет ключ из 2-3 дерева с корнем p.
func (n *Node) remove(p *Node, k int) *Node {
	var item *Node = n.search(p, k) // Ищем узел, где находится ключ k

	if item == nil {
		return p
	}

	var minNode *Node = nil
	if item.key[0] == k {
		minNode = n.searchMin(item.second) // Ищем Node с эквивалентным ключом (минимальный элемент в правом поддереве)
	} else {
		minNode = n.searchMin(item.third)
	}

	if minNode != nil { // меняем ключи местами
		var z *int
		if k == item.key[0] {
			z = &item.key[0]
		} else {
			z = &item.key[1]
		}
		item.swap(z, &minNode.key[0])
		item = minNode // Перемещаем указатель на лист, т.к. min всегда лист
	}
	item.removeFromNode(k) // И удаляем требуемый ключ из листа
	return p.fix(item)
}

// Исправление дерева после удаления ключа

func (n *Node) fix(leaf *Node) *Node {
	// Случай 1, когда удаляется единственный ключ в дереве
	if leaf.size == 0 && leaf.parent == nil {
		leaf = nil
		return nil
	}
	// Случай 2, когда вершина, в которой удалили ключ, имела два ключа
	if leaf.size != 0 {
		if leaf.parent != nil {
			return n.fix(leaf.parent)
		} else {
			return leaf
		}
	}

	parent := leaf.parent
	// Случай 3, когда достаточно перераспределить ключи в дереве
	if parent.first.size == 2 || parent.second.size == 2 || parent.size == 2 {
		leaf = n.redistribute(leaf)
	} else if parent.size == 2 && parent.third.size == 2 {
		leaf = n.redistribute(leaf)
	} else {
		// Случай 4, когда нужно произвести склеивание и пройтись вверх по дереву как минимум на еще одну вершину
		leaf = n.merge(leaf)
	}
	return n.fix(leaf)
}

// -------- Перераспределение ---------------------------------

func (n *Node) redistribute(leaf *Node) *Node {
	parent := leaf.parent
	first := parent.first
	second := parent.second
	third := parent.third

	if parent.size == 2 && first.size < 2 && second.size < 2 && third.size < 2 {
		if first == leaf {
			parent.first = parent.second
			parent.second = parent.third
			parent.third = nil

			parent.first.insertToNode(parent.key[0])
			parent.first.third = parent.first.second
			parent.first.second = parent.first.first

			if leaf.first != nil {
				parent.first.first = leaf.first
			} else if leaf.second != nil {
				parent.first.first = leaf.second
			}

			if parent.first.first != nil {
				parent.first.first.parent = parent.first
			}

			parent.removeFromNode(parent.key[0])
			first = nil
		} else if second == leaf {
			first.insertToNode(parent.key[0])
			parent.removeFromNode(parent.key[0])
			if leaf.first != nil {
				first.third = leaf.first
			} else if leaf.second != nil {
				first.third = leaf.second
			}

			if first.third != nil {
				first.third.parent = first
			}

			parent.second = parent.third
			parent.third = nil

			second = nil
		} else if third == leaf {
			second.insertToNode(parent.key[1])
			parent.third = nil
			parent.removeFromNode(parent.key[1])
			if leaf.first != nil {
				second.third = leaf.first
			} else if leaf.second != nil {
				second.third = leaf.second
			}

			if second.third != nil {
				second.third.parent = second
			}

			third = nil
		}
	} else if parent.size == 2 && (first.size == 2 || second.size == 2 || third.size == 2) {
		if third == leaf {
			if leaf.first != nil {
				leaf.second = leaf.first
				leaf.first = nil
			}

			leaf.insertToNode(parent.key[1])
			if second.size == 2 {
				parent.key[1] = second.key[1]
				second.removeFromNode(second.key[1])
				leaf.first = second.third
				second.third = nil
				if leaf.first != nil {
					leaf.first.parent = leaf
				}
			} else if first.size == 2 {
				parent.key[1] = second.key[0]
				leaf.first = second.second
				second.second = second.first
				if leaf.first != nil {
					leaf.first.parent = leaf
				}

				second.key[0] = parent.key[0]
				parent.key[0] = first.key[1]
				first.removeFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		} else if second == leaf {
			if third.size == 2 {
				if leaf.first == nil {
					leaf.first = leaf.second
					leaf.second = nil
				}
				second.insertToNode(parent.key[1])
				parent.key[1] = third.key[0]
				third.removeFromNode(third.key[0])
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			} else if first.size == 2 {
				if leaf.second == nil {
					leaf.second = leaf.first
					leaf.first = nil
				}
				second.insertToNode(parent.key[0])
				parent.key[0] = first.key[1]
				first.removeFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		} else if first == leaf {
			if leaf.first == nil {
				leaf.first = leaf.second
				leaf.second = nil
			}
			first.insertToNode(parent.key[0])
			if second.size == 2 {
				parent.key[0] = second.key[0]
				second.removeFromNode(second.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = second.third
				second.third = nil
			} else if third.size == 2 {
				parent.key[0] = second.key[0]
				second.key[0] = parent.key[1]
				parent.key[1] = third.key[0]
				third.removeFromNode(third.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			}
		}
	} else if parent.size == 1 {
		leaf.insertToNode(parent.key[0])

		if first == leaf && second.size == 2 {
			parent.key[0] = second.key[0]
			second.removeFromNode(second.key[0])

			if leaf.first == nil {
				leaf.first = leaf.second
			}

			leaf.second = second.first
			second.first = second.second
			second.second = second.third
			second.third = nil
			if leaf.second != nil {
				leaf.second.parent = leaf
			}
		} else if second == leaf && first.size == 2 {
			parent.key[0] = first.key[1]
			first.removeFromNode(first.key[1])

			if leaf.second == nil {
				leaf.second = leaf.first
			}

			leaf.first = first.third
			first.third = nil
			if leaf.first != nil {
				leaf.first.parent = leaf
			}
		}
	}
	return parent
}

// ------------ Склеивание ---------------------------------
func (n *Node) merge(leaf *Node) *Node {
	parent := leaf.parent

	if parent.first == leaf {
		parent.second.insertToNode(parent.key[0])
		parent.second.third = parent.second.second
		parent.second.second = parent.second.first

		if leaf.first != nil {
			parent.second.first = leaf.first
		} else if leaf.second != nil {
			parent.second.first = leaf.second
		}

		if parent.second.first != nil {
			parent.second.first.parent = parent.second
		}

		parent.removeFromNode(parent.key[0])
		parent.first = nil

	} else if parent.second == leaf {
		parent.first.insertToNode(parent.key[0])

		if leaf.first != nil {
			parent.first.third = leaf.first
		} else if leaf.second != nil {
			parent.first.third = leaf.second
		}

		if parent.first.third != nil {
			parent.first.third.parent = parent.first
		}

		parent.removeFromNode(parent.key[0])
		parent.second = nil
	}

	if parent.parent == nil {
		var tmp *Node = nil
		if parent.first != nil {
			tmp = parent.first
		} else {
			tmp = parent.second
		}
		tmp.parent = nil
		parent = nil
		return tmp
	}
	return parent
}

// -------- Tree ------------------------------------

type Tree struct {
	root *Node
}

func (t *Tree) insert(k int) {
	t.root = t.root.insert(t.root, k)
}

func (t *Tree) remove(k int) {
	t.root = t.root.remove(t.root, k)
}

func (t *Tree) search(k int) bool {
	if t.root.search(t.root, k) != nil {
		return true
	}
	return false
}

func (t *Tree) printTree() {
	var printNode func(*Node)
	printNode = func(p *Node) {
		if p != nil {
			printNode(p.first)
			fmt.Printf("%d ", p.key[0])
			printNode(p.second)
			if p.size == 2 {
				fmt.Printf("%d ", p.key[1])
				printNode(p.third)
			}
		}
	}
	printNode(t.root)
	fmt.Println()
}

func main() {
	tree := &Tree{}

	tree.insert(10)
	tree.insert(20)
	tree.insert(5)
	tree.insert(6)
	tree.insert(12)
	tree.insert(30)
	tree.insert(7)
	tree.insert(17)

	fmt.Print("Tree after insertion: ")
	tree.printTree()

	fmt.Println("Tree has key 6: ", tree.search(6))
	tree.remove(6)
	fmt.Print("Tree after removing 6: ")
	tree.printTree()
	fmt.Println("Tree has key 6: ", tree.search(6))

	tree.remove(13)
	fmt.Print("Tree after removing 13: ")
	tree.printTree()

	tree.remove(7)
	fmt.Print("Tree after removing 7: ")
	tree.printTree()

	tree.remove(4)
	fmt.Print("Tree after removing 4: ")
	tree.printTree()
}
