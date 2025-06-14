import { useEffect, useState, useRef } from "react";
import { useLocation } from "react-router-dom";
import { useDispatch } from "react-redux";
import ChessBoard from "../components/board/ChessBoard";
import Header from "../components/Header";
import Footer from "../components/Footer";
import useChessWebSocket from "../hooks/useChessWebSocket";
import { PositionOutcome } from "../utils/outcomeEnum";
import { setCanNotDrag, setCanDrag } from "../store/actions";
import Clocks from "../components/clocks/Clock";
import GameResultPopup from "../components/GameResultPopup";

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

    const [timeControlType, setTimeControlType] = useState("fisher"); // 'fisher' or 'standard'
    const [baseTime, setBaseTime] = useState(180);
    const [timeIncrement, setTimeIncrement] = useState(2);
    const [selectedPreset, setSelectedPreset] = useState("rapid");

    const timeControls = {
        bullet: { base: 60, increment: 1, name: "(1+1) Куля" },
        blitz: { base: 180, increment: 2, name: "(3+2) Блискавичний контроль" },
        rapid: { base: 600, increment: 5, name: "(10+5) Швидкий контроль" },
        classical: {
            base: 1800,
            increment: 10,
            name: "Класичний контроль (30+10)",
        },
        custom: { base: 180, increment: 2, name: "Custom" },
    };

    const location = useLocation();
    const dispatch = useDispatch();

    const handleMessage = (data) => {
        console.log(data);
        if (data.outcome === PositionOutcome.DRAW) {
            setGameOver(true);
            setGameResult("draw");
        } else if (data.outcome === PositionOutcome.LOSS) {
            setGameOver(true);
            setGameResult("loss");
        } else if (data.outcome === PositionOutcome.WIN) {
            setGameOver(true);
            setGameResult("win");
        } else {
            dispatch(setCanDrag());
            setTurnToPlay("User");
            setBlackTime((prev) => prev + increment);
            setPossibleMoves(data.possibleMoves || []);
            setFen(data.fen);
        }
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
        }
        if (Number.isInteger(inc)) {
            setIncrement(inc);
        }
        if (model) {
            setSelectedModel(model);
        }
        return () => {
            // Clean up
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
            };
            sendMessage(engineMove);
        }
    }, [isBoardFlipped]);

    // --- Called when the user drops a piece on the board ---
    const handlePieceDrop = (fromSquare, toSquare, promotion = null) => {
        dispatch(setCanNotDrag());
        setTurnToPlay("Computer");

        setWhiteTime((prev) => prev + increment);

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
        setTimeIncrement(preset.increment);
        setSelectedPreset(presetKey);
        if (!gameStarted) {
            setWhiteTime(preset.base);
            setBlackTime(preset.base);
        }
    };

    const startGame = () => {
        setGameStarted(true);
        setWhiteTime(baseTime);
        setBlackTime(baseTime);
        setIncrement(timeIncrement);
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

                                    <div className="time-control-options">
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
                                                />
                                            </div>
                                            {timeControlType === "fisher" && (
                                                <div className="time-input">
                                                    <label>
                                                        Додавання (в сек.):
                                                    </label>
                                                    <input
                                                        type="number"
                                                        value={timeIncrement}
                                                        onChange={(e) =>
                                                            setTimeIncrement(
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
                                        (думає ...)
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
                                onTimeUpWhite={() => {
                                    setGameOver(true);
                                    alert("Час білих вичерпано!");
                                }}
                                onTimeUpBlack={() => {
                                    setGameOver(true);
                                    alert("Час чорних вичерпано!");
                                }}
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
