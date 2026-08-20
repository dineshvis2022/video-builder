package renderer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"video-builder/models"
)

type FFmpegRenderer struct {
	FFmpegPath string
	OutputDir  string
}

func NewFFmpegRenderer(
	ffmpegPath string,
	outputDir string,
) *FFmpegRenderer {

	return &FFmpegRenderer{
		FFmpegPath: ffmpegPath,
		OutputDir:  outputDir,
	}
}

func (r *FFmpegRenderer) Render(
	invitation models.Invitation,
) (string, error) {

	// 1. Build filter graph
	filterBuilder := NewFilterBuilder()

	filterComplex, imageCount, err :=
		filterBuilder.Build(invitation)

	if err != nil {
		return "", err
	}

	// 2. Create output directory
	if err := os.MkdirAll(
		r.OutputDir,
		0755,
	); err != nil {
		return "", err
	}

	// 3. Calculate total duration
	totalDuration := calculateTotalDuration(
		invitation,
	)

	// 4. Create output file
	outputFile := fmt.Sprintf(
		"%s/%s.mp4",
		r.OutputDir,
		invitation.ProjectID,
	)

	// 5. Build FFmpeg arguments
	args := []string{
		"-y",

		"-i",
		invitation.Background.URL,
	}

	// 6. Add media inputs
	addMediaInputs(
		&args,
		invitation,
	)

	// 7. Add filter complex
	args = append(
		args,
		"-filter_complex",
		filterComplex,
	)

	// 8. Map rendered video
	args = append(
		args,
		"-map",
		"[vout]",
	)

	// 9. Map background audio if available
	args = append(
		args,
		"-map",
		"0:a?",
	)

	// 10. Video codec
	args = append(
		args,
		"-c:v",
		"libx264",
	)

	// 11. Encoding preset
	args = append(
		args,
		"-preset",
		"veryfast",
	)

	// 12. Video quality
	args = append(
		args,
		"-crf",
		"23",
	)

	// 13. Audio codec
	args = append(
		args,
		"-c:a",
		"aac",
	)

	// 14. Audio bitrate
	args = append(
		args,
		"-b:a",
		"128k",
	)

	// 15. Output duration
	args = append(
		args,
		"-t",
		fmt.Sprintf(
			"%d",
			totalDuration,
		),
	)

	// 16. Output file
	args = append(
		args,
		outputFile,
	)

	fmt.Println("FFmpeg command:")
	fmt.Println(
		r.FFmpegPath,
		strings.Join(args, " "),
	)

	// 17. Execute FFmpeg
	cmd := exec.Command(
		r.FFmpegPath,
		args...,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {

		fmt.Println(
			"FFmpeg output:",
		)

		fmt.Println(
			string(output),
		)

		return "",
			fmt.Errorf(
				"ffmpeg failed: %w",
				err,
			)
	}

	fmt.Println(
		"FFmpeg completed successfully",
	)

	_ = imageCount

	return outputFile, nil
}

func addMediaInputs(
	args *[]string,
	invitation models.Invitation,
) {

	for _, slide := range invitation.Slides {

		for _, element := range slide.Elements {

			switch element.Type {

			case "image":

				*args = append(
					*args,
					"-loop",
					"1",
					"-i",
					element.AssetPath,
				)

			case "gif":

				*args = append(
					*args,
					"-stream_loop",
					"-1",
					"-i",
					element.AssetPath,
				)
			}
		}
	}
}

func calculateTotalDuration(
	invitation models.Invitation,
) int {

	totalDuration := 0

	for _, slide := range invitation.Slides {

		totalDuration +=
			slide.DisplayDuration
	}

	return totalDuration
}

func escapeFFmpegText(text string) string {
	text = strings.ReplaceAll(
		text,
		"\\",
		"\\\\",
	)

	text = strings.ReplaceAll(
		text,
		"'",
		"\\'",
	)

	text = strings.ReplaceAll(
		text,
		":",
		"\\:",
	)

	return text
}

// func buildImageFilter(
// 	videoLabel string,
// 	imageLabel string,
// 	outputLabel string,
// 	element models.Element,
// 	start int,
// 	end int,
// ) string {

// 	scale := fmt.Sprintf(
// 		"scale=%d:%d:force_original_aspect_ratio=decrease",
// 		element.Width,
// 		element.Height,
// 	)

//		return fmt.Sprintf(
//			"%s%s%s%s;"+
//				"%s%s%s",
//			imageLabel,
//			scale,
//			"[img]",
//			"",
//			videoLabel,
//			"[img]overlay=%d:%d:enable='between(t,%d,%d)'%s",
//			element.X,
//			element.Y,
//			start,
//			end,
//			outputLabel,
//		)
//	}
func buildImageFilter(
	videoLabel string,
	imageLabel string,
	outputLabel string,
	element models.Element,
	start int,
	end int,
) string {

	imageLabelName :=
		strings.Trim(
			outputLabel,
			"[]",
		) + "_img"

	opacity := getOpacity(element)

	rotation := getRotation(element)

	scaleFilter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=decrease",
		element.Width,
		element.Height,
	)

	opacityFilter := fmt.Sprintf(
		"format=rgba,colorchannelmixer=aa=%f",
		opacity,
	)

	imageFilters := scaleFilter

	if rotation != 0 {

		imageFilters += fmt.Sprintf(
			",rotate=%f*PI/180:"+
				"ow=rotw(iw):"+
				"oh=roth(ih):"+
				"c=none",
			rotation,
		)
	}

	imageFilters += "," + opacityFilter

	return fmt.Sprintf(
		"%s%s[%s];"+
			"%s[%s]overlay=%d:%d:"+
			"enable='between(t,%d,%d)'%s",

		imageLabel,

		imageFilters,

		imageLabelName,

		videoLabel,

		imageLabelName,

		element.X,

		element.Y,

		start,

		end,

		outputLabel,
	)
}
