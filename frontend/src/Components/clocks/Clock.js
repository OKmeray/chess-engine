import React, { useEffect, useRef } from "react";

const Clocks = ({
    whiteTime,
    blackTime,
    turnToPlay,
    gameOver,
    setWhiteTime,
    setBlackTime,
    onTimeUpWhite,
    onTimeUpBlack,
    playerColor,
}) => {
    const intervalRef = useRef(null);

    useEffect(() => {
        // Clear any old interval whenever turn changes or gameOver changes
        clearInterval(intervalRef.current);

        if (gameOver) {
            // If the game is over, do NOT start any new timer
            return;
        }

        if (turnToPlay === "User") {
            // Start White’s clock
            intervalRef.current = setInterval(() => {
                setWhiteTime((prev) => {
                    if (prev <= 1) {
                        clearInterval(intervalRef.current);
                        onTimeUpWhite?.();
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        } else if (turnToPlay === "Computer") {
            // Start Black’s clock
            intervalRef.current = setInterval(() => {
                setBlackTime((prev) => {
                    if (prev <= 1) {
                        clearInterval(intervalRef.current);
                        onTimeUpBlack?.();
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }

        // Cleanup on unmount or re-run
        return () => clearInterval(intervalRef.current);
    }, [
        turnToPlay,
        gameOver,
        setWhiteTime,
        setBlackTime,
        onTimeUpWhite,
        onTimeUpBlack,
    ]);

    const formatTime = (seconds) => {
        const minutes = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${minutes}:${secs < 10 ? "0" : ""}${secs}`;
    };

    return (
        <div className="clocks-container">
            <div
                className={`clock ${
                    turnToPlay === "Computer" && !gameOver ? "active" : ""
                }`}
            >
                <h2>{playerColor === "white" ? "Час чорних" : "Час білих"}</h2>
                <p>{formatTime(blackTime)}</p>
            </div>
            <div
                className={`clock ${
                    turnToPlay === "User" && !gameOver ? "active" : ""
                }`}
            >
                <h2>{playerColor === "white" ? "Час білих" : "Час чорних"}</h2>
                <p>{formatTime(whiteTime)}</p>
            </div>
        </div>
    );
};

export default Clocks;
