package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/MrGru/ChineseChess/res"
	"github.com/hajimehoshi/ebiten"
)

var square *ebiten.Image

const (
	screenWidth  = 375
	screenHeight = 667

	cellSize = float64(40)
)

var (
	backgroundFrame *ebiten.Image
	colorTable      = color.NRGBA{0x32, 0x32, 0x32, 0xff}
	board           *Board
)

func update(screen *ebiten.Image) error {

	Input()
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, 0)
	screen.DrawImage(backgroundFrame, opts)
	DrawTable(screen)
	board.Draw(screen)
	return nil
}

func Input() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		board.HandleReleaseEvent(ebiten.CursorPosition())
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		board.HandleTouchEvent(ebiten.CursorPosition())
	}
}

func DrawTable(screen *ebiten.Image) {
	for i := 0; i < 10; i++ {
		ebitenutil.DrawLine(screen, screenWidth/2-4*cellSize, screenHeight/2+(float64(i)-4)*cellSize-cellSize/2, screenWidth/2+4*cellSize, screenHeight/2+(float64(i)-4)*cellSize-cellSize/2, colorTable)
	}
	for i := 0; i < 9; i++ {
		if i > 0 && i < 8 {
			ebitenutil.DrawLine(screen, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(0-4)*cellSize-cellSize/2, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(4-4)*cellSize-cellSize/2, colorTable)
			ebitenutil.DrawLine(screen, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(5-4)*cellSize-cellSize/2, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(9-4)*cellSize-cellSize/2, colorTable)
		} else {
			ebitenutil.DrawLine(screen, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(0-4)*cellSize-cellSize/2, screenWidth/2+(float64(i)-4)*cellSize, screenHeight/2+(9-4)*cellSize-cellSize/2, colorTable)
		}
	}
	x1, y1 := getPosition(0, 3)
	x2, y2 := getPosition(2, 5)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, colorTable)
	x1, y1 = getPosition(0, 5)
	x2, y2 = getPosition(2, 3)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, colorTable)

	x1, y1 = getPosition(9, 3)
	x2, y2 = getPosition(7, 5)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, colorTable)
	x1, y1 = getPosition(9, 5)
	x2, y2 = getPosition(7, 3)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, colorTable)
}

func getPosition(row, col float64) (float64, float64) {
	startPosX := screenWidth/2 - 4*cellSize
	startPosY := screenHeight/2 + (0-4)*float64(cellSize) - cellSize/2
	x := startPosX + col*cellSize
	y := startPosY + row*cellSize
	return x, y
}
func main() {

	bg, _, err := image.Decode(bytes.NewReader(res.Background_png))
	if err != nil {
		log.Fatal(err)
	}
	backgroundFrame, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
	board = &Board{}
	board.InitPosition(screenWidth, screenHeight, cellSize)
	board.CreatePiece()

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Chinese Chess"); err != nil {
		panic(err)
	}
}
