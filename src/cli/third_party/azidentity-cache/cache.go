package cache

// Cache is a tiny local stub for the optional persistent cache module.
// Sprint 1 doesn't use the persistent-cache surface, so this keeps the build offline.
type Cache struct{}

// New returns an empty stub cache.
func New(_ any) (Cache, error) {
	return Cache{}, nil
}
