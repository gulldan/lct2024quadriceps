"use client";
import { useMemo, useState, useEffect, useContext } from "react";

import { Col, Container, Image, Row } from "react-bootstrap";
import { PageContext } from "./context";
import { Player } from "./player";
import { CurTrack } from "./cur_track";
import TaskResult from "./task_result";
import PrettyStatus from "@/components/status";
import { getTask } from "@/utils/web";

function ActiveWait() {
  const { data, setData } = useContext(PageContext);

  useEffect(() => {
    const interval = setInterval(() => {
      getTask(data.id).then((t) => {
        if (t.status != t.status) {
          setData(t);
        }
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex flex-col justify-evenly md:flex-row">
      <Image
        className="self-center lg:w-[20%] md:w-[20%] w-[40%]"
        src="/ani2.gif"
      />
    </div>
  );
}
function Page({ init }) {
  const [data, setData] = useState(init);
  const value = useMemo(() => ({ data, setData }), [data]);

  return (
    <PageContext.Provider value={value}>
      <Container fluid>
        <Row>
          <Col className="py-2 text-center">
            <div className="p-2">
              <PrettyStatus status={data.status} />
            </div>

            {data.status == "TASK_STATUS_DONE" && (
              <TaskResult is_good={data.copyright.length == 0} />
            )}

            {data.status == "TASK_STATUS_IN_PROGRESS" && (
              <ActiveWait />
            )}
          </Col>
        </Row>
        <Row className="py-2  ">
          <Col xs={12} lg={data.copyright.length == 0 ? 12 : 6}>
            <div className="text-center">
              <h2>{data.name}</h2>
            </div>
            <Player />
          </Col>
          <Col xs={12} lg={6}>
            <CurTrack />
          </Col>
        </Row>
      </Container>
    </PageContext.Provider>
  );
}

export default Page;
