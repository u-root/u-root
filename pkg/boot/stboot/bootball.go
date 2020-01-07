package stboot

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/bootconfig"
	"github.com/u-root/u-root/pkg/uzip"
)

const (
	signaturesDirName string = "signatures"
	rootCertName      string = "root.cert"
	bootFilesDirName  string = "bootconfig"
)

// BootBall contains data to operate on the system transparency
// bootball archive. There is an underlaying temporary directory
// representing the extracted archive.
type BootBall struct {
	Archive        string
	dir            string
	config         *Stconfig
	numBootConfigs int
	bootFiles      map[string][]string
	rootCert       *x509.CertPool
	signatures     map[string][]signature
	NumSignatures  int
	hashes         map[string][]byte
	Signer         Signer
}

// BootBallFromArchie constructs a BootBall zip file at archive
func BootBallFromArchie(archive string) (*BootBall, error) {
	var ball = new(BootBall)

	dir, err := ioutil.TempDir("", "bootball")
	if err != nil {
		return ball, fmt.Errorf("BootBall: cannot create tmp dir: %v", err)
	}

	err = uzip.FromZip(archive, dir)
	if err != nil {
		return ball, fmt.Errorf("BootBall: cannot unzip %s: %v", archive, err)
	}

	cfg, err := getConfig(filepath.Join(dir, ConfigName))
	if err != nil {
		return ball, fmt.Errorf("BootBall: getting configuration faild: %v", err)
	}

	ball.Archive = archive
	ball.dir = dir
	ball.config = cfg

	err = ball.init()
	if err != nil {
		return ball, err
	}

	return ball, nil
}

// BootBallFromConfig constructs a BootBall from a stconfig.json at configFile.
// the underlaying tmporary directory is created with standardized path and an
// updated copy of stconfig.json
func BootBallFromConfig(configFile string) (*BootBall, error) {
	var ball = new(BootBall)

	archive := filepath.Join(filepath.Dir(configFile), BallName)

	cfg, err := getConfig(configFile)
	if err != nil {
		return ball, fmt.Errorf("BootBall: getting configuration faild: %v", err)
	}

	dir, err := makeConfigDir(cfg, filepath.Dir(configFile))
	if err != nil {
		return ball, fmt.Errorf("BootBall: creating standard configuration directory faild: %v", err)
	}

	ball.Archive = archive
	ball.dir = dir
	ball.config = cfg

	err = ball.init()
	if err != nil {
		return ball, err
	}

	return ball, nil
}

func (ball *BootBall) init() error {
	cert, err := getRootCert(filepath.Join(ball.dir, signaturesDirName, rootCertName))
	if err != nil {
		return fmt.Errorf("BootBall: getting configuration faild: %v", err)
	}

	bootFiles, err := getBootFiles(ball.config, ball.dir)
	if err != nil {
		return fmt.Errorf("BootBall: getting boot files faild: %v", err)
	}

	ball.rootCert = cert
	ball.numBootConfigs = len(ball.config.BootConfigs)
	ball.bootFiles = bootFiles
	ball.Signer = Sha512PssSigner{}

	err = ball.getSignatures()
	if err != nil {
		return fmt.Errorf("BootBall: getting signatures: %v", err)
	}

	var x int = 0
	for _, sigPool := range ball.signatures {
		if x == 0 {
			x = len(sigPool)
			continue
		}
		if len(sigPool) != x {
			return errors.New("BootBall: invalid map of signatures")
		}
	}
	ball.NumSignatures = x
	return nil
}

// Clean removes the underlaying temporary directory.
func (ball *BootBall) Clean() error {
	err := os.RemoveAll(ball.dir)
	if err != nil {
		return err
	}
	ball.dir = ""
	return nil
}

// Pack archives the all contents of the underlaying temporary
// directory using zip.
func (ball *BootBall) Pack() error {
	if ball.Archive == "" || ball.dir == "" {
		return errors.New("BootBall.Pacstandak: booball.archive and bootball.dir must be set")
	}
	return uzip.ToZip(ball.dir, ball.Archive)
}

// Dir returns the temporary directory associated with BootBall.
func (ball *BootBall) Dir() string {
	return ball.dir
}

// GetBootConfigByIndex returns the Bootconfig at index from the BootBall's configs arrey.
func (ball *BootBall) GetBootConfigByIndex(index int) (*bootconfig.BootConfig, error) {
	bc, err := ball.config.getBootConfig(index)
	if err != nil {
		return nil, err
	}
	bc.SetFilePathsPrefix(ball.dir)
	return bc, nil
}

// Hash calculates hashes of all boot configurations in BootBall using the
// BootBall.Signer's hash function
func (ball *BootBall) Hash() error {
	ball.hashes = make(map[string][]byte)
	for i, files := range ball.bootFiles {
		hash, herr := ball.Signer.Hash(files...)
		if herr != nil {
			return herr
		}
		ball.hashes[i] = hash
	}
	return nil
}

// Sign signes the hashes of all boot configurations in BootBall using the
// BootBall.Signer's hash function with the provided privKeyFile. The signature
// is stored along with the provided certFile inside the BootBall.
func (ball *BootBall) Sign(privKeyFile, certFile string) error {
	err := validateFiles("", privKeyFile, certFile)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}

	cert, err := parseCertificate(buf)
	if err != nil {
		return err
	}

	err = validateCertificate(cert, ball.rootCert)
	if err != nil {
		return err
	}

	log.Printf("Signing with: %s", privKeyFile)

	if ball.hashes == nil {
		err = ball.Hash()
		if err != nil {
			return err
		}
	}

	sigs := make([]signature, 0)
	for _, hash := range ball.hashes {
		s, err := ball.Signer.Sign(privKeyFile, hash)
		if err != nil {
			return err
		}
		sigs = append(sigs, signature{
			Bytes: s,
			Cert:  cert})
	}

	if err = writeSignatures(sigs, certFile, ball.dir); err != nil {
		return err
	}

	ball.NumSignatures++
	return nil
}

// VerifyBootconfigs validates the certificates stored together with the
// signatures of each boot configuration in BootBall and verifies the
// signatures. A map is returned with the BootConfig's name as key and the
// according number of valid signatures of this BootConfig.
func (ball *BootBall) VerifyBootconfigs() (map[string]int, error) {
	verified := make(map[string]int)
	for i := 0; 1 < ball.NumSignatures; i++ {
		n, err := ball.VerifyBootconfigByIndex(i)
		if err != nil {
			return nil, err
		}
		verified[ball.config.BootConfigs[i].Name] = n
	}
	return verified, nil
}

// VerifyBootconfigByIndex validates the certificates stored together with the
// signatures of BootConfig index at BootBall.Config.BootConfigs[] and verifies
// the signatures. The number of valid signatures is returned.
func (ball *BootBall) VerifyBootconfigByIndex(index int) (int, error) {
	bcName := ball.config.BootConfigs[index].Name
	return ball.VerifyBootconfigByName(bcName)
}

// VerifyBootconfigByName validates the certificates stored together with the
// signatures of BootConfig name and verifies the signatures. The number of
// valid signatures is returned.
func (ball *BootBall) VerifyBootconfigByName(name string) (int, error) {
	if ball.hashes == nil {
		err := ball.Hash()
		if err != nil {
			return 0, err
		}
	}

	sigs := ball.signatures[name]
	var verified int = 0
	for _, sig := range sigs {
		err := validateCertificate(sig.Cert, ball.rootCert)
		if err != nil {
			return verified, err
		}
		err = ball.Signer.Verify(sig, ball.hashes[name])
		if err != nil {
			return verified, err
		}
		verified++
	}
	return verified, nil
}

func getConfig(dest string) (*Stconfig, error) {
	cfgBytes, err := ioutil.ReadFile(dest)
	if err != nil {
		return nil, err
	}
	cfg, err := stconfigFromBytes(cfgBytes)
	if err != nil {
		return nil, err
	}
	if !(cfg.IsValid()) {
		return nil, errors.New("invalid configuration")
	}
	return cfg, nil
}

func getRootCert(dest string) (*x509.CertPool, error) {
	certBytes, err := ioutil.ReadFile(dest)
	if err != nil {
		return nil, err
	}
	cert, err := certPool(certBytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func getBootFiles(cfg *Stconfig, prefix string) (map[string][]string, error) {
	bootFiles := make(map[string][]string)
	for _, bc := range cfg.BootConfigs {
		files := make([]string, 0)
		for _, file := range bc.FileNames() {
			file = filepath.Join(prefix, file)
			files = append(files, file)
			if err := validateFiles("", files...); err != nil {
				return nil, err
			}
		}
		bootFiles[bc.Name] = files
	}
	return bootFiles, nil
}

func (ball *BootBall) getSignatures() error {
	ball.signatures = make(map[string][]signature)
	path := filepath.Join(ball.dir, signaturesDirName)

	sigPool := make([]signature, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(info.Name())

		if !info.IsDir() && (ext == ".signature") {
			sigBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			dir := filepath.Dir(path)
			index, err := strconv.Atoi(dir[len(dir)-1:])
			if err != nil {
				return err
			}

			certFile := strings.TrimSuffix(path, filepath.Ext(path)) + ".cert"
			certBytes, err := ioutil.ReadFile(certFile)
			if err != nil {
				return err
			}

			cert, err := parseCertificate(certBytes)
			if err != nil {
				return err
			}

			sig := signature{
				Bytes: sigBytes,
				Cert:  cert,
			}
			sigPool = append(sigPool, sig)
			key := ball.config.BootConfigs[index].Name
			ball.signatures[key] = sigPool
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func validateFiles(prefix string, files ...string) error {
	for _, file := range files {
		_, err := os.Stat(filepath.Join(prefix, file))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeSignatures(sigs []signature, certFile, dir string) error {
	for i, sig := range sigs {
		d := fmt.Sprintf("%s%d", bootFilesDirName, i)
		path := filepath.Join(dir, signaturesDirName, d)
		os.Mkdir(path, os.ModePerm)

		id := fmt.Sprintf("%x", sig.Cert.PublicKey)[2:18]
		sigName := fmt.Sprintf("%s.signature", id)
		sigPath := filepath.Join(path, sigName)
		werr := ioutil.WriteFile(sigPath, sig.Bytes, 0644)
		if werr != nil {
			return werr
		}

		certName := fmt.Sprintf("%s.cert", id)
		certPath := filepath.Join(path, certName)
		cerr := copyFile(certFile, certPath)
		if cerr != nil {
			return cerr
		}
	}
	return nil
}

func makeConfigDir(cfg *Stconfig, origDir string) (string, error) {
	if err := validateFiles(cfg.RootCertPath); err != nil {
		return "", err
	}

	for _, bc := range cfg.BootConfigs {
		if err := validateFiles(origDir, bc.FileNames()...); err != nil {
			return "", err
		}
	}

	dir, err := ioutil.TempDir(os.TempDir(), "bootball")
	if err != nil {
		return "", err
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

		bc.ChangeFilePaths(dirName)
		cfg.BootConfigs[i] = bc
	}

	dstPath = filepath.Join(dir, ConfigName)
	bytes, err := cfg.bytes()
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(dstPath, bytes, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()

	return nil
}
