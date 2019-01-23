package main

import (
	"bytes"
	"image"
	"log"

	"github.com/MrGru/ChineseChess/res"
	"github.com/hajimehoshi/ebiten"
)

const (
	NON = iota
	RED_CHARIOT
	RED_ADVISOR
	RED_ELEPHANT
	RED_CANNON
	RED_HORSE
	RED_SOLDIER
	RED_GENERAL

	BLACK_CHARIOT
	BLACK_ADVISOR
	BLACK_ELEPHANT
	BLACK_CANNON
	BLACK_HORSE
	BLACK_SOLDIER
	BLACK_GENERAL
)

type Piece struct {
	id  int
	img *ebiten.Image
	row float64
	col float64
}

func (p *Piece) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	w, h := p.img.Size()
	opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	posX, posY := getPosition(p.row, p.col)
	opts.GeoM.Translate(posX, posY)
	screen.DrawImage(p.img, opts)
}

func NewPiece(id int, row, col float64) Piece {
	piece := Piece{}
	piece.id = id
	piece.row = row
	piece.col = col
	var im image.Image
	var err error
	switch id {
	case RED_ADVISOR:
		im, _, err = image.Decode(bytes.NewReader(res.Red_advisor_png))
	case RED_CANNON:
		im, _, err = image.Decode(bytes.NewReader(res.Red_canon_png))
	case RED_CHARIOT:
		im, _, err = image.Decode(bytes.NewReader(res.Red_chariot_png))
	case RED_ELEPHANT:
		im, _, err = image.Decode(bytes.NewReader(res.Red_elephant_png))
	case RED_GENERAL:
		im, _, err = image.Decode(bytes.NewReader(res.Red_general_png))
	case RED_HORSE:
		im, _, err = image.Decode(bytes.NewReader(res.Red_horse_png))
	case RED_SOLDIER:
		im, _, err = image.Decode(bytes.NewReader(res.Red_soldier_png))

	case BLACK_ADVISOR:
		im, _, err = image.Decode(bytes.NewReader(res.Black_advisor_png))
	case BLACK_CANNON:
		im, _, err = image.Decode(bytes.NewReader(res.Black_cannon_png))
	case BLACK_CHARIOT:
		im, _, err = image.Decode(bytes.NewReader(res.Black_chariot_png))
	case BLACK_ELEPHANT:
		im, _, err = image.Decode(bytes.NewReader(res.Black_elephant_png))
	case BLACK_GENERAL:
		im, _, err = image.Decode(bytes.NewReader(res.Black_general_png))
	case BLACK_HORSE:
		im, _, err = image.Decode(bytes.NewReader(res.Black_horse_png))
	case BLACK_SOLDIER:
		im, _, err = image.Decode(bytes.NewReader(res.Black_soldier_png))
	}
	if err != nil {
		log.Fatal(err)
	}
	piece.img, _ = ebiten.NewImageFromImage(im, ebiten.FilterDefault)
	return piece
}
