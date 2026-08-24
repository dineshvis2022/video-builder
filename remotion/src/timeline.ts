import { Invitation, Slide } from "./types";

export type TimelineSlide = Slide & {
  startFrame: number;
  endFrame: number;
  durationInFrames: number;
};

export const buildTimeline = (
  invitation: Invitation,
  fps: number,
): TimelineSlide[] => {
  let currentFrame = 0;

  return invitation.slides.map((slide) => {
    const durationInFrames = Math.round(
      slide.displayDuration * fps,
    );

    const startFrame = currentFrame;

    const endFrame =
      startFrame + durationInFrames;

    currentFrame = endFrame;

    return {
      ...slide,

      startFrame,
      endFrame,
      durationInFrames,
    };
  });
};

export const getTotalDurationInFrames = (
  invitation: Invitation,
  fps: number,
): number => {
  return invitation.slides.reduce(
    (total, slide) =>
      total +
      Math.round(slide.displayDuration * fps),
    0,
  );
};