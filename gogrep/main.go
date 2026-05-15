package main

import (
	"bufio"         //leitura de arquivos
	"fmt"           //print de dados no terminal
	"os"            // acessar dados do sistema operacional
	"path/filepath" //manipular caminhos
	"strings"       // manipulação de strings
	"sync"
)

func buscarNovoArquivo(arquivo string, termo string, ch chan int, wy *sync.WaitGroup) {

	defer wy.Done() // informa que terminou a função

	ocorrencias := 0

	f, err := os.Open(arquivo)

	if err != nil {
		fmt.Println("Erro ao abrir o Arquivo")
		ch <- 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	linha := 1

	for scanner.Scan() {
		texto := scanner.Text()

		if strings.Contains(strings.ToLower(texto), strings.ToLower(termo)) {
			fmt.Printf("[%s] Linha %d: %s\n", arquivo, linha, texto)
			ocorrencias++
		}

		linha++

	}

	ch <- ocorrencias // envia o resultado pelo cannal

}

func main() {
	ch := make(chan int) // canal para resceber os resultados

	var wy sync.WaitGroup // controla o numero de gorountines rodando

	var lista []string

	if len(os.Args) < 3 {
		fmt.Println("Uso: gogrep <termo> <arquivo>")
		os.Exit(1)
	}

	termo := os.Args[1]
	arquivos := os.Args[2:]
	ocorrencias := 0

	for _, caminho := range arquivos {

		info, err := os.Stat(caminho)

		if err != nil {
			fmt.Println("Erro", err)
			continue
		}

		if info.IsDir() {
			filepath.Walk(caminho, func(path string, info os.FileInfo, err error) error {
				if strings.HasSuffix(path, ".go") {
					lista = append(lista, path)
				}

				return nil
			})
		} else {
			lista = append(lista, caminho)
		}
	}

	for _, arquivo := range lista {
		wy.Add(1)
		go buscarNovoArquivo(arquivo, termo, ch, &wy)
	}

	go func() {
		wy.Wait()
		close(ch)
	}()

	for resultado := range ch {
		ocorrencias += resultado
	}

	fmt.Printf("%d ocorrencias encontradas.\n", ocorrencias)

}
