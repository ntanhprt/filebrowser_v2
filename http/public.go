package fbhttp

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/crypto/bcrypt"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
)

var withHashFile = func(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, ifPath := ifPathWithName(r)
		link, err := d.store.Share.GetByHash(id)
		if err != nil {
			return errToStatus(err), err
		}

		status, err := authenticateShareRequest(r, link)
		if status != 0 || err != nil {
			return status, err
		}

		user, err := d.store.Users.Get(d.server.Root, link.UserID)
		if err != nil {
			return errToStatus(err), err
		}

		d.user = user

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       link.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
			Token:      link.Token,
		})
		if err != nil {
			return errToStatus(err), err
		}

		// share base path
		basePath := link.Path

		// file relative path
		filePath := ""

		if file.IsDir {
			basePath = filepath.Dir(basePath)
			filePath = ifPath
		}

		// set fs root to the shared file/folder
		d.user.Fs = afero.NewBasePathFs(d.user.Fs, basePath)

		file, err = files.NewFileInfo(&files.FileOptions{
			Fs:      d.user.Fs,
			Path:    filePath,
			Modify:  d.user.Perm.Modify,
			Expand:  true,
			Checker: d,
			Token:   link.Token,
		})
		if err != nil {
			return errToStatus(err), err
		}

		if file.IsDir {
			// extract name from the last directory in the path
			name := filepath.Base(strings.TrimRight(link.Path, string(filepath.Separator)))
			file.Name = name
		}

		d.raw = file
		return fn(w, r, d)
	}
}

// ref to https://github.com/filebrowser/filebrowser/pull/727
// `/api/public/dl/MEEuZK-v/file-name.txt` for old browsers to save file with correct name
func ifPathWithName(r *http.Request) (id, filePath string) {
	pathElements := strings.Split(r.URL.Path, "/")
	// prevent maliciously constructed parameters like `/api/public/dl/XZzCDnK2_not_exists_hash_name`
	// len(pathElements) will be 1, and golang will panic `runtime error: index out of range`

	switch len(pathElements) {
	case 1:
		return r.URL.Path, "/"
	default:
		joinedPath := path.Join("/", path.Join(pathElements[1:]...))
		// Preserve trailing slash for directory creation
		if strings.HasSuffix(r.URL.Path, "/") && !strings.HasSuffix(joinedPath, "/") {
			joinedPath += "/"
		}
		return pathElements[0], joinedPath
	}
}

// PublicShareResponse wraps FileInfo with share permission info
type PublicShareResponse struct {
	*files.FileInfo
	ShareInfo *SharePermissionInfo `json:"shareInfo,omitempty"`
}

// SharePermissionInfo contains permission info for the share
type SharePermissionInfo struct {
	Hash       string `json:"hash"`
	IsPublic   bool   `json:"isPublic"`
	Permission string `json:"permission"`
	CanEdit    bool   `json:"canEdit"`
}

var publicShareHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	id, ifPath := ifPathWithName(r)
	link, err := d.store.Share.GetByHash(id)
	if err != nil {
		return errToStatus(err), err
	}

	status, err := authenticateShareRequest(r, link)
	if status != 0 || err != nil {
		return status, err
	}

	user, err := d.store.Users.Get(d.server.Root, link.UserID)
	if err != nil {
		return errToStatus(err), err
	}

	d.user = user

	file, err := files.NewFileInfo(&files.FileOptions{
		Fs:         d.user.Fs,
		Path:       link.Path,
		Modify:     d.user.Perm.Modify,
		Expand:     false,
		ReadHeader: d.server.TypeDetectionByHeader,
		Checker:    d,
		Token:      link.Token,
	})
	if err != nil {
		return errToStatus(err), err
	}

	// share base path
	basePath := link.Path

	// file relative path
	filePath := ""

	if file.IsDir {
		basePath = filepath.Dir(basePath)
		filePath = ifPath
	}

	// set fs root to the shared file/folder
	d.user.Fs = afero.NewBasePathFs(d.user.Fs, basePath)

	file, err = files.NewFileInfo(&files.FileOptions{
		Fs:      d.user.Fs,
		Path:    filePath,
		Modify:  d.user.Perm.Modify,
		Expand:  true,
		Checker: d,
		Token:   link.Token,
	})
	if err != nil {
		return errToStatus(err), err
	}

	if file.IsDir {
		// extract name from the last directory in the path
		name := filepath.Base(strings.TrimRight(link.Path, string(filepath.Separator)))
		file.Name = name
		file.Sorting = files.Sorting{By: "name", Asc: false}
		file.ApplySort()
	}

	// Create response with share info
	response := &PublicShareResponse{
		FileInfo: file,
		ShareInfo: &SharePermissionInfo{
			Hash:       link.Hash,
			IsPublic:   link.IsPublic,
			Permission: link.Permission,
			CanEdit:    link.IsPublic && link.Permission == "change",
		},
	}

	return renderJSON(w, r, response)
}

var publicDlHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)
	if !file.IsDir {
		return rawFileHandler(w, r, file)
	}

	return rawDirHandler(w, r, d, file)
})

func authenticateShareRequest(r *http.Request, l *share.Link) (int, error) {
	if l.PasswordHash == "" {
		return 0, nil
	}

	if r.URL.Query().Get("token") == l.Token {
		return 0, nil
	}

	password := r.Header.Get("X-SHARE-PASSWORD")
	password, err := url.QueryUnescape(password)
	if err != nil {
		return 0, err
	}
	if password == "" {
		return http.StatusUnauthorized, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(l.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return http.StatusUnauthorized, nil
		}
		return 0, err
	}

	return 0, nil
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"OK"}`))
}

// withHashFileAuth is similar to withHashFile but requires user authentication
// and checks for edit permission on the share
var withHashFileAuth = func(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, ifPath := ifPathWithName(r)
		link, err := d.store.Share.GetByHash(id)
		if err != nil {
			return errToStatus(err), err
		}

		// Check if this is a public share with change permission
		if !link.IsPublic || link.Permission != "change" {
			return http.StatusForbidden, fmt.Errorf("share does not allow editing")
		}

		// Check if current user is allowed to edit
		// AllowedUsers empty means everyone can edit
		if len(link.AllowedUsers) > 0 {
			allowed := false
			for _, uid := range link.AllowedUsers {
				if uid == d.user.ID {
					allowed = true
					break
				}
			}
			if !allowed && !d.user.Perm.Admin {
				return http.StatusForbidden, fmt.Errorf("user not allowed to edit this share")
			}
		}

		// Get the owner user to access their filesystem
		owner, err := d.store.Users.Get(d.server.Root, link.UserID)
		if err != nil {
			return errToStatus(err), err
		}

		// Get file info from owner's fs
		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         owner.Fs,
			Path:       link.Path,
			Modify:     true,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})
		if err != nil {
			return errToStatus(err), err
		}

		// share base path
		basePath := link.Path

		if file.IsDir {
			basePath = filepath.Dir(basePath)
		}

		// set fs root to the shared file/folder using owner's fs
		sharedFs := afero.NewBasePathFs(owner.Fs, basePath)
		
		// Store original user and set owner's fs for operations
		d.raw = &sharedResourceData{
			link:     link,
			filePath: ifPath,
			sharedFs: sharedFs,
			owner:    owner,
		}
		
		return fn(w, r, d)
	})
}

// sharedResourceData holds data for shared resource operations
type sharedResourceData struct {
	link     *share.Link
	filePath string
	sharedFs afero.Fs
	owner    interface{}
}

// publicSharePostHandler handles file/folder creation in shared folders
func publicSharePostHandler(fileCache FileCache) handleFunc {
	return withHashFileAuth(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		srd := d.raw.(*sharedResourceData)
		
		filePath := srd.filePath
		if filePath == "" || filePath == "/" {
			// Get path from URL after hash
			_, filePath = ifPathWithName(r)
		}

		// Directories creation on POST.
		if strings.HasSuffix(filePath, "/") {
			err := srd.sharedFs.MkdirAll(filePath, d.settings.DirMode)
			return errToStatus(err), err
		}

		// Check if file exists
		_, err := srd.sharedFs.Stat(filePath)
		if err == nil {
			if r.URL.Query().Get("override") != "true" {
				return http.StatusConflict, nil
			}
		}

		// Write the file
		info, err := writeFile(srd.sharedFs, filePath, r.Body, d.settings.FileMode, d.settings.DirMode)
		if err != nil {
			return errToStatus(err), err
		}

		etag := fmt.Sprintf(`"%x%x"`, info.ModTime().UnixNano(), info.Size())
		w.Header().Set("ETag", etag)
		return http.StatusOK, nil
	})
}

// publicShareDeleteHandler handles file/folder deletion in shared folders
func publicShareDeleteHandler(fileCache FileCache) handleFunc {
	return withHashFileAuth(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
		srd := d.raw.(*sharedResourceData)
		
		filePath := srd.filePath
		if filePath == "" || filePath == "/" {
			return http.StatusForbidden, fmt.Errorf("cannot delete root of shared folder")
		}

		err := srd.sharedFs.RemoveAll(filePath)
		if err != nil {
			return errToStatus(err), err
		}

		return http.StatusNoContent, nil
	})
}

// publicSharePutHandler handles file modification in shared folders
func publicSharePutHandler() handleFunc {
	return withHashFileAuth(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		srd := d.raw.(*sharedResourceData)
		
		filePath := srd.filePath
		if filePath == "" || filePath == "/" {
			return http.StatusBadRequest, fmt.Errorf("cannot modify root")
		}

		// Only allow PUT for files
		if strings.HasSuffix(filePath, "/") {
			return http.StatusMethodNotAllowed, nil
		}

		// Check if file exists
		exists, err := afero.Exists(srd.sharedFs, filePath)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if !exists {
			return http.StatusNotFound, nil
		}

		info, err := writeFile(srd.sharedFs, filePath, r.Body, d.settings.FileMode, d.settings.DirMode)
		if err != nil {
			return errToStatus(err), err
		}

		etag := fmt.Sprintf(`"%x%x"`, info.ModTime().UnixNano(), info.Size())
		w.Header().Set("ETag", etag)
		return http.StatusOK, nil
	})
}

// publicSharePatchHandler handles file move/rename in shared folders
func publicSharePatchHandler(fileCache FileCache) handleFunc {
	return withHashFileAuth(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
		srd := d.raw.(*sharedResourceData)
		
		src := srd.filePath
		dst := r.URL.Query().Get("destination")
		action := r.URL.Query().Get("action")
		
		dst, err := url.QueryUnescape(dst)
		if err != nil {
			return errToStatus(err), err
		}
		
		if dst == "/" || src == "/" || src == "" {
			return http.StatusForbidden, nil
		}

		err = checkParent(src, dst)
		if err != nil {
			return http.StatusBadRequest, err
		}

		override := r.URL.Query().Get("override") == "true"
		rename := r.URL.Query().Get("rename") == "true"
		if !override && !rename {
			if _, err = srd.sharedFs.Stat(dst); err == nil {
				return http.StatusConflict, nil
			}
		}
		if rename {
			dst = addVersionSuffix(dst, srd.sharedFs)
		}

		switch action {
		case "copy":
			err = copyFile(srd.sharedFs, src, dst, d.settings.FileMode, d.settings.DirMode)
		case "rename":
			err = moveFile(srd.sharedFs, src, dst, d.settings.FileMode, d.settings.DirMode)
		default:
			return http.StatusBadRequest, fmt.Errorf("unsupported action: %s", action)
		}

		return errToStatus(err), err
	})
}

// copyFile copies a file or directory
func copyFile(afs afero.Fs, src, dst string, fileMode, dirMode fs.FileMode) error {
	srcInfo, err := afs.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(afs, src, dst, fileMode, dirMode)
	}

	return copySingleFile(afs, src, dst, fileMode)
}

// copySingleFile copies a single file
func copySingleFile(afs afero.Fs, src, dst string, fileMode fs.FileMode) error {
	// Read content from source
	content, err := afero.ReadFile(afs, src)
	if err != nil {
		return err
	}

	// Create destination directory if needed
	dstDir := path.Dir(dst)
	if err := afs.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	
	// Write to destination
	return afero.WriteFile(afs, dst, content, fileMode)
}

// copyDir copies a directory recursively
func copyDir(afs afero.Fs, src, dst string, fileMode, dirMode fs.FileMode) error {
	err := afs.MkdirAll(dst, dirMode)
	if err != nil {
		return err
	}

	entries, err := afero.ReadDir(afs, src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := path.Join(src, entry.Name())
		dstPath := path.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(afs, srcPath, dstPath, fileMode, dirMode)
		} else {
			err = copySingleFile(afs, srcPath, dstPath, fileMode)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// moveFile moves/renames a file or directory
func moveFile(afs afero.Fs, src, dst string, fileMode, dirMode fs.FileMode) error {
	return afs.Rename(src, dst)
}
