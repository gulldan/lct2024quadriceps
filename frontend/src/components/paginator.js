"use client";
import Pagination from "react-bootstrap/Pagination";

function Firts({ page }) {
  if (page > 3) {
    return (
      <>
        <Pagination.First href="/pages/1" />
        <Pagination.Item href="/pages/1">{1}</Pagination.Item>
        <Pagination.Ellipsis />
      </>
    );
  } else {
    return <></>;
  }
}

function End({ page, per_page, total  }) {
  const total_pages = Math.ceil(total / per_page);


  if (parseInt(page) <= total_pages - 3) {
    return (
      <>
         <Pagination.Ellipsis />
        <Pagination.Item href={`/pages/${total_pages}`}>{total_pages}</Pagination.Item>
        <Pagination.Last href={`/pages/${total_pages}`} />
       
      </>
    );
  } else {
    return <></>;
  }
}


function Prev({ page }) {
  if (page > 1) {
    return <Pagination.Prev href={`/pages/${page - 1}`} />;
  }
  return <></>;
}

function Next({ page, per_page, total }) {
  const total_pages = Math.ceil(total / per_page);
  if (page < total_pages) {
    return <Pagination.Next href={`/pages/${parseInt(page) + 1}`} />;
  }
  return <></>;
}

function Cur({ page, per_page, total }) {
  const pa = parseInt(page);

  const total_pages = Math.ceil(total / per_page);
  let items = [];

  for (let i = 1; i <= 2; i++)
    if (pa - i >= 1) {
      items = [pa - i, ...items];
    }

  items.push(pa);

  for (let i = 1; i <= 2; i++)
    if (pa + i <= total_pages) {
      items = [...items, pa + i];
    }

  
  return items.map((p, i) => (
    <Pagination.Item key={i} href={`/pages/${p}`} active={p == page}>
      {p}
    </Pagination.Item>
  ));
}

export default function Paginator({ page, per_page, total }) {
  return (
    <Pagination>
      <Firts page={page} />
      <Prev page={page} />
      <Cur page={page} per_page={per_page} total={total} />
      <Next page={page} per_page={per_page} total={total} />
      <End page={page} per_page={per_page} total={total} />
    </Pagination>
  );
}
