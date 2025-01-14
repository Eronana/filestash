package backend

import (
	. "github.com/mickael-kerjean/filestash/server/common"
	"io/ioutil"
	"io"
	"os"
	"strings"
	"golang.org/x/crypto/bcrypt"
)


type FileSystem struct {
	rootPath string;
}

func init() {
	Backend.Register("filesystem", FileSystem{})
}

func (s FileSystem) Init(params map[string]string, app *App) (IBackend, error) {
  admin := Config.Get("auth.admin").String()
  if admin == "" {
    return nil, NewError("Missing admin account, please contact your administrator", 500)
  }

  if err := bcrypt.CompareHashAndPassword([]byte(admin), []byte(params["password"])); err != nil {
    return nil, ErrInvalidPassword
  }

	s.rootPath = strings.TrimRight(params["rootPath"], "/")
	return &s, nil
}

func (b FileSystem) LoginForm() Form {
	return Form{
		Elmnts: []FormElement{
			FormElement{
				Name:        "type",
				Type:        "hidden",
				Value:       "filesystem",
			},
			FormElement{
				Name:        "password",
				Type:        "password",
				Placeholder: "Passphrase",
			},
			FormElement{
				Name:        "rootPath",
				Type:        "text",
				Placeholder: "root path*",
			},
		},
	}
}

func (b FileSystem) GetFullPath(path string) (string) {
	return b.rootPath + path
}

func (b FileSystem) Ls(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(b.GetFullPath(path))
}

func (b FileSystem) Cat(path string) (io.ReadCloser, error) {
	return os.OpenFile(b.GetFullPath(path), os.O_RDONLY, os.ModePerm)
}

func (b FileSystem) Mkdir(path string) error {
	return os.Mkdir(b.GetFullPath(path), os.ModePerm)
}

func (b FileSystem) Rm(path string) error {
	return os.RemoveAll(b.GetFullPath(path))
}

func (b FileSystem) Mv(from string, to string) error {
	return os.Rename(b.GetFullPath(from), b.GetFullPath(to))
}

func (b FileSystem) Touch(path string) error {
	file, err := os.Create(b.GetFullPath(path))
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (b FileSystem) Save(path string, file io.Reader) error {
	newFile, err := os.Create(b.GetFullPath(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(newFile, file)
	newFile.Close()
	return err
}

func (b FileSystem) Stat(path string) (os.FileInfo, error) {
	f, err := os.Stat(b.GetFullPath(path))
	return f, err
}