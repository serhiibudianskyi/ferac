const apiURL = import.meta.env.VITE_API_COSMOS ?? "https://budianskyi-s.com/ferac";
const rpcURL = import.meta.env.VITE_WS_TENDERMINT ?? "https://budianskyi-s.com/ferac/";
const prefix = import.meta.env.VITE_ADDRESS_PREFIX ?? "f8-";

export const env = {
  apiURL,
  rpcURL,
  prefix,
};
