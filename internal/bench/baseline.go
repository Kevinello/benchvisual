package bench

// Baseline compare benchmark result with baseline
//
//	@param sets []Set
//	@param baselines []float64
//	@return []Set
//	@author kevineluo
//	@update 2023-04-18 04:14:32
func Baseline(sets []Set, baselines []float64) {
	for setIdx, set := range sets {
		for target, benchList := range set.Targets {
			for idx, benchmark := range benchList {
				var (
					fast  bool
					small bool
					less  bool
				)
				if baselines[0] > 0 {
					// consider cpu core nums when compare
					fast = benchmark.NsPerOp*float64(benchmark.CPUCores) < baselines[0]
				}
				if baselines[1] > 0 {
					small = benchmark.Mem.BytesPerOp < baselines[1]
				}
				if baselines[2] > 0 {
					less = benchmark.Mem.AllocsPerOp < baselines[2]
				}
				sets[setIdx].Targets[target][idx].ReachBaseline = fast && small && less
			}
		}
	}
}
