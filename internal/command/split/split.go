package split

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type SplitParameters struct {
	File          string
	SizeStr       string
	BufferSizeStr string
	size          int64
	bufferSize    int64
}

type FichierI interface {
	Write(b []byte) error
	Fin() error
}

type FichierS struct {
	compteurNomFichier   int
	nomFichier           string
	nomFichierCourant    string
	fichier              *os.File
	compteurDebutFichier int64
	tailleMax            int64
}

func SplitCommand(param SplitParameters) error {

	err := calculParametre(&param)
	if err != nil {
		return err
	}

	err = splitFile(param)

	return err
}

func calculParametre(param *SplitParameters) error {
	if param.BufferSizeStr == "" {
		param.bufferSize = 1024 * 8
	} else {
		var err error
		param.bufferSize, err = strconv.ParseInt(param.BufferSizeStr, 10, 64)
		if err != nil {
			return fmt.Errorf("buffer size must be an integer : %w", err)
		}
	}
	if param.SizeStr == "" {
		param.size = 1024 * 1024
	} else {
		var err error
		param.size, err = strconv.ParseInt(param.SizeStr, 10, 64)
		if err != nil {
			return fmt.Errorf("size must be an integer : %w", err)
		}
	}
	if param.bufferSize > param.size {
		param.bufferSize = param.size
	}
	return nil
}

func splitFile(param SplitParameters) error {

	BUFFERSIZE := param.bufferSize

	var tailleFichier int64 = 0
	src := param.File

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	fichier, err := creationFichier(src, param.size)
	if err != nil {
		return err
	}
	defer fichier.Fin()

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if err := fichier.Write(buf[:n]); err != nil {
			return err
		}
		tailleFichier += int64(n)

	}

	return nil
}

func creationFichier(nomFichier string, tailleMax int64) (*FichierS, error) {
	res := &FichierS{compteurNomFichier: 1, nomFichier: nomFichier,
		compteurDebutFichier: 0, tailleMax: tailleMax}
	err := res.creationFichier()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *FichierS) Write(b []byte) error {
	if len(b) > 0 && v.fichier == nil {
		err := v.creationFichier()
		if err != nil {
			return err
		}
	}

	if v.tailleMax > v.compteurDebutFichier+int64(len(b)) {

		if _, err := v.fichier.Write(b); err != nil {
			return err
		}
		v.compteurDebutFichier += int64(len(b))
	} else {
		n2 := v.tailleMax - v.compteurDebutFichier
		if n2 <= 0 {
			panic("n2<=0:" + strconv.FormatInt(n2, 10) + ", tailleMax=" +
				strconv.FormatInt(v.tailleMax, 10) + ",compteurDebutFichier=" +
				strconv.FormatInt(v.compteurDebutFichier, 10))
		}
		if _, err := v.fichier.Write(b[:n2]); err != nil {
			return err
		}
		err := v.fichier.Close()
		if err != nil {
			return err
		}
		v.fichier = nil

		n3 := int64(len(b)) - n2
		if n3 > 0 {
			if n3 < 0 {
				panic("n3<=0:" + strconv.FormatInt(n3, 10) + ", tailleMax=" +
					strconv.FormatInt(v.tailleMax, 10) + ",compteurDebutFichier=" +
					strconv.FormatInt(v.compteurDebutFichier, 10) + ",n=" +
					strconv.FormatInt(int64(len(b)), 10))
			}
			if n3 > int64(len(b)) {
				panic("n3>n:" + strconv.FormatInt(n3, 10))
			}
			err = v.creationFichier()
			if err != nil {
				return err
			}
			if _, err := v.fichier.Write(b[n2:]); err != nil {
				return err
			}
		}
		v.compteurDebutFichier = n3
	}
	return nil
}

func (v *FichierS) creationFichier() error {
	dst := fmt.Sprintf("%s.%03d", v.nomFichier, v.compteurNomFichier)
	fmt.Printf("creation de %s\n", dst)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	v.fichier = destination
	v.compteurNomFichier++
	v.nomFichierCourant = dst
	return nil
}

func (v *FichierS) Fin() error {
	if v.fichier != nil {
		err := v.fichier.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
