package application

// SyncResponseHook is a no-op placeholder for Java Guice SyncResponseHook extensions (Phase 8).
type SyncResponseHook interface {
	Handle(deviceID int64, resp any) (any, error)
}
