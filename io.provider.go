package elio

// ProvideIo provide service
func ProvideIo(h *IoHost, c ConfigIo, s Service) *Io {
	return NewIo(h, c, s)
}
