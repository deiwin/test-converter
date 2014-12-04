package convert

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type field struct {
	Name   string
	Type   string
	Values []string
}

var (
	typeMatcher     = regexp.MustCompile(`\[\][^{]+`)
	commaWhitespace = regexp.MustCompile(`, `)
)

// Test function converts an array-driven test to a table-driven test
func Test(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "tests") {
			fields, values, testFuncName := getTestData(scanner)
			printStructuredTestData(output, fields, values, testFuncName)
		} else {
			fmt.Fprintln(output, line)
		}
	}
}

func printStructuredTestData(output io.Writer, fields []field, values []string, testFuncName string) {
	fmt.Fprintln(output, "\tvar tests = []struct{")
	for _, f := range fields {
		fmt.Fprintf(output, "\t\t%s\t\t%s\n", f.Name, f.Type)
	}
	fmt.Fprintf(output, "\t\t%s\t%s\n", "expected", "bool")
	fmt.Fprintln(output, "\t}{")
	minNumberOfValues := getMinNumberOfValues(fields, values)
	for i, value := range values {
		if i >= minNumberOfValues {
			break
		}
		fieldValues := getNthFieldValues(fields, i)
		fieldPlusExpectedValues := append(fieldValues, value)
		fmt.Fprintf(output, "\t\t{%s},\n", strings.Join(fieldPlusExpectedValues, ", "))
	}
	fmt.Fprintln(output, "\t}")
	fmt.Fprintln(output, "\tfor i, test := range tests {")
	fmt.Fprintf(output, "\t\tactual := %s(%s)\n", testFuncName, getCommaSeparatedFieldNames(fields))
	fmt.Fprintln(output, "\t\tif actual != test.expected {")
	fmt.Fprintln(output, "\t\t\t"+`t.Log("Case ", i, ": expected ", test.expected, ", but result was ", actual)`)
}

func getMinNumberOfValues(fields []field, values []string) (min int) {
	min = len(values)
	for _, f := range fields {
		nrOfValues := len(f.Values)
		if nrOfValues < min {
			min = nrOfValues
		}
	}
	return
}

func getCommaSeparatedFieldNames(fields []field) string {
	fieldNames := make([]string, len(fields))
	for i, f := range fields {
		fieldNames[i] = "test." + f.Name
	}
	return strings.Join(fieldNames, ", ")
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
	// also skip the log line, it will be rewritten
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
	fieldType := typeMatcher.FindString(line)[2:] // [2:] to remove brackets
	fieldName := strings.Replace(strings.Fields(line)[0], "tests", "param", 1)
	values := findAllValues(scanner)

	return field{
		Name:   fieldName,
		Type:   fieldType,
		Values: values,
	}
}

func findAllValues(scanner *bufio.Scanner) []string {
	stringBetweenBraces := getStringBetweenBraces(scanner)
	return parseValues(stringBetweenBraces)
}

func getStringBetweenBraces(scanner *bufio.Scanner) (allValuesString string) {
	line := scanner.Text()
	allValuesString = ""
	afterOpeningBracket := line[strings.Index(line, "{")+1:]
	for line = afterOpeningBracket; !strings.HasSuffix(line, "}"); line = scanner.Text() {
		allValuesString = allValuesString + line + "\n"
		scanner.Scan()
	}
	beforeClosingBracket := line[:len(line)-1]
	allValuesString = allValuesString + beforeClosingBracket
	return
}

func parseValues(s string) (allValues []string) {
	allValues = []string(nil)
	numberOfQuotes := 0
	lastSplit := -1
	for i, rune := range s {
		if rune == '`' || rune == '"' {
			numberOfQuotes = numberOfQuotes + 1
		} else if numberOfQuotes%2 == 0 && rune == ',' {
			value := strings.TrimSpace(s[lastSplit+1 : i])
			lastSplit = i
			allValues = append(allValues, value)
		}
	}
	if lastSplit+1 < len(s) {
		value := strings.TrimSpace(s[lastSplit+1:])
		allValues = append(allValues, value)
	}
	return
}
