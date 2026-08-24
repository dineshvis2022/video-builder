import React from "react";

import {
  AbsoluteFill,
  useCurrentFrame,
  useVideoConfig,
  Video as RemotionVideo,
} from "remotion";

import {
  Invitation,
  Slide,
  Element,
} from "./types";

import {
  buildTimeline,
} from "./timeline";

import {getAnimationStyle} from "./animations/animations";


type Props = {
  invitation: Invitation;
};

export const Video: React.FC<Props> = ({
  invitation,
}) => {
  const frame = useCurrentFrame();

  const {
    fps,
    width,
    height,
  } = useVideoConfig();

  const timeline = buildTimeline(
    invitation,
    fps,
  );

  const currentSlide =
    timeline.find(
      (slide) =>
        frame >= slide.startFrame &&
        frame < slide.endFrame,
    );

  if (!currentSlide) {
    return (
      <AbsoluteFill
        style={{
          backgroundColor: "black",
        }}
      />
    );
  }

  const slideFrame =
    frame - currentSlide.startFrame;

 return (
    <AbsoluteFill
      style={{
        scale: 0.997,
      }}
    >
      <RemotionVideo
        src={invitation.background.url}
        muted
        loop
        style={{
          width: "100%",
          height: "100%",
          objectFit: "cover",
        }}
      />

      <SlideView
        slide={currentSlide}
        frame={slideFrame}
        width={width}
        height={height}
      />
    </AbsoluteFill>
  );
};

type SlideViewProps = {
  slide: Slide;
  frame: number;
  width: number;
  height: number;
  
};

const SlideView: React.FC<SlideViewProps> = ({
  slide,
  frame,
  width,
  height,
}) => {
  return (
    <AbsoluteFill>
      {slide.elements.map(
        (element) => (
         <ElementView
            key={element.id}
            element={element}
            frame={frame}
            fps={30}
            width={width}
            height={height}
            slideDurationInFrames={
              slide.displayDuration
            }
          />
        ),
      )}

      <div
        style={{
          position: "absolute",

          left: 30,
          bottom: 30,

          color: "white",

          fontSize: 20,
        }}
      >
        Slide {slide.slideId}
      </div>
    </AbsoluteFill>
  );
};

type ElementViewProps = {
  element: Element;
  frame: number;
  fps: number;
  width: number;
  height: number;
  slideDurationInFrames: number;
};

const ElementView: React.FC<
  ElementViewProps
> = ({
  element,
  frame,
  fps,
  width,
  height,
  slideDurationInFrames
}) => {
  if (element.type === "text") {
    return (
      <TextElement
        element={element}
        frame={frame}
        fps={fps}
        width={width}
        height={height}
        slideDurationInFrames={slideDurationInFrames}
      />
    );
  }

  if (element.type === "image") {
    return (
      <ImageElement
        element={element}
        frame={frame}
        fps={fps}
        width={width}
        height={height}
        slideDurationInFrames={slideDurationInFrames}
      />
    );
  }

  return null;
};

const TextElement: React.FC<
  ElementViewProps
> = ({
  element,
  frame,
  fps,
  slideDurationInFrames
}) => {
  const textAlign =
    element.align === "center"
      ? "center"
      : element.align === "right"
        ? "right"
        : "left";

  const baseTransform =
    element.align === "center"
      ? "translateX(-50%)"
      : element.align === "right"
        ? "translateX(-100%)"
        : "none";

  const animationStyle =
    getAnimationStyle(
      element.animation,
      frame,
      fps,
      slideDurationInFrames,
    );
  const transform =
    animationStyle.transform === "none"
      ? baseTransform
      : `${baseTransform} ${animationStyle.transform}`;

  return (
    <div
      style={{
        position: "absolute",

        left: element.x,
        top: element.y,

        transform,

        fontSize:
          element.fontSize || 40,

        color:
          element.fontColor ||
          "white",

        textAlign,

        opacity:
          (element.opacity ?? 1) *
          animationStyle.opacity,

        whiteSpace:
          "pre-wrap",
      }}
    >
      {element.text}
    </div>
  );
};

const ImageElement: React.FC<
  ElementViewProps
> = ({
  element,
  frame,
  fps,
  slideDurationInFrames
}) => {
  if (!element.assetPath) {
    return null;
  }

 const animationStyle =
  getAnimationStyle(
    element.animation,
    frame,
    fps,
    slideDurationInFrames,
  );
  const rotation =
    element.rotation
      ? `rotate(${element.rotation}deg)`
      : "";

  const transform =
    animationStyle.transform === "none"
      ? rotation || undefined
      : `${rotation} ${animationStyle.transform}`;

  return (
    <img
      src={element.assetPath}
      style={{
        position: "absolute",

        left: element.x,
        top: element.y,

        width: element.width,
        height: element.height,

        objectFit: "contain",

        opacity:
          (element.opacity ?? 1) *
          animationStyle.opacity,

        transform,
      }}
    />
  );
};