package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Move struct {
	Color int
	Col   int
}

type Board struct {
	BoardString string
}

type ConnectGame int

// Board size. Connect 4 is normally 6 rows by 7 columns.
const rows = 6
const cols = 7

// -1 means the cell is empty.
// 0 means a white piece, 1 means a black piece.
const empty = -1

// The board is a global 2D slice: gameBoard[row][col].
// Row 0 is the TOP, row 5 is the BOTTOM.
var gameBoard [][]int

// Remembers the color that moved last.
// -1 means no one has moved yet (so the first move is always allowed).
var lastColorMoved = -1

func (t *ConnectGame) Move(args *Move, reply *int) error {
	log.Println("Received Move -> Color:", args.Color, "Col:", args.Col)

	// Turn-order rule: the same color cannot move twice in a row.
	// If this color is the same one that just moved, reject the move.
	if args.Color == lastColorMoved {
		return errors.New("Turn Order Violation")
	}

	// Drop the piece into the column.
	// Start at the bottom row and go up until we find an empty cell.
	for r := rows - 1; r >= 0; r-- {
		if gameBoard[r][args.Col] == empty {
			gameBoard[r][args.Col] = args.Color
			break // the piece landed, so stop
		}
	}

	// Remember who just moved, so the next move must be the other color.
	lastColorMoved = args.Color

	return nil
}

func (t *ConnectGame) Get(args *int, reply *Board) error {
	// Build one big string, row by row.
	text := ""
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if gameBoard[r][c] == 0 {
				text = text + "W "
			} else if gameBoard[r][c] == 1 {
				text = text + "B "
			} else {
				text = text + ". "
			}
		}
		text = text + "\n"
	}

	reply.BoardString = text
	return nil
}

func main() {
	// Make the outer slice (the rows).
	gameBoard = make([][]int, rows)
	for r := 0; r < rows; r++ {
		// Make each inner row (the columns).
		gameBoard[r] = make([]int, cols)
		// Start every cell as empty.
		for c := 0; c < cols; c++ {
			gameBoard[r][c] = empty
		}
	}

	cg := new(ConnectGame)
	rpc.Register(cg)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Println("Serving on PORT 1234")
	http.Serve(l, nil)
}
