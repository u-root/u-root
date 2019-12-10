package stboot

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	urootcrypto "github.com/u-root/u-root/pkg/crypto"
)

// memoryZipReader is used to unpack a zip file from a byte sequence in memory.
type memoryZipReader struct {
	Content []byte
}

func (r *memoryZipReader) ReadAt(p []byte, offset int64) (n int, err error) {
	cLen := int64(len(r.Content))
	if offset > cLen {
		return 0, io.EOF
	}
	if cLen-offset >= int64(len(p)) {
		n = len(p)
		err = nil
	} else {
		err = io.EOF
		n = int(int64(cLen) - offset)
	}
	copy(p, r.Content[offset:int(offset)+n])
	return n, err
}

// FromZip tries to extract a Stconfig from provided bootball. The returned
// string argument is the temporary directory where the files were extracted,
// otherwise an error is returned
func FromZip(bootball string) (*Stconfig, string, error) {
	// load the whole zip file in memory - we need it anyway for the signature
	// matching.
	// TODO refuse to read if too big?
	data, err := ioutil.ReadFile(bootball)
	if err != nil {
		return nil, "", err
	}
	urootcrypto.TryMeasureData(urootcrypto.BlobPCR, data, bootball)
	zipbytes := data

	r, err := zip.NewReader(&memoryZipReader{Content: zipbytes}, int64(len(zipbytes)))
	if err != nil {
		return nil, "", err
	}
	tempDir, err := ioutil.TempDir(os.TempDir(), "bootconfig")
	if err != nil {
		return nil, "", err
	}
	log.Printf("Created temporary directory %s", tempDir)
	var cfg *Stconfig
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			// Dont care - will be handled later
			continue
		}

		destination := path.Join(tempDir, f.Name)
		if len(f.Name) == 0 {
			log.Printf("Warning: skipping zero-length file name (flags: %d, mode: %s)", f.Flags, f.Mode())
			continue
		}
		// Check if folder exists
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(destination), os.ModeDir|os.FileMode(0700)); err != nil {
				return nil, "", err
			}
		}
		fd, err := f.Open()
		if err != nil {
			return nil, "", err
		}
		buf, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, "", err
		}
		if f.Name == ConfigName {
			// make sure it's not a duplicate manifest within the ZIP file
			// and inform the user otherwise
			if cfg != nil {
				log.Printf("Warning: duplicate config file found, the last found wins")
			}
			// parse the configuration
			cfg, err = stconfigFromBytes(buf)
			if err != nil {
				return nil, "", err
			}
		}
		if err := ioutil.WriteFile(destination, buf, f.Mode()); err != nil {
			return nil, "", err
		}
		log.Printf("Extracted file %s (flags: %d, mode: %s)", f.Name, f.Flags, f.Mode())
	}
	if cfg == nil {
		return nil, "", errors.New("no manifest found")
	}
	return cfg, tempDir, nil
}

// ToZip tries to pack all files specified in the the provided config file
// into a zip archive. A copy of the config file is included in the resulting
// archive with adopted paths. The archive is created at output. An error is
// returned, if the files listed inside the file don't exist.
func ToZip(output string, config string) error {
	// Make sure the file is named according to the guidelines
	if base := path.Base(config); base != ConfigName {
		return fmt.Errorf("expect '%s', got: %s", ConfigName, base)
	}
	configBytes, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}
	cfg, err := stconfigFromBytes(configBytes)
	if err != nil {
		return fmt.Errorf("error parsing config: %v", err)
	} else if !cfg.IsValid() {
		return errors.New("config is not valid")
	}

	// Create a buffer to write the archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	z := zip.NewWriter(buf)

	var dest, origin string
	//Archive boot files
	for i, bc := range cfg.BootConfigs {
		dir := fmt.Sprintf("bootconfig_%d/", i)
		z.Create(dir)
		if bc.Kernel != "" {
			dest = path.Join(dir, path.Base(bc.Kernel))
			origin = path.Join(path.Dir(config), bc.Kernel)
			if err := tozip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack kernel: %v", err)
			}
			bc.Kernel = dest
		}
		if bc.Initramfs != "" {
			dest = path.Join(dir, path.Base(bc.Initramfs))
			origin = path.Join(path.Dir(config), bc.Initramfs)
			if err := tozip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack initramfs: %v", err)
			}
			bc.Initramfs = dest
		}
		if bc.DeviceTree != "" {
			dest = path.Join(dir, path.Base(bc.DeviceTree))
			origin = path.Join(path.Dir(config), bc.DeviceTree)
			if err := tozip(z, dest, origin); err != nil {
				return fmt.Errorf("cant pack device tree: %v", err)
			}
			bc.DeviceTree = dest
		}
		cfg.BootConfigs[i] = bc
	}

	// Archive root certificate
	z.Create("certs/")
	dest = "certs/root.cert"
	origin = path.Join(path.Dir(config), cfg.RootCertPath)
	if err = tozip(z, dest, origin); err != nil {
		return fmt.Errorf("cannot pack certificate: %v", err)
	}
	cfg.RootCertPath = dest

	// Archive config
	newConfig, err := cfg.bytes()
	if err != nil {
		return fmt.Errorf("cannot serialize modified config: %v", err)
	}
	dst, err := z.Create(path.Base(config))
	if err != nil {
		return fmt.Errorf("cannot create modified config file in archive: %v", err)
	}
	_, err = io.Copy(dst, bytes.NewReader(newConfig))
	if err != nil {
		return fmt.Errorf("cannot write modified config: %v", err)
	}

	// Write central directory of archive
	err = z.Close()
	if err != nil {
		return fmt.Errorf("cannot write central archive directory: %v", err)
	}

	err = ioutil.WriteFile(output, buf.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("cannot write archive to filesystem: %v", err)
	}
	return nil
}

func tozip(w *zip.Writer, newPath, originPath string) error {
	dst, err := w.Create(newPath)
	if err != nil {
		return err
	}
	// Copy content from inputpath to new file
	src, err := os.Open(originPath)
	if err != nil {
		return fmt.Errorf("cannot find %s specified in config", originPath)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return src.Close()
}
