"use client";

import { useEffect, useState } from "react";
import { Fireplace } from "./fire";
import { Alert, Button, Image } from "react-bootstrap";

import "./styles.css";

export default function HomePage() {
  const [fire, setFire] = useState(null);
  const [play, setPlay] = useState(false);

  useEffect(() => {
    setFire(new Fireplace());
  }, []);
  useEffect(() => {
    if (fire) {
      if (!play) {
        fire.config.heat = 0;
      }
      if (fire.inited) {
        return;
      }
      fire.init();
      fire.particlesSpawnRate = 10;
    }
  }, [fire, play]);

  return (
    <>
      <canvas id="cvs"></canvas>
      {play && <Image className="left-screen" src="ani2.gif" />}
      {play && <Image className="right-screen" src="ani1.gif" />}
      <div className="center-screen">
        <Alert variant="dark">
          <Button variant="dark"
            disabled={play}
            onClick={() => {
              fire.config.heat = 1000;
              setPlay(true);
            }}
          >
            {play ? "Вот и думай" : "Включить светлую тему"}

            {play && (
              <audio autoPlay loop>
                <source src="/sound_grabej.mp4" type="audio/mp4" />
                Your browser does not support the audio element.
              </audio>
            )}
          </Button>
        </Alert>
      </div>
    </>
  );
}
