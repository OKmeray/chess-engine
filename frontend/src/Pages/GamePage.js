import { useEffect, useState, useRef } from "react";
import { useLocation } from "react-router-dom";
import { useDispatch } from "react-redux";
import ChessBoard from "../Components/board/ChessBoard";
import Header from "../Components/Header";
import Footer from "../Components/Footer";
import useChessWebSocket from "../Hooks/useChessWebSocket";
import { PositionOutcome } from "../utils/outcomeEnum";
import { setCanNotDrag, setCanDrag } from "../Store/actions";
import Clocks from "../Components/clocks/Clock";
import GameResultPopup from "../Components/GameResultPopup";

import "./GamePage.css";

const GamePage = () => {
    const [fen, setFen] = useState(
        "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    );
    const [possibleMoves, setPossibleMoves] = useState([]);
    const [turnToPlay, setTurnToPlay] = useState("");
    const [gameOver, setGameOver] = useState(false);
    const [isBoardFlipped, setIsBoardFlipped] = useState(false);
    const [gameStarted, setGameStarted] = useState(false);
    const [playerColor, setPlayerColor] = useState("white");
    const [selectedModel, setSelectedModel] = useState("");
    const [gameResult, setGameResult] = useState(null);

    // Timers
    const [whiteTime, setWhiteTime] = useState(180);
    const [blackTime, setBlackTime] = useState(180);
    const [increment, setIncrement] = useState(2);
    const timerRef = useRef(null);
    const incrementRef = useRef(increment);

    const [timeControlType, setTimeControlType] = useState("fisher"); // 'fisher' or 'standard'
    const [baseTime, setBaseTime] = useState(180);
    const [selectedPreset, setSelectedPreset] = useState("custom");
    const [areTimeInputsDisabled, setAreTimeInputsDisabled] = useState(false);

    const timeControls = {
        bullet: { base: 60, increment: 1, name: "(1+1) Куля" },
        blitz: { base: 180, increment: 2, name: "(3+2) Блискавичний контроль" },
        rapid: { base: 600, increment: 5, name: "(10+5) Швидкий контроль" },
        classical: {
            base: 2100,
            increment: 30,
            name: "(35+30) Класичний контроль",
        },
        custom: {
            base: baseTime,
            increment: increment,
            name: "Індивідуальне налаштування",
        },
    };

    const location = useLocation();
    const dispatch = useDispatch();

    const handleMessage = (data) => {
        console.log(data);
        if (data.outcome === PositionOutcome.DRAW) {
            if (data.fen) setFen(data.fen);
            setFen(data.fen); // Update the board first
            setTimeout(() => {
                setGameOver(true);
                setGameResult("draw");
            }, 500);
        } else if (data.outcome === PositionOutcome.LOSS) {
            setFen(data.fen); // Update the board first
            setTimeout(() => {
                setGameOver(true);
                setGameResult("loss");
            }, 500);
        } else if (data.outcome === PositionOutcome.WIN) {
            setFen(data.fen); // Update the board first
            setTimeout(() => {
                setGameOver(true);
                setGameResult("win");
            }, 500);
        } else {
            dispatch(setCanDrag());
            setTurnToPlay("User");
            if (timeControlType === "fisher" && incrementRef.current > 0) {
                setBlackTime((prev) => prev + incrementRef.current);
            }
            setPossibleMoves(data.possibleMoves || []);
            setFen(data.fen);
        }
    };

    const handleTimeOut = (result) => {
        setGameOver(true);
        setGameResult(result === "loss" ? "timeoutWin" : "timeoutLoss"); // This will trigger the GameResultPopup
    };

    const resetGame = () => {
        setFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
        setPossibleMoves([]);
        setTurnToPlay("");
        setGameOver(false);
        setGameStarted(false);
        setGameResult(null);
        setWhiteTime(baseTime);
        setBlackTime(baseTime);
        clearInterval(timerRef.current);
    };

    const { sendMessage } = useChessWebSocket(
        `${process.env.REACT_APP_SOCKET_URL}move`,
        handleMessage
    );

    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const baseTime = parseInt(params.get("baseTime"), 10);
        const inc = parseInt(params.get("increment"), 10);
        const model = params.get("model");

        if (Number.isInteger(baseTime)) {
            setWhiteTime(baseTime);
            setBlackTime(baseTime);
            setBaseTime(baseTime);

            const matchedPreset = Object.keys(timeControls).find((key) => {
                return (
                    key !== "custom" &&
                    timeControls[key].base === baseTime &&
                    timeControls[key].increment === inc
                );
            });

            setSelectedPreset(matchedPreset || "custom");
            setAreTimeInputsDisabled(!!matchedPreset);
        }

        if (Number.isInteger(inc)) {
            setIncrement(inc);
        }
        if (model) {
            setSelectedModel(model);
        }
        return () => {
            clearInterval(timerRef.current);
        };
    }, []);

    useEffect(() => {
        fetch(`${process.env.REACT_APP_BACKEND_URL}Game?fen=${fen}`)
            .then((res) => {
                return res.json();
            })
            .then((res) => {
                setPossibleMoves(res["possibleMoves"]);
            });
    }, [location, fen]);

    useEffect(() => {
        if (isBoardFlipped && turnToPlay === "User") {
            dispatch(setCanNotDrag());
            const engineMove = {
                fen: fen,
                from: "0",
                to: "0",
                time: whiteTime,
                promotion: null,
                requestFirstMove: true,
                model: selectedModel,
            };
            sendMessage(engineMove);
        }
    }, [isBoardFlipped]);

    // --- Called when the user drops a piece on the board ---
    const handlePieceDrop = (fromSquare, toSquare, promotion = null) => {
        if (!gameStarted) return; // Prevent moves if game hasn't started

        dispatch(setCanNotDrag());
        setTurnToPlay("Computer");

        // setWhiteTime((prev) => prev + increment);
        if (timeControlType === "fisher" && incrementRef.current > 0) {
            setWhiteTime((prev) => prev + incrementRef.current);
        }

        const moveData = {
            fen,
            from: fromSquare,
            to: toSquare,
            time: blackTime,
            promotion: promotion, // e.g., 'q', 'r', 'b', 'n'
            model: selectedModel,
        };
        sendMessage(moveData);
        console.log("Sending move to server:", moveData);
    };

    const handlePresetSelect = (presetKey) => {
        const preset = timeControls[presetKey];
        setBaseTime(preset.base);
        setIncrement(preset.increment);
        setSelectedPreset(presetKey);
        setAreTimeInputsDisabled(presetKey !== "custom");
        if (!gameStarted) {
            setWhiteTime(preset.base);
            setBlackTime(preset.base);
        }
    };

    const startGame = () => {
        setGameStarted(true);
        setWhiteTime(baseTime);
        setBlackTime(baseTime);
        setIncrement(increment);
        setTurnToPlay(playerColor === "white" ? "User" : "Computer");

        if (playerColor === "black") {
            // Request engine's first move
            dispatch(setCanNotDrag());
            const engineMove = {
                fen: fen,
                from: "0",
                to: "0",
                time: whiteTime,
                promotion: null,
                requestFirstMove: true,
                model: selectedModel,
            };
            sendMessage(engineMove);
        } else {
            dispatch(setCanDrag());
        }
    };

    const handleColorSelection = (color) => {
        if (!gameStarted) {
            setPlayerColor(color);
            setIsBoardFlipped(color === "black");
        }
    };

    useEffect(() => {
        incrementRef.current = increment; // Keep ref in sync with state
    }, [increment]);

    return (
        <div>
            <Header />
            <section
                className={`game-section ${gameStarted ? "game-started" : ""}`}
            >
                <div className="container">
                    <div className="game-page-content">
                        <div className="game-details-wrapper">
                            <div className="game-settings-controls">
                                <div className="game-settings">
                                    <h3>Налаштування часу</h3>

                                    <div className="time-presets">
                                        {Object.keys(timeControls).map(
                                            (key) => (
                                                <button
                                                    key={key}
                                                    onClick={() =>
                                                        handlePresetSelect(key)
                                                    }
                                                    className={`base-button time-preset ${
                                                        selectedPreset === key
                                                            ? "selected"
                                                            : ""
                                                    }`}
                                                >
                                                    {timeControls[key].name}
                                                </button>
                                            )
                                        )}
                                    </div>

                                    <div
                                        className={`time-control-options ${
                                            areTimeInputsDisabled
                                                ? "disabled"
                                                : ""
                                        }`}
                                    >
                                        <div className="time-control-type">
                                            <label>
                                                <input
                                                    type="radio"
                                                    name="timeControlType"
                                                    checked={
                                                        timeControlType ===
                                                        "standard"
                                                    }
                                                    onChange={() =>
                                                        setTimeControlType(
                                                            "standard"
                                                        )
                                                    }
                                                />
                                                Стандартний
                                            </label>
                                            <label>
                                                <input
                                                    type="radio"
                                                    name="timeControlType"
                                                    checked={
                                                        timeControlType ===
                                                        "fisher"
                                                    }
                                                    onChange={() =>
                                                        setTimeControlType(
                                                            "fisher"
                                                        )
                                                    }
                                                />
                                                Фішера (з додаванням)
                                            </label>
                                        </div>
                                        <div className="time-inputs">
                                            <div className="time-input">
                                                <label>Час (в секундах):</label>
                                                <input
                                                    type="number"
                                                    value={baseTime}
                                                    onChange={(e) => {
                                                        const value =
                                                            parseInt(
                                                                e.target.value
                                                            ) || 0;
                                                        setBaseTime(value);
                                                        if (!gameStarted) {
                                                            setWhiteTime(value);
                                                            setBlackTime(value);
                                                        }
                                                    }}
                                                    min="30"
                                                    onBlur={(e) => {
                                                        // Ensure value is at least 30 when input loses focus
                                                        if (baseTime < 30) {
                                                            setBaseTime(30);
                                                            if (!gameStarted) {
                                                                setWhiteTime(
                                                                    30
                                                                );
                                                                setBlackTime(
                                                                    30
                                                                );
                                                            }
                                                        }
                                                    }}
                                                />
                                            </div>
                                            {timeControlType === "fisher" && (
                                                <div className="time-input">
                                                    <label>
                                                        Додавання (в сек.):
                                                    </label>
                                                    <input
                                                        type="number"
                                                        value={increment}
                                                        onChange={(e) =>
                                                            setIncrement(
                                                                parseInt(
                                                                    e.target
                                                                        .value
                                                                ) || 0
                                                            )
                                                        }
                                                        min="0"
                                                    />
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                                <div className="game-settings">
                                    <div className="game-start-controls">
                                        <h3>Виберіть колір:</h3>
                                        <div className="color-choice-wrapper">
                                            <button
                                                onClick={() =>
                                                    handleColorSelection(
                                                        "white"
                                                    )
                                                }
                                                className={`base-button color-select-button ${
                                                    playerColor === "white"
                                                        ? "selected"
                                                        : ""
                                                }`}
                                            >
                                                Грати за білих
                                            </button>
                                            <button
                                                onClick={() =>
                                                    handleColorSelection(
                                                        "black"
                                                    )
                                                }
                                                className={`base-button color-select-button ${
                                                    playerColor === "black"
                                                        ? "selected"
                                                        : ""
                                                }`}
                                            >
                                                Грати за чорних
                                            </button>
                                        </div>
                                    </div>
                                </div>
                                <div className="game-settings">
                                    <button
                                        onClick={() => {
                                            setFen(
                                                "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
                                            );
                                            setPossibleMoves([]);
                                        }}
                                        className="base-button hover-selected"
                                    >
                                        Початкова позиція
                                    </button>
                                    <button
                                        onClick={startGame}
                                        className="base-button start-game-button"
                                    >
                                        Почати гру
                                    </button>
                                </div>
                            </div>
                        </div>

                        <div className="chess-board-wrapper">
                            <div className="player-headings-wrapper">
                                <p className="player-heading">
                                    Шаховий рушій&nbsp;
                                    <span
                                        className="player-subheading"
                                        style={{
                                            opacity:
                                                turnToPlay === "Computer"
                                                    ? "100%"
                                                    : "0%",
                                        }}
                                    >
                                        (
                                        <div className="clarification-icon">
                                            <svg
                                                width="14"
                                                height="14"
                                                viewBox="0 0 18 18"
                                                fill="none"
                                                xmlns="http://www.w3.org/2000/svg"
                                            >
                                                <path
                                                    d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                    fill="black"
                                                />
                                                <circle
                                                    cx="9"
                                                    cy="9"
                                                    r="8.5"
                                                    stroke="black"
                                                />
                                            </svg>
                                            <span className="tooltip-text">
                                                Шаховий рушій думаю 5% від часу,
                                                який був після останнього
                                                зробленого ходу
                                            </span>
                                        </div>
                                        думає ...)
                                    </span>
                                </p>
                            </div>

                            {/* The board & pieces */}
                            <div className="chess-board">
                                <ChessBoard
                                    fen={fen}
                                    possibleMoves={possibleMoves}
                                    onPieceDrop={handlePieceDrop}
                                    isFlipped={isBoardFlipped}
                                    isDraggable={gameStarted}
                                />
                            </div>

                            <div className="player-headings-wrapper">
                                <p className="player-heading">
                                    Гість&nbsp;
                                    <span
                                        className="player-subheading"
                                        style={{
                                            opacity:
                                                turnToPlay === "User"
                                                    ? "100%"
                                                    : "0%",
                                        }}
                                    >
                                        (ваш хід)
                                    </span>
                                </p>
                            </div>
                        </div>
                        <div className="right-panel">
                            <Clocks
                                whiteTime={whiteTime}
                                blackTime={blackTime}
                                turnToPlay={turnToPlay}
                                gameOver={gameOver}
                                setWhiteTime={setWhiteTime}
                                setBlackTime={setBlackTime}
                                onTimeOut={handleTimeOut}
                                playerColor={playerColor}
                            />
                            <button
                                onClick={() => {
                                    setGameOver(true);
                                    setGameResult("win");
                                }}
                                className="base-button resign-button"
                                disabled={!gameStarted}
                            >
                                Здатися
                            </button>
                        </div>
                    </div>
                </div>
            </section>
            <Footer />

            {/* Add this at the bottom of your component */}
            {gameResult && (
                <GameResultPopup result={gameResult} onClose={resetGame} />
            )}
        </div>
    );
};

export default GamePage;
