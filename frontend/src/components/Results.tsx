import { useState } from "react";
import { useAppSelector } from "../app/hooks";
import ResultCard from "./ResultCard";
import RichResultCard from "./RichResultCard";
import { isHtmlResult, isRichResult } from "../app/store";
import Loading from "./Loading";
import HtmlResultCard from "./HtmlResultCard";

function Results() {
  const results = useAppSelector((state) => state.results.value);
  const query = useAppSelector((state) => state.query.value);
  const loading = useAppSelector((state) => state.loading.value);

  const [extended, setExtended] = useState(-1);

  return (
    <div>
      {loading ? (
        <Loading />
      ) : results.length === 0 ? (
        query.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-500 text-4xl m-10 font-semibold text-center">
            Your results will appear here...
          </p>
        ) : (
          <p className="text-gray-600 dark:text-gray-500 text-4xl m-10 font-semibold text-center">
            No results
          </p>
        )
      ) : null}
      <div className="flex flex-wrap justify-center overflow-x-hidden max-w-screen md:mx-10 transition-all duration-500">
        {results.map((result, i) =>
          isRichResult(result) ? (
            <RichResultCard
              key={i}
              result={result}
              extended={i === extended}
              callback={() => {
                i === extended ? setExtended(-1) : setExtended(i);
              }}
            />
          ) : isHtmlResult(result) ? (
            <HtmlResultCard key={i} result={result} extended={i === extended} callback={() => {i === extended ? setExtended(-1) : setExtended(i);}} />
          ) : (
            <ResultCard key={i} result={result} />
          )
        )}
      </div>
    </div>
  );
}

export default Results;
