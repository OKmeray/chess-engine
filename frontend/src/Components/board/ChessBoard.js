import React, { useState, useEffect } from "react";
import squareSize from "../variables";
import Pieces from "./Pieces";
import queenWhite from "../../assets/images/queen_white.svg";
import rookWhite from "../../assets/images/rook_white.svg";
import bishopWhite from "../../assets/images/bishop_white.svg";
import knightWhite from "../../assets/images/knight_white.svg";
import queenBlack from "../../assets/images/queen_black.svg";
import rookBlack from "../../assets/images/rook_black.svg";
import bishopBlack from "../../assets/images/bishop_black.svg";
import knightBlack from "../../assets/images/knight_black.svg";

const ChessBoard = ({ fen, possibleMoves, onPieceDrop, isFlipped }) => {
    const [pieces, setPieces] = useState([]);
    const [promotionPending, setPromotionPending] = useState(null);
    const [promotionMovingPiece, setPromotionMovingPiece] = useState(null);
    const [justPromoted, setJustPromoted] = useState(false);
    const [errorMsg, setErrorMsg] = useState("");

    const showError = (msg) => {
        setErrorMsg(msg);
        setTimeout(() => setErrorMsg(""), 700);
    };

    const boardWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "relative",
        align: "center",
        display: "grid",
        gridTemplateColumns: "repeat(8, 1fr)",
        gridTemplateRows: "repeat(8, 1fr)",
    };

    function isDigit(char) {
        if (char >= "0" && char <= "9") {
            return true;
        }
        return false;
    }

    useEffect(() => {
        if (justPromoted && fen === promotionPending?.previousFen) {
            // setJustPromoted(false); // Reset flag and skip update
            return;
        }

        const piecePlacement = fen.split(" ")[0];
        let parsedPieces = [];
        let squareIndex = 0;
        let indexCounter = 0;

        for (let i = 0; i < piecePlacement.length; i++) {
            if (isDigit(piecePlacement[i])) {
                squareIndex += parseInt(piecePlacement[i]);
            } else if (piecePlacement[i] === "/") {
                continue;
            } else {
                parsedPieces.push({
                    id: indexCounter,
                    piece: piecePlacement[i],
                    offset: squareIndex,
                });
                squareIndex += 1;
                indexCounter += 1;
            }
        }

        setPieces(parsedPieces);
        setJustPromoted(false);
    }, [fen]);
    const handleDrop = (fromIndex, toIndex) => {
        if (
            !possibleMoves[fromIndex] ||
            !possibleMoves[fromIndex].includes(parseInt(toIndex))
        ) {
            if (fromIndex === toIndex) {
                return;
            }
            showError("Неправильний хід");
            return;
        }

        // === CASTLING DETECTION ===
        const castlingMoves = {
            whiteKingside: {
                kingFrom: 60,
                kingTo: 62,
                rookFrom: 63,
                rookTo: 61,
            },
            whiteQueenside: {
                kingFrom: 60,
                kingTo: 58,
                rookFrom: 56,
                rookTo: 59,
            },
            blackKingside: { kingFrom: 4, kingTo: 6, rookFrom: 7, rookTo: 5 },
            blackQueenside: { kingFrom: 4, kingTo: 2, rookFrom: 0, rookTo: 3 },
        };

        let updatedPieces = [...pieces];
        const movingPiece = pieces.find(
            (p) => p.offset.toString() === fromIndex.toString()
        );

        // Detect and apply castling rook movement
        for (const move of Object.values(castlingMoves)) {
            if (fromIndex == move.kingFrom && toIndex == move.kingTo) {
                updatedPieces = updatedPieces.map((p) => {
                    if (p.offset === move.rookFrom) {
                        return { ...p, offset: move.rookTo };
                    }
                    return p;
                });
            }
        }

        // === EN PASSANT DETECTION ===
        if (
            movingPiece &&
            (movingPiece.piece === "P" || movingPiece.piece === "p")
        ) {
            const diff = Math.abs(fromIndex - toIndex);
            if (diff === 7 || diff === 9) {
                const capturedPawnSquare =
                    parseInt(toIndex) + (movingPiece.piece === "P" ? 8 : -8);
                const isEmptySquare = !pieces.find(
                    (p) => p.offset.toString() === toIndex.toString()
                );
                if (isEmptySquare) {
                    updatedPieces = updatedPieces.filter(
                        (p) => p.offset !== capturedPawnSquare
                    );
                }
            }
        }

        // Remove any piece on the target square
        updatedPieces = updatedPieces.filter(
            (p) => p.offset.toString() !== toIndex.toString()
        );

        // Check if this is a promotion move
        const targetRank = Math.floor(toIndex / 8);
        const isPromotion =
            (movingPiece.piece === "P" && targetRank === 0) ||
            (movingPiece.piece === "p" && targetRank === 7);

        if (isPromotion) {
            setPromotionPending({
                from: fromIndex,
                to: toIndex,
                previousFen: fen,
            });
            setPromotionMovingPiece(movingPiece);
            setPieces(
                updatedPieces.map((p) =>
                    p.offset.toString() === fromIndex.toString()
                        ? { ...p, offset: toIndex }
                        : p
                )
            );
        } else {
            updatedPieces = updatedPieces.map((p) =>
                p.offset.toString() === fromIndex.toString()
                    ? { ...p, offset: toIndex }
                    : p
            );
            setPieces(updatedPieces);
            onPieceDrop(fromIndex, toIndex);
        }
    };

    const promotionImages = {
        P: {
            q: queenWhite,
            r: rookWhite,
            b: bishopWhite,
            n: knightWhite,
        },
        p: {
            q: queenBlack,
            r: rookBlack,
            b: bishopBlack,
            n: knightBlack,
        },
    };

    return (
        <div className="board-image-wrapper" style={boardWrapperStyles}>
            {errorMsg && <div className="error-toast">{errorMsg}</div>}
            {Array.from({ length: 8 }, (_, row) =>
                Array.from({ length: 8 }, (_, col) => {
                    const displayRow = isFlipped ? 7 - row : row;
                    const displayCol = isFlipped ? 7 - col : col;
                    return (
                        <div
                            key={`${displayRow}-${displayCol}`}
                            style={{
                                background:
                                    (row + col) % 2 === 0
                                        ? "#f0d9b5"
                                        : "#b58863",
                                width: `${squareSize}px`,
                                height: `${squareSize}px`,
                                gridRow: isFlipped ? 8 - row : row + 1,
                                gridColumn: isFlipped ? 8 - col : col + 1,
                            }}
                        ></div>
                    );
                })
            )}
            <Pieces
                pieces={pieces}
                onPieceDrop={handleDrop}
                isFlipped={isFlipped}
            />
            {promotionPending && (
                <div className="promotion-dialog">
                    {["q", "r", "b", "n"].map((pieceType) => (
                        <button
                            key={pieceType}
                            onClick={() => {
                                const promotedPiece =
                                    promotionMovingPiece.piece === "P"
                                        ? pieceType.toUpperCase()
                                        : pieceType;

                                // Create new piece with promotion
                                const promotedPieceObj = {
                                    ...promotionMovingPiece,
                                    offset: promotionPending.to,
                                    piece: promotedPiece,
                                };

                                // Remove the original piece and any piece on the target square
                                const updated = pieces
                                    .filter(
                                        (p) =>
                                            p.offset !==
                                                promotionPending.from &&
                                            p.offset !== promotionPending.to
                                    )
                                    .concat([promotedPieceObj]);

                                setPieces(updated);
                                setPromotionPending(null);
                                setPromotionMovingPiece(null);
                                setJustPromoted(true); // Skip the next FEN update

                                onPieceDrop(
                                    promotionPending.from,
                                    promotionPending.to,
                                    promotedPiece
                                );
                            }}
                            style={{
                                border: "none",
                                background: "none",
                                padding: "6px",
                            }}
                        >
                            <img
                                src={
                                    promotionImages[promotionMovingPiece.piece][
                                        pieceType
                                    ]
                                }
                                alt={pieceType}
                                width={squareSize}
                                height={squareSize}
                            />
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
};
export default ChessBoard;
