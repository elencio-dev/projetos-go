package main

import (
	"crypto/rand"
	"flag" // ler argumentos da linha de comando
	"fmt"
	"strings"
)

func gerarSenha(tamanho int, maiusculas bool, numeros bool, simbolos bool) string {
	chars := []byte("abcdefghijkmnloqrstuvwxyz")
	senha := []byte{}
	b := make([]byte, 1)

	if maiusculas {
		chars = append(chars, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")...)
	}

	if numeros {
		chars = append(chars, []byte("0123456789")...)
	}

	if simbolos {
		chars = append(chars, []byte("@!$%&*(#)")...)
	}

	for i := 0; i < tamanho; i++ {
		rand.Read(b)
		char := chars[int(b[0])%len(chars)]
		senha = append(senha, char)
	}

	return string(senha)

}

func avaliarSenha(senha string) string {
	points := 0

	if strings.ContainsAny(senha, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		points++
	}

	if strings.ContainsAny(senha, "@!$%&*(#)") {
		points++
	}

	if strings.ContainsAny(senha, "0123456789") {
		points++
	}

	if len(senha) >= 8 {
		points++
	}

	if len(senha) >= 16 {
		points++
	}

	if points <= 2 {
		return "Fraca"
	} else if points <= 4 {
		return "Média"
	}
	return "Forte"

}

func main() {

	tamanho := flag.Int("tamanho", 16, "comprimento da senha")
	maiusculas := flag.Bool("Maiusculas", true, "incluir maiusculas na senha")
	numeros := flag.Bool("numeros", true, "incluir números na senha")
	simbolos := flag.Bool("simbolos", true, "incluir simbolos na senha")

	flag.Parse()

	senha := gerarSenha(*tamanho, *maiusculas, *numeros, *simbolos)
	forca := avaliarSenha(senha)
	fmt.Println("Senha: ", senha)
	fmt.Println("Força da Senha: ", forca)

}
