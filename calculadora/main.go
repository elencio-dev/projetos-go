package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func sumTwoNumbers(num01 float64, num02 float64) float64 {
	return num01 + num02
}

func divideTwoNumbers(num01 float64, num02 float64) (float64, error) {

	if num02 == 0 {
		return 0, errors.New("Não é possivel dividir por zero")
	}
	return num01 / num02, nil
}

func multiplyTwoNumbers(num01 float64, num02 float64) float64 {
	return num01 * num02
}

func subtractTwoNumbers(num01 float64, num02 float64) float64 {
	return num01 - num02
}

func salvarHistorico(f *os.File, registro string) {
	fmt.Fprintln(f, registro)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	f, err := os.OpenFile("historico.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao Abrir histórico", err)
		os.Exit(1)
	}

	defer f.Close()

	fmt.Print(">")
	for scanner.Scan() {
		linha := scanner.Text()

		if linha == "sair" {
			fmt.Println("Até logo")
			break
		}

		partes := strings.Fields(linha)

		if len(partes) != 3 {
			fmt.Println("Formato inválido. Use: numero operador numero")
			fmt.Println(">")
			continue
		}

		operator := partes[1]
		number01, err := strconv.ParseFloat(partes[0], 64)
		if err != nil {
			fmt.Println("Número invalido", err)
			continue
		}

		number02, err := strconv.ParseFloat(partes[2], 64)
		if err != nil {
			fmt.Println("Número invalido", err)
			continue
		}

		switch operator {
		case "+":
			result := sumTwoNumbers(number01, number02)
			fmt.Println("Resultado:", result)
			salvarHistorico(f, fmt.Sprintf("%g + %g = %g", number01, number02, result))
		case "-":
			result := subtractTwoNumbers(number01, number02)
			fmt.Println("Resultado:", result)
			salvarHistorico(f, fmt.Sprintf("%g - %g = %g", number01, number02, result))
		case "*":
			result := multiplyTwoNumbers(number01, number02)
			fmt.Println("Resultado:", result)
			salvarHistorico(f, fmt.Sprintf("%g * %g = %g", number01, number02, result))
		case "/":
			result, err := divideTwoNumbers(number01, number02)
			if err != nil {
				fmt.Println("Erro:", err)
			} else {
				fmt.Println("Resultado:", result)
				salvarHistorico(f, fmt.Sprintf("%g / %g = %g", number01, number02, result))
			}
		default:
			fmt.Println("Operador Inválido:", operator)
		}

		fmt.Println(">")

	}

}
