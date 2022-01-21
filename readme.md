### 使用示例
```golang
package main

import (
    "fmt"

    "test/swf"
)

func main() {
	// 1.创建前缀树
	tree := swf.NewTree()

	// 2.添加敏感词
	tree.AddWords("she", "his", "hers", "he")

	// 3.构建失配指针
	tree.BuildFail()

	// 打印 tree
	swf.PrintTree(tree, true)

	// 4.匹配
	fmt.Println(tree.Match("bshes"))
	fmt.Println(tree.MatchAll("bshes"))
}
```
