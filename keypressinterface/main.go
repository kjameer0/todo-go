package keypressinterface

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// okay I can render a list of items broken up by line
// what is the next part?
// I need to be able to generate unique ids for each rendered task that are two characters in length

// a menu that can be navigated by moving cursor around when you press keys
// generating a menu should entail submitting
// a list of strings
// number of cols
// number of rows
type MatrixMenu struct {
	matrixData    [][]string
	Items         []string
	Rows          int
	Cols          int
	fd            int
	cursorPos     [2]int
	originalState *term.State
}

const up = "\033[A"
const down = "\033[B"
const right = "\033[C"
const left = "\033[D"

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
func NewMatrixMenu(items []string, cols int, fd int) (*MatrixMenu, error) {
	matrix, err := generateMatrix(cols, items)
	if err != nil {
		return nil, err
	}
	return &MatrixMenu{Cols: cols, matrixData: matrix, fd: fd, cursorPos: [2]int{0, 0}}, nil
}

//	func (m *MatrixMenu) handleControls() error {
//		row := m.cursorPos[0]
//		col := m.cursorPos[1]
//		return nil
//	}
func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}
func (m *MatrixMenu) RenderInterface() error {
	if m.matrixData == nil {
		return errors.New("no data to create menu from")
	}
	oldState, err := term.MakeRaw(m.fd)
	if err != nil {
		return err
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
			return err
		}
		userInput := string(buf)
		if buf[0] == 3 {
			term.Restore(m.fd, m.originalState)
			return errors.New("interrupt triggered")
		}
		clearLines(len(m.matrixData))
		if userInput == up {
			nextCursorPos := m.cursorPos[0] - 1
			if nextCursorPos < 0 {
				nextCursorPos = len(m.matrixData) - 1
			}
			m.cursorPos[0] = nextCursorPos
		} else if userInput == down {
			nextCursorPos := m.cursorPos[0]
			nextCursorPos = (nextCursorPos + 1) % len(m.matrixData)
			m.cursorPos[0] = nextCursorPos
		} else if userInput == right {
			nextCursorPos := m.cursorPos[1]
			nextCursorPos = (nextCursorPos + 1) % m.Cols
			// skip to front of row if we see a blank
			if m.matrixData[m.cursorPos[0]][nextCursorPos] == "" {
				m.cursorPos[1] = 0
			} else {
				m.cursorPos[1] = nextCursorPos
			}
		} else if userInput == left {
			nextCursorPos := m.cursorPos[1] - 1
			if nextCursorPos < 0 {
				nextCursorPos = len(m.matrixData[0]) - 1
				for m.matrixData[m.cursorPos[0]][nextCursorPos] == "" {
					nextCursorPos--
				}
			}
			m.cursorPos[1] = nextCursorPos
		}
		fmt.Print("\r")
	}
	return nil
}

func RenderInterface() {
	fmt.Println("aa: Send morning greeting", "bb. Review pull requests", "cc. Fix database migration issue", "dd. Update API documentation", "ee. Attend stand-up meeting", "ff. Refactor authentication logic")
	fmt.Println("aa: Check emails", "bb. Implement caching strategy", "cc. Debug failing test cases", "dd. Optimize SQL queries", "ee. Deploy staging environment", "ff. Write unit tests")
	fmt.Println("aa: Plan sprint tasks", "bb. Code review feedback", "cc. Investigate performance bottleneck", "dd. Push feature branch", "ee. Write end-to-end tests", "ff. Prepare for demo meeting")

}
