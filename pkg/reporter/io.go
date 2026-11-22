package reporter

import "os"

func writeBytes(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
