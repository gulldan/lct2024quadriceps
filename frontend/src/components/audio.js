"use client";

import { useEffect, useRef } from "react";
export default function Audio() {
  const a = useRef();
  useEffect(() => {
    //audio.play();
    console.log(a);
    if (a.current) {
      try {
        a.current.play();
      } catch (err) {}
    }
    return () => {
      // --> componentWillUnmount
      if (a.current) {
        a.current.pause();
      }
      //audio.release();
    };
  }, [a]);

  return <audio ref={a} loop src="/komanda.mp3"></audio>;
}
