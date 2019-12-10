package stboot

import (
	"archive/zip"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	signaturesDirName string = "signatures"
	rootCertName      string = "root.cert"
	bootFilesDirName  string = "bootconfig"
)

type signature struct {
	Bytes []byte
	Cert  *x509.Certificate
}

type BootBall struct {
	archive        string
	dir            string
	config         *Stconfig
	numBootConfigs int
	bootFiles      [][]string
	rootCert       *x509.CertPool
	signatures     []signature
	hashes         [][]byte
}

func BootBallFromArchie(archive string) (BootBall, error) {
	var ball BootBall

	dir, err := ioutil.TempDir(os.TempDir(), "bootball")
	if err != nil {
		return ball, fmt.Errorf("bootball: cannot create tmp dir: %v", err)
	}

	err = fromZip(archive, dir)
	if err != nil {
		return ball, fmt.Errorf("bootball: cannot unzip %s: %v", archive, err)
	}

	cfg, err := getConfig(filepath.Join(dir, ConfigName))
	if err != nil {
		return ball, fmt.Errorf("bootball: getting configuration faild: %v", err)
	}

	cert, err := getRootCert(filepath.Join(dir, signaturesDirName, rootCertName))
	if err != nil {
		return ball, fmt.Errorf("bootball: getting configuration faild: %v", err)
	}

	bootFiles, err := getBootFiles(cfg, dir)
	if err != nil {
		return ball, fmt.Errorf("bootball: getting boot files faild: %v", err)
	}

	ball.archive = archive
	ball.dir = dir
	ball.config = cfg
	ball.rootCert = cert
	ball.numBootConfigs = len(ball.config.BootConfigs)
	ball.bootFiles = bootFiles

	return ball, nil
}

func BootBallFromConfig(configFile string) (BootBall, error) {
	var ball BootBall

	archive := filepath.Join(filepath.Dir(configFile), BallName)

	cfg, err := getConfig(configFile)
	if err != nil {
		return ball, fmt.Errorf("bootball: getting configuration faild: %v", err)
	}

	dir, err := makeConfigDir(cfg, filepath.Dir(configFile))
	if err != nil {
		return ball, fmt.Errorf("bootball: creating standard configuration directory faild: %v", err)
	}

	cert, err := getRootCert(filepath.Join(dir, signaturesDirName, rootCertName))
	if err != nil {
		return ball, fmt.Errorf("bootball: getting configuration faild: %v", err)
	}

	bootFiles, err := getBootFiles(cfg, dir)
	if err != nil {
		return ball, fmt.Errorf("bootball: getting boot files faild: %v", err)
	}

	ball.archive = archive
	ball.dir = dir
	ball.config = cfg
	ball.rootCert = cert
	ball.numBootConfigs = len(ball.config.BootConfigs)
	ball.bootFiles = bootFiles

	return ball, nil
}

func getConfig(dest string) (cfg *Stconfig, err error) {
	cfgBytes, err := ioutil.ReadFile(dest)
	if err != nil {
		return
	}
	cfg, err = stconfigFromBytes(cfgBytes)
	if err != nil {
		return
	}
	if !(cfg.IsValid()) {
		return cfg, errors.New("invalid configuration")
	}
	return
}

func getRootCert(dest string) (cert *x509.CertPool, err error) {
	certBytes, err := ioutil.ReadFile(dest)
	if err != nil {
		return
	}
	cert, err = certPool(certBytes)
	if err != nil {
		return
	}
	return
}

func getBootFiles(cfg *Stconfig, prefix string) (bootFiles [][]string, err error) {
	bootFiles = make([][]string, 0)
	for _, bc := range cfg.BootConfigs {
		files := make([]string, 0)
		for _, file := range bc.FileNames() {
			file = filepath.Join(prefix, file)
			files = append(files, file)
			if err = validateFiles("", files...); err != nil {
				return
			}
		}
		bootFiles = append(bootFiles, files)
	}
	return
}

func validateFiles(prefix string, files ...string) (err error) {
	for _, file := range files {
		_, err = os.Stat(filepath.Join(prefix, file))
		if err != nil {
			return
		}
	}
	return
}

func makeConfigDir(cfg *Stconfig, origDir string) (dir string, err error) {
	if err = validateFiles(cfg.RootCertPath); err != nil {
		return
	}

	for _, bc := range cfg.BootConfigs {
		if err = validateFiles(origDir, bc.FileNames()...); err != nil {
			return
		}
	}

	dir, err = ioutil.TempDir(os.TempDir(), "bootball")
	if err != nil {
		return
	}

	dstPath := filepath.Join(dir, signaturesDirName, rootCertName)
	srcPath := filepath.Join(origDir, cfg.RootCertPath)
	copyFile(srcPath, dstPath)

	for i, bc := range cfg.BootConfigs {
		dirName := fmt.Sprintf("%s%d", bootFilesDirName, i)
		for _, file := range bc.FileNames() {
			fileName := filepath.Base(file)
			dstPath := filepath.Join(dir, dirName, fileName)
			srcPath := filepath.Join(origDir, file)
			copyFile(srcPath, dstPath)
		}

		if bc.Kernel != "" {
			newPath := filepath.Join(dirName, path.Base(bc.Kernel))
			bc.Kernel = newPath
		}
		if bc.Initramfs != "" {
			newPath := filepath.Join(dirName, path.Base(bc.Initramfs))
			bc.Initramfs = newPath
		}
		if bc.DeviceTree != "" {
			newPath := filepath.Join(dirName, path.Base(bc.DeviceTree))
			bc.DeviceTree = newPath
		}
		if bc.Multiboot != "" {
			newPath := filepath.Join(dirName, path.Base(bc.Multiboot))
			bc.Multiboot = newPath
		}
		for j, mod := range bc.Modules {
			if mod != "" {
				newPath := filepath.Join(dirName, path.Base(mod))
				bc.Modules[j] = newPath
			}
		}
		cfg.BootConfigs[i] = bc
	}

	dstPath = filepath.Join(dir, ConfigName)
	bytes, err := cfg.bytes()
	if err != nil {
		return
	}
	ioutil.WriteFile(dstPath, bytes, os.ModePerm)

	return
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return
	}

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()

	return
}

func toZip(srcDir, dest string) (err error) {
	info, err := os.Stat(srcDir)
	if err != nil {
		return
	}
	if !(info.IsDir()) {
		return fmt.Errorf("%s is not a directory", srcDir)
	}
	archive, err := os.Create(dest)
	if err != nil {
		return
	}
	defer archive.Close()

	z := zip.NewWriter(archive)
	defer z.Close()

	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// do not include srcDir into archive
		if strings.Compare(info.Name(), filepath.Base(srcDir)) == 0 {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// adjust header.Name to preserve folder strulture
		header.Name = strings.TrimPrefix(path, srcDir)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := z.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return
}

func fromZip(src, destDir string) (err error) {
	z, err := zip.OpenReader(src)
	if err != nil {
		return
	}

	if err = os.MkdirAll(destDir, 0755); err != nil {
		return
	}

	for _, file := range z.File {
		path := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return
}
