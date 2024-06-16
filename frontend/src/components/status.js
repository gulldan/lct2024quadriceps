import { Spinner, Badge } from "react-bootstrap";

export default function PrettyStatus({status}) {
    let status_pretty = "Неизвестный статус";

  switch (status) {
    case "TASK_STATUS_UNSPECIFIED":
      status_pretty = (
        <Badge bg="warning" text="dark">
          Неопределенный статус
        </Badge>
      );
      break;
    case "TASK_STATUS_FAIL":
      status_pretty = <Badge bg="danger">Анализ сорвался</Badge>;
      break;
    case "TASK_STATUS_IN_PROGRESS":
      status_pretty = (
        <>
          Идет анализ.. <Spinner size="sm" />
        </>
      );
      break;
    case "TASK_STATUS_DONE":
      status_pretty = <Badge bg="success">Анализ завершен</Badge>;
      break;
  }


  return <>{status_pretty}</>
}