package pkg

import (
	"os/exec"

	"streaming/internal/shared"
)

func EncodeH264(source, ext string) error {
	// ffmpeg -i input.flv -vcodec libx264 -acodec aac output.mp4
	cmd := exec.Command("ffmpeg", "-i", source+"_"+ext,
		"-vcodec", "libx264", "-acodec", "aac", source+shared.PreprocessedFileExt)
	return cmd.Run()
}
