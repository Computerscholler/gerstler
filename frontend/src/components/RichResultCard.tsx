import { MouseEvent, useState } from "react";
import { Part, RichResult } from "../app/store";

interface RichResultCardProps {
  result: RichResult;
  extended: boolean;
  callback: Function;
}

function PartComponent({ part }: { part: Part }) {
  return part.highlight ? (
    <span
      className="bg-yellow-700 rounded-md bg-opacity-60 px-1"
      style={{ paddingBottom: "0.08rem", paddingTop: "0.08rem" }}
    >
      {part.content}
    </span>
  ) : (
    <span className="">{part.content}</span>
  );
}
function RichResultCard(props: RichResultCardProps) {
  const { result, extended, callback } = props;
  const [click, setClick] = useState(false);

  const mouseDown = (event: MouseEvent) => {
    setClick(true);
    setTimeout(() => setClick(false), 200);
  };

  const mouseUp = (event: MouseEvent) => {
    if (click) {
      callback();
      setClick(false);
    }
  };

  const lookaround = 10;

  const parts = result.parts.map((part: Part, i) => {
    if (part.highlight) {
      return part;
    } else if (
      i === 0 &&
      i + 1 < result.parts.length &&
      result.parts[i + 1].highlight
    ) {
      if (part.content.length - lookaround <= 0) {
        return {
          content: part.content,
          highlight: false,
          index: part.index,
        };
      } else {
        return {
          content:
            "..." +
            part.content
              .substring(part.content.length - lookaround)
              .trimStart(),
          highlight: false,
          index: part.index,
        };
      }
    } else if (i + 1 === result.parts.length) {
      return {
        content:
          part.content
            .substring(0, Math.min(lookaround, part.content.length))
            .trimEnd() + "...",
        highlight: false,
        index: part.index,
      };
    } else {
      if (part.content.length < lookaround * 2 + 3) {
        return {
          content: part.content,
          highlight: false,
          index: part.index,
        };
      }
      return {
        content:
          part.content.substring(0, lookaround).trimEnd() +
          "..." +
          part.content.substring(part.content.length - lookaround).trimStart(),
        highlight: false,
        index: part.index,
      };
    }
  });

  return (
    <div
      className={`rounded-md p-4 transition-all duration-300 dark:bg-gray-700 m-4 dark:text-gray-100 divide-gray-200 ${
        extended ? "max-w-full md:max-w-6xl" : "max-w-lg"
      } shadow-lg`}
    >
      <p
        className="py-2 text-center font-semibold"
        onMouseUp={mouseUp}
        onMouseDown={mouseDown}
      >
        {result.title}
      </p>
      <div
        className={`py-3 font-thin transition-all duration-300 text-justify rounded-md dark:bg-gray-600 px-3 ${
          extended ? "text-lg" : "text-xs"
        }`}
        onMouseUp={mouseUp}
        onMouseDown={mouseDown}
      >
        {extended
          ? result.parts.map((part, i) => <PartComponent part={part} key={i} />)
          : parts.map((part, i) => <PartComponent part={part} key={i} />)}
      </div>
      <p className="text-sm  text-right pt-2 pb-1 underline">
        <a
          className="hover:font-semibold"
          href={result.link}
          rel="noopener noreferrer"
          target="_blank"
        >
          {result.provider}
        </a>
      </p>
    </div>
  );
}

export default RichResultCard;
