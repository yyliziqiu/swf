package swf

import (
	"bufio"
	"io"
	"os"
)

/**
node
*/
type Node struct {
	Index     int            // 节点序号。不一定连续
	Character rune           // 节点字符
	Fail      *Node          // 失配节点
	Parent    *Node          // 父节点
	Children  map[rune]*Node // 子节点
	IsRoot    bool           // 是否为根节点
	IsEnd     bool           // 是否为单词尾部字符节点
	LnEnd     bool           //	是否本节点或失配节点链中存在单词尾部字符节点
	Depth     int            // 当前节点深度。从 0 开始
}

/**
添加子节点
*/
func (node *Node) AddChild(index int, character rune, isEnd bool) (*Node, bool) {
	if child := node.GetChild(character); child != nil {
		if isEnd {
			child.IsEnd = isEnd
			child.LnEnd = isEnd
		}
		return child, false
	}

	node.Children[character] = &Node{
		Index:     index,
		Character: character,
		Fail:      nil,
		Parent:    node,
		Children:  map[rune]*Node{},
		IsRoot:    false,
		IsEnd:     isEnd,
		LnEnd:     isEnd,
		Depth:     node.Depth + 1,
	}

	return node.GetChild(character), true
}

/**
获取子节点
*/
func (node *Node) GetChild(character rune) *Node {
	if child, ok := node.Children[character]; ok {
		return child
	}
	return nil
}

/**
获取子节点，包括失配节点的子节点
*/
func (node *Node) GetChildOrFailChild(character rune) *Node {
	child := node.GetChild(character)
	if child != nil {
		return child
	}

	if node.IsRoot {
		return node
	}

	return node.Fail.GetChildOrFailChild(character)
}

/**
获取失配节点链中首个单词尾部字符节点
*/
func (node *Node) GetFirstEndNode() *Node {
	if node.IsEnd {
		return node
	}

	fail := node
	for !fail.IsRoot {
		if fail.Fail.IsEnd {
			return fail.Fail
		}
		fail = fail.Fail
	}

	return nil
}

/**
获取失配节点链中单词尾部字符节点
*/
func (node *Node) GetEndNodes() []*Node {
	fails := make([]*Node, 0)

	if node.IsEnd {
		fails = append(fails, node)
	}

	fail := node
	for !fail.IsRoot {
		if fail.Fail.IsEnd {
			fails = append(fails, fail.Fail)
		}
		fail = fail.Fail
	}

	return fails
}

/**
返回当前节点到根节点之间的字符
*/
func (node *Node) GetString() string {
	var (
		parent     = node
		characters = make([]rune, node.Depth, node.Depth)
	)

	for !parent.IsRoot {
		characters[parent.Depth-1] = parent.Character
		parent = parent.Parent
	}

	return string(characters)
}

/**
tree
*/
type Tree struct {
	Root *Node // 根节点
	CCnt int   // 添加字符次数
}

/**
构建失配指针。BFS
*/
func (tree *Tree) BuildFail() {
	q := &Queue{}
	q.Push(tree.Root)
	tree.Root.Fail = tree.Root
	for !q.IsEmpty() {
		node := q.Pop()
		for _, child := range node.Children {
			if node.IsRoot {
				child.Fail = node
			} else {
				fail := node.Fail
				for child.Fail == nil {
					if valid := fail.GetChild(child.Character); valid != nil {
						child.Fail = valid
						if valid.LnEnd {
							child.LnEnd = true
						}
					} else if fail.IsRoot {
						child.Fail = fail
					} else {
						fail = fail.Fail
					}
				}
			}
			q.Push(child)
		}
	}
}

/**
匹配敏感词
*/
func (tree *Tree) Match(text string) (bool, string) {
	var (
		node       = tree.Root
		child      *Node
		characters = []rune(text)
	)

	for _, character := range characters {
		child = node.GetChildOrFailChild(character)
		if child.LnEnd {
			if endNode := child.GetFirstEndNode(); endNode != nil {
				return true, endNode.GetString()
			}
		}
		node = child
	}

	return false, ""
}

/**
匹配敏感词
*/
func (tree *Tree) MatchAll(text string) (bool, []string) {
	var (
		node    = tree.Root
		child   *Node
		result  []string
		existed = make(map[*Node]struct{})
	)

	for _, character := range []rune(text) {
		child = node.GetChildOrFailChild(character)
		if child.LnEnd {
			endNodes := child.GetEndNodes()
			for _, fail := range endNodes {
				if _, ok := existed[fail]; !ok {
					result = append(result, fail.GetString())
					existed[fail] = struct{}{}
				}
			}
		}
		node = child
	}

	if len(result) == 0 {
		return false, result
	}

	return true, result
}

/**
向 tree 添加多个单词
*/
func (tree *Tree) AddWords(words ...string) {
	for _, word := range words {
		tree.AddWord(word)
	}
}

/**
向 tree 添加单词
*/
func (tree *Tree) AddWord(word string) {
	node := tree.Root
	characters := []rune(word)
	for i, character := range characters {
		tree.CCnt++
		node, _ = node.AddChild(tree.CCnt, character, i == len(characters)-1)
	}
}

/**
从文件中添加敏感词
*/
func (tree *Tree) AddWordsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		tree.AddWord(string(line))
	}

	return nil
}

/**
队列。用于 BFS
*/
type Queue struct {
	s []*Node
}

/**
向队列中添加一个元素
*/
func (q *Queue) Push(node *Node) {
	q.s = append(q.s, node)
}

/**
取出队列首部元素
*/
func (q *Queue) Pop() *Node {
	n := q.s[0]
	q.s = q.s[1:]

	return n
}

/**
获取队列长度
*/
func (q *Queue) Len() int {
	return len(q.s)
}

/**
判断队列是否为空
*/
func (q *Queue) IsEmpty() bool {
	return q.Len() == 0
}
