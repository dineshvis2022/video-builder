import React from "react";

import {
  getFontDefinition,
} from "./fonts";

import {
  AbsoluteFill,
  Img,
  useCurrentFrame,
  useVideoConfig,
  staticFile,
} from "remotion";

import {
  Invitation,
  Slide,
  Element,
} from "./types";

import {
  buildTimeline,
  TimelineSlide,
} from "./timeline";

import {
  getAnimationStyle,
} from "./animations/animations";

import {
  getTransitionStyle,
} from "./transitions";

import {
  GifElement,
} from "./media";

import {
  AudioLayer,
} from "./audio";

// {fontDefinitions.map(
//   (font) => (
//     <style
//       key={font.family}
//     >{`
//       @font-face {
//         font-family: '${font.family}';
//         src: url('${font.src}');
//       }
//     `}</style>
//   ),
// )}

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

const activeSlides =
  getActiveSlides(
    timeline,
    frame,
    fps,
  );
  

  const getFontDefinitions = (
  ) => {
    const fonts = new Map<
      string,
      ReturnType<
        typeof getFontDefinition
      >
    >();

    for (
      const slide of invitation.slides
    ) {
      for (
        const element of slide.elements
      ) {
        if (
          element.type === "text" &&
          element.fontFile
        ) {
          const definition =
            getFontDefinition(
              element.fontFile,
            );

          if (definition) {
            fonts.set(
              definition.family,
              definition,
            );
          }
        }
      }
    }

  return Array.from(
    fonts.values(),
  ).filter(
    (
      font,
    ): font is NonNullable<
      typeof font
    > => Boolean(font),
  );
};


  return (
    <AbsoluteFill
      style={{
        backgroundColor: "black",
      }}
    >
      <BackgroundVideo
        invitation={invitation}
      />

      {activeSlides.map(
        ({
          slide,
          transitionFrame,
          transitionDuration,
          isEntering,
        }) => (
          <TransitionSlide
            key={slide.slideId}
            slide={slide}
            frame={
              frame -
              slide.startFrame
            }
            transitionFrame={
              transitionFrame
            }
            transitionDuration={
              transitionDuration
            }
            isEntering={isEntering}
            width={width}
            height={height}
            fps={fps}
          />
        ),
      )}

      <AudioLayer
        invitation={invitation}
      />
    </AbsoluteFill>
  );
};

type ActiveSlide = {
  slide: TimelineSlide;
  transitionFrame: number;
  transitionDuration: number;
  isEntering: boolean;
};

const getActiveSlides = (
  timeline: TimelineSlide[],
  frame: number,
  fps: number,
): ActiveSlide[] => {
  if (timeline.length === 0) {
    return [];
  }

  const activeIndex =
    timeline.findIndex(
      (slide) =>
        frame >= slide.startFrame &&
        frame < slide.endFrame,
    );

  if (activeIndex < 0) {
    return [];
  }

  const current =
    timeline[activeIndex];

  const previous =
    activeIndex > 0
      ? timeline[activeIndex - 1]
      : undefined;

  const transition =
    current.transition;

  const transitionDuration =
    transition
      ? transition.duration
      : 0;

  const transitionFrames =
    Math.round(
      transitionDuration *
        fps,
    );

  const transitionStart =
    current.startFrame;

  const isEntering =
    transitionFrames > 0 &&
    frame <
      transitionStart +
        transitionFrames;

  if (
    isEntering &&
    previous
  ) {
    return [
      {
        slide: previous,
        transitionFrame:
          frame -
          transitionStart,
        transitionDuration:
          transitionFrames,
        isEntering: false,
      },
      {
        slide: current,
        transitionFrame:
          frame -
          transitionStart,
        transitionDuration:
          transitionFrames,
        isEntering: true,
      },
    ];
  }

  return [
    {
      slide: current,
      transitionFrame: 0,
      transitionDuration: 0,
      isEntering: false,
    },
  ];
};

type BackgroundVideoProps = {
  invitation: Invitation;
};

const BackgroundVideo: React.FC<
  BackgroundVideoProps
> = ({ invitation }) => {
  if (!invitation.background.url) {
    return null;
  }

  const hasExternalAudio =
    invitation.audio &&
    invitation.audio.url !== "";

  const mixBackground =
    invitation.audio
      ?.mixWithBackground === true;

  const muted =
    Boolean(hasExternalAudio) &&
    !mixBackground;

  const backgroundVolume =
    invitation.audio
      ?.backgroundVolume || 1;

  return (
    <video
      src={getAssetPath(
        invitation.background.url,
      )}
      muted={muted}
      // volume={
      //     muted
      //       ? 0
      //       : backgroundVolume
      //   }
        loop
      autoPlay
      style={{
        position: "absolute",
        width: "100%",
        height: "100%",
        objectFit: "cover",
      }}
    />
  );
};

type TransitionSlideProps = {
  slide: TimelineSlide;
  frame: number;
  transitionFrame: number;
  transitionDuration: number;
  isEntering: boolean;
  width: number;
  height: number;
  fps: number;
};

const TransitionSlide: React.FC<
  TransitionSlideProps
> = ({
  slide,
  frame,
  transitionFrame,
  transitionDuration,
  isEntering,
  width,
  height,
  fps,
}) => {
  const transitionStyle =
    getTransitionStyle(
      slide.transition?.type,
      transitionFrame,
      transitionDuration,
      isEntering,
    );

  
  return (
    <AbsoluteFill
      style={{
      opacity:
        transitionStyle.opacity,

      transform:
        transitionStyle.transform,

      clipPath:
        transitionStyle.clipPath,

      zIndex:
        isEntering ? 2 : 1,
    }}
    >
      <SlideView
        slide={slide}
        frame={frame}
        width={width}
        height={height}
        fps={fps}
      />
    </AbsoluteFill>
  );
};

type SlideViewProps = {
  slide: Slide;
  frame: number;
  width: number;
  height: number;
  fps: number;
};

const SlideView: React.FC<
  SlideViewProps
> = ({
  slide,
  frame,
  width,
  height,
  fps,
}) => {
  return (
    <AbsoluteFill>
      {slide.elements.map(
        (element) => (
          <ElementView
            key={element.id}
            element={element}
            frame={frame}
            fps={fps}
            width={width}
            height={height}
            slideDurationInFrames={
              slide.durationInFrames
            }
          />
        ),
      )}
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
  slideDurationInFrames,
}) => {
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
    animationStyle.transform ===
    "none"
      ? rotation || undefined
      : `${rotation} ${animationStyle.transform}`;

  const opacity =
    (element.opacity ?? 1) *
    animationStyle.opacity;

  if (element.type === "text") {
    return (
      <TextElement
        element={element}
        transform={transform}
        opacity={opacity}
      />
    );
  }

  if (element.type === "image") {
    return (
      <ImageElement
        element={element}
        transform={transform}
        opacity={opacity}
      />
    );
  }

  if (element.type === "gif") {
    return (
      <GifElement
        element={element}
        transform={transform}
        opacity={opacity}
        style={animationStyle}
      />
    );
  }

  return null;
};

type TextElementProps = {
  element: Element;
  transform: string | undefined;
  opacity: number;
};

const TextElement: React.FC<
  TextElementProps
> = ({
  element,
  transform,
  opacity,
}) => {
  const align =
    element.align || "left";

  const xTransform =
    align === "center"
      ? "translateX(-50%)"
      : align === "right"
        ? "translateX(-100%)"
        : "";

  const finalTransform =
    transform
      ? `${xTransform} ${transform}`
      : xTransform || undefined;

  return (
    <div
      style={{
        position: "absolute",

        left: element.x,
        top: element.y,

        transform:
          finalTransform,

        fontSize:
          element.fontSize || 60,

        color:
          element.fontColor ||
          "white",

        opacity,

        textAlign: align,

        whiteSpace:
          "pre-wrap",

        fontFamily:
          element.fontFile
            ? element.fontFile
                .split("/")
                .pop()
                ?.replace(
                  /\.[^/.]+$/,
                  "",
                )
            : "Arial",
      }}
    >
      {element.text}
    </div>
  );
};

type ImageElementProps = {
  element: Element;
  transform: string | undefined;
  opacity: number;
};

const ImageElement: React.FC<
  ImageElementProps
> = ({
  element,
  transform,
  opacity,
}) => {
  if (!element.assetPath) {
    return null;
  }

  return (
    <Img
      src={getAssetPath(
        element.assetPath,
      )}
      style={{
        position: "absolute",

        left: element.x,
        top: element.y,

        width:
          element.width,

        height:
          element.height,

        objectFit:
          "contain",

        opacity,

        transform,
      }}
    />
  );
};

const getAssetPath = (
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
    return path;
  }

  return staticFile(`/${path}`);
};

const getFontFamily = (
  path: string,
): string => {
  return path
    .split("/")
    .pop()
    ?.split(".")[0] || "CustomFont";
};

