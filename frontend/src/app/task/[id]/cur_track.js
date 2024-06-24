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

export function CurTrack() {
  const [selected, setSelected] = useState(null);

  const { data } = useContext(PageContext);

  useEffect(() => {
    if (data.current) {
      setSelected(data.current);
    } else {
      setSelected(null);
    }
  }, [data.current, data.copyright]);

  if (!selected) {
    return <></>;
  }

  if (!selected.origUrl) {
    <div>
      <div className="text-center">
        <h2>
          Извините, пока почему-то нет ссылки на оригинал, наша команда
          специалистов разбирается с данной проблемой
        </h2>
      </div>
    </div>;
  }

  return (
    <div>
      <div className="text-center">
        <h2>Оригинальное видео:</h2>
      </div>
      <MediaPlayer
        viewType={"video"}
        src={selected.origUrl}
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
