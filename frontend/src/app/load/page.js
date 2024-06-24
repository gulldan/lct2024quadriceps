"use client";
import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button, Alert, ButtonGroup, Form, Image } from "react-bootstrap";
import { MediaPlayer, MediaProvider } from "@vidstack/react";
import {
  DefaultAudioLayout,
  defaultLayoutIcons,
  DefaultVideoLayout,
} from "@vidstack/react/player/layouts/default";
import "@vidstack/react/player/styles/default/theme.css";
import "@vidstack/react/player/styles/default/layouts/video.css";
import { create_task, upload } from "@/utils/web";

// <MediaProvider />
function PreviewVideo({ file }) {
  const [url, setUrl] = useState(null);
  useEffect(() => {
    if (file == null) {
      setUrl(null);
    } else {
      setUrl(URL.createObjectURL(file));
    }
  }, [file]);
  if (url == null) {
    return <Form.Label>ВИДЯШКУ СЮДА</Form.Label>;
  } else {
    return (
      <>
        <MediaPlayer
          title={"Копирайтшная"}
          controls={false}
          viewType="video"
          src={[{ src: url, type: file.type }]}
          load="visible"
          posterLoad="visible"
        >
          <MediaProvider />
          <DefaultAudioLayout icons={defaultLayoutIcons} />
          <DefaultVideoLayout icons={defaultLayoutIcons} />
        </MediaPlayer>
      </>
    );
  }
}
export default function Home() {
  const [file, setFile] = useState(null);
  const [error, setError] = useState(null);
  const [statusMsg, setStatusMsg] = useState(null);
  const [isLoad, setIsLoad] = useState(false);
  const router = useRouter();

  function handleChange(event) {
    setFile(event.target.files[0]);
  }

  function check() {
    setIsLoad(true);
    setStatusMsg(null);
    create_task(file)
      .catch((e) => {
        setError(e.message);
      })
      .then((t) => {
        router.push(`/task/${t.id}`)
      })
      .finally(() => {
        setIsLoad(false);
      });
  }
  function load() {
    setStatusMsg(null);
    setIsLoad(true);
    upload(file)
      .catch((e) => {
        setError(e.message);
      })
      .then(() => {
        setStatusMsg("Видео загружено в базу оригиналов")
        setFile(null)
      })
      .finally(() => {
        setIsLoad(false);
      });
  }

  function handleSubmit(event) {
    event.preventDefault();
    const buttonType = event.nativeEvent.submitter.name;

    if (file.type != "video/mp4") {
      setError("Поддерживаются только видео формата mp4");
      return;
    } else {
      setError(null);
    }

    if (buttonType == "check") {
      check();
    } else if (buttonType == "load") {
      load();
    } else {
      setError("Ошибка загрузки");
    }
  }

  return (
    <div className="d-flex justify-content-center align-items-center text-center fullscreen">
      <div className="md:w-[50%] xs:w-[100%]">
        {!isLoad && (
          <Form
            onSubmit={handleSubmit}
            style={{ visibility: isLoad ? "hidden" : "" }}
          >
            <Form.Group className="mb-3">
              <PreviewVideo file={file} />
              <Form.Control
                type="file"
                accept="video/mp4"
                onChange={handleChange}
              />
            </Form.Group>

            
            {file && (
              <Form.Group className="mb-3">
                <ButtonGroup>
                  <Button
                    name="check"
                    title="сравнить с текущими видео и если копирайта нет, то внести в базу"
                    type="submit"
                  >
                    ПРОВЕРИТЬ КОПИРАЙТ
                  </Button>
                  <Button
                    name="load"
                    title="занести в базу без проверки"
                    type="submit"
                  >
                    ЗАНЕСТИ В БАЗУ
                  </Button>
                </ButtonGroup>
              </Form.Group>
            )}
            {error && <Alert variant={"danger"}>{error}</Alert>}
            {!error && statusMsg && <Alert variant={"success"}>{statusMsg}</Alert>}
          </Form>
        )}

        {isLoad && (
          <div>
            <p>Видео загружается на сервер...</p>
            <Image src={"/ani1.gif"} />
          </div>
        )}
      </div>
    </div>
  );
}

{
  /* <Form>
  <fieldset>
    <Form.Label>ВИДЯШКУ СЮДА</Form.Label>
    <Form.Control type="file"></Form.Control>
  </fieldset>
</Form>; */
}
