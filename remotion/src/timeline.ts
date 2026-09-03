import {
  Invitation,
  Slide,
} from "./types";

export type TimelineSlide =
  Slide & {
    startFrame: number;
    endFrame: number;
    durationInFrames: number;
  };

export const buildTimeline = (
  invitation: Invitation,
  fps: number,
): TimelineSlide[] => {
  let currentFrame = 0;

  return invitation.slides.map(
    (slide) => {
      const durationInFrames =
        Math.round(
          slide.displayDuration *
            fps,
        );

      const startFrame =
        currentFrame;

      const endFrame =
        startFrame +
        durationInFrames;

      const transitionFrames =
        Math.round(
          (slide.transition
            ?.duration || 0) *
            fps,
        );

      currentFrame =
        endFrame -
        transitionFrames;

      return {
        ...slide,

        startFrame,

        endFrame,

        durationInFrames,
      };
    },
  );
};

export const getTotalDurationInFrames =
  (
    invitation: Invitation,
    fps: number,
  ): number => {
    let totalFrames = 0;

    invitation.slides.forEach(
      (slide) => {
        totalFrames +=
          Math.round(
            slide.displayDuration *
              fps,
          );

        totalFrames -=
          Math.round(
            (slide.transition
              ?.duration || 0) *
              fps,
          );
      },
    );

    return Math.max(
      1,
      totalFrames,
    );
  };