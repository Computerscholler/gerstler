import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface QueryState {
  value: string;
}

const initalState = { value: "" } as QueryState;

export const querySlice = createSlice({
  name: "query",
  initialState: initalState,
  reducers: {
    clearQuery: (state) => {
      state.value = "";
    },
    updateQuery: (state, action: PayloadAction<string>) => {
      state.value = action.payload;
    },
  },
});

export const { clearQuery, updateQuery } = querySlice.actions;

export default querySlice.reducer;
