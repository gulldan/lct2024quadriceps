function root() {
  const base = `http://bff:8888`;
  if (typeof window === "undefined") {
    return base;
  }
  return "";
}

// от сервера не всегда приходит json))
// return json \ or raise error
async function GetJsonOrDie(res) {
  let data = await res.text();
  let jdata = {};

  try {
    jdata = JSON.parse(data);
  } catch (e) {
    throw new Error(data);
  }
  if (res.status != 200) {
    throw new Error(jdata["message"]);
  }
  return jdata;
}

export async function upload(file) {
  let form = new FormData();
  form.append("file", file);

  let res = await fetch(`${root()}/api/v1/original/upload`, {
    method: "POST",
    cache: "no-store",
    body: form,
  });
  return await GetJsonOrDie(res);
}

export async function create_task(file) {
  let form = new FormData();
  form.append("file", file);

  let res = await fetch(`${root()}/api/v1/tasks/create/upload`, {
    method: "POST",
    cache: "no-store",
    body: form,
  });
  return await GetJsonOrDie(res);
}

export async function getTask(id) {
  let res = await fetch(`${root()}/api/v1/tasks/${id}`, {
    method: "GET",
    cache: "no-store",
    "Content-Type": "multipart/form-data",
    accept: "application/json",
  });
  return await GetJsonOrDie(res);
}

export async function getTasks(page, limit) {
  let res = await fetch(
    `${root()}/api/v1/tasks_preview?page=${page - 1}&limit=${limit}`,
    {
      method: "GET",
      cache: "no-store",
      "Content-Type": "multipart/form-data",
      accept: "application/json",
    }
  );
  return await GetJsonOrDie(res);
}
