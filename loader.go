package dotenv

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

const maxStackLen = 50

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

	var pc [maxStackLen]uintptr
	n := runtime.Callers(1, pc[:])

	override := false

	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])

		if f.Name() == withOverridePackagePath {
			override = true

			break
		}
	}

	if !override {
		return godotenv.Overload(file)
	}

	tuples, err := godotenv.Read(file)
	if err != nil {
		return err
	}

	for k, v := range tuples {
		if _, defined := os.LookupEnv(k); !defined {
			if sErr := os.Setenv(k, v); sErr != nil {
				return nil
			}
		}
	}

	return nil
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
