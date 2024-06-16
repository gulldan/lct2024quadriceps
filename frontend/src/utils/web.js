function root() {
  const base = `http://bff:8888`
  if (typeof window === "undefined") {
    return base;
  }
  return "";
}

export async function upload(file) {
  let form = new FormData();
  form.append("file", file);
  

  let res = await fetch(`${root()}/api/v1/original/upload`, {
    method: "POST",
    cache: "no-cache",
    body: form,
  });
  let data = await res.json();

  if (res.status != 200) {
    throw new Error(data["message"]);
  }
  return data;
}

export async function create_task(file) {
  let form = new FormData();
  form.append("file", file);
  
  let res = await fetch(`${root()}/api/v1/tasks/create/upload`, {
    method: "POST",
    cache: "no-cache",
    body: form,
  });
  let data = await res.json();

  if (res.status != 200) {
    throw new Error(data["message"]);
  }
  return data;
}

export async function getTask(id) {
  let res = await fetch(`${root()}/api/v1/tasks/${id}`, {
    method: "GET",
    cache: "no-cache",
    "Content-Type": "multipart/form-data",
    accept: "application/json",
  });
  let data = await res.json();
  if (res.status != 200) {
    throw new Error(data["message"]);
  }
  return data;
}

export async function getTasks(page, limit) {
  let res = await fetch(
    `${root()}/api/v1/tasks_preview?page=${page - 1}&limit=${limit}`,
    {
      method: "GET",
      cache: "no-cache",
      "Content-Type": "multipart/form-data",
      accept: "application/json",
    }
  );

  let data = await res.json();
  console.log(data);
  if (res.status != 200) {
    throw new Error(data["message"]);
  }
  return data;
}
