


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

  slides: [
    {
      slideId: 1,

      displayDuration: 5,

      elements: [
        {
          id: "image-1",
          type: "image",

          assetPath: staticFile("/assets/couple.png"),

          x: 710,
          y: 150,

          width: 500,
          height: 500,

          animation: {
            type: "fadeIn",
            duration: 2,
          },
        },

        {
          id: "text-1",
          type: "text",
          text: "Sadhna ❤️ Kartik",

          x: 960,
          y: 750,

          fontSize: 70,

          fontColor: "white",

          align: "center",

          animation: {
            type: "slideLeft",
            duration: 2,
            distance: 500,
          },
        },
      ],
    },

    {
      slideId: 2,

      displayDuration: 5,

      elements: [
        {
          id: "text-2",
          type: "text",

          text: "The Grand Palace",

          x: 960,
          y: 400,

          fontSize: 70,

          fontColor: "white",

          align: "center",
        },
      ],
    },

    {
      slideId: 3,

      displayDuration: 5,

      elements: [
        {
          id: "text-3",
          type: "text",

          text: "Save The Date",

          x: 960,
          y: 400,

          fontSize: 70,

          fontColor: "white",

          align: "center",
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