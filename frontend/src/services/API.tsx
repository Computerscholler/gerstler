import React, { createContext, useMemo } from "react";
import axios, { AxiosInstance } from "axios";

import { clearResults, updateResults } from "../features/resultsSlice";
import { setEmailLoading, setFilesystemLoading, setGoogleDriveLoading, setLoading, setNotionLoading } from "../features/loadingSlice";
import { ConfigurationState } from "../features/configurationSlice";
import { useAppDispatch } from "../app/hooks";
import { store } from "../app/store";
export interface API {
  ws: WebSocket;
  api: AxiosInstance;
}

export function createAPI(dispatch: any) {
  var base_url: string;
  if (process.env.REACT_APP_GERSTLER_PORT) {
    base_url = `${window.location.hostname}:${process.env.REACT_APP_GERSTLER_PORT}`
  } else {
    if (process.env.NODE_ENV === "development") {
      base_url = "localhost:5000";
    } else {
      if (window.location.port === "80") {
        base_url = `${window.location.hostname}:80`;
      } else {
        base_url = `${window.location.hostname}:5000`;
      }
    }
  }

  const secure = window.location.protocol === "https" ? true : false;

  let api = axios.create({
    baseURL: `${window.location.protocol}://${base_url}/api/`,
  });

  const connect = () => {
    let ws = new WebSocket(`${secure ? "wss" : "ws"}://${base_url}/api/ws`);

    ws.onmessage = (message) => {
      let obj = JSON.parse(message.data);
      if ("action" in obj) {
        switch (obj.action) {
          case "results": {
            const query = store.getState().query.value;
            if (obj.query === query) {
              if (obj.data) {
                dispatch(updateResults(obj.data));
              }
            }
            break;
          }
          case "loading_status": {
            const query = store.getState().query.value;
            if(obj.query === query) {
              switch (obj.data.provider.toLowerCase()) {
                case "email": dispatch(setEmailLoading(obj.data.loading)); break;
                case "notion": dispatch(setNotionLoading(obj.data.loading)); break;
                case "gdrive": dispatch(setGoogleDriveLoading(obj.data.loading)); break;
                case "filesystem": dispatch(setFilesystemLoading(obj.data.loading)); break;
              }
            }
          }
        }
      } else {
        console.error("Invalid message %s", message);
      }
    };

    ws.onclose = () => {
      setTimeout(() => {
        connect()
      }, 1000)
    }

    return ws;
  };

  let ws = connect();

  return { ws, api };
}

export function send(api: API, payload: string) {
  if (
    api.ws.readyState === api.ws.CLOSED ||
    api.ws.readyState === api.ws.CLOSING
  ) {
    console.error("Websocket connection not active");
    // TODO: retry
  } else if (api.ws.readyState === api.ws.CONNECTING) {
      setTimeout(() => send(api, payload), 1000)
  } else {
    api.ws.send(payload);
  }
}

export function updateConfig(api: API, config: ConfigurationState) {}

export function getConfig(api: API) {
  axios.get("config").then((response) => {});
}

export function sendQuery(api: API, query: String) {
  send(
    api,
    JSON.stringify({
      action: "query",
      query: query,
    })
  );
}

export const APIContext = createContext<API>({} as API);

export default function APIProvider({ children }: { children: any }) {
  const dispatch = useAppDispatch();
  const api = useMemo(() => {
    return createAPI(dispatch);
  }, []);
  return <APIContext.Provider value={api}>{children}</APIContext.Provider>;
}
