package yang

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
	_ "github.com/mattn/go-sqlite3"
	"github.com/openconfig/goyang/pkg/yang"
)

// pathInfo stores the data for each discovered schema path
type pathInfo struct {
	path   string      // path is the full schema path of the node
	entry  *yang.Entry // entry is a pointer to the Entry struct corresponding to the node
	module string      // module is the name of the module from which the node was instantiated
	descr  string
}

func processModules(srcFiles []string) ([]*yang.Entry, error) {
	moduleSet := yang.NewModules()
	var errs []error
	for _, name := range srcFiles {
		err := moduleSet.Read(name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		return nil, errors.Join(errs...)
	}

	errs = moduleSet.Process()
	if errs != nil {
		return nil, errors.Join(errs...)
	}

	// since multiple modules may depend on a given module,
	// some modules appear multiple times in  Modules after
	// processing.  These are de-duplicated before generating
	// the Entries (this code is reused from ygot)
	var modNames []string
	mods := make(map[string]*yang.Module)
	for _, m := range moduleSet.Modules {
		if mods[m.Name] == nil {
			mods[m.Name] = m
			modNames = append(modNames, m.Name)
		}
	}
	// sort the deduped module names to ensure deterministic
	// order when processing paths
	sort.Strings(modNames)

	var entries []*yang.Entry
	for _, name := range modNames {
		entries = append(entries, yang.ToEntry(mods[name]))
	}
	return entries, nil
}

func processEntries(entries []*yang.Entry) ([]pathInfo, error) {
	var pathInfos []pathInfo
	for _, entry := range entries {
		if entry.Errors != nil {
			return nil, fmt.Errorf("errors in %s", entry.Name)
		}
		more, err := processEntry(entry)
		if err != nil {
			return nil, err
		}
		pathInfos = append(pathInfos, more...)
	}
	return pathInfos, nil
}

func processEntry(entry *yang.Entry) ([]pathInfo, error) {
	var pathInfos []pathInfo
	var module string
	if m, err := entry.InstantiatingModule(); err == nil {
		module = m
	}
	path := getEntryPath(entry)
	pathInfos = append(pathInfos, pathInfo{
		path:   path,
		entry:  entry,
		module: module,
		descr:  entry.Description,
	})
	for _, direntry := range entry.Dir {
		more, err := processEntry(direntry)
		if err != nil {
			return nil, err
		}
		pathInfos = append(pathInfos, more...)
	}
	return pathInfos, nil
}

// getEntryPath returns the schema path for the given Entry by following
// parent pointers
func getEntryPath(entry *yang.Entry) string {
	var path []string
	if entry.Parent == nil { // root entry
		return "/"
	}
	for e := entry; e.Parent != nil; e = e.Parent {
		path = append(path, e.Name)
	}
	// need to reverse the path created by the loop above
	for i := len(path)/2 - 1; i >= 0; i-- {
		o := len(path) - 1 - i
		path[i], path[o] = path[o], path[i]
	}
	return strings.Join(append([]string{""}, path...), "/")
}

func initDB(infos []pathInfo) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE VIRTUAL TABLE path_index USING fts5(path, description, type);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		var t string
		if info.entry.Type != nil {
			t = info.entry.Type.Name
		}
		sqlStmt := `INSERT INTO path_index (path, description, type) VALUES (?, ?, ?)`
		_, err = db.Exec(sqlStmt, info.path, info.descr, t)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func filePaths(dirPaths []string) ([]string, error) {
	var r fs.FS
	var srcFiles []string
	r, err := runfiles.New()
	if err != nil {
		r = os.DirFS(".")
	}

	for _, dirPath := range dirPaths {
		err := fs.WalkDir(r, dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".yang") {
				realPath, err := runfiles.Rlocation(path)
				if err != nil {
					srcFiles = append(srcFiles, path)
				} else {
					srcFiles = append(srcFiles, realPath)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return srcFiles, nil
}
