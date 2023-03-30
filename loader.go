package dotenv

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const maxStackLen = 20

// Load loads the environment.
func Load() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	file := findDotEnv(cwd)
	if file == "" {
		return nil
	}

	return godotenv.Overload(file)
}

func findDotEnv(dir string) string {
	for {
		file := filepath.Join(dir, ".env")
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file
		}

		parent := "../"
		next := filepath.Clean(filepath.Join(dir, parent))

		if next == dir {
			return ""
		}

		dir = next
	}
}
