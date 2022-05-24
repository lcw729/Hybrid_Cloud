package resourceinfo

// 전체 스케줄링 레이어에 대한 기록
type CycleStatus struct {
	AvailableNodes int
	TotalNumNodes  int
}

func NewCycleStatus(totalNum int) *CycleStatus {
	return &CycleStatus{
		TotalNumNodes:  totalNum,
		AvailableNodes: totalNum,
	}
}

func (c *CycleStatus) MinusOneAvailableNodes() {
	c.AvailableNodes--
}

func (c *CycleStatus) IsAnyClusters() bool {
	if c.AvailableNodes == 0 {
		return false
	} else {
		return true
	}
}
