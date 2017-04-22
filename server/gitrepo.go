package main

import (
	git "gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetRepo(dir, repo string) error {

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               repo,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		logger.Error(err.Error())
		return err
	} else {
		logger.Info("gitrepo: clone of repo completed")
		return nil
	}

}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	logger.Info("gitrepo: repo removed")
	return nil
}

func GetAppConfig(filename string) ([]byte, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Info("config " + filename + " read")
	return file, nil
}
