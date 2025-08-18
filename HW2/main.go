package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var raw = `
<!DOCTYPE html>
<html>
<body>
<h1>My First Heading</h1>
<p>My first paragraph.</p>
<p>HTML images are defined with the img tag:</p>
<img src="xxx.jpg" width="104" height="142">
</body>
</html>`

func countTree(t *html.Node) (int, int) {
	var tc, ic int

	for node := range t.Descendants() {
		if node.Type == html.TextNode {
			tc += len(strings.Split(node.Data, " "))

		} else if node.Type == html.ElementNode && node.Data == "img" {
			ic++

		}
	}

	return tc, ic

}

func countTreeRecur(t *html.Node, words, pics *int)  {
	

	{
		for node := range t.ChildNodes() {
			if node.Type == html.TextNode {
				
				*words += len(strings.Split(node.Data, " "))

			} else if node.Type == html.ElementNode && node.Data == "img" {
				*pics ++

			}
			countTreeRecur(node, words, pics)

		}
	}

}



func main() {

	var tc, ic int

	raw = strings.ReplaceAll(raw, "\n", "")
	s := strings.NewReader(raw)
	tree, err := html.Parse(s)
	if err != nil {
		fmt.Printf("Unable to parse html, error : %v", err)
		os.Exit(-1)
	}

	tc, ic = countTree(tree)

	fmt.Printf("words count: %d\n", tc)
	fmt.Printf("image count: %d\n ", ic)

	tc, ic = 0,0
	countTreeRecur(tree, &tc, &ic)
	fmt.Printf("words count: %d\n", tc)
		fmt.Printf("image count: %d\n ", ic)
}
