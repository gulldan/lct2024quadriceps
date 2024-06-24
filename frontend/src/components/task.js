"use client";
import Card from "react-bootstrap/Card";
import { Spinner } from "react-bootstrap";
import Badge from "react-bootstrap/Badge";
import CardImg from "react-bootstrap/CardImg";
import CardBody from "react-bootstrap/CardBody";
import CardTitle from "react-bootstrap/CardTitle";
import CardText from "react-bootstrap/CardText";

import { useRouter } from "next/navigation";
import { useCallback } from "react";
import PrettyStatus from "@/components/status";

export default function Task({ id, name, previewUrl, status }) {
  const router = useRouter();

  const click = useCallback(() => {
    router.push("/task/" + id);
  }, [id]);

  return (
    <Card onClick={click} style={{ cursor: "pointer", maxHeight: 300 }}>
      <CardImg style={{"objectFit": "contain", height: "100%", width: "100%"}} variant="top" src={previewUrl} />
      <CardBody>
        <CardTitle>{name}</CardTitle>
        <CardTitle>
          <PrettyStatus status={status} />
        </CardTitle>
      </CardBody>
    </Card>
  );
}
