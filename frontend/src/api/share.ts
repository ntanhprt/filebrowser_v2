import { fetchURL, fetchJSON, removePrefix, createURL } from "./utils";

export async function list() {
  return fetchJSON<Share[]>("/api/shares");
}

export async function get(url: string) {
  url = removePrefix(url);
  return fetchJSON<Share>(`/api/share${url}`);
}

export async function remove(hash: string) {
  await fetchURL(`/api/share/${hash}`, {
    method: "DELETE",
  });
}

export async function create(
  url: string,
  password = "",
  expires = "",
  unit = "hours"
) {
  url = removePrefix(url);
  url = `/api/share${url}`;
  if (expires !== "") {
    url += `?expires=${expires}&unit=${unit}`;
  }
  let body = "{}";
  if (password != "" || expires !== "" || unit !== "hours") {
    body = JSON.stringify({
      password: password,
      expires: expires.toString(), // backend expects string not number
      unit: unit,
    });
  }
  return fetchJSON(url, {
    method: "POST",
    body: body,
  });
}

export async function listPublic() {
  return fetchJSON<Share[]>("/api/publicshares");
}

export async function makePublic(
  hash: string,
  permission: "view" | "change",
  allowedUsers: number[] = []
) {
  return fetchJSON(`/api/share/${hash}/public`, {
    method: "PUT",
    body: JSON.stringify({
      permission: permission,
      allowed_users: allowedUsers,
    }),
  });
}

export async function makePrivate(hash: string) {
  return fetchJSON(`/api/share/${hash}/private`, {
    method: "PUT",
    body: JSON.stringify({}),
  });
}

export function getShareURL(share: Share) {
  return createURL("share/" + share.hash, {});
}
