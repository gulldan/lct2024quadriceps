import { createContext } from "react";


export const PageContext = createContext({
  data: null,
  setData: () => {}
});
