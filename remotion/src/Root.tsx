


import React from "react";
import { Composition } from "remotion";

import { Video } from "./Video";
import { Invitation } from "./types";
import { getTotalDurationInFrames } from "./timeline";
import {staticFile} from 'remotion';

const defaultInvitation: Invitation = {
  projectId: "TEST_V5",

  background: {
    url: staticFile("/assets/background.mp4"),
  },
audio: {
  url: "/assets/music.mp3",
  volume: 0.6,
  loop: true,
  fadeIn: 1,
  fadeOut: 2,
  mixWithBackground: true,
  backgroundVolume: 0.2,
},

slides: [
  {
    slideId: 1,
    durationInFrames:5,
    displayDuration: 5,

    elements: [
      {
        id: "couple",
        type: "image",

        assetPath: staticFile("/assets/couple.png"),

        x: 700,
        y: 150,

        width: 500,
        height: 500,

        animation: {
          type: "zoomIn",
          duration: 2,
        },
      },

      {
        id: "title",
        type: "text",

        text:
          "Sadhna ❤️ Kartik",

        x: 960,
        y: 750,

        fontSize: 70,

        fontColor: "white",

        align: "center",

        animation: {
          type: "fadeIn",
          duration: 2,
        },
      },

      {
        id: "gif",
        type: "gif",

        assetPath:
          "/assets/animation.gif",

        x: 1200,
        y: 200,

        width: 300,
        height: 300,
      },
    ],
  },

  {
    slideId: 2,
    durationInFrames:5,
    displayDuration: 5,

    transition: {
      type: "fade",
      duration: 1,
    },

    elements: [
      {
        id: "venue",
        type: "text",

        text:
          "The Grand Palace",

        x: 960,
        y: 400,

        fontSize: 70,

        fontColor: "white",

        align: "center",

        animation: {
          type: "slideUp",
          duration: 2,
        },
      },
    ],
  },

  {
    slideId: 3,
    durationInFrames: 5,
    displayDuration: 5,

    transition: {
      type: "wipeleft",
      duration: 1,
    },

    elements: [
      {
        id: "save",
        type: "text",

        text:
          "Save The Date",

        x: 960,
        y: 400,

        fontSize: 80,

        fontColor: "white",

        align: "center",

        animation: {
          type: "zoomOut",
          duration: 2,
        },
      },
    ],
  },
],
};

const FPS = 30;

const durationInFrames =
  getTotalDurationInFrames(
    defaultInvitation,
    FPS,
  );

export const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="VideoBuilder"

      component={Video}

      width={1920}
      height={1080}

      fps={FPS}

      durationInFrames={
        durationInFrames
      }

      defaultProps={{
        invitation: defaultInvitation,
      }}
    />
  );
};