package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type tarefasDiurnas struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
	Concluido bool   `json:"concluido"`
}

const arquivoJson = "tarefas.json"
const arquivoTxt = "tarefas.txt"

func salvarJson(tarefas []tarefasDiurnas) {
	dados, err := json.MarshalIndent(tarefas, "", "")
	if err != nil {
		fmt.Println("Erro ao converter para Json", err)
		return
	}

	err = os.WriteFile(arquivoJson, dados, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar arquivo", err)
		return
	}
}

func carregarJson() []tarefasDiurnas {
	dados, err := os.ReadFile(arquivoJson)
	if err != nil {
		return []tarefasDiurnas{}
	}
	var tarefas []tarefasDiurnas

	err = json.Unmarshal(dados, &tarefas)

	if err != nil {
		fmt.Println("Erro ao ler Json", err)
		return []tarefasDiurnas{}
	}

	return tarefas

}

func lerInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func exportarTxt(tarefas []tarefasDiurnas) {

	f, err := os.Create(arquivoTxt)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo Txt", err)
		os.Exit(1)
	}

	defer f.Close()

	if len(tarefas) == 0 {
		fmt.Println("Nenhuma Tarefa cadastrada!")
		return
	}

	for _, t := range tarefas {
		status := "[]"
		if t.Concluido {
			status = "[x]"
		}

		linha := fmt.Sprintf("%d. %s %s - %s", t.ID, t.Titulo, status, t.Descricao)
		fmt.Fprintln(f, linha)
	}

	fmt.Println("Tarefas exportadas para", arquivoTxt)
}

func adicionarTarefas(scanner *bufio.Scanner, tarefas []tarefasDiurnas, proximoID *int) []tarefasDiurnas {
	titulo := lerInput(scanner, "Escreva o titulo da tarefa: \n")
	descricao := lerInput(scanner, "Escreva a Descrição da Tarefa: \n")

	tarefas = append(tarefas, tarefasDiurnas{
		ID:        *proximoID,
		Titulo:    titulo,
		Descricao: descricao,
		Concluido: false,
	})

	*proximoID++

	return tarefas

}

func salvarTarefa(task *os.File, tarefa string) {
	fmt.Fprintln(task, tarefa)
}

func marcarTarefaConcluida(scanner *bufio.Scanner, tarefas []tarefasDiurnas) []tarefasDiurnas {
	ID := lerInput(scanner, "ID da tarefa para marcar com Concluida: \n")

	id, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("Erro", err)
		return tarefas
	}

	for i := range tarefas {
		if tarefas[i].ID == id {
			tarefas[i].Concluido = true
			return tarefas
		}
	}

	return tarefas
}
func removerTarefa(scanner *bufio.Scanner, tarefas []tarefasDiurnas) []tarefasDiurnas {
	ID := lerInput(scanner, "ID da tarefa para remover: \n")

	id, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("Erro", err)
		return tarefas
	}

	for i := range tarefas {
		if tarefas[i].ID == id {
			tarefas = append(tarefas[:i], tarefas[i+1:]...)

			return tarefas
		}
	}

	return tarefas
}

func exibirMenu() {
	fmt.Println("====MENU====")
	fmt.Println("1. Adicionar Tarefa:")
	fmt.Println("2. Exportar Txt")
	fmt.Println("3. Marcar com concluida Tarefa:")
	fmt.Println("4. Remover Tarefa:")
	fmt.Println("5. Escolha: ")
}

func main() {
	tarefas := carregarJson()

	scanner := bufio.NewScanner(os.Stdin)
	proximoID := 1
	for _, t := range tarefas {
		if t.ID >= proximoID {
			proximoID = t.ID + 1
		}
	}

	for {
		exibirMenu()
		scanner.Scan()
		opcao := strings.TrimSpace(scanner.Text())

		switch opcao {
		case "1":
			tarefas = adicionarTarefas(scanner, tarefas, &proximoID)
			salvarJson(tarefas)
		case "2":
			exportarTxt(tarefas)
		case "3":
			tarefas = marcarTarefaConcluida(scanner, tarefas)
			salvarJson(tarefas)
		case "4":
			tarefas = removerTarefa(scanner, tarefas)
			salvarJson(tarefas)
		case "5":
			salvarJson(tarefas)
			return
		}
	}

}
