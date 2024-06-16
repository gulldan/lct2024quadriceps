"use client"

import { Alert } from "react-bootstrap";

 

const style = {
    position: "fixed",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
  };
  
  export default function Error({error}) {
    return (
      <div style={style}>
        <Alert variant="danger">
            {error.message}
        </Alert>
      </div>
    );
  }
  