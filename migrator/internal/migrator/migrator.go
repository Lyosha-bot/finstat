package migrator

import (
	"fmt"
	"log"
	"path/filepath"
	ewrap "postgres-migrator/internal/lib"
	"postgres-migrator/internal/repository"
	"postgres-migrator/migrations"
	"strconv"
	"strings"
)

const SEPARATOR = "_"

func ApplyMigrations(creds repository.Credentials) error {
	conn, err := repository.NewConnection(creds)
	if err != nil {
		return ewrap.Wrap("Couldn't establish connection", err)
	}
	defer conn.Close()

	existingMigrations, err := conn.Migrations()
	if err != nil {
		return ewrap.Wrap("Couldn't get existing migrations", err)
	}

	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return ewrap.Wrap("Couldn't read directory", err)
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if filepath.Ext(fileName) != ".sql" {
			continue
		}

		log.Println("Migrating", fileName)

		fileName = strings.TrimSuffix(fileName, ".sql")
		split := strings.Split(fileName, SEPARATOR)
		verStr := split[0]
		version, err := strconv.ParseUint(verStr, 10, 64)
		if err != nil {
			return ewrap.Wrap(fmt.Sprintf("Couldn't parse %s to uint", verStr), err)
		}

		if !existingMigrations[version] {
			contentBytes, err := migrations.FS.ReadFile(entry.Name())
			if err != nil {
				return ewrap.Wrap(fmt.Sprintf("Couldn't open %s", fileName), err)
			}

			name := ""
			if len(split) > 1 {
				name = split[1]
			}

			err = conn.ApplyMigration(version, name, string(contentBytes))
			if err != nil {
				return ewrap.Wrap(fmt.Sprintf("Couldn't apply %d migration", version), err)
			}

			existingMigrations[version] = true
		}
	}

	return nil
}
