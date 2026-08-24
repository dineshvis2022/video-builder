import { interpolate } from "remotion";

import { Animation } from "../types";

export type AnimationStyle = {
  opacity: number;
  transform: string;
};

export const getAnimationStyle = (
  animation: Animation | undefined,
  frame: number,
  fps: number,
  totalFrames: number,
): AnimationStyle => {
  if (!animation) {
    return {
      opacity: 1,
      transform: "none",
    };
  }

  const durationInFrames = Math.max(
    1,
    Math.round(
      (animation.duration || 1) * fps,
    ),
  );

  const type = animation.type;

  if (type === "fadeIn") {
    return {
      opacity: interpolate(
        frame,
        [0, durationInFrames],
        [0, 1],
        {
          extrapolateLeft: "clamp",
          extrapolateRight: "clamp",
        },
      ),

      transform: "none",
    };
  }

  if (type === "fadeOut") {
    const startFrame = Math.max(
      0,
      totalFrames - durationInFrames,
    );

    return {
      opacity: interpolate(
        frame,
        [
          startFrame,
          totalFrames,
        ],
        [1, 0],
        {
          extrapolateLeft: "clamp",
          extrapolateRight: "clamp",
        },
      ),

      transform: "none",
    };
  }

  if (type === "slideLeft") {
    const distance =
      animation.distance || 500;

    const x = interpolate(
      frame,
      [0, durationInFrames],
      [distance, 0],
      {
        extrapolateLeft: "clamp",
        extrapolateRight: "clamp",
      },
    );

    return {
      opacity: 1,
      transform: `translateX(${x}px)`,
    };
  }

  if (type === "slideRight") {
    const distance =
      animation.distance || 500;

    const x = interpolate(
      frame,
      [0, durationInFrames],
      [-distance, 0],
      {
        extrapolateLeft: "clamp",
        extrapolateRight: "clamp",
      },
    );

    return {
      opacity: 1,
      transform: `translateX(${x}px)`,
    };
  }

  if (type === "slideUp") {
    const distance =
      animation.distance || 500;

    const y = interpolate(
      frame,
      [0, durationInFrames],
      [distance, 0],
      {
        extrapolateLeft: "clamp",
        extrapolateRight: "clamp",
      },
    );

    return {
      opacity: 1,
      transform: `translateY(${y}px)`,
    };
  }

  if (type === "slideDown") {
    const distance =
      animation.distance || 500;

    const y = interpolate(
      frame,
      [0, durationInFrames],
      [-distance, 0],
      {
        extrapolateLeft: "clamp",
        extrapolateRight: "clamp",
      },
    );

    return {
      opacity: 1,
      transform: `translateY(${y}px)`,
    };
  }

  return {
    opacity: 1,
    transform: "none",
  };
};