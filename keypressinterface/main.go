package keypressinterface

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// i need the ability to have pressing enter on an entry return the underlying data(not the string version) to the user
type StringerSignaler[T any] interface {
	String() string
	Signal() T
}

type MatrixMenu[T fmt.Stringer] struct {
	matrixData           [][]string
	underlyingDataMatrix [][]T
	Items                []T
	fd                   int
	cursorPos            [2]int
	originalState        *term.State
}

const up = "\033[A"
const down = "\033[B"
const right = "\033[C"
const left = "\033[D"
const TERMINAL_WINDOW_WIDTH = 80
const enter = 13

func generateRows(items []string, windowWidth int) [][]string {
	curWidth := 0
	m := [][]string{}
	curRow := []string{}
	for _, item := range items {
		if len(item) >= windowWidth-5 {
			truncatedRow := []string{(item[0:windowWidth-5] + "...")}
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

// TODO write a function to generate lines of 80 char width
// a task can be length l < 80 or l >= 80
//write a function that
// func (m *MatrixMenu) generateMatrix(n int) [][]int {

// }

func generateMatrix(cols int, items []string) ([][]string, error) {
	matrix := make([][]string, 0)
	itemIdx := 0
	for itemIdx < len(items) {
		if itemIdx == len(items) {
			return matrix, nil
		}
		itemList := make([]string, cols)
		matrix = append(matrix, itemList)
		for col := 0; col < cols; col++ {
			if itemIdx == len(items) {
				itemList[col] = ""
				continue
			}
			itemList[col] = items[itemIdx]
			itemIdx++
		}
	}
	return matrix, nil
}
func NewMatrixMenu[T fmt.Stringer](items []T, fd int) (*MatrixMenu[T], error) {
	stringItems := []string{}
	for _, item := range items {
		stringItems = append(stringItems, item.String())
	}

	matrix := generateRows(stringItems, TERMINAL_WINDOW_WIDTH)
	//make the data matrix match the string matrix
	underlyingDataMatrix := [][]T{}
	itemsIdx := 0
	for _, row := range matrix {
		// iterate cols
		curRow := []T{}
		for i := 0; i < len(row); i++ {
			curRow = append(curRow, items[itemsIdx])
			itemsIdx++
		}
		underlyingDataMatrix = append(underlyingDataMatrix, curRow)
	}
	return &MatrixMenu[T]{matrixData: matrix, underlyingDataMatrix: underlyingDataMatrix, fd: fd, cursorPos: [2]int{0, 0}}, nil
}

func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}
func (m *MatrixMenu[T]) RenderInterface() (T, error) {
	var zeroValue T
	if m.matrixData == nil {
		return zeroValue, errors.New("no data to create menu from")
	}
	oldState, err := term.MakeRaw(m.fd)
	if err != nil {
		return zeroValue, err
	}
	defer term.Restore(m.fd, oldState)
	m.originalState = oldState

	for {
		for rowIdx, row := range m.matrixData {
			rowText := ""
			for colIdx, val := range row {
				if rowIdx == m.cursorPos[0] && colIdx == m.cursorPos[1] {
					rowText += " >"
				} else {
					rowText += "  "
				}
				rowText += val
			}
			fmt.Println(rowText)
			fmt.Print("\r")
		}
		fmt.Print("\r")

		buf := make([]byte, 3)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return zeroValue, err
		}
		curMatrixRow := m.cursorPos[0]
		curMatrixCol := m.cursorPos[1]
		userInput := string(buf)
		if buf[0] == 3 {
			term.Restore(m.fd, m.originalState)
			return zeroValue, errors.New("interrupt triggered")
		}
		clearLines(len(m.matrixData))
		if userInput == up {
			nextCursorPos := curMatrixRow - 1
			if nextCursorPos < 0 {
				nextCursorPos = len(m.matrixData) - 1
			}
			m.cursorPos[0] = nextCursorPos
		} else if userInput == down {
			nextCursorPos := curMatrixRow
			nextCursorPos = (nextCursorPos + 1) % len(m.matrixData)
			if curMatrixCol >= len(m.matrixData[nextCursorPos]) {
				m.cursorPos[1] = len(m.matrixData[nextCursorPos]) - 1
			}
			m.cursorPos[0] = nextCursorPos
		} else if userInput == right {
			nextCursorPos := curMatrixCol
			nextCursorPos = (nextCursorPos + 1) % len(m.matrixData[curMatrixRow])
			m.cursorPos[1] = nextCursorPos
		} else if userInput == left {
			nextCursorPos := curMatrixCol - 1
			if nextCursorPos < 0 {
				nextCursorPos = len(m.matrixData[curMatrixRow]) - 1
			}
			m.cursorPos[1] = nextCursorPos
		} else if buf[0] == enter {
			return m.underlyingDataMatrix[curMatrixRow][curMatrixCol], nil
		}
		fmt.Print("\r")
	}
}

func RenderInterface() {
	fmt.Println("aa: Send morning greeting", "bb. Review pull requests", "cc. Fix database migration issue", "dd. Update API documentation", "ee. Attend stand-up meeting", "ff. Refactor authentication logic")
	fmt.Println("aa: Check emails", "bb. Implement caching strategy", "cc. Debug failing test cases", "dd. Optimize SQL queries", "ee. Deploy staging environment", "ff. Write unit tests")
	fmt.Println("aa: Plan sprint tasks", "bb. Code review feedback", "cc. Investigate performance bottleneck", "dd. Push feature branch", "ee. Write end-to-end tests", "ff. Prepare for demo meeting")

}
