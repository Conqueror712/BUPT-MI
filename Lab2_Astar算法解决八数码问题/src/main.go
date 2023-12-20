package main

import (
	"container/heap"
	"net/http"

	"github.com/gin-gonic/gin"
)

const n int = 3

var start [n * n]int
var target [n * n]int

// FindZeroPosition 找到0的位置
// 返回值为0的位置，如果没有0则返回-1
func FindZeroPosition(a [n * n]int) int {
	for i, v := range a {
		if v == 0 {
			return i
		}
	}
	return -1
}

// HeuristicSearch 启发式搜索
// 返回值为当前状态和目标状态不同的数的个数
func HeuristicSearch(a [n * n]int) int {
	cnt := 0
	for i := 0; i < n*n; i++ {
		if a[i] != target[i] {
			cnt++
		}
	}
	// 由于每次只能交换两个数，所以返回值除以2
	return cnt / 2
}

// Node 节点
// state 为当前状态
// parent 为父节点
// cost 为当前节点的代价
// direction 为当前节点的方向
type Node struct {
	state     [n * n]int
	parent    *Node
	cost      int
	direction byte
}

// NodeHeap 优先队列
// 优先队列的比较函数为当前状态和目标状态不同的数的个数加上当前节点的代价
type NodeHeap []Node

func (h NodeHeap) Len() int { return len(h) }
func (h NodeHeap) Less(i, j int) bool {
	return HeuristicSearch(h[i].state)+h[i].cost < HeuristicSearch(h[j].state)+h[j].cost
}
func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(Node))
}
func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Solve(a [n * n]int) *Node {
	var state map[[n * n]int]bool = make(map[[n * n]int]bool)
	open := &NodeHeap{}
	heap.Init(open)
	var node = new(Node)
	node.parent = nil
	node.cost = 0
	node.state = a
	heap.Push(open, *node)
	state[a] = true
	for open.Len() > 0 {
		node = new(Node)
		*node = heap.Pop(open).(Node)
		if node.state == target {
			return node
		}
		zeroPos := FindZeroPosition(node.state)
		var newNode Node
		newNode.parent = node
		newNode.cost = node.cost + 1
		dx := [4]int{1, 0, -1, 0}
		dy := [4]int{0, 1, 0, -1}
		dp := [4]int{n, 1, -n, -1}
		direction := [4]byte{'U', 'L', 'D', 'R'}

		for i := 0; i < 4; i++ {
			if zeroPos/n+dx[i] >= 0 && zeroPos/n+dx[i] < n && zeroPos%n+dy[i] >= 0 && zeroPos%n+dy[i] < n {
				newNode.state = node.state
				newNode.direction = direction[i]
				newNode.state[zeroPos], newNode.state[zeroPos+dp[i]] = newNode.state[zeroPos+dp[i]], newNode.state[zeroPos]
				if _, ok := state[newNode.state]; ok {
					continue
				}
				heap.Push(open, newNode)
				state[newNode.state] = true
			}
		}
	}
	return nil
}

func Calc(input string) int {
	for i := 0; i < 8; i++ {
		start[i] = int(input[i]) - '0'
	}
	for i := 9; i < 17; i++ {
		target[i-9] = int(input[i]) - '0'
	}
	ans := Solve(start)
	return ans.cost
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	// Serve the HTML form at the root route
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "This data is come from Go background.",
		})
	})

	// var a [9]int

	// Handle form submissions
	router.POST("/submit", func(c *gin.Context) {
		// Retrieve the input value from the form
		inputValue := string(c.PostForm("input"))
		ans := Calc(inputValue)
		// Return a response to the user
		c.HTML(http.StatusOK, "result.html", gin.H{
			"inputValue": inputValue,
			"ans":        ans,
		})
	})

	// Start the web server
	router.Run(":8080")
}
