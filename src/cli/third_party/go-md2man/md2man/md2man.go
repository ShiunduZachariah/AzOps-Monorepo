package md2man

// Render is a minimal local stub that satisfies Cobra's docs dependency.
// Sprint 1 doesn't generate man pages, so a pass-through implementation is enough.
func Render(input string) []byte {
	return []byte(input)
}
