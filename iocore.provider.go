package elio

// ProvideIoCore provide service
func ProvideIoCore(c ConfigIo, s Service) *IoCore {
	return NewIoCore(c, s)
}
