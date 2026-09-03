import {
  staticFile,
} from "remotion";

export type FontDefinition = {
  family: string;
  src: string;
};

export const getFontDefinition = (
  fontFile?: string,
): FontDefinition | null => {
  if (!fontFile) {
    return null;
  }

  const filename =
    fontFile
      .split("/")
      .pop() || "";

  const family =
    filename
      .replace(/\.[^/.]+$/, "");

  return {
    family,
    src: getFontPath(
      fontFile,
    ),
  };
};

const getFontPath = (
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