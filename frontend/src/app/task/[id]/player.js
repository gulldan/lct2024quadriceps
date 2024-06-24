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
  const { data, setData } = useContext(PageContext);

  useEffect(() => {
    if (activeCues.length > 0) {
      setData((prevdata) => ({
        ...prevdata,
        current: activeCues[0],
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
    console.log(data)
    setContent(
      data.copyright.map((t, i) => ({
        id: t.origId,
        origUrl: t.origUrl,
        startTime: parseInt(t.copyrightStart),
        endTime: parseInt(t.copyrightEnd),
        origStart: parseInt(t.origStart),
        origEnd: parseInt(t.origEnd),
      }))
    );
  }, [data.copyright]);

  return (<div>
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
    </MediaPlayer></div>
  );
}
