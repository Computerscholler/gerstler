import { useEffect, useState } from "react";
import "./App.css";
import { useAppDispatch } from "./app/hooks";
import Results from "./components/Results";
import SearchBar from "./components/SearchBar";
import { updateResult } from "./features/resultsSlice";
import APIProvider from "./services/API";
import { FaCog } from "react-icons/fa";
import { CgCloseO } from "react-icons/cg";
import Configuration from "./components/Configuration";

function App() {
  const dispatch = useAppDispatch();
  // useEffect(() => {
  //   dispatch(
  //     updateResult({
  //       title: "lorem_ipsum.pdf",
  //       provider: "Google Drive",
  //       content: "",
  //       link: "http://example.com",
  //       matches: 2
  //     })
  //   )
  //   dispatch(
  //     updateResult({
  //       title: "test.pdf",
  //       provider: "Email",
  //       content: "",
  //       link: "",
  //       matches: 1
  //     })
  //   )
  //   dispatch(
  //     updateResult({
  //       parts: [
  //         {
  //           content:
  //             "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris tellus tellus, elementum eget convallis vel, commodo vitae tortor. Morbi in leo at",
  //           highlight: false,
  //           index: 0,
  //         },
  //         { content: "justo lacinia congue", highlight: true, index: 1 },
  //         {
  //           content:
  //             "at ut urna. Vivamus eget sapien eget nisi semper congue vitae ut nisl.",
  //           highlight: false,
  //           index: 2,
  //         },
  //         {
  //           content: "Sed imperdiet diam",
  //           highlight: true,
  //           index: 3,
  //         },
  //         {
  //           content:
  //             "in lectus ultrices, in aliquet magna elementum. Vestibulum euismod magna scelerisque diam semper, eget facilisis urna pharetra. Quisque molestie dignissim consectetur. Integer porta tortor et vehicula malesuada. Nam enim leo, rutrum ullamcorper scelerisque elementum, hendrerit et leo. Suspendisse placerat sapien quis vehicula mollis. Ut vulputate magna vitae lorem aliquam suscipit. Maecenas ultrices sem et tortor lobortis lacinia. Mauris pretium felis id lacus convallis, non dapibus arcu maximus. Cras nec ultricies nulla. Duis posuere sem blandit sem semper, quis imperdiet erat blandit.",
  //           highlight: false,
  //           index: 4,
  //         },
  //       ],
  //       provider: "Notion.so",
  //       title: "Notes",
  //       link: "http://example.com"
  //     })
  //   );
  // });

  const [showConfig, setShowConfig] = useState(false);

  return (
    <APIProvider>
      <div className="w-screen h-full min-h-screen font-mono dark:bg-gray-900 overflow-x-hidden">
        <div className="fixed w-full top-0 dark:bg-gray-900 shadow-xl h-24">
          <div className="flex ml-6 md:m-0 md:justify-center h-full">
            <SearchBar />
          </div>
          <div className="fixed top-0 h-24 right-4 md:right-8">
            <div className="h-full flex align-center ">
              {showConfig ? (
                <CgCloseO
                  onClick={() => setShowConfig(false)}
                  className="text-4xl cursor-pointer dark:text-gray-600 my-auto transition-all duration-200"
                />
              ) : (
                <FaCog
                  onClick={() => setShowConfig(true)}
                  className="text-4xl cursor-pointer dark:text-gray-600 my-auto transition-all duration-200"
                />
              )}
            </div>
          </div>
        </div>
        <div className="mt-40">{showConfig ? <Configuration /> : <Results />}</div>
      </div>
    </APIProvider>
  );
}

export default App;
