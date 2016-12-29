package compress

import (
	"github.com/goreleaser/releaser/config"
	"os"
	"github.com/goreleaser/releaser/uname"
	"compress/gzip"
	"archive/tar"
	"io"
	"fmt"
)

func ArchiveAll(version string, config config.ProjectConfig) error {
	for _, system := range config.Build.Oses {
		for _, arch := range config.Build.Arches {
			if err := create(system, arch, config); err != nil {
				return err
			}
		}
	}
	return nil
}

func create(system, arch string, config config.ProjectConfig) error {
	file, err := os.Create("dist/" + nameFor(system, arch, config.BinaryName) + ".tar.gz")
	fmt.Println("Creating", file.Name(), "...")
	if err != nil {
		return err
	}
	defer file.Close()
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, f := range config.Files {
		if err := addFile(tw, f, f); err != nil {
			return err
		}
	}
	if err := addFile(tw, config.BinaryName, binaryName(system, arch, config.BinaryName)); err != nil {
		return err
	}
	return nil
}

func addFile(tw *tar.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := new(tar.Header)
	header.Name = name
	header.Size = stat.Size()
	header.Mode = int64(stat.Mode())
	header.ModTime = stat.ModTime()
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tw, file); err != nil {
		return err
	}
	return nil
}

func nameFor(system, arch, binary string) string {
	return binary + "_" + uname.FromGo(system) + "_" + uname.FromGo(arch)
}

func binaryName(system, arch, binary string) string {
	return "dist/" + nameFor(system, arch, binary) + "/" + binary
}
