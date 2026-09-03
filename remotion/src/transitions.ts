import { interpolate } from "remotion";

export type TransitionStyle = {
  opacity: number;
  transform: string;
  clipPath?: string;
};

export const getTransitionStyle = (
  type: string | undefined,
  frame: number,
  durationInFrames: number,
  isEntering: boolean,
): TransitionStyle => {
  if (
    !type ||
    durationInFrames <= 0
  ) {
    return {
      opacity: 1,
      transform: "none",
    };
  }

  const progress = interpolate(
    frame,
    [0, durationInFrames],
    [0, 1],
    {
      extrapolateLeft: "clamp",
      extrapolateRight: "clamp",
    },
  );

  const normalized =
    type.toLowerCase();

  if (normalized === "fade") {
    return {
      opacity: isEntering
        ? progress
        : 1 - progress,

      transform: "none",
    };
  }

  if (
    normalized ===
    "slideleft"
  ) {
    return {
      opacity: 1,

      transform: isEntering
        ? `translateX(${100 * (1 - progress)}%)`
        : `translateX(${-100 * progress}%)`,
    };
  }

  if (
    normalized ===
    "slideright"
  ) {
    return {
      opacity: 1,

      transform: isEntering
        ? `translateX(${-100 * (1 - progress)}%)`
        : `translateX(${100 * progress}%)`,
    };
  }

  if (
    normalized === "slideup"
  ) {
    return {
      opacity: 1,

      transform: isEntering
        ? `translateY(${100 * (1 - progress)}%)`
        : `translateY(${-100 * progress}%)`,
    };
  }

  if (
    normalized === "slidedown"
  ) {
    return {
      opacity: 1,

      transform: isEntering
        ? `translateY(${-100 * (1 - progress)}%)`
        : `translateY(${100 * progress}%)`,
    };
  }

  if (
    normalized === "wipeleft"
  ) {
    return {
      opacity: 1,

      transform: "none",

      clipPath: isEntering
        ? `inset(0 ${100 * (1 - progress)}% 0 0)`
        : `inset(0 0 0 ${100 * progress}%)`,
    };
  }

  if (
    normalized ===
    "wiperight"
  ) {
    return {
      opacity: 1,

      transform: "none",

      clipPath: isEntering
        ? `inset(0 0 0 ${100 * (1 - progress)}%)`
        : `inset(0 ${100 * progress}% 0 0)`,
    };
  }

  if (
    normalized === "wipeup"
  ) {
    return {
      opacity: 1,

      transform: "none",

      clipPath: isEntering
        ? `inset(${100 * (1 - progress)}% 0 0 0)`
        : `inset(0 0 ${100 * progress}% 0)`,
    };
  }

  if (
    normalized === "wipedown"
  ) {
    return {
      opacity: 1,

      transform: "none",

      clipPath: isEntering
        ? `inset(0 0 ${100 * (1 - progress)}% 0)`
        : `inset(${100 * progress}% 0 0 0)`,
    };
  }

  return {
    opacity: 1,
    transform: "none",
  };
};


// style={{
//   opacity:
//     transitionStyle.opacity,

//   transform:
//     transitionStyle.transform,

//   clipPath:
//     transitionStyle.clipPath,

//   zIndex:
//     isEntering ? 2 : 1,
// }}