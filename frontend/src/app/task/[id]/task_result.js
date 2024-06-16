"use client";
import { useState } from "react";
import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";
import { TypeAnimation } from "react-type-animation";

function good() {
  return (
    <>
      Молодец!!! В видео не найдено копирайтов
      <video autoPlay loop src="/+soc.mp4"></video>
    </>
  );
}

function bad() {
  return (
    <>
      <p>Ужас, тут есть копирайты!</p>

      <TypeAnimation
        sequence={[
          ">sudd",
          500,
          ">sudo ", //  Continuing previous Text
          500,
          ">sudo rm /",
          500,
          ">sudo rm -rf /",
          5000,
          ">sudo rm -rf /\n>",
          2000,
          ">sudo rm -rf /\n>ls",
          2000,
          ">sudo rm -rf /\n>ls\n>zsh: command not found: ls\n>",
          5000,
        ]}
        style={{
          font: "consolas",
          background: "black",
          color: "white",
          whiteSpace: "pre-line",
          height: "110px",
          width: "100%",
          display: "block",
        }}
      />
      <video autoPlay loop src="/-soc.mp4"></video>
    </>
  );
}

export default function TaskResult({ is_good }) {
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);

  return (
    <>
      <Button variant="primary" onClick={handleShow}>
        Подробнее о моей судьбе
      </Button>

      <Modal show={show} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>ИТОГО:</Modal.Title>
        </Modal.Header>
        <Modal.Body>{is_good==true ? good() : bad()}</Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Понятно
          </Button>
          <Button variant="primary" onClick={handleClose}>
            Закрыть
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}
