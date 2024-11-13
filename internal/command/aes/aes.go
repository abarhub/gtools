package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"gtools/internal/command/base64"
	"io"
	"os"
)

type AesParameters struct {
	InputFile   string
	OutpoutFile string
	Encrypt     bool
	Password    string
}

const chunkSize = 64 * 1024 // 64KB chunks

func AesCommand(param AesParameters) error {
	return aesCmd(param, os.Stdout)
}

func aesCmd(param AesParameters, out io.Writer) error {

	if param.Encrypt {

		key := make([]byte, 32) // Clé AES-256
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("Erreur génération clé: %v\n", err)
		}

		err := encrypt(key, param.InputFile, param.OutpoutFile)
		if err != nil {
			return err
		}

		s, err := base64.EncodeStr(key)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(out, "%s\n", s)
		if err != nil {
			return err
		}

	} else {
		pwd := param.Password
		key, err := base64.DecodeStr([]byte(pwd))
		if err != nil {
			return err
		}
		err = decrypt(key, param.InputFile, param.OutpoutFile)
		if err != nil {
			return err
		}
	}

	return nil
}

//func encrypt0(plaintext string, secretKey []byte) string {
//	aes, err := aes.NewCipher(secretKey)
//	if err != nil {
//		panic(err)
//	}
//
//	gcm, err := cipher.NewGCM(aes)
//	if err != nil {
//		panic(err)
//	}
//
//	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
//	// A nonce should always be randomly generated for every encryption.
//	nonce := make([]byte, gcm.NonceSize())
//	_, err = rand.Read(nonce)
//	if err != nil {
//		panic(err)
//	}
//
//	// ciphertext here is actually nonce+ciphertext
//	// So that when we decrypt, just knowing the nonce size
//	// is enough to separate it from the ciphertext.
//	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
//
//	return string(ciphertext)
//}
//
//func decrypt0(ciphertext string, secretKey []byte) string {
//	aes, err := aes.NewCipher(secretKey)
//	if err != nil {
//		panic(err)
//	}
//
//	gcm, err := cipher.NewGCM(aes)
//	if err != nil {
//		panic(err)
//	}
//
//	// Since we know the ciphertext is actually nonce+ciphertext
//	// And len(nonce) == NonceSize(). We can separate the two.
//	nonceSize := gcm.NonceSize()
//	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
//
//	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
//	if err != nil {
//		panic(err)
//	}
//
//	return string(plaintext)
//}

func encrypt(key []byte, inputFile, outputFile string) error {
	// Ouvrir le fichier source
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("erreur ouverture fichier source: %v", err)
	}
	defer in.Close()

	// Créer le fichier de destination
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("erreur création fichier destination: %v", err)
	}
	defer out.Close()

	// Créer le cipher block AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("erreur création cipher AES: %v", err)
	}

	// Créer le mode GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("erreur création GCM: %v", err)
	}

	// Générer un nonce aléatoire
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("erreur génération nonce: %v", err)
	}

	// Écrire le nonce au début du fichier chiffré
	if _, err := out.Write(nonce); err != nil {
		return fmt.Errorf("erreur écriture nonce: %v", err)
	}

	// Lire et chiffrer par morceaux
	buf := make([]byte, chunkSize)
	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lecture fichier: %v", err)
		}

		// Chiffrer le morceau
		ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)

		// Écrire la taille du morceau chiffré suivi du morceau
		if err := writeChunk(out, ciphertext); err != nil {
			return fmt.Errorf("erreur écriture morceau chiffré: %v", err)
		}
	}

	return nil
}

func decrypt(key []byte, inputFile, outputFile string) error {
	// Ouvrir le fichier chiffré
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("erreur ouverture fichier chiffré: %v", err)
	}
	defer in.Close()

	// Créer le fichier de sortie
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("erreur création fichier déchiffré: %v", err)
	}
	defer out.Close()

	// Créer le cipher block AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("erreur création cipher AES: %v", err)
	}

	// Créer le mode GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("erreur création GCM: %v", err)
	}

	// Lire le nonce du début du fichier
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(in, nonce); err != nil {
		return fmt.Errorf("erreur lecture nonce: %v", err)
	}

	// Lire et déchiffrer par morceaux
	for {
		// Lire la taille du prochain morceau
		chunk, err := readChunk(in)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lecture morceau chiffré: %v", err)
		}

		// Déchiffrer le morceau
		plaintext, err := gcm.Open(nil, nonce, chunk, nil)
		if err != nil {
			return fmt.Errorf("erreur déchiffrement: %v", err)
		}

		// Écrire le morceau déchiffré
		if _, err := out.Write(plaintext); err != nil {
			return fmt.Errorf("erreur écriture morceau déchiffré: %v", err)
		}
	}

	return nil
}

func readChunk(r io.Reader) ([]byte, error) {
	// Lire la taille du morceau (4 bytes)
	sizeBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, sizeBytes); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(sizeBytes)

	// Lire le morceau
	chunk := make([]byte, size)
	if _, err := io.ReadFull(r, chunk); err != nil {
		return nil, err
	}

	return chunk, nil
}

func writeChunk(w io.Writer, chunk []byte) error {
	// Écrire la taille du morceau sur 4 bytes
	size := uint32(len(chunk))
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, size)

	if _, err := w.Write(sizeBytes); err != nil {
		return err
	}

	// Écrire le morceau
	_, err := w.Write(chunk)
	return err
}
