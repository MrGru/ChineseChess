package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

const (
	MAX_COL = 9
	MAX_ROW = 10
)

type Board struct {
	pieces [][]*Piece
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

func (b *Board) CreatePiece() {
	b.pieces = make([][]*Piece, MAX_COL)
	for i := range b.pieces {
		b.pieces[i] = make([]*Piece, MAX_ROW)
	}

	fmt.Println("len: ", len(initBoard))
	for i := 0; i < len(initBoard); i++ {
		for j := 0; j < len(initBoard[i]); j++ {
			if initBoard[i][j] != 0 {
				p := NewPiece(initBoard[i][j], float64(i), float64(j))
				b.pieces[j][i] = &p
			}
		}
	}

}

func (b *Board) Draw(screen *ebiten.Image) {
	for i := 0; i < len(b.pieces); i++ {
		for j := 0; j < len(b.pieces[i]); j++ {
			if b.pieces[i][j] != nil {
				b.pieces[i][j].Draw(screen)
			}
		}
	}
}
