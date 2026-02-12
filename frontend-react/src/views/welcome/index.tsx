import { FC, useState, useEffect } from 'react';
import { Link } from 'react-router';

const texts = [
  "Ready to explore further?",
  "Keep growing and learning",
  "With us, let's grow faster"
] as const;

const typingSpeed = 80;
const pauseDuration = 1000;

const Welcome: FC = () => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [displayedText, setDisplayedText] = useState("");
  const [isTyping, setIsTyping] = useState(true);

  const currentText = texts[currentIndex];

  useEffect(() => {
    let timeout: any;
    if (isTyping) {
      if (displayedText.length < currentText.length) {
        timeout = setTimeout(() => {
          setDisplayedText(currentText.slice(0, displayedText.length + 1));
        }, typingSpeed);
      } else {
        timeout = setTimeout(() => {
          setIsTyping(false);
        }, pauseDuration);
      }
    } else {
      if (displayedText.length > 0) {
        timeout = setTimeout(() => {
          setDisplayedText(displayedText.slice(0, -1));
        }, typingSpeed / 1.5);
      } else {
        setCurrentIndex((prev) => (prev + 1) % texts.length);
        setIsTyping(true);
      }
    }
    return () => clearTimeout(timeout);
  }, [displayedText, isTyping, currentIndex]);

  useEffect(() => {
    if (isTyping && displayedText === "") {
      setDisplayedText(currentText.slice(0, 1));
    }
  }, [currentIndex, isTyping]);

  return (
    <>
      {/* CSS */}
      <style>{`
        @keyframes gradientMove {
          0%, 100% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
        }
        /* Untuk MSRooot */
        .animated-text-gradient {
          background-image: linear-gradient(
            90deg,
            #6366f1, /* indigo-500 */
            #a855f7, /* purple-500 */
            #7c3aed, /* indigo-600 */
            #6366f1
          );
          background-size: 200% 100%;
          -webkit-background-clip: text;
          background-clip: text;
          color: transparent;
          animation: gradientMove 7s ease infinite;
        }
        /* Untuk button Try Now */
        .animated-gradient-btn {
          background-image: linear-gradient(
            to right,
            #eab30888, /* yellow-500 semi-transparan */
            #d9770688,
            #fbbf2488,
            #eab30888
          );
          background-size: 200% 200%;
          animation: gradientMove 6s ease infinite;
          backdrop-filter: blur(4px); /* efek kaca tipis, optional */
        }
        .animated-gradient-btn:hover {
          background-image: linear-gradient(
            to right,
            #facc1588,
            #fbbf2488,
            #fcd34d88,
            #facc1588
          );
          transform: scale(1.05);
        }
      `}</style>

      <div className="min-h-svh w-full bg-gradient-to-br from-white via-gray-300 to-white dark:from-black dark:via-gray-950 dark:to-black flex flex-col justify-between px-6 sm:px-12 md:px-16 lg:px-36 relative">
        <div className="flex-1 flex items-center">
          <div className="text-left space-y-2 max-w-2xl z-10">
            <h1 className="text-5xl sm:text-6xl md:text-7xl font-extrabold tracking-tight whitespace-nowrap">
              <span className="text-gray-600 dark:text-white">Welcome to </span>
              <span className="animated-text-gradient">
                MSRooot.ID
              </span>
            </h1>
            <p className="text-gray-600 dark:text-gray-400 text-xl sm:text-2xl font-light tracking-wide max-w-lg min-h-[3rem] font-mono">
              {displayedText}
              <span className="animate-pulse">|</span>
            </p>
            <div className="pt-8">
              <Link
                to="/home"
                className={`
                  animated-gradient-btn
                  inline-flex items-center justify-center
                  px-10 py-4 text-xl font-semibold
                  text-gray-600 hover:text-gray-700 dark:text-gray-300  hover:dark:text-gray-200
                  rounded-full
                  shadow-lg shadow-yellow-900/20
                  transition-all duration-300
                  hover:shadow-xl hover:shadow-yellow-800/30
                  active:scale-95
                  overflow-hidden
                  border border-yellow-500/30
                `}
              >
                Try Now
              </Link>
            </div>
          </div>
        </div>

        {/* Footer */}
        <footer className="py-6 text-center text-gray-600 text-sm z-10">
          <p>Version 1.0.0</p>
          <p>© {new Date().getFullYear()} MSRooot. Made with React.JS and Golang .</p>
        </footer>

        {/* Glow effect */}
        <div className="absolute inset-0 pointer-events-none overflow-hidden">
          <div className="absolute -top-20 -right-20 sm:-right-40 w-64 h-64 sm:w-96 sm:h-96 bg-indigo-900/10 rounded-full blur-3xl"></div>
          <div className="absolute -bottom-20 -left-20 sm:-left-40 w-64 h-64 sm:w-96 sm:h-96 bg-purple-900/10 rounded-full blur-3xl"></div>
        </div>
      </div>
    </>
  );
};

export default Welcome;