import React from "react";

import {
  AbsoluteFill,
  interpolate,
  staticFile,
  useVideoConfig,
} from "remotion";

import {
  Html5Audio,
} from "remotion";

import {
  Invitation,
} from "./types";

type Props = {
  invitation: Invitation;
};

export const AudioLayer: React.FC<
  Props
> = ({ invitation }) => {
  const {
    fps,
    durationInFrames,
  } = useVideoConfig();

  const audio =
    invitation.audio;

  if (
    !audio ||
    !audio.url
  ) {
    return null;
  }

  const volume =
    audio.volume && audio.volume > 0
      ? audio.volume
      : 1;

  const fadeInFrames =
    Math.round(
      (audio.fadeIn || 0) *
        fps,
    );

  const fadeOutFrames =
    Math.round(
      (audio.fadeOut || 0) *
        fps,
    );

  const fadeOutStart =
    Math.max(
      0,
      durationInFrames -
        fadeOutFrames,
    );

  const getVolume = (
    frame: number,
  ) => {
    let result = volume;

    if (
      fadeInFrames > 0
    ) {
      result *= interpolate(
        frame,
        [0, fadeInFrames],
        [0, 1],
        {
          extrapolateLeft:
            "clamp",
          extrapolateRight:
            "clamp",
        },
      );
    }

    if (
      fadeOutFrames > 0
    ) {
      result *= interpolate(
        frame,
        [
          fadeOutStart,
          durationInFrames,
        ],
        [1, 0],
        {
          extrapolateLeft:
            "clamp",
          extrapolateRight:
            "clamp",
        },
      );
    }

    return result;
  };

  return (
    <AbsoluteFill
      style={{
        pointerEvents:
          "none",
      }}
    >
      <Html5Audio
        src={getAudioPath(
          audio.url,
        )}
        volume={getVolume}
        loop={audio.loop}
      />
    </AbsoluteFill>
  );
};

const getAudioPath = (
  path: string,
): string => {
  if (
    path.startsWith(
      "http://",
    ) ||
    path.startsWith(
      "https://",
    )
  ) {
    return path;
  }

  if (path.startsWith("/")) {
    return staticFile(path);
  }

  return staticFile(`/${path}`);
};