package swf

import (
	"fmt"
	"testing"
)

func TestTree_Match(t *testing.T) {
	// 1.创建前缀树
	tree := NewTree()

	// 2.添加敏感词
	tree.AddWords("she", "his", "hers", "he")

	// 3.构建失配指针
	tree.BuildFail()

	// 打印 tree
	PrintTree(tree, true)

	// 4.匹配
	fmt.Println(tree.Match("ashes"))
	fmt.Println(tree.MatchAll("ashes"))
}
