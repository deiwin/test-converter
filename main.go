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
			printStructuredTestData(fields, values, testFuncName)
		} else {
			fmt.Println(line)
		}
	}
}

func printStructuredTestData(fields []field, values []string, testFuncName string) {
	fmt.Println("\tvar tests = []struct{")
	for _, f := range fields {
		fmt.Printf("\t\t%s\t%s\n", f.Name, f.Type)
	}
	fmt.Printf("\t\t%s\t%s\n", "expected", "bool")
	fmt.Println("\t}{")
	for i, value := range values {
		fieldValues := getNthFieldValues(fields, i)
		fieldPlusExpectedValues := append(fieldValues, value)
		fmt.Printf("\t\t{%s},\n", strings.Join(fieldPlusExpectedValues, ", "))
	}
	fmt.Println("\tfor _, test := range tests {")
	fmt.Printf("\t\tactual := %s(%s)\n", testFuncName, getCommaSeparatedFieldTypes(fields))
	fmt.Println("\t\tif actual != test.expected {")
}

func getCommaSeparatedFieldTypes(fields []field) string {
	fieldTypes := make([]string, len(fields))
	for i, f := range fields {
		fieldTypes[i] = f.Type
	}
	return strings.Join(fieldTypes, ", ")
}

func getNthFieldValues(fields []field, n int) (fieldValues []string) {
	fieldValues = make([]string, len(fields))
	for i, f := range fields {
		fieldValues[i] = f.Values[n]
	}
	return
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
