package main

type ConsensusData struct {
	hashCount       map[string]int
	totalCount      int
	maxCount        int
	maxHash         string
	myHash          string
	onConsensusDone func(bool)
}

func (cd *ConsensusData) AddVote(hash string) {
	if cd.hashCount == nil {
		cd.hashCount = make(map[string]int)
	}
	if cnt, ok := cd.hashCount[hash]; ok {
		cd.hashCount[hash] = cnt + 1
	} else {
		cd.hashCount[hash] = 1
	}
	newCnt := cd.hashCount[hash]
	if newCnt > cd.maxCount {
		cd.maxCount = newCnt
		cd.maxHash = hash
	}
	cd.totalCount += 1
}
