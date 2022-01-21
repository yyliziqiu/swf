package swf

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/**
创建一个空的 tree
*/
func NewTree() *Tree {
	return &Tree{
		Root: &Node{
			Index:     0,
			Character: '/',
			Fail:      nil,
			Parent:    nil,
			Children:  map[rune]*Node{},
			IsRoot:    true,
			IsEnd:     false,
			LnEnd:     false,
			Depth:     0,
		},
		CCnt: 0,
	}
}

/**
输出 tree，调试用
*/
func PrintTree(tree *Tree, pretty bool) {
	printTree(tree.Root, "", pretty)
}

/**
输出 tree，调试用
*/
func printTree(node *Node, f string, pretty bool) {
	if node.IsRoot {
		fmt.Println(string(node.Character))
	}

	if len(node.Children) == 0 {
		return
	}

	i := 0
	for _, child := range node.Children {
		line := string(child.Character) + " " + strconv.Itoa(child.Index)
		if child.Fail != nil && !child.Fail.IsRoot {
			if pretty {
				line += " \033[1;47;31m" + strconv.Itoa(child.Fail.Index) + "\033[0m"
			} else {
				line += " " + strconv.Itoa(child.Fail.Index)
			}
		}
		if child.IsEnd {
			if pretty {
				line += " \033[1;47;31mEND\033[0m"
			} else {
				line += " END"
			}
		}

		if i = i + 1; i == len(node.Children) {
			fmt.Println(f + "└──" + line)
			printTree(child, f+"   ", pretty)
		} else {
			fmt.Println(f + "│──" + line)
			printTree(child, f+"│  ", pretty)
		}
	}
}

/**
优化字符串
*/
func PromoteStr(str string) (string, string) {
	var (
		characters = []rune(str)
		sb1        strings.Builder
		sb2        strings.Builder
	)

	for key, value := range characters {
		// 大写转小写
		if unicode.IsUpper(value) {
			value = unicode.ToLower(value)
		}

		// 合并相邻空格
		if key > 0 && value == ' ' && characters[key-1] == ' ' {
			continue
		}

		// 提取中文、英文、空格
		IsHan := unicode.Is(unicode.Han, value)
		IsEng := 'a' <= value && value <= 'z'
		if !IsHan && !IsEng && value != ' ' {
			continue
		}

		sb1.WriteRune(value)
		if IsHan {
			sb2.WriteRune(value)
		}
	}

	return sb1.String(), sb2.String()
}
