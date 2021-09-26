import {
  ChangeEvent,
  useEffect,
  useState,
  MouseEvent,
  useContext,
} from "react";
import { MdClear } from "react-icons/md";
import { useDebounce } from "use-debounce/lib";
import { useAppDispatch, useAppSelector } from "../app/hooks";
import { setLoading } from "../features/loadingSlice";
import { clearQuery, updateQuery } from "../features/querySlice";
import { clearResults } from "../features/resultsSlice";

import { APIContext, sendQuery } from "../services/API";

function SearchBar() {
  const globalQuery = useAppSelector((state) => state.query.value);
  const [query, setQuery] = useState(globalQuery);
  const [lastUpdate, setLastUpdate] = useState(new Date());
  const dispatch = useAppDispatch();

  const api = useContext(APIContext);

  const clear = () => {
    setQuery("");
    dispatch(clearQuery())
    dispatch(clearResults())
    dispatch(setLoading(false));
  }

  const [debouncedQuery] = useDebounce(query, 300)

  useEffect(() => {
    if (debouncedQuery.length > 0) {
      dispatch(clearResults())
      dispatch(updateQuery(debouncedQuery))
      sendQuery(api, debouncedQuery);
    }
  }, [debouncedQuery]);

  const onChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (event.target.value) {
      setQuery(event.target.value);
      dispatch(updateQuery(event.target.value));
    } else {
      clear(); 
    }
  };

  const onClick = (event: MouseEvent<HTMLElement>) => {
    event.preventDefault();
    clear();
  };

  return (
    <div className="my-auto w-4/5 lg:w-1/3 flex align-center justify-center border-2 rounded-xl shadow-md p-2 pr-4 pl-6 focus:outline-none  border-gray-300 focus:border-gray-500 focus:shadow-xl dark:bg-gray-700 dark:border-gray-900 dark:focus:border-gray-600 dark:text-gray-200">
      <input
        type="text"
        autoFocus
        className="w-full text-gray-800 md:text-center text-sm md:text-xl dark:bg-gray-700 dark:text-gray-200 focus:outline-none "
        value={query}
        onChange={onChange}
        placeholder="Type something to start..."
      />
      <span className="dark:text-gray-400 my-auto" onClick={onClick}>
        <MdClear />
      </span>
    </div>
  );
}

export default SearchBar;
