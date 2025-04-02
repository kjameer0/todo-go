package navmenu

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
	"todo.com/ansi"
)

const ctrlD = 4
const ctrlC = 3
const enter = 13
const delete = 127

// TODO: pressing ctrlD should cancel the operation but return to the current menu
// TODO:  ctrl C should end the menu
// NavMenu holds a 2D slice of menu items
type NavMenu[T fmt.Stringer] struct {
	menu            [][]string
	keyLookup       map[string]T
	fd              int
	originalState   *term.State
	numPrintedLines int
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func indexOf(s string, item byte) int {
	itemRune := rune(item)
	for i, v := range s {
		if v == itemRune {
			return i
		}
	}
	return -1 // Return -1 if not found
}
func createLookupKeys(size int) []string {
	keys := []string{}
	alphabetIdx := 0
	curKey := []byte{alphabet[alphabetIdx]}
	for range size {
		keys = append(keys, string(curKey))
		alphabetIdx = (alphabetIdx + 1) % len(alphabet)
		//if we are on a 9
		//roll over everything behind
		i := len(curKey) - 1
		for ; i >= 0; i-- {
			char := curKey[i]
			if char == alphabet[len(alphabet)-1] {
				curKey[i] = alphabet[0]
			} else {
				nextIdx := (indexOf(alphabet, curKey[i]) + 1) % len(alphabet)
				curKey[i] = alphabet[nextIdx]
				break
			}
		}
		if i <= -1 {
			curKey = append(curKey, alphabet[0])
		}
	}
	return keys
}
func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}

func generateRows(items []string, windowWidth int) [][]string {
	curWidth := 0
	m := [][]string{}
	curRow := []string{}
	for _, item := range items {
		if len(item) >= windowWidth-6 {
			truncatedRow := []string{(item[0:windowWidth-6] + "...")}
			if len(curRow) > 0 {
				m = append(m, curRow)
				curRow = []string{}
			}
			m = append(m, truncatedRow)
		} else if len(item)+curWidth >= windowWidth-5 {
			curWidth = 0
			m = append(m, curRow)
			curRow = []string{item}
		} else if len(item)+curWidth < windowWidth {
			curWidth += len(item)
			curRow = append(curRow, item)
		}
	}
	if len(curRow) > 0 {
		m = append(m, curRow)
	}
	return m
}

func NewMenu[T fmt.Stringer](items []T, fd int) *NavMenu[T] {
	keys := createLookupKeys(len(items))
	menuItems := []string{}
	keyLookup := make(map[string]T)
	for idx, item := range items {
		keyLookup[keys[idx]] = item
		menuItem := fmt.Sprintf("%s. %s", keys[idx], item.String())
		menuItems = append(menuItems, menuItem)
	}
	return &NavMenu[T]{keyLookup: keyLookup, menu: generateRows(menuItems, 80), fd: fd}
}

func (m *NavMenu[T]) Render() (T, error) {
	var zeroValue T
	oldState, err := term.MakeRaw(m.fd)
	if err != nil {
		return zeroValue, err
	}
	defer term.Restore(m.fd, oldState)
	m.originalState = oldState
	userInput := []byte{}
	for {
		for _, row := range m.menu {
			rowText := ""
			for _, entry := range row {
				rowText += entry
			}
			fmt.Println(rowText)
			m.numPrintedLines += 1
			fmt.Print("\r")
		}
		fmt.Println("\r")
		fmt.Print("Input: ", string(userInput))
		buf := make([]byte, 3)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return zeroValue, err
		}
		if buf[0] == ctrlC {
			term.Restore(m.fd, m.originalState)
			log.Fatal("interrupted\r")
		} else if buf[0] == ctrlD {
			fmt.Println("cancelled\r")
			return zeroValue, errors.New("user cancelled action")
		} else if buf[0] == 27 {
			clearLines(len(m.menu) + 1)
			continue
		} else if buf[0] == enter {
			fmt.Print(ansi.ClearScreen)
			fmt.Print(ansi.Home)
			fmt.Print("\r")
			return m.keyLookup[string(userInput)], nil
		} else if buf[0] == delete {
			if len(userInput) == 0 {
				clearLines(len(m.menu) + 1)
				continue
			}
			userInput = userInput[0 : len(userInput)-1]
			fmt.Print(ansi.Left)
			fmt.Print(ansi.ClearLine)
		} else {
			userInput = append(userInput, buf[0])
		}
		clearLines(len(m.menu) + 1)
		fmt.Print("\r")
	}
}
