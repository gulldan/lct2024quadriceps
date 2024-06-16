import { Spinner } from "react-bootstrap";

const style = {
  position: "fixed",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
};

export default function Load() {
  return (
    <div style={style}>
      <Spinner />
    </div>
  );
}
