import React, { useEffect, useState, useRef } from "react";
import { useLocation } from "react-router-dom";
import { useDispatch } from "react-redux";
import ChessBoard from "../components/board/ChessBoard";
import Header from "../components/Header";
import Footer from "../components/Footer";
import OpeningsPanel from "../components/OpeningsPanel";
import useChessWebSocket from "../hooks/useChessWebSocket";
import { PositionOutcome } from "../utils/outcomeEnum";
import { setCanNotDrag, setCanDrag } from "../store/actions";
import Clocks from "../components/clocks/Clock";

import "./GamePage.css";

const GamePage = () => {
    const [fen, setFen] = useState("4k2r/4ppbp/8/8/1N6/8/5PPP/4K2R w Kk - 0 9");
    const [possibleMoves, setPossibleMoves] = useState([]);
    const [turnToPlay, setTurnToPlay] = useState("");
    const [gameOver, setGameOver] = useState(false);
    const [selectedVariations, setSelectedVariations] = useState([]);

    // Timers
    const [whiteTime, setWhiteTime] = useState(160);
    const [blackTime, setBlackTime] = useState(170);
    const [increment, setIncrement] = useState(2);
    const timerRef = useRef(null);

    const location = useLocation();
    const dispatch = useDispatch();

    const handleMessage = (data) => {
        console.log(data);
        if (data.outcome === PositionOutcome.DRAW) {
            setGameOver(true);
            console.log("DRAW");
            alert("Draw");
        } else if (data.outcome === PositionOutcome.LOSS) {
            setGameOver(true);
            console.log("Computer lost, you won!");
            alert("Computer lost, you won!");
        } else if (data.outcome === PositionOutcome.WIN) {
            setGameOver(true);
            console.log("Computer won, you lost");
            alert("Computer won, you lost");
        } else {
            dispatch(setCanDrag());
            setTurnToPlay("User");

            setBlackTime((prev) => prev + increment);

            setPossibleMoves(data.possibleMoves || []);

            // Update the FEN from the server
            setFen(data.fen);
        }
    };

    // --- Use the custom WebSocket hook ---
    const { sendMessage } = useChessWebSocket(
        `${process.env.REACT_APP_SOCKET_URL}move`,
        handleMessage
    );

    // --- On mount, parse query params, fetch initial data, etc. ---
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const baseTime = parseInt(params.get("baseTime"), 10);
        const inc = parseInt(params.get("increment"), 10);

        if (Number.isInteger(baseTime)) {
            setWhiteTime(baseTime);
            setBlackTime(baseTime);
        }
        if (Number.isInteger(inc)) {
            setIncrement(inc);
        }

        fetch(`${process.env.REACT_APP_BACKEND_URL}Game?fen=${fen}`)
            .then((res) => {
                return res.json();
            })
            .then((res) => {
                setPossibleMoves(res["possibleMoves"]);
            });

        return () => {
            // Clean up
            clearInterval(timerRef.current);
        };
    }, [location, fen]);

    // --- Called when the user drops a piece on the board ---
    const handlePieceDrop = (fromSquare, toSquare) => {
        dispatch(setCanNotDrag());
        setTurnToPlay("Computer");

        setWhiteTime((prev) => prev + increment);

        // Send the move to the server
        const moveData = {
            fen,
            from: fromSquare,
            to: toSquare,
            selectedVariations: selectedVariations, // TODO:
            time: blackTime,
        };
        sendMessage(moveData);
        console.log("Sending move to server:", moveData);
    };

    return (
        <div>
            <Header />
            <section className="game-section">
                <div className="game-page-content">
                    <div className="game-details-wrapper">
                        {/* <OpeningsPanel
                            setSelectedVariations={setSelectedVariations}
                            selectedVariations={selectedVariations}
                        /> */}
                    </div>

                    <div className="chess-board-wrapper">
                        <div className="player-headings-wrapper">
                            <p className="player-heading">Шаховий рушій</p>
                            <p
                                className="player-subheading"
                                style={{
                                    opacity:
                                        turnToPlay === "Computer"
                                            ? "100%"
                                            : "0%",
                                }}
                            >
                                Рушій думає ...
                            </p>
                        </div>

                        {/* The board & pieces */}
                        <div className="chess-board">
                            <ChessBoard
                                fen={fen}
                                possibleMoves={possibleMoves}
                                onPieceDrop={handlePieceDrop}
                            />
                        </div>

                        <div className="player-headings-wrapper">
                            <p className="player-heading">Гість</p>
                            <p
                                className="player-subheading"
                                style={{
                                    opacity:
                                        turnToPlay === "User" ? "100%" : "0%",
                                }}
                            >
                                Ваш хід
                            </p>
                        </div>
                    </div>

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
                    />
                </div>
            </section>
            <Footer />
        </div>
    );
};

export default GamePage;
