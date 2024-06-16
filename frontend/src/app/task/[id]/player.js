"use client";

import { useState, useContext, useEffect } from "react";

import {
  MediaPlayer,
  MediaProvider,
  Track,
  useActiveTextTrack,
  useActiveTextCues,
} from "@vidstack/react";

import {
  DefaultAudioLayout,
  defaultLayoutIcons,
  DefaultVideoLayout,
} from "@vidstack/react/player/layouts/default";
import "@vidstack/react/player/styles/default/theme.css";
import "@vidstack/react/player/styles/default/layouts/video.css";

import { PageContext } from "./context";

function ChaptersObserver() {
  const activeChaptersTrack = useActiveTextTrack("chapters");
  const activeCues = useActiveTextCues(activeChaptersTrack);
  const { setData } = useContext(PageContext);

  useEffect(() => {
    if (activeCues.length > 0) {
      setData((prevdata) => ({
        ...prevdata,
        current: activeCues[0].id,
      }));
    } else {
      setData((prevdata) => ({ ...prevdata, current: null }));
    }
  }, [activeCues]);
  return <></>;
}

export function Player() {
  const [content, setContent] = useState([]);
  const { data } = useContext(PageContext);

  useEffect(() => {
    setContent(
      data.copyright.map((t) => ({
        id: t.origId,
        startTime: t.copyrightStart,
        endTime: t.copyrightEnd,
      }))
    );
  }, [data.copyright]);

  return (
    <MediaPlayer
      viewType={"video"}
      src={data.videoUrl}
      load={"visible"} // play
      posterLoad={"visible"}
      aspectRatio="16 / 9"
    >
      <MediaProvider>
        <ChaptersObserver />
        <Track
          content={content}
          kind="chapters"
          language="en-US"
          type="json"
          default={true}
        />
      </MediaProvider>
      <DefaultAudioLayout icons={defaultLayoutIcons} />
      <DefaultVideoLayout icons={defaultLayoutIcons} />
    </MediaPlayer>
  );
}
