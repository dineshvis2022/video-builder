export type Invitation = {
  projectId: string;
  background: Background;
  audio?: AudioConfig;
  slides: Slide[];
};

export type Background = {
  url: string;
};

export type AudioConfig = {
  url: string;
  volume?: number;
  loop?: boolean;
  fadeIn?: number;
  fadeOut?: number;
  mixWithBackground?: boolean;
  backgroundVolume?: number;
};

export type Slide = {
  slideId: number;

  displayDuration: number;

  start?: number;
  end?: number;

  transition?: Transition;
  durationInFrames: number;
  elements: Element[];
};

export type Transition = {
  type: string;
  duration: number;
};

export type Element = {
  id: string;

  type: "text" | "image" | "gif";

  text?: string;

  assetPath?: string;

  x: number;
  y: number;

  width?: number;
  height?: number;

  fontSize?: number;

  fontColor?: string;

  fontFile?: string;

  align?: "left" | "center" | "right";

  opacity?: number;

  rotation?: number;

  animation?: Animation;
  
};

export type Animation = {
  type:
    | "fadeIn"
    | "fadeOut"
    | "slideLeft"
    | "slideRight"
    | "slideUp"
    | "slideDown"
    | "zoomIn"
    | "zoomOut";

  duration?: number;

  distance?: number;
};

