package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// ProcessConfig populates the specified struct based on local environment variables and env files within a supplied directory
func ProcessConfig(envDir string, cfg interface{}) error {
	if envDir != "" {
		files, err := os.ReadDir(envDir)
		if err != nil {
			return err
		}
		for i := range files {
			f := files[i]

			fullFilePath := filepath.Join(envDir, f.Name())
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".env") {
				continue
			}

			if err := godotenv.Load(fullFilePath); err != nil {
				return err
			}
		}
	}

	return envconfig.Process("", cfg)
}
