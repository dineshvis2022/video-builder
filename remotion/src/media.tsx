import React from "react";
import { Img, staticFile } from "remotion";

import { Gif } from "@remotion/gif";

import { Element } from "./types";

type MediaElementProps = {
  element: Element;
  style: React.CSSProperties;
  transform: string | undefined;
  opacity: number;
};

export const ImageElement: React.FC<
  MediaElementProps
> = ({ element, style }) => {
  if (!element.assetPath) {
    return null;
  }

  return (
    <Img
      src={getAssetPath(element.assetPath)}
      style={{
        ...style,
        objectFit: "contain",
      }}
    />
  );
};

export const GifElement: React.FC<
  MediaElementProps
> = ({ element, style }) => {
  if (!element.assetPath) {
    return null;
  }

  return (
    <Gif
      src={getAssetPath(element.assetPath)}
      width={element.width}
      height={element.height}
      fit="contain"
      style={style}
    />
  );
};

const getAssetPath = (
  assetPath: string,
): string => {
  if (
    assetPath.startsWith("http://") ||
    assetPath.startsWith("https://")
  ) {
    return assetPath;
  }

  if (assetPath.startsWith("/")) {
    return staticFile(assetPath);
  }

  return staticFile(`/${assetPath}`);
};