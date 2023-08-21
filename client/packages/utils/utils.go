package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	c "ticoma/client/packages/constants"
	intf "ticoma/client/packages/interfaces"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Get current hash as string
func GetCommitHash() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

// Gen random number in range - (inclusive, exclusive)
func RandRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

func GetScreenInfo() *intf.ScreenInfo {
	w := rl.GetScreenWidth()
	h := rl.GetScreenHeight()
	return &intf.ScreenInfo{
		Width:       w,
		Height:      h,
		RefreshRate: 60, // Tmp
	}
}

// Get a texture with specific Id from a spritesheet.
// Spritesheet must be a .png file
func GetTextureFromId(textureId int, spritePath string) rl.Texture2D {

	spriteFile, err := os.Open(spritePath)
	if err != nil {
		fmt.Println("[UTILS] - Couldn't open spritesheet. Err: ", err)
	}

	spriteImg, err := png.Decode(spriteFile)
	if err != nil {
		fmt.Println("[UTILS] - Couldn't decode spritesheet. Err: ", err)
	}

	texturePosX := textureId * c.BLOCK_SIZE
	texturePosY := (textureId / 8) * c.BLOCK_SIZE

	img := spriteImg.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(texturePosX, texturePosY, texturePosX+c.BLOCK_SIZE, texturePosY+c.BLOCK_SIZE))

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		fmt.Println("[UTILS] - Couldn't encode image. Err: ", err)
	}

	txtBytes := buf.Bytes()

	textureImg := rl.NewImage(txtBytes, c.BLOCK_SIZE, c.BLOCK_SIZE, 0, rl.UncompressedR32g32b32a32)
	texture2D := rl.LoadTextureFromImage(textureImg)

	return texture2D

}
