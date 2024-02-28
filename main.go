package main

import (
	"bufio"
	"errors"
	"fmt"
	rom "github.com/brandenc40/romannumeral"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Calculator struct {
	numbers [2]int
	action  string
}

func main() {
	r := bufio.NewReader(os.Stdin)

	for {
		str, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		res, err := process(str)
		if err != nil {
			panic(err)
		}

		fmt.Println(res)
	}
}

func process(input string) (string, error) {
	var calculator Calculator

	var numbers [2]int
	var romans [2]bool

	for action := range calculator.getActions() {
		if strings.Contains(input, action) {
			items := strings.Split(input, action)
			if len(items) != 2 {
				return "", errors.New("failed to parse input parameter")
			}

			for i, item := range items {
				var err error
				numbers[i], romans[i], err = prepareParam(item)
				if err != nil {
					return "", err
				}
			}

			calculator.numbers = numbers
			calculator.action = action

			break
		}
	}

	_, err := validate(numbers, romans)
	if err != nil {
		return "", err
	}

	intResult := calculator.Calculate()

	if romans[0] && romans[1] {
		return rom.IntToString(intResult)
	}

	return strconv.Itoa(intResult), nil
}

func prepareParam(param string) (int, bool, error) {
	param = strings.TrimSpace(param)

	matched, err := isRomanNumber(param)
	if err != nil {
		return 0, false, err
	}

	number, err := toInteger(param, matched)
	if err != nil {
		return 0, false, err
	}

	return number, matched, nil
}

func isRomanNumber(str string) (bool, error) {
	return regexp.MatchString("^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$", str)
}

func toInteger(str string, isRomanNumber bool) (int, error) {
	if isRomanNumber {
		num, err := rom.StringToInt(str)
		if err != nil {
			return num, err
		}

		return num, nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return num, err
	}

	return num, nil
}

func validate(numbers [2]int, romans [2]bool) (bool, error) {
	if (romans[0] == true && romans[1] == false) || (romans[0] == false && romans[1] == true) {
		return false, errors.New("input numbers have different format")
	}

	if numbers[0] > 10 || numbers[1] > 10 {
		return false, errors.New("input numbers must be from 0 to 10 inclusive")
	}

	return true, nil
}

func (c *Calculator) Calculate() int {
	return c.getCallback()(c.numbers)
}

func (c *Calculator) getActions() map[string]func([2]int) int {
	return map[string]func([2]int) int{
		"+": func(nums [2]int) int { return nums[0] + nums[1] },
		"-": func(nums [2]int) int { return nums[0] - nums[1] },
		"*": func(nums [2]int) int { return nums[0] * nums[1] },
		"/": func(nums [2]int) int { return nums[0] / nums[1] },
	}
}

func (c *Calculator) getCallback() func([2]int) int {
	return c.getActions()[c.action]
}
