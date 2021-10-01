package hooks

func Load() {
	// since all hooks use init method to add themselves this function does nothing
	// calling this function will just the code optimizer from annoyingly removing the "unused" package import
}