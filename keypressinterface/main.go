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
func generateMatrix(rows int, cols int, items []string) ([][]string, error) {
	matrix := make([][]string, 0)
	if len(items) > rows*cols {
		return nil, errors.New("too many items for provided rows and columns")
	}
	itemIdx := 0
	for row := 0; row < rows; row++ {
		if itemIdx == len(items) {
			return matrix, nil
		}
		itemList := make([]string, cols)
		matrix = append(matrix, itemList)
		for col := 0; col < cols; col++ {
			if itemIdx == len(items) {
				return matrix, nil
			}
			matrix[row][col] = items[itemIdx]
			itemIdx++
		}
	}
	return matrix, nil
}
func NewMatrixMenu(items []string, rows int, cols int, fd int) (*MatrixMenu, error) {
	matrix, err := generateMatrix(rows, cols, items)
	if err != nil {
		return nil, err
	}
	return &MatrixMenu{Rows: rows, Cols: cols, matrixData: matrix, fd: fd, cursorPos: [2]int{0, 0}}, nil
}
// func (m *MatrixMenu) handleControls() error {
// 	row := m.cursorPos[0]
// 	col := m.cursorPos[1]
// 	return nil
// }

func (m *MatrixMenu) RenderInterface() error {
	if m.matrixData == nil {
		return errors.New("no data to create menu from")
	}
	oldState, err := term.MakeRaw(m.fd)
	if err != nil {
		return err
	}
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
			fmt.Println(rowText)
		}
	}
	return nil
}

func RenderInterface() {
	fmt.Println("aa: Send morning greeting", "bb. Review pull requests", "cc. Fix database migration issue", "dd. Update API documentation", "ee. Attend stand-up meeting", "ff. Refactor authentication logic")
	fmt.Println("aa: Check emails", "bb. Implement caching strategy", "cc. Debug failing test cases", "dd. Optimize SQL queries", "ee. Deploy staging environment", "ff. Write unit tests")
	fmt.Println("aa: Plan sprint tasks", "bb. Code review feedback", "cc. Investigate performance bottleneck", "dd. Push feature branch", "ee. Write end-to-end tests", "ff. Prepare for demo meeting")

}
