package renderer

import (
	"fmt"
	"strings"
	"video-builder/models"
)

func buildAudioFilter(
	invitation models.Invitation,
	audioInputIndex int,
	totalDuration int,
) (string, bool) {

	if invitation.Audio == nil ||
		invitation.Audio.URL == "" {
		return "", false
	}

	audio := invitation.Audio

	volume := audio.Volume
	if volume <= 0 {
		volume = 1
	}

	filters := []string{
		fmt.Sprintf(
			"[%d:a]atrim=duration=%d",
			audioInputIndex,
			totalDuration,
		),
		"asetpts=PTS-STARTPTS",
		fmt.Sprintf(
			"volume=%f",
			volume,
		),
	}

	if audio.FadeIn > 0 {
		filters = append(
			filters,
			fmt.Sprintf(
				"afade=t=in:st=0:d=%f",
				audio.FadeIn,
			),
		)
	}

	if audio.FadeOut > 0 {

		start :=
			float64(totalDuration) -
				audio.FadeOut

		if start < 0 {
			start = 0
		}

		filters = append(
			filters,
			fmt.Sprintf(
				"afade=t=out:st=%f:d=%f",
				start,
				audio.FadeOut,
			),
		)
	}

	return strings.Join(filters, ",") +
		"[music]", true
}

func buildBackgroundAudioFilter(
	invitation models.Invitation,
	totalDuration int,
) string {

	volume := 1.0

	if invitation.Audio != nil &&
		invitation.Audio.BackgroundVolume > 0 {

		volume =
			invitation.Audio.BackgroundVolume
	}

	return fmt.Sprintf(
		"[0:a]atrim=duration=%d,"+
			"asetpts=PTS-STARTPTS,"+
			"volume=%f[bgAudio]",
		totalDuration,
		volume,
	)
}
