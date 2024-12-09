const BASE_URL = "http://localhost:9091";

export const API_ENDPOINTS = {
  CHECK_FILE_HASH: `${BASE_URL}/checkFileHash`,
  UPLOAD_FILES: `${BASE_URL}/upload`,
  GET_ALL_FILES: `${BASE_URL}/getAllFiles`,
  DOWNLOAD_FILES: `${BASE_URL}/download`,
};

export function setBaseUrl(newBaseUrl) {
  Object.keys(API_ENDPOINTS).forEach(key => {
    API_ENDPOINTS[key] = API_ENDPOINTS[key].replace(BASE_URL, newBaseUrl);
  });
}