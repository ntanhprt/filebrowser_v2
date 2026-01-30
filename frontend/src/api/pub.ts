import { fetchURL, removePrefix, createURL } from "./utils";
import { baseURL } from "@/utils/constants";

export interface ShareInfo {
  hash: string;
  path: string;
  isPublic: boolean;
  permission: string;
  canEdit: boolean;
}

export async function fetch(url: string, password: string = "") {
  url = removePrefix(url);

  const res = await fetchURL(
    `/api/public/share${url}`,
    {
      headers: { "X-SHARE-PASSWORD": encodeURIComponent(password) },
    },
    false
  );

  const data = (await res.json()) as Resource & { shareInfo?: ShareInfo };
  data.url = `/share${url}`;

  if (data.isDir) {
    if (!data.url.endsWith("/")) data.url += "/";
    data.items = data.items.map((item: any, index: any) => {
      item.index = index;
      item.url = `${data.url}${encodeURIComponent(item.name)}`;

      if (item.isDir) {
        item.url += "/";
      }

      return item;
    });
  }

  // Include share info in the response
  if (data.shareInfo) {
    (data as any).canEdit = data.shareInfo.canEdit;
    (data as any).sharePermission = data.shareInfo.permission;
  }

  return data;
}

export function download(
  format: DownloadFormat,
  hash: string,
  token: string,
  ...files: string[]
) {
  let url = `${baseURL}/api/public/dl/${hash}`;

  if (files.length === 1) {
    url += encodeURIComponent(files[0]) + "?";
  } else {
    let arg = "";

    for (const file of files) {
      arg += encodeURIComponent(file) + ",";
    }

    arg = arg.substring(0, arg.length - 1);
    arg = encodeURIComponent(arg);
    url += `/?files=${arg}&`;
  }

  if (format) {
    url += `algo=${format}&`;
  }

  if (token) {
    url += `token=${token}&`;
  }

  window.open(url);
}

export function getDownloadURL(res: Resource, inline = false) {
  const params = {
    ...(inline && { inline: "true" }),
    ...(res.token && { token: res.token }),
  };

  return createURL("api/public/dl/" + res.hash + res.path, params);
}

// Upload file to shared folder (requires authentication and edit permission)
export async function upload(
  hash: string,
  filePath: string,
  content: Blob | string,
  override = false
) {
  const url = `/api/public/share/${hash}${filePath}${override ? "?override=true" : ""}`;
  
  const res = await fetchURL(url, {
    method: "POST",
    body: content,
  });

  return res;
}

// Create folder in shared folder (requires authentication and edit permission)
export async function createFolder(hash: string, folderPath: string) {
  // Ensure path ends with /
  const path = folderPath.endsWith("/") ? folderPath : folderPath + "/";
  const url = `/api/public/share/${hash}${path}`;
  
  const res = await fetchURL(url, {
    method: "POST",
  });

  return res;
}

// Delete file/folder from shared folder (requires authentication and edit permission)
export async function remove(hash: string, filePath: string) {
  const url = `/api/public/share/${hash}${filePath}`;
  
  const res = await fetchURL(url, {
    method: "DELETE",
  });

  return res;
}

// Rename/move file in shared folder (requires authentication and edit permission)
export async function rename(
  hash: string,
  srcPath: string,
  dstPath: string
) {
  const url = `/api/public/share/${hash}${srcPath}?action=rename&destination=${encodeURIComponent(dstPath)}`;
  
  const res = await fetchURL(url, {
    method: "PATCH",
  });

  return res;
}

// Copy file in shared folder (requires authentication and edit permission)
export async function copy(
  hash: string,
  srcPath: string,
  dstPath: string
) {
  const url = `/api/public/share/${hash}${srcPath}?action=copy&destination=${encodeURIComponent(dstPath)}`;
  
  const res = await fetchURL(url, {
    method: "PATCH",
  });

  return res;
}

// Update file content in shared folder (requires authentication and edit permission)
// Uses PUT method to update existing file content
export async function save(
  hash: string,
  filePath: string,
  content: string
) {
  const url = `/api/public/share/${hash}${filePath}`;
  
  const res = await fetchURL(url, {
    method: "PUT",
    body: content,
  });

  return res;
}
