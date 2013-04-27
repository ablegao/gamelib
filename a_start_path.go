package gamelib

import (
	"container/list"
)

// a * 移动算方法 。
const (
	COST_HORIZONTAL = 20 //横向移动评分
	COST_VERTICAL   = 5  //纵向移动评分
	COST_DIAGONAL   = 12 //斜向移动评分。
)

type Point struct {
	X int
	Y int
}

type Node struct {
	X          int
	Y          int
	IsInOpen   bool
	IsInClose  bool
	G          int
	H          int
	F          int
	ParentNode *Node
}

func NewNode(x int, y int, ParentNode *Node) *Node {
	return &Node{
		X:          x,
		Y:          y,
		IsInOpen:   false,
		IsInClose:  false,
		G:          0,
		H:          0,
		F:          0,
		ParentNode: ParentNode,
	}
}

func NewAStartPathFinder(MoveArray [][]int, BigMapWidth int, BigMapHeight int, ScroW int, ScroH int) *AStartPathFinder {

	asp := &AStartPathFinder{xMapStart: 0,
		yMapStart: 0,
		wMap:      BigMapWidth / ScroW,
		hMap:      BigMapHeight / (ScroH / 2),
		moveArray: MoveArray,
		openList:  list.New(),
		closeList: list.New(),
		clousMap:  []int{},
	}
	asp.InitMap()
	return asp
}

func Abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}

type AStartPathFinder struct {
	//地图起始网格坐标
	xMapStart int //x 开始坐标
	yMapStart int //y

	wMap int ///  地图列数（每行格点数）
	hMap int //	//地图行数（每列格点数

	mapData     [][]*Node  //map[int]map[int]Node
	openList    *list.List //开放列表
	closeList   *list.List //关闭列表
	isFinded    bool       // 能否找到路径。
	RunTimeInMs int        //寻路时间
	moveArray   [][]int

	clousMap []int
}

//初始化一个地图坐标。矩阵
func (this *AStartPathFinder) InitMap() {
	this.mapData = [][]*Node{}
	var xindex int
	for y := this.yMapStart; y < this.hMap; y += 1 {
		if y&1 == 0 {
			xindex = this.wMap
		} else {
			xindex = this.wMap + 1
		}
		xData := []*Node{}
		for x := this.xMapStart; x < xindex; x++ {
			//this.mapData[y][x] = Node{x, y}
			xData = append(xData, NewNode(x, y, nil))
			//fmt.Printf("%d:%d\n", x, y)
		}
		this.mapData = append(this.mapData, xData)
	}
}

/**
 * 判断坐标是否是在数组里面的
 * @param po 坐标位置
 * @return 真 假
 *
 */
func (this *AStartPathFinder) whetherMovable(po *Point) bool {

	if len(this.mapData) > po.Y && len(this.mapData[po.Y]) > po.X {
		return false
	} else {
		return true
	}
}

/**
 * 开始计算 并返回 路径数组
 * @param startPoint 初始移动坐标
 * @param endPoint 终点坐标
 * @return 所移动路径数组
 *
 */
func (this *AStartPathFinder) Find(startp *Point, endp *Point) []*Point {
	if this.whetherMovable(startp) || this.whetherMovable(endp) {
		return nil
	}

	currentNode := this.mapData[startp.Y][startp.X]
	endNode := this.mapData[endp.Y][endp.X]
	//.this.openList.
	this.openList.PushBack(currentNode)

	//只要openList有数据， 就一只循环， 知道openList里面的所有内容消失 。
	for this.openList.Len() > 0 {
		NowNode := this.openList.Front()
		currentNode = NowNode.Value.(*Node)

		//元素用完需要删除
		this.openList.Remove(NowNode)

		currentNode.IsInOpen = false
		currentNode.IsInClose = true
		this.closeList.PushBack(currentNode)

		//已经到达终点
		if currentNode.X == endp.X && currentNode.Y == endp.Y {
			this.isFinded = true
			break
		} else {

			aroundNodes := this.getAroundsNode(currentNode.X, currentNode.Y)

			for _, node := range aroundNodes {
				if node == nil {
					continue
				}
				g := this.getGValue(currentNode, node)
				h := this.getHValue(currentNode, endNode, node)
				if node.IsInOpen == true {
					if g < node.G {

						node.G = g
						node.H = h
						node.F = g + h
						node.ParentNode = currentNode
						this.findAndSort(node)

					}
				} else {
					node.G = g
					node.H = h
					node.F = g + h
					node.ParentNode = currentNode

					this.insertAndSort(node)

				}
			}
		}

	}
	if this.isFinded {
		path := this.createPath(startp.X, startp.Y)
		this.destroyLists()
		return path
	} else {
		this.destroyLists()
		return nil
	}
	return nil
}

/**
 * 生成路径数组
 */
func (this *AStartPathFinder) createPath(xStart int, yStart int) []*Point {

	path := list.New()
	if nodeByList := this.closeList.Back(); nodeByList != nil {
		node := nodeByList.Value.(*Node)
		for node.X != xStart || node.Y != yStart {
			path.PushFront(&Point{node.X, node.Y})
			node = node.ParentNode
		}
		path.PushFront(&Point{node.X, node.Y})

		/*
			for node := this.closeList.Back(); node != nil && (node.X != xStart || node.Y != yStart); node = node.ParentNode {

			}*/

	} else {

	}

	retPath := []*Point{}
	for nn := path.Front(); nn != nil; nn = nn.Next() {
		retPath = append(retPath, nn.Value.(*Point))
	}

	return retPath

}

/**
 * 删除列表中原来的node，并将node放到正确顺序
 */

func (this *AStartPathFinder) findAndSort(node *Node) {
	listLength := this.openList.Len()
	if listLength < 1 {
		return
	}
	for nn := this.openList.Front(); nn != nil; nn = nn.Next() {
		if node.F <= nn.Value.(*Node).F {
			this.openList.InsertBefore(node, nn)
		}
		if node.X == nn.Value.(*Node).X && node.Y == nn.Value.(*Node).Y {
			this.openList.Remove(nn)
		}
	}

}

/**
 * 按由小到大顺序将节点插入到列表
 */
func (this *AStartPathFinder) insertAndSort(node *Node) {
	node.IsInOpen = true

	listLength := this.openList.Len()

	if listLength == 0 {
		this.openList.PushBack(node)
	} else {
		for nn := this.openList.Front(); nn != nil; nn = nn.Next() {
			if node.F <= nn.Value.(*Node).F {
				this.openList.InsertBefore(node, nn)
				return
			}
		}

		this.openList.PushBack(node)
	}
}
func (this *AStartPathFinder) getGValue(currentNode *Node, n *Node) int {
	g := 0
	if currentNode.Y == n.Y { // 横向  左右

		g = currentNode.G + COST_HORIZONTAL //this.COST_HORIZONTAL;
	} else if currentNode.Y+2 == n.Y || currentNode.Y-2 == n.Y { // 竖向  上下

		g = currentNode.G + COST_VERTICAL*2
	} else { // 斜向  左上 左下 右上 右下

		g = currentNode.G + COST_DIAGONAL
	}

	return g

}
func (this *AStartPathFinder) getHValue(currentNode *Node, endNode *Node, node *Node) int {
	var dx int
	var dy int
	//节点到0，0点的x轴距离
	var dxNodeTo0 int = node.X*COST_HORIZONTAL + (node.Y&1)*COST_HORIZONTAL/2
	//终止节点到0，0点的x轴距离
	var dxEndNodeTo0 int = endNode.X*COST_HORIZONTAL + (endNode.Y&1)*COST_HORIZONTAL/2
	dx = Abs(dxEndNodeTo0 - dxNodeTo0)
	dy = Abs(endNode.Y-node.Y) * COST_VERTICAL
	return dx + dy
}

/**
 * 检查点在地图上是否可通过
 */
func (this *AStartPathFinder) getAroundsNode(x int, y int) []*Node {
	//var aroundNodes Array = new Array();
	aroundNodes := []*Node{}

	var checkX int
	var checkY int
	/**
	 * 菱形组合的地图八方向与正常不同
	 */
	//左
	checkX = x - 1
	checkY = y
	//左边是否可以动  左上 左下是否可移动
	if this.isWalkable(checkX, checkY) && this.isWalkable(x-1+(y&1), y-1) && this.isWalkable(x-1+(y&1), y+1) && !this.mapData[checkY][checkX].IsInClose {
		//trace("左:::"+checkY+"##"+checkX);
		//aroundNodes.push(this.map[checkY][checkX]);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//右
	checkX = x + 1
	checkY = y
	//右边是否可以动  右上 右下是否可移动
	if this.isWalkable(checkX, checkY) && this.isWalkable(x+(y&1), y-1) && this.isWalkable(x+(y&1), y+1) && !this.mapData[checkY][checkX].IsInClose {
		//trace("右:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//上
	checkX = x
	checkY = y - 2
	//上边是否可移动 左上 右上是否可移动
	if this.isWalkable(checkX, checkY) && this.isWalkable(x-1+(y&1), y-1) && this.isWalkable(x+(y&1), y-1) && !this.mapData[checkY][checkX].IsInClose {
		//trace("上:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//下
	checkX = x
	checkY = y + 2
	//下边是否可移动 左下 右下是否可移动
	if this.isWalkable(checkX, checkY) && this.isWalkable(x-1+(y&1), y+1) && this.isWalkable(x+(y&1), y+1) && !this.mapData[checkY][checkX].IsInClose {
		//trace("下:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//左上
	checkX = x - 1 + (y & 1)
	checkY = y - 1
	if this.isWalkable(checkX, checkY) && !this.mapData[checkY][checkX].IsInClose {
		//trace("左上:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//左下
	checkX = x - 1 + (y & 1)
	checkY = y + 1
	if this.isWalkable(checkX, checkY) && !this.mapData[checkY][checkX].IsInClose {
		//trace("左下:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//右上
	checkX = x + (y & 1)
	checkY = y - 1
	if this.isWalkable(checkX, checkY) && !this.mapData[checkY][checkX].IsInClose {
		//trace("右上:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}
	//右下
	checkX = x + (y & 1)
	checkY = y + 1
	if this.isWalkable(checkX, checkY) && !this.mapData[checkY][checkX].IsInClose {
		//trace("右下:::"+checkY+"##"+checkX);
		aroundNodes = append(aroundNodes, this.mapData[checkY][checkX])
	}

	return aroundNodes
}

/**
 * 检查点在地图上是否可通过
 */
func (this *AStartPathFinder) isWalkable(x int, y int) bool {
	// 1. 是否是有效的地图上点（数组边界检查）
	if x < this.xMapStart || x >= this.wMap {
		return false
	}
	if y < this.yMapStart || y >= this.hMap {

		return false
	}
	// 2. 是否是walkable
	return this.getWalkable(x, y)
}

/**
 * 判断是否是地图上的数值
 * @param xtile X坐标
 * @param ytile Y坐标
 * @return 真 假
 *
 */
func (this *AStartPathFinder) getWalkable(xtitle int, ytitle int) bool {
	yArrLength := len(this.moveArray)
	if yArrLength == 0 {
		return false
	}
	xArrLength := len(this.moveArray[ytitle])
	if xArrLength <= 0 {
		return false
	}

	if ytitle >= yArrLength {
		ytitle = yArrLength - 1
	} else if xtitle >= xArrLength {
		xtitle = xArrLength - 1
	}

	if it := this.moveArray[ytitle][xtitle]; it == 1 {

		return false
	} else {
		return true
	}
	return false
}

func (this *AStartPathFinder) destroyLists() {
	this.isFinded = false
	//this.closeList = new Array();
	//this.openList = new Array();
	this.closeList = list.New()
	this.openList = list.New()
	this.xMapStart = 0
	this.yMapStart = 0
	this.InitMap()
}
