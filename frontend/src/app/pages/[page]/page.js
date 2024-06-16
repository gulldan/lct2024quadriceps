import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import { notFound } from "next/navigation";

import { getTasks } from "@/utils/web";

import Task from "@/components/task";
import Paginator from "@/components/paginator";

export default async function Home({ params }) {
  const per_page_limit = 12;

  const page = parseInt(params.page);
  console.debug("page", page);
  if (isNaN(page) || page < 1) {
    return notFound();
  }

  // try {
   
    let res = await getTasks(params.page, per_page_limit);
  // } catch (err) {
  //   return <div>{err.message}</div>;
  // }

  const total = parseInt(res["total"]);
  const tasks_preview = res["tasksPreview"];

  if (tasks_preview.length == 0) {
    return notFound();
  }

  return (
    <Container className="p-2">
      <Row className="justify-content-md-center">
        {tasks_preview.map((t, i) => (
          <Col key={i} className="p-2" id={i} xs={12} md={6} lg={4}>
            <Task {...t} />
          </Col>
        ))}
      </Row>
      <Row>
        <Col>
          <Paginator
            page={params.page}
            per_page={per_page_limit}
            total={total}
          />
        </Col>
      </Row>
    </Container>
  );
}
