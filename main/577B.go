package main

import (
	"bufio"
	. "fmt"
	"io"
)

// github.com/EndlessCheng/codeforces-go
func Sol577B(reader io.Reader, writer io.Writer) {
	max := func(a, b int) int {
		if a >= b {
			return a
		}
		return b
	}
	const inf int = 1e8

	in := bufio.NewReader(reader)
	out := bufio.NewWriter(writer)
	defer out.Flush()

	var n, m int
	Fscan(in, &n, &m)
	if n >= m {
		Fprint(out, "YES")
		return
	}
	arr := make([]int, n)
	for i := range arr {
		Fscan(in, &arr[i])
		arr[i] %= m
	}

	// dp[i][j]>0: 前i个数组合出了j
	dp := [2][]int{}
	for i := 0; i < 2; i++ {
		dp[i] = make([]int, m)
		for j := range dp[i] {
			dp[i][j] = -inf
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		cur, prev := i&1, (i-1)&1
		dp[cur][0] = dp[prev][0]
		for j := range dp[prev] {
			dp[cur][j] = max(dp[cur][j], dp[prev][(j-arr[i-1]+m)%m]+1)
		}
	}
	if dp[n&1][0] > 0 {
		Fprint(out, "YES")
	} else {
		Fprint(out, "NO")
	}
}

//func main() {
//	Sol577B(os.Stdin, os.Stdout)
//}
