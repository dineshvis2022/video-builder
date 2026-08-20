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

	// 1. Build video filter graph
	filterBuilder := NewFilterBuilder()

	filterComplex, _, err :=
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

	// 3. Calculate final rendered duration
	// Transitions overlap two slides, so their duration
	// is subtracted from the total duration.
	totalDuration := calculateRenderedDuration(
		invitation,
	)

	if totalDuration <= 0 {
		return "", fmt.Errorf(
			"invalid rendered duration: %f",
			totalDuration,
		)
	}

	// 4. Create output file
	outputFile := fmt.Sprintf(
		"%s/%s.mp4",
		r.OutputDir,
		invitation.ProjectID,
	)

	// 5. Build FFmpeg arguments
	args := []string{
		"-y",

		// Background video
		"-i",
		invitation.Background.URL,
	}

	// 6. Add image and GIF inputs
	addMediaInputs(
		&args,
		invitation,
	)

	// 7. Calculate external audio input index
	//
	// Input 0 = background video
	// Input 1..N = image/GIF inputs
	// Input N+1 = external audio
	audioInputIndex := -1

	if invitation.Audio != nil &&
		invitation.Audio.URL != "" {

		audioInputIndex =
			1 +
				countMediaInputs(
					invitation,
				)

		// Loop audio when requested
		if invitation.Audio.Loop {

			args = append(
				args,
				"-stream_loop",
				"-1",
			)
		}

		// External music
		args = append(
			args,
			"-i",
			invitation.Audio.URL,
		)
	}

	// 8. Build audio filters
	//
	// Video filter already contains [vout].
	// Here we append audio filters when external
	// audio is configured.
	if audioInputIndex >= 0 {

		audioFilter,
			hasAudio :=
			buildAudioFilter(
				invitation,
				audioInputIndex,
				int(totalDuration),
			)

		if !hasAudio {
			return "", fmt.Errorf(
				"failed to build audio filter",
			)
		}

		filterComplex += ";" +
			audioFilter

		// 9. Mix external music with background audio
		if invitation.Audio.MixWithBackground {

			backgroundAudioFilter :=
				buildBackgroundAudioFilter(
					invitation,
					int(totalDuration),
				)

			filterComplex += ";" +
				backgroundAudioFilter

			// Mix background audio and music
			filterComplex +=
				"[bgAudio][music]" +
					"amix=inputs=2:" +
					"duration=first:" +
					"dropout_transition=0" +
					"[aout]"

		} else {

			// Use external music only
			filterComplex +=
				"[music]anull[aout]"
		}
	}

	// 10. Add filter complex
	args = append(
		args,
		"-filter_complex",
		filterComplex,
	)

	// 11. Map rendered video
	args = append(
		args,
		"-map",
		"[vout]",
	)

	// 12. Map audio
	if audioInputIndex >= 0 {

		// External/mixed audio
		args = append(
			args,
			"-map",
			"[aout]",
		)

	} else {

		// Keep background video's original audio
		args = append(
			args,
			"-map",
			"0:a?",
		)
	}

	// 13. Video codec
	args = append(
		args,
		"-c:v",
		"libx264",
	)

	// 14. Encoding preset
	args = append(
		args,
		"-preset",
		"veryfast",
	)

	// 15. Video quality
	args = append(
		args,
		"-crf",
		"23",
	)

	// 16. Pixel format
	args = append(
		args,
		"-pix_fmt",
		"yuv420p",
	)

	// 17. Audio codec
	args = append(
		args,
		"-c:a",
		"aac",
	)

	// 18. Audio bitrate
	args = append(
		args,
		"-b:a",
		"128k",
	)

	// 19. Output duration
	args = append(
		args,
		"-t",
		fmt.Sprintf(
			"%.3f",
			totalDuration,
		),
	)

	// 20. Output file
	args = append(
		args,
		outputFile,
	)

	// 21. Print FFmpeg command
	fmt.Println("FFmpeg command:")
	fmt.Println(
		r.FFmpegPath,
		strings.Join(args, " "),
	)

	// 22. Execute FFmpeg
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

	// 23. Print success
	fmt.Println(
		"FFmpeg completed successfully",
	)

	fmt.Println(
		"Output:",
		outputFile,
	)

	return outputFile, nil
}

func countMediaInputs(
	invitation models.Invitation,
) int {

	count := 0

	for _, slide := range invitation.Slides {

		for _, element := range slide.Elements {

			switch element.Type {

			case "image", "gif":
				count++
			}
		}
	}

	return count
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

func calculateRenderedDuration(
	invitation models.Invitation,
) float64 {

	duration := 0.0

	for i, slide := range invitation.Slides {

		duration +=
			float64(
				slide.DisplayDuration,
			)

		if i > 0 &&
			slide.Transition != nil {

			duration -=
				slide.Transition.Duration
		}
	}

	return duration
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
