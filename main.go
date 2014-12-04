package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type field struct {
	Name   string
	Type   string
	Values []string
}

var (
	typeMatcher = regexp.MustCompile(`\[\][^{]+`)
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "tests") {
			fields, values, testFuncName := getTestData(scanner)
			fmt.Println("var tests = []struct{")
			fmt.Println("var tests = []struct{")
		} else {
			fmt.Println(line)
		}
	}
}

func getTestData(scanner *bufio.Scanner) (fields []field, values []string, testFuncName string) {
	fields = parseFields(scanner)
	values = findAllValues(scanner)
	// Scanner should be at the for ... line, just skip it
	scanner.Scan()
	scanner.Scan()
	testFuncName = getTestFuncName(scanner.Text())
	// also skip the if line
	scanner.Scan()
	return
}

func getTestFuncName(line string) string {
	testFunctionCall := line[strings.Index(line, ":=")+3:]
	return testFunctionCall[:strings.Index(testFunctionCall, "(")]
}

func parseFields(scanner *bufio.Scanner) (fields []field) {
	fields = []field(nil)
	for !strings.HasPrefix(strings.TrimSpace(scanner.Text()), "expected") {
		fields = append(fields, parseField(scanner))
		scanner.Scan()
	}
	return
}

func parseField(scanner *bufio.Scanner) field {
	line := scanner.Text()
	fieldType := typeMatcher.FindString(line)
	fieldName := strings.Replace(strings.Fields(line)[0], "tests", "param", 1)
	values := findAllValues(scanner)

	return field{
		Name:   fieldName,
		Type:   fieldType,
		Values: values,
	}
}

func findAllValues(scanner *bufio.Scanner) (allValues []string) {
	line := scanner.Text()
	allValues = []string(nil)
	afterOpeningBracket := line[strings.Index(line, "{"):]
	for line = afterOpeningBracket; !strings.Contains(line, "}"); line = scanner.Text() {
		values := splitCommaAndSpaceSeparatedString(line)
		allValues = append(allValues, values...)

		scanner.Scan()
	}
	beforeClosingBracket := line[:len(line)-1]
	values := splitCommaAndSpaceSeparatedString(beforeClosingBracket)
	allValues = append(allValues, values...)
	return
}

func splitCommaAndSpaceSeparatedString(s string) (values []string) {
	valuesWithCommas := strings.Fields(s)
	values = make([]string, len(valuesWithCommas))
	for i, valueWithComma := range valuesWithCommas {
		values[i] = valueWithComma[:len(valueWithComma)-1]
	}
	return
}
