package score

const (
	Accuracy = 10
)

func UInt64ToFloat32Score(score uint64) float32 {
	return float32(score) / Accuracy
}
func Float32ToUInt64Score(score float32) uint64 {
	return uint64(score * Accuracy)
}
