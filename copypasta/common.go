package copypasta

import (
	. "fmt"
	"io"
	"math"
	"math/rand"
	"sort"
	"time"
)

// General ideas https://codeforces.com/blog/entry/48417
// 从特殊到一般：尝试修改条件或缩小题目的数据范围，先研究某个特殊情况下的思路，然后再逐渐扩大数据范围来思考怎么改进算法

// 异类双变量：固定某变量统计另一变量的 [0,n)
//     EXTRA: 值域上的双变量，见 https://codeforces.com/contest/486/problem/D
// 同类双变量①：固定 i 统计 [0,n)
// 同类双变量②：固定 i 统计 [0,i-1]
// 套路：预处理数据（按照某种顺序排序/优先队列/BST/...），或者边遍历边维护，
//      然后固定变量 i，用均摊 O(1)~O(logn) 的复杂度统计范围内的另一变量 j
// 这样可以将复杂度从 O(n^2) 降低到 O(n) 或 O(nlogn)
// 进阶：https://codeforces.com/contest/1483/problem/D

// 利用前缀和实现巧妙的构造 https://www.luogu.com.cn/blog/duyi/qian-zhui-he
// 邻项修改->前缀和->单项修改 https://codeforces.com/problemset/problem/1254/B2 https://ac.nowcoder.com/acm/contest/7612/C

/* 横看成岭侧成峰
转换为距离的众数 https://codeforces.com/problemset/problem/1365/C
转换为差分数组的变化 https://codeforces.com/problemset/problem/1110/E
转换为差 http://www.51nod.com/Challenge/Problem.html#problemId=1217
考虑每个点产生的贡献 https://codeforces.com/problemset/problem/1009/E
考虑每条边产生的负贡献 https://atcoder.jp/contests/abc173/tasks/abc173_f
考虑符合范围要求的贡献 https://codeforces.com/problemset/problem/1151/E
和式的另一视角。若每一项的值都在一个范围，不妨考虑另一个问题：值为 x 的项有多少个？https://atcoder.jp/contests/abc162/tasks/abc162_e
对所有排列考察所有子区间的性质，可以转换成对所有子区间考察所有排列。将子区间内部的排列和区间外部的排列进行区分，内部的性质单独研究，外部的当作 (n-(r-l))! 个排列 https://codeforces.com/problemset/problem/1284/C
从最大值入手 https://codeforces.com/problemset/problem/1381/B
等效性 https://leetcode-cn.com/contest/biweekly-contest-8/problems/maximum-number-of-ones/
逆向思维 https://leetcode-cn.com/contest/biweekly-contest-9/problems/minimum-time-to-build-blocks/
https://leetcode-cn.com/contest/biweekly-contest-31/problems/minimum-number-of-increments-on-subarrays-to-form-a-target-array/
*/

/* 奇偶性
https://codeforces.com/problemset/problem/763/B
https://codeforces.com/problemset/problem/1270/E
https://codeforces.com/problemset/problem/1332/E 配对法：将合法局面与非法局面配对
*/

/* 归纳：solve(n)->solve(n-1) 或者 solve(n-1)->solve(n)
https://codeforces.com/problemset/problem/1517/C
https://codeforces.com/problemset/problem/412/D
https://codeforces.com/problemset/problem/266/C
*/

/* 正难则反：小学奥数告诉我们，不可行方案永远比可行方案好求
https://codeforces.com/problemset/problem/621/C
https://codeforces.com/problemset/problem/571/A
*/

/* 见微知著：考察单个点的规律，从而推出全局规律
https://codeforces.com/problemset/problem/1510/K
https://leetcode-cn.com/problems/minimum-number-of-operations-to-reinitialize-a-permutation/
*/

// 栈+懒删除 https://codeforces.com/problemset/problem/1000/F
// 栈的应用 https://codeforces.com/problemset/problem/1092/D1
//         https://codeforces.com/problemset/problem/1092/D2

// 锻炼分类讨论能力 https://codeforces.com/problemset/problem/356/C

/* Golang 注意事项
for range array 会拷贝一份 array，这种情况可以用 for range array[:]
for-switch 内的 break 跳出的是该 switch，而不是其外部的 for 循环
对于存在海量小对象的情况（如 trie, treap 等），使用 debug.SetGCPercent(-1) 来禁用 GC，不去扫描大量对象，能明显减少耗时
对于可以回收的情况（如 append 在超过 cap 时），使用 debug.SetGCPercent(-1) 虽然会减少些许耗时，但若有大量内存没被回收，会有 MLE 的风险
其他情况下使用 debug.SetGCPercent(-1) 对耗时和内存使用无明显影响
对于多组数据的情况，禁用 GC 若 MLE，可在每组数据的开头或者末尾调用 debug.FreeOSMemory() 手动 GC
参考 https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/
    https://zhuanlan.zhihu.com/p/77943973
*/
func commonCollection() {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	pow10 := func(x int) int64 { return int64(math.Pow10(x)) } // 不需要 round

	// TIPS: dir4[i] 和 dir4[i^1] 互为相反方向
	type pair struct{ x, y int }
	dir4 := []pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 上下左右
	dir4C := []pair{ // 西东南北
		'W': {-1, 0},
		'E': {1, 0},
		'S': {0, -1},
		'N': {0, 1},
	}
	dir4c := []pair{ // 左右下上
		'L': {-1, 0},
		'R': {1, 0},
		'D': {0, -1},
		'U': {0, 1},
	}
	dir4R := []pair{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}
	dir8 := []pair{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}
	perm3 := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	perm4 := [][]int{
		{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 2, 1, 3}, {0, 2, 3, 1}, {0, 3, 1, 2}, {0, 3, 2, 1},
		{1, 0, 2, 3}, {1, 0, 3, 2}, {1, 2, 0, 3}, {1, 2, 3, 0}, {1, 3, 0, 2}, {1, 3, 2, 0},
		{2, 0, 1, 3}, {2, 0, 3, 1}, {2, 1, 0, 3}, {2, 1, 3, 0}, {2, 3, 0, 1}, {2, 3, 1, 0},
		{3, 0, 1, 2}, {3, 0, 2, 1}, {3, 1, 0, 2}, {3, 1, 2, 0}, {3, 2, 0, 1}, {3, 2, 1, 0},
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	mins := func(a ...int) int {
		res := a[0]
		for _, v := range a[1:] {
			if v < res {
				res = v
			}
		}
		return res
	}
	maxs := func(a ...int) int {
		res := a[0]
		for _, v := range a[1:] {
			if v > res {
				res = v
			}
		}
		return res
	}
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	ceil := func(a, b int) int {
		// assert a >= 0 && b > 0
		if a == 0 {
			return 0
		}
		return (a-1)/b + 1
	}
	// 另一种写法，无需考虑 a 为 0 的情况
	ceil = func(a, b int) int {
		return (a + b - 1) / b
	}
	bin := func(v int) []byte {
		const maxLen = 30 // 62 for int64
		s := make([]byte, maxLen+1)
		for i := range s {
			s[i] = byte(v >> (maxLen - i) & 1)
		}
		return s
	}

	sort3 := func(a ...int) (x, y, z int) { sort.Ints(a); return a[0], a[1], a[2] }
	minString := func(a, b string) string {
		if len(a) != len(b) {
			if len(a) < len(b) {
				return a
			}
			return b
		}
		if a < b {
			return a
		}
		return b
	}
	ternaryI := func(cond bool, r1, r2 int) int {
		if cond {
			return r1
		}
		return r2
	}
	ternaryS := func(cond bool, r1, r2 string) string {
		if cond {
			return r1
		}
		return r2
	}
	zip := func(a, b []int) {
		n := len(a)
		type pair struct{ x, y int }
		ps := make([]pair, n)
		for i := range ps {
			ps[i] = pair{a[i], b[i]}
		}
	}
	zipI := func(a []int) {
		n := len(a)
		type pair struct{ x, y int }
		ps := make([]pair, n)
		for i := range ps {
			ps[i] = pair{a[i], i}
		}
	}

	// 顺时针旋转矩阵 90°
	rotate := func(a [][]int) [][]int {
		n, m := len(a), len(a[0])
		b := make([][]int, m)
		for i := range b {
			b[i] = make([]int, n)
		}
		for i, r := range a {
			for j, v := range r {
				b[j][n-1-i] = v
			}
		}
		return b
	}
	// 转置
	transpose := func(a [][]int) [][]int {
		n, m := len(a), len(a[0])
		b := make([][]int, m)
		for i := range b {
			b[i] = make([]int, n)
			for j, r := range a {
				b[i][j] = r[i]
			}
		}
		return b
	}

	// 适用于 mod 超过 int32 范围的情况
	// 还有一种用浮点数的写法，此略
	mul := func(a, b, mod int64) (res int64) {
		for ; b > 0; b >>= 1 {
			if b&1 == 1 {
				res = (res + a) % mod
			}
			a = (a + a) % mod
		}
		return
	}

	// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
	pow := func(x, n, mod int64) int64 {
		x %= mod
		res := int64(1) % mod
		for ; n > 0; n >>= 1 {
			if n&1 == 1 {
				res = res * x % mod
			}
			x = x * x % mod
		}
		return res
	}

	// 从低位到高位
	toAnyBase := func(x, base int) (res []int) {
		for ; x > 0; x /= base {
			res = append(res, x%base)
		}
		return
	}
	digits := func(x int) (res []int) {
		for ; x > 0; x /= 10 {
			res = append(res, x%10)
		}
		return
	}

	// 合并有序数组，保留重复元素
	// a b 必须是有序的（可以为空）
	merge := func(a, b []int) []int {
		i, n := 0, len(a)
		j, m := 0, len(b)
		res := make([]int, 0, n+m)
		for {
			if i == n {
				return append(res, b[j:]...)
			}
			if j == m {
				return append(res, a[i:]...)
			}
			if a[i] < b[j] { // 改成 > 为降序
				res = append(res, a[i])
				i++
			} else {
				res = append(res, b[j])
				j++
			}
		}
	}

	// 返回 a 的各个子集的元素和
	// https://codeforces.com/contest/1209/problem/E2
	subSum := func(a []int) []int {
		sum := make([]int, 1<<len(a)) // int64
		for p, v := range a {
			for s := 0; s < 1<<p; s++ {
				sum[1<<p|s] = sum[s] + v
				// NOTE: 若要直接在此写循环遍历 sum，注意别漏了 sum[0] = 0 的情况
			}
		}
		return sum
	}

	// 返回 a 的各个子集的元素和的排序后的结果
	// 若已求出前 i-1 个数的有序子集和 b，那么前 i 个数的有序子集和可以由 b 和 {b 的每个数加上 a[i]} 归并得到
	// 复杂度为 O(1+2+4+...+2^(n-1)) = O(2^n)
	// 参考 https://leetcode-cn.com/problems/closest-subsequence-sum/solution/o2n2de-zuo-fa-by-heltion-0yn7/
	subSumSorted := func(a []int) []int {
		sum := []int{0}
		for _, v := range a {
			b := make([]int, len(sum))
			for i, w := range sum {
				b[i] = w + v
			}
			sum = merge(sum, b)
		}
		return sum
	}

	// 分组前缀和（具体见 query 上的注释）
	// LC1664/周赛216C https://leetcode-cn.com/contest/weekly-contest-216/problems/ways-to-make-a-fair-array/
	groupPrefixSum := func(a []int, k int) {
		// 补 0 简化后续逻辑
		n := len(a)
		for len(a)%k > 0 {
			a = append(a, 0)
		}
		sum := make([]int, len(a)+k) // int64
		for i, v := range a {
			sum[i+k] = sum[i] + v
		}
		pre := func(x, m int) int {
			if x%k <= m {
				return sum[x/k*k+m]
			}
			return sum[(x+k-1)/k*k+m]
		}
		// 求下标在 [l,r) 范围内且下标同余于 m 的元素和 (0<=m<k)
		query := func(l, r, m int) int {
			return pre(r, m) - pre(l, m)
		}
		a = a[:n] // 如果要枚举等，可能需要复原

		_ = query
	}

	// 环形区间和 [l,r) 0<=l<r
	circularRangeSum := func(a []int) {
		n := len(a)
		sum := make([]int64, n+1)
		for i, v := range a {
			sum[i+1] = sum[i] + int64(v)
		}
		pre := func(p int) int64 {
			return sum[n]*int64(p/n) + sum[p%n]
		}
		query := func(l, r int) int64 {
			return pre(r) - pre(l)
		}

		_ = query
	}

	// 带权(等差数列)前缀和
	{
		var n int // read
		a := make([]int64, n)
		// read a ...

		sum := make([]int64, n+1)
		iSum := make([]int64, n+1)
		for i, v := range a {
			sum[i+1] = sum[i] + v
			iSum[i+1] = iSum[i] + int64(i+1)*v
		}
		query := func(l, r int) int64 { return iSum[r] - iSum[l] - int64(l)*(sum[r]-sum[l]) } // [l,r)

		_ = query
	}

	// 二维前缀和
	var sum2d [][]int
	initSum2D := func(a [][]int) {
		n, m := len(a), len(a[0])
		sum2d = make([][]int, n+1)
		sum2d[0] = make([]int, m+1)
		for i, row := range a {
			sum2d[i+1] = make([]int, m+1)
			for j, v := range row {
				sum2d[i+1][j+1] = sum2d[i+1][j] + sum2d[i][j+1] - sum2d[i][j] + v
			}
		}
	}
	// r1<=r<=r2 && c1<=c<=c2
	querySum2D := func(r1, c1, r2, c2 int) int {
		r2++
		c2++
		return sum2d[r2][c2] - sum2d[r2][c1] - sum2d[r1][c2] + sum2d[r1][c1]
	}

	// 矩阵每行每列的前缀和
	rowColSum := func(a [][]int) (sumR, sumC [][]int) {
		n, m := len(a), len(a[0])
		sumR = make([][]int, n) // int64
		for i, row := range a {
			sumR[i] = make([]int, m+1)
			for j, v := range row {
				sumR[i][j+1] = sumR[i][j] + v
			}
		}
		sumC = make([][]int, n+1) // int64
		for i := range sumC {
			sumC[i] = make([]int, m)
		}
		for j := 0; j < m; j++ {
			for i, row := range a {
				sumC[i+1][j] = sumC[i][j] + row[j]
			}
		}
		return
	}

	// 矩阵每条主对角线、反对角线的前缀和
	// https://leetcode-cn.com/problems/get-biggest-three-rhombus-sums-in-a-grid/
	diagonalSum := func(a [][]int) {
		n, m := len(a), len(a[0])

		ds := make([][]int, n+1) // 主对角线前缀和
		as := make([][]int, n+1) // 反对角线前缀和
		for i := range ds {
			ds[i] = make([]int, m+1)
			as[i] = make([]int, m+1)
		}
		for i, r := range a {
			for j, v := range r {
				ds[i+1][j+1] = ds[i][j] + v // ↘
				as[i+1][j] = as[i][j+1] + v // ↙
			}
		}
		// 从 x,y 开始，向 ↘，连续的 k 个数的和（需要保证至少有 k 个数）
		queryDiagonal := func(x, y, k int) int { return ds[x+k][y+k] - ds[x][y] }
		// 从 x,y 开始，向 ↙，连续的 k 个数的和（需要保证至少有 k 个数）
		queryAntiDiagonal := func(x, y, k int) int { return as[x+k][y+1-k] - as[x][y+1] }

		_, _ = queryDiagonal, queryAntiDiagonal
	}

	// 利用每个数产生的贡献计算 Σ|ai-aj|, i!=j
	// 相关题目 https://codeforces.com/contest/1311/problem/F
	contributionSum := func(a []int) (sum int64) {
		n := len(a)
		sort.Ints(a)
		for i, v := range a {
			sum += int64(v) * int64(2*i+1-n)
		}
		return
	}

	// 二维差分
	// todo https://blog.csdn.net/weixin_43914593/article/details/113782108
	//      https://www.luogu.com.cn/problem/P3397

	reverse := func(a []byte) []byte {
		n := len(a)
		b := make([]byte, n)
		for i, v := range a {
			b[n-1-i] = v
		}
		return b
	}
	reverseInPlace := func(a []byte) {
		for i, n := 0, len(a); i < n/2; i++ {
			a[i], a[n-1-i] = a[n-1-i], a[i]
		}
	}

	equal := func(a, b []int) bool {
		// assert len(a) == len(b)
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	// 求差集 A-B, B-A 和交集 A∩B
	// EXTRA: 求并集 union: A∪B = A-B+A∩B = merge(differenceA, intersection) 或 merge(differenceB, intersection)
	// EXTRA: 求对称差 symmetric_difference: A▲B = A-B ∪ B-A = merge(differenceA, differenceB)
	// a b 必须是有序的（可以为空）
	// 与图论结合 https://codeforces.com/problemset/problem/243/B
	splitDifferenceAndIntersection := func(a, b []int) (differenceA, differenceB, intersection []int) {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n {
				differenceB = append(differenceB, b[j:]...)
				return
			}
			if j == m {
				differenceA = append(differenceA, a[i:]...)
				return
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				differenceA = append(differenceA, x)
				i++
			} else if x > y { // 改成 < 为降序
				differenceB = append(differenceB, y)
				j++
			} else {
				intersection = append(intersection, x)
				i++
				j++
			}
		}
	}

	// 求交集简洁写法
	intersection := func(a, b []int) []int {
		mp := map[int]bool{}
		for _, v := range a {
			mp[v] = true
		}
		mp2 := map[int]bool{}
		for _, v := range b {
			if mp[v] {
				mp2[v] = true
			}
		}
		mp = mp2

		keys := make([]int, 0, len(mp))
		for k := range mp {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		return keys
	}

	// a 是否为 b 的子集（相当于 differenceA 为空）
	// a b 需要是有序的
	isSubset := func(a, b []int) bool {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n {
				return true
			}
			if j == m {
				return false
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				return false
			} else if x > y { // 改成 < 为降序
				j++
			} else {
				i++
				j++
			}
		}
	}

	// EXTRA: a 是否为 b 的子序列
	// https://codeforces.com/problemset/problem/778/A
	isSubSequence := func(a, b []int) bool {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n {
				return true
			}
			if j == m {
				return false
			}
			if a[i] == b[j] {
				i++
				j++
			} else {
				j++
			}
		}
	}

	// 是否为不相交集合（相当于 intersection 为空）
	// a b 需要是有序的
	isDisjoint := func(a, b []int) bool {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n || j == m {
				return true
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				i++
			} else if x > y { // 改成 < 为降序
				j++
			} else {
				return false
			}
		}
	}

	// 去重
	// a 必须是有序的
	unique := func(a []int) (res []int) {
		for i, v := range a {
			if i == 0 || v != a[i-1] {
				res = append(res, v)
			}
		}
		//n = len(res)
		return
	}

	// 原地去重
	// a 必须是有序的
	uniqueInPlace := func(a []int) []int {
		n := len(a)
		if n == 0 {
			return nil
		}
		k := 0
		for _, w := range a[1:] {
			if a[k] != w {
				k++
				a[k] = w
			}
		}
		//n = k + 1
		return a[:k+1]
	}

	// 离散化，不保留原始数据（保留原始数据的版本见下面的 discreteMap）
	// discrete([]int{100,20,50,50}, 1) => []int{3,1,2,2}
	// https://leetcode-cn.com/contest/biweekly-contest-18/problems/rank-transform-of-an-array/
	discrete := func(a []int, startIndex int) (kth []int) {
		type pair struct{ v, i int }
		ps := make([]pair, len(a))
		for i, v := range a {
			ps[i] = pair{v, i}
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v }) // or SliceStable
		kth = make([]int, len(a))

		// a 有重复元素
		k := startIndex
		for i, p := range ps {
			if i > 0 && p.v != ps[i-1].v {
				k++
			}
			kth[p.i] = k
		}

		// 若有需要，求出 kth 后还可以对 ps 进行去重，这样可以用 kth 值访问原始值

		// a 无重复元素，或者给相同元素也加上顺序（例如主席树的离散化写法）
		for i, p := range ps {
			kth[p.i] = i + startIndex
		}

		return
	}

	// 简化版，不要求值连续 [10,30,20,20] => [0,3,1,1]
	discrete2 := func(a []int, startIndex int) []int {
		b := append([]int(nil), a...)
		sort.Ints(b)
		for i, v := range a {
			a[i] = sort.SearchInts(b, v) + startIndex
		}
		return a
	}

	// 保留原始数据的离散化
	// 返回一个名次 map
	// discreteMap([]int{100,20,20,50}, 1) => map[int]int{20:1, 50:2, 100:3}
	// 例题：LC327 https://leetcode-cn.com/problems/count-of-range-sum/
	discreteMap := func(a []int, startIndex int) (kth map[int]int) {
		sorted := append([]int(nil), a...)
		sort.Ints(sorted)

		// 有重复元素
		kth = map[int]int{}
		k := startIndex
		for i, v := range sorted {
			if i == 0 || v != sorted[i-1] {
				kth[v] = k
				k++
			}
		}

		// 无重复元素
		kth = make(map[int]int, len(sorted))
		for i, v := range sorted {
			kth[v] = i + startIndex
		}

		// EXTRA: 第 k 小元素在原数组中的下标 kthPos
		pos := make(map[int][]int, k-startIndex)
		for i, v := range a {
			pos[v] = append(pos[v], i)
		}
		kthPos := make([][]int, k+1)
		for v, k := range kth {
			kthPos[k] = pos[v]
		}

		return
	}

	// 哈希编号，也可以理解成另一种离散化（无序）
	// 编号从 0 开始
	indexMap := func(a []string) map[string]int {
		mp := map[string]int{}
		for _, v := range a {
			if _, ok := mp[v]; !ok {
				mp[v] = len(mp)
			}
		}
		return mp
	}

	allSame := func(a ...int) bool {
		for _, v := range a[1:] {
			if v != a[0] {
				return false
			}
		}
		return true
	}

	// a 相对于 [0,n) 的补集
	// a 必须是升序且无重复元素
	complement := func(n int, a []int) (res []int) {
		j := 0
		for i := 0; i < n; i++ {
			if j == len(a) || i < a[j] {
				res = append(res, i)
			} else {
				j++
			}
		}
		return
	}

	// 数组第 k 小 (Quick Select)       kthElement nthElement
	// 0 <= k < len(a)
	// 调用会改变数组中元素顺序
	// 代码实现参考算法第四版 p.221
	// 算法的平均比较次数为 ~2n+2kln(n/k)+2(n-k)ln(n/(n-k))
	// https://en.wikipedia.org/wiki/Quickselect
	// https://www.geeksforgeeks.org/quickselect-algorithm/
	// 模板题 LC215 https://leetcode-cn.com/problems/kth-largest-element-in-an-array/
	//       LC973 https://leetcode-cn.com/problems/k-closest-points-to-origin/submissions/
	// 模板题 https://codeforces.com/contest/977/problem/C
	quickSelect := func(a []int, k int) int {
		//k = len(a) - 1 - k // 求第 k 大
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for l, r := 0, len(a)-1; l < r; {
			v := a[l] // 切分元素
			i, j := l, r+1
			for {
				for i++; i < r && a[i] < v; i++ { // less(i, l)
				}
				for j--; j > l && a[j] > v; j-- { // less(l, j)
				}
				if i >= j {
					break
				}
				a[i], a[j] = a[j], a[i]
			}
			a[l], a[j] = a[j], v
			if j == k {
				break
			} else if j < k {
				l = j + 1
			} else {
				r = j - 1
			}
		}
		return a[k] //  a[:k+1]  a[k:]
	}

	contains := func(a []int, x int) bool {
		for _, v := range a {
			if v == x {
				return true
			}
		}
		return false
	}

	// x 是否包含 y 中的所有元素，且顺序一致
	containsAll := func(x, y []int) bool {
		for len(y) < len(x) {
			if len(y) == 0 {
				return true
			}
			if x[0] == y[0] {
				y = y[1:]
			}
			x = x[1:]
		}
		return false
	}

	// 扫描线
	// 某些题目需要配合线段树
	// https://cses.fi/book/book.pdf 30.1
	// TODO 窗口的星星 https://www.luogu.com.cn/problem/P1502
	// TODO 矩形周长 https://www.luogu.com.cn/problem/P1856
	// 天际线问题 LC218 https://leetcode-cn.com/problems/the-skyline-problem/
	// TODO 矩形面积并 LC850 https://leetcode-cn.com/problems/rectangle-area-ii/ 《算法与实现》5.4.3
	// 经典题 https://codeforces.com/problemset/problem/1000/C
	// https://codeforces.com/problemset/problem/1379/D
	// 线段相交统计（栈） https://codeforces.com/contest/1278/problem/D
	// 统计水平方向的线段与垂直方向的线段的交点个数 https://codeforces.com/problemset/problem/610/D
	// LC 套题 https://leetcode-cn.com/tag/line-sweep/
	// todo CF652D
	sweepLine := func(in io.Reader, n int) {
		type event struct{ pos, delta int }
		events := make([]event, 0, 2*n)
		for i := 0; i < n; i++ {
			var l, r int
			Fscan(in, &l, &r)
			events = append(events, event{l, 1}, event{r, -1})
		}
		sort.Slice(events, func(i, j int) bool {
			a, b := events[i], events[j]
			return a.pos < b.pos || a.pos == b.pos && a.delta < b.delta // 先出后进。改成 a.delta > b.delta 为先进后出
		})

		for _, e := range events {
			if e.delta > 0 {

			} else {

			}
		}
	}

	// 扫描线另一种写法，把 delta 压缩进 pos
	// 这样可以避免写一个复杂的 sort.Slice
	sweepLine2 := func(in io.Reader, n int) {
		events := make([]int, 0, 2*n)
		for i := 0; i < n; i++ {
			var l, r int
			Fscan(in, &l, &r)
			// 注意移位后是否溢出
			events = append(events, l<<1|1, r<<1) // 先出后进
			//events = append(events, l<<1, r<<1|1) // 先进后出
		}
		sort.Ints(events)

		for _, e := range events {
			pos, delta := e>>1, e&1
			_ = pos
			if delta > 0 { // 根据上面的写法来定义何为出何为进

			} else {

			}
		}
	}

	// 扫描线：一维格点刷漆，返回被刷到的格点数
	countCoveredPoints := func(in io.Reader, m int) int {
		type pair struct{ p, d int }
		es := make([]pair, 0, 2*m)
		for i := 0; i < m; i++ {
			var l, r int
			Fscan(in, &l, &r)
			es = append(es, pair{l, 1}, pair{r, -1})
		}
		// assert len(es) > 0
		sort.Slice(es, func(i, j int) bool { return es[i].p < es[j].p })
		ans := es[len(es)-1].p - es[0].p + 1
		// 减去没被刷到的格点
		eventCnt, st := 0, es[0].p
		for _, e := range es {
			if eventCnt == 0 {
				if d := e.p - st - 1; d > 0 {
					ans -= d
				}
			}
			eventCnt += e.d
			if eventCnt == 0 {
				st = e.p
			}
		}
		return ans
	}

	// 二维离散化
	// 代码来源 https://atcoder.jp/contests/abc168/tasks/abc168_f
	discrete2D := func(n, m int) (ans int) {
		type line struct{ a, b, c int }
		lr := make([]line, n)
		du := make([]line, m)
		// read ...

		xs := []int{-2e9, 0, 2e9}
		ys := []int{-2e9, 0, 2e9}
		for _, l := range lr {
			a, b, c := l.a, l.b, l.c
			xs = append(xs, a, b)
			ys = append(ys, c)
		}
		for _, l := range du {
			a, b, c := l.a, l.b, l.c
			xs = append(xs, a)
			ys = append(ys, b, c)
		}
		sort.Ints(xs)
		xs = unique(xs)
		xi := discreteMap(xs, 0)
		sort.Ints(ys)
		ys = unique(ys)
		yi := discrete(ys, 0)

		lx, ly := len(xi), len(yi)
		glr := make([][]int, lx)
		gdu := make([][]int, lx)
		vis := make([][]bool, lx)
		for i := range glr {
			glr[i] = make([]int, ly)
			gdu[i] = make([]int, ly)
			vis[i] = make([]bool, ly)
		}
		for _, p := range lr {
			glr[xi[p.a]][yi[p.c]]++
			glr[xi[p.b]][yi[p.c]]--
		}
		for _, p := range du {
			gdu[xi[p.a]][yi[p.b]]++
			gdu[xi[p.a]][yi[p.c]]--
		}
		for i := 1; i < lx-1; i++ {
			for j := 1; j < ly-1; j++ {
				glr[i][j] += glr[i-1][j]
				gdu[i][j] += gdu[i][j-1]
			}
		}

		type pair struct{ x, y int }
		q := []pair{{xi[0], yi[0]}}
		for len(q) > 0 {
			p := q[0]
			q = q[1:]
			x, y := p.x, p.y
			if x == 0 || x == lx-1 || y == 0 || y == ly-1 {
				return -1
			} // 无穷大
			if !vis[x][y] {
				vis[x][y] = true
				ans += (xs[x+1] - xs[x]) * (ys[y+1] - ys[y])
				if glr[x][y] == 0 {
					q = append(q, pair{x, y - 1})
				}
				if glr[x][y+1] == 0 {
					q = append(q, pair{x, y + 1})
				}
				if gdu[x][y] == 0 {
					q = append(q, pair{x - 1, y})
				}
				if gdu[x+1][y] == 0 {
					q = append(q, pair{x + 1, y})
				}
			}
		}
		return
	}

	_ = []interface{}{
		pow10, dir4, dir4C, dir4c, dir4R, dir8, perm3, perm4,
		min, mins, max, maxs, abs, ceil, bin,
		ternaryI, ternaryS, zip, zipI, rotate, transpose, minString,
		pow, mul, toAnyBase, digits,
		subSum, subSumSorted, groupPrefixSum, circularRangeSum, initSum2D, querySum2D, rowColSum, diagonalSum,
		contributionSum,
		sort3, reverse, reverseInPlace, equal,
		merge, splitDifferenceAndIntersection, intersection, isSubset, isSubSequence, isDisjoint,
		unique, uniqueInPlace, discrete, discrete2, discreteMap, indexMap, allSame, complement, quickSelect, contains, containsAll,
		sweepLine, sweepLine2, countCoveredPoints,
		discrete2D,
	}
}
