import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Result, RichResult } from "../app/store";

export const resultsSlice = createSlice({
  name: "results",
  initialState: {
    value: new Array<Result|RichResult>(),
  },
  reducers: {
    clearResults: (state) => {
      state.value = []
    },
    updateResult: (state, action: PayloadAction<Result|RichResult>) => {
      state.value.push(action.payload);
    },
    updateResults: (state, action: PayloadAction<Result|RichResult>) => {
      state.value = state.value.concat(action.payload)
    },
    replaceResults: (state, action: PayloadAction<Array<Result|RichResult>>) => {
      state.value = action.payload
    }
  },
});

export const { clearResults, updateResult, updateResults, replaceResults } = resultsSlice.actions;

export default resultsSlice.reducer;
