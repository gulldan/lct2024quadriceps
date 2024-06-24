var incr = (function () {
  var i = 1;

  return function () {
    return i++;
  };
})();

export async function getTask(id) {
  const i = incr();

  let out = {
    id: id,
    name: `name name ${id}`,
    videoUrl: "/test/video2.mp4",
    status: i > 5 ? "TASK_STATUS_DONE" : "TASK_STATUS_IN_PROGRESS",
    copyright: [],
  };

  if (i > 5) {
    out.timestamps.push({
      start: 3,
      end: 6,
      orig_id: 2,
      origStart: 5,
      origEnd: 8,
    });
  }
  return out;
}

function randomTask(id) {
  const s = [
    "TASK_STATUS_UNSPECIFIED",
    "TASK_STATUS_FAIL",
    "TASK_STATUS_IN_PROGRESS",
    "TASK_STATUS_DONE",
  ];

  return {
    id: id,
    name: "stub",
    previewUrl: "/test/gool.png",
    status: s[id % s.length],
  };
}

export async function getTasks(page, size) {
  let out = [];

  for (let i = 1; i <= 64; i++) {
    out = [...out, randomTask(i)];
  }

  return {
    tasksPreview: out.slice((page - 1) * size, (page - 1) * size + size),
    total: 64,
  };
}
