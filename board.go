package main

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/MrGru/ChineseChess/res"

	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_COL    = 9
	MAX_ROW    = 10
	RED_TURN   = 0
	BLACK_TURN = 1
)

type Board struct {
	pieces         [][]*Piece
	turn           int
	selectedPiece  *Piece
	isCheckedPiece bool
	startPosX      float64
	startPosY      float64
	cellSize       float64
	selectedRow    int
	selectedCol    int
	moveImage      [][]*ebiten.Image
}

var initBoard = [][]int{
	{BLACK_CHARIOT, BLACK_HORSE, BLACK_ELEPHANT, BLACK_ADVISOR,
		BLACK_GENERAL, BLACK_ADVISOR, BLACK_ELEPHANT,
		BLACK_HORSE, BLACK_CHARIOT},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, BLACK_CANNON, 0, 0, 0, 0, 0, BLACK_CANNON, 0},
	{BLACK_SOLDIER, 0, BLACK_SOLDIER, 0, BLACK_SOLDIER, 0,
		BLACK_SOLDIER, 0, BLACK_SOLDIER},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{RED_SOLDIER, 0, RED_SOLDIER, 0, RED_SOLDIER, 0, RED_SOLDIER,
		0, RED_SOLDIER},
	{0, RED_CANNON, 0, 0, 0, 0, 0, RED_CANNON, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{RED_CHARIOT, RED_HORSE, RED_ELEPHANT, RED_ADVISOR,
		RED_GENERAL, RED_ADVISOR, RED_ELEPHANT, RED_HORSE,
		RED_CHARIOT},
}

var moveList = [][]int{{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0}}

const (
	MOVE = 1
	EAT  = 2
)

func (b *Board) InitPosition(screenWidth, screenHeight, cellSize float64) {
	b.startPosX = screenWidth/2 - 4*cellSize
	b.startPosY = screenHeight/2 + (0-4)*float64(cellSize) - cellSize/2
	b.cellSize = cellSize
}

func (b *Board) CreatePiece() {
	b.pieces = make([][]*Piece, MAX_ROW)
	for i := range b.pieces {
		b.pieces[i] = make([]*Piece, MAX_COL)
	}

	fmt.Println("len: ", len(initBoard))
	for i := 0; i < len(initBoard); i++ {
		for j := 0; j < len(initBoard[i]); j++ {
			if initBoard[i][j] != 0 {
				p := NewPiece(initBoard[i][j], float64(i), float64(j))
				b.pieces[i][j] = &p
			}
		}
	}
	moveImg, _, err := image.Decode(bytes.NewReader(res.Move_png))
	if err != nil {
		log.Fatal(err)
	}

	b.moveImage = make([][]*ebiten.Image, MAX_ROW)
	for i := range b.moveImage {
		b.moveImage[i] = make([]*ebiten.Image, MAX_COL)
	}
	for i := 0; i < len(b.moveImage); i++ {
		for j := 0; j < len(b.moveImage[i]); j++ {
			b.moveImage[i][j], _ = ebiten.NewImageFromImage(moveImg, ebiten.FilterDefault)
		}
	}
}

func (b *Board) HandleTouchEvent(x, y int) {
	row, col := b.getRowColFromPosition(x, y)
	if row >= 0 && col >= 0 && row < 10 && col < 9 {
		if b.pieces[row][col] != nil {
			b.selectedCol = col
			b.selectedRow = row
		}
	}
}

func (b *Board) HandleReleaseEvent(x, y int) {
	row, col := b.getRowColFromPosition(x, y)
	if  row >= 0 && col >= 0 && row < 10 && col < 9 && b.isCheckedPiece && moveList[row][col] > 0 {
		b.ClearMoveList()
		oldRow, oldCol := int(b.selectedPiece.row), int(b.selectedPiece.col)
		b.pieces[oldRow][oldCol] = nil

		b.selectedPiece.col = float64(col)
		b.selectedPiece.row = float64(row)
		b.pieces[row][col] = b.selectedPiece

		b.isCheckedPiece = false
		b.selectedPiece = nil

	} else {
		b.ClearMoveList()
		if row >= 0 && col >= 0 && row < 10 && col < 9 {
			if b.pieces[row][col] != nil && b.selectedRow == b.selectedRow && b.selectedCol == b.selectedCol {
				b.selectedPiece = b.pieces[row][col]
				b.checkListMove()
				b.isCheckedPiece = true
			} else {
				b.selectedPiece = nil
				b.isCheckedPiece = false
			}

		} else {
			b.selectedPiece = nil
			b.isCheckedPiece = false
		}
	}
}
func (b *Board) ClearMoveList() {
	for i := 0; i < len(moveList); i++ {
		for j := 0; j < len(moveList[i]); j++ {
			moveList[i][j] = 0
		}
	}
}

func (b *Board) checkListMove() {
	row := int(b.selectedPiece.row)
	col := int(b.selectedPiece.col)
	switch b.selectedPiece.id {
	case BLACK_SOLDIER:
		if row < MAX_ROW-1 {
			if b.pieces[row+1][col] == nil {
				moveList[row+1][col] = MOVE
			} else if b.pieces[row+1][col].id > 20 {
				moveList[row+1][col] = EAT
			}
		}
		if row > 4 {
			if col > 0 {
				if b.pieces[row][col-1] == nil {
					moveList[row][col-1] = MOVE
				}
				if b.pieces[row][col-1].id > 20 {
					moveList[row][col-1] = EAT
				}
			}
			if col < 8 {
				if b.pieces[row][col+1] == nil {
					moveList[row][col+1] = MOVE
				} else if b.pieces[row][col+1].id > 20 {
					moveList[row][col+1] = EAT
				}
			}
		}
	case BLACK_CANNON:
		if row+1 < MAX_ROW {
			for i := row + 1; i < MAX_ROW; i++ {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if i+1 < MAX_ROW {
						for j := i + 1; j < MAX_ROW; j++ {
							fmt.Println("Check row: ", j, col)
							if b.pieces[j][col] != nil && b.pieces[j][col].id > 20 {
								moveList[j][col] = EAT
								break
							}
						}
					}
					break
				}
			}
		}
		if row-1 >= 0 {
			for i := row - 1; i >= 0; i-- {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if i-1 > 0 {
						for j := i - 1; j >= 0; j-- {
							if b.pieces[j][col] != nil && b.pieces[j][col].id > 20 {
								moveList[j][col] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
		if col+1 < MAX_COL {
			for i := col + 1; i < MAX_COL; i++ {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if i+1 < MAX_COL {
						for j := i + 1; j < MAX_COL; j++ {
							if b.pieces[row][j] != nil && b.pieces[row][j].id > 20 {
								moveList[row][j] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
		if col-1 >= 0 {
			for i := col - 1; i >= 0; i-- {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if i-1 > 0 {
						for j := i - 1; j >= 0; j-- {
							if b.pieces[row][j] != nil && b.pieces[row][j].id > 20 {
								moveList[row][j] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
	case BLACK_CHARIOT:
		if row+1 < MAX_ROW {
			for i := row + 1; i < MAX_ROW; i++ {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if b.pieces[i][col] != nil && b.pieces[i][col].id > 20 {
						moveList[i][col] = EAT
					}
					break
				}
			}
		}
		if row-1 >= 0 {
			for i := row - 1; i >= 0; i-- {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if b.pieces[i][col] != nil && b.pieces[i][col].id > 20 {
						moveList[i][col] = EAT
					}
					break
				}
			}
		}
		if col+1 < MAX_COL {
			for i := col + 1; i < MAX_COL; i++ {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if b.pieces[row][i] != nil && b.pieces[row][i].id > 20 {
						moveList[row][i] = EAT
					}
					break
				}
			}
		}
		if col-1 >= 0 {
			for i := col - 1; i >= 0; i-- {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if b.pieces[row][i] != nil && b.pieces[row][i].id > 20 {
						moveList[row][i] = EAT
					}
					break
				}
			}
		}
	case BLACK_HORSE:
		nextRow, nextCol := row-2, col-1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-2, col+1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col+1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col-1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col-2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col-2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col+2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col+2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case BLACK_ELEPHANT:
		nextRow, nextCol := row-2, col-2
		if nextRow >= 0 && nextRow <= 4 && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col+2
		if nextRow >= 0 && nextRow <= 4 && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-2, col+2
		if nextRow >= 0 && nextRow <= 4 && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col-2
		if nextRow >= 0 && nextRow <= 4 && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case BLACK_ADVISOR:
		nextRow, nextCol := row-1, col-1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col+1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col-1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col+1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case BLACK_GENERAL:
		nextRow, nextCol := row-1, col
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row, col-1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row, col+1
		if nextRow >= 0 && nextRow <= 2 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id > 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}

	case RED_SOLDIER:
		if row > 0 {
			if b.pieces[row-1][col] == nil {
				moveList[row-1][col] = MOVE
			} else if b.pieces[row-1][col].id < 20 {
				moveList[row-1][col] = EAT
			}
		}
		if row <= 4 {
			if col > 0 {
				if b.pieces[row][col-1] == nil {
					moveList[row][col-1] = MOVE
				} else if b.pieces[row][col-1].id < 20 {
					moveList[row][col-1] = EAT
				}
			}
			if col < 8 {
				if b.pieces[row][col+1] == nil {
					moveList[row][col+1] = MOVE
				} else if b.pieces[row][col+1].id < 20 {
					moveList[row][col+1] = EAT
				}
			}
		}
	case RED_CANNON:
		if row+1 < MAX_ROW {
			for i := row + 1; i < MAX_ROW; i++ {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if i+1 < MAX_ROW {
						for j := i + 1; j < MAX_ROW; j++ {
							if b.pieces[j][col] != nil && b.pieces[j][col].id < 20 {
								moveList[j][col] = EAT
								break
							}
						}
					}
					break
				}
			}
		}
		if row-1 >= 0 {
			for i := row - 1; i >= 0; i-- {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if i-1 > 0 {
						for j := i - 1; j >= 0; j-- {
							if b.pieces[j][col] != nil && b.pieces[j][col].id < 20 {
								moveList[j][col] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
		if col+1 < MAX_COL {
			for i := col + 1; i < MAX_COL; i++ {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if i+1 < MAX_COL {
						for j := i + 1; j < MAX_COL; j++ {
							if b.pieces[row][j] != nil && b.pieces[row][j].id < 20 {
								moveList[row][j] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
		if col-1 >= 0 {
			for i := col - 1; i >= 0; i-- {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if i-1 > 0 {
						for j := i - 1; j >= 0; j-- {
							if b.pieces[row][j] != nil && b.pieces[row][j].id < 20 {
								moveList[row][j] = EAT
								break
							}

						}
					}
					break
				}
			}
		}
	case RED_CHARIOT:
		if row+1 < MAX_ROW {
			for i := row + 1; i < MAX_ROW; i++ {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if b.pieces[i][col] != nil && b.pieces[i][col].id < 20 {
						moveList[i][col] = EAT
					}
					break
				}
			}
		}
		if row-1 >= 0 {
			for i := row - 1; i >= 0; i-- {
				if b.pieces[i][col] == nil {
					moveList[i][col] = MOVE
				} else {
					if b.pieces[i][col] != nil && b.pieces[i][col].id < 20 {
						moveList[i][col] = EAT
					}
					break
				}
			}
		}
		if col+1 < MAX_COL {
			for i := col + 1; i < MAX_COL; i++ {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if b.pieces[row][i] != nil && b.pieces[row][i].id < 20 {
						moveList[row][i] = EAT
					}
					break
				}
			}
		}
		if col-1 >= 0 {
			for i := col - 1; i >= 0; i-- {
				if b.pieces[row][i] == nil {
					moveList[row][i] = MOVE
				} else {
					if b.pieces[row][i] != nil && b.pieces[row][i].id < 20 {
						moveList[row][i] = EAT
					}
					break
				}
			}
		}
	case RED_HORSE:
		nextRow, nextCol := row-2, col-1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-2, col+1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col+1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col-1
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col-2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col-2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col+2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col+2
		if nextRow >= 0 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case RED_ELEPHANT:
		nextRow, nextCol := row-2, col-2
		if nextRow >= 5 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col+2
		if nextRow >= 5 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-2, col+2
		if nextRow >= 5 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row-1][col+1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+2, col-2
		if nextRow >= 5 && nextRow < MAX_ROW && nextCol >= 0 && nextCol < MAX_COL && b.pieces[row+1][col-1] == nil {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case RED_ADVISOR:
		nextRow, nextCol := row-1, col-1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col+1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col-1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row-1, col+1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
	case RED_GENERAL:
		nextRow, nextCol := row-1, col
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row+1, col
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row, col-1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}
		nextRow, nextCol = row, col+1
		if nextRow >= 7 && nextRow <= 9 && nextCol >= 3 && nextCol <= 5 {
			if b.pieces[nextRow][nextCol] == nil {
				moveList[nextRow][nextCol] = MOVE
			} else if b.pieces[nextRow][nextCol].id < 20 {
				moveList[nextRow][nextCol] = EAT
			}
		}

	}
}

func (b *Board) getRowColFromPosition(x, y int) (int, int) {
	if b.startPosX-b.cellSize/2 > float64(x) || b.startPosY-b.cellSize/2 > float64(y) {
		return -1, -1
	}
	return int((float64(y) - b.startPosY + b.cellSize/2) / cellSize), int((float64(x) - b.startPosX + b.cellSize/2) / cellSize)
}

func (b *Board) Draw(screen *ebiten.Image) {
	for i := 0; i < len(b.pieces); i++ {
		for j := 0; j < len(b.pieces[i]); j++ {
			if b.pieces[i][j] != nil {
				b.pieces[i][j].Draw(screen)
			}
		}
	}

	// fmt.Println(moveList)
	for i := 0; i < len(moveList); i++ {
		for j := 0; j < len(moveList[i]); j++ {
			if moveList[i][j] > 0 {
				// fmt.Println("Draw: ", i, j)
				opts := &ebiten.DrawImageOptions{}
				w, h := b.moveImage[i][j].Size()
				opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
				posX, posY := getPosition(float64(i), float64(j))
				opts.GeoM.Translate(posX, posY)
				screen.DrawImage(b.moveImage[i][j], opts)
			}
		}
	}
}
