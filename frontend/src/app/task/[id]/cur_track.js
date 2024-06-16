import { useContext, useEffect, useState } from "react";

import { MediaPlayer, MediaProvider, Track } from "@vidstack/react";

import {
  DefaultAudioLayout,
  defaultLayoutIcons,
  DefaultVideoLayout,
} from "@vidstack/react/player/layouts/default";
import "@vidstack/react/player/styles/default/theme.css";
import "@vidstack/react/player/styles/default/layouts/video.css";

import { PageContext } from "./context";
import { getTask } from "@/utils/web";

export function CurTrack() {
  const [selected, setSelected] = useState(null);
  const [orig, setOrig] = useState(null);
  const { data } = useContext(PageContext);

  useEffect(() => {
    if (data.current) {
      setSelected(data.copyright[data.current]);
    } else {
      setSelected(null);
    }
  }, [data.current, data.copyright]);

  useEffect(() => {
    if (selected) {
      getTask(selected.orig_id).then((res) => {
        setOrig(res);
      });
    }
  }, [selected]);

  if (!selected || !orig) {
    return <></>;
  }

  return (
    <div>
      <div className="text-center">
        <h2>Оригинальное видео:</h2>
      </div>
      <MediaPlayer
        viewType={"video"}
        src={orig.videoUrl}
        load={"eager"} // play
        posterLoad={"play"}
        aspectRatio="16 / 9"
        currentTime={selected.origStart}
        autoPlay={true}
        muted={true}
      >
        <MediaProvider>
          <Track
            content={[
              { startTime: selected.origStart, endTime: selected.origEnd },
            ]}
            kind="chapters"
            language="en-US"
            type="json"
            default={true}
          />
        </MediaProvider>
        <DefaultAudioLayout icons={defaultLayoutIcons} />
        <DefaultVideoLayout icons={defaultLayoutIcons} />
      </MediaPlayer>
    </div>
  );
}
