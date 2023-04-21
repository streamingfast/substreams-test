package validator

type Stats struct {
	successCount uint64
	failedCount  uint64
}

func (s *Stats) GetSuccessCount() uint64 {
	return s.successCount
}

func (s *Stats) GetFailedCount() uint64 {
	return s.failedCount
}

func (s *Stats) GetTestCount() uint64 {
	return s.successCount + s.failedCount
}
