package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/grafana-tools/sdk"
)

// Dashboard contains the bare minimum information necessary to
// uniquely identify a Grafana dashboard and sync it disk.
type Dashboard struct {
	FolderTitle     string `json:"-"`
	FolderDirectory string ""
	Title           string `json:"title"`
	UID             string `json:"uid"`
	Filename        string `json:"-"`
}

func NewFromFoundBoard(db sdk.FoundBoard) Dashboard {
	folderDirectory := ""
	if db.FolderTitle == "" {
		folderDirectory = "General"
	} else {
		folderDirectory = db.FolderTitle
	}
	return Dashboard{
		FolderTitle:     db.FolderTitle,
		FolderDirectory: folderDirectory,
		Title:           db.Title,
		UID:             db.UID,
		Filename:        filepath.Join(folderDirectory, filepath.Base(db.URL)+".json"),
	}
}

func NewFromFile(baseDir, path string) (Dashboard, error) {
	rPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return Dashboard{}, err
	}
	dirName := filepath.Dir(rPath)
	if dirName == "." {
		dirName = ""
	}
	d := Dashboard{
		Filename:    rPath,
		FolderTitle: dirName,
	}
	f, err := os.Open(path)
	if err != nil {
		return Dashboard{}, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return Dashboard{}, err
	}
	return d, nil
}
