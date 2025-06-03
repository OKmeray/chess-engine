import React, { useState, useEffect } from "react";
import squareSize from "../variables";
import Pieces from "./Pieces";

const ChessBoard = ({ fen, possibleMoves, onPieceDrop }) => {
    const [pieces, setPieces] = useState([]);

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
    }, [fen]);

    const handleDrop = (fromIndex, toIndex) => {
        console.log(possibleMoves[fromIndex]);
        if (
            !possibleMoves[fromIndex] ||
            (possibleMoves[fromIndex] &&
                !possibleMoves[fromIndex].includes(parseInt(toIndex)))
        ) {
            alert("Illegal move!");
            return;
        }

        setPieces((prevPieces) =>
            prevPieces
                .filter((p) => p.offset.toString() !== toIndex.toString()) // Remove captured piece
                .map((p) =>
                    p.offset.toString() === fromIndex.toString()
                        ? { ...p, offset: toIndex }
                        : p
                )
        );

        onPieceDrop(fromIndex, toIndex);
    };

    return (
        <div className="board-image-wrapper" style={boardWrapperStyles}>
            {Array.from({ length: 8 }, (_, row) =>
                Array.from({ length: 8 }, (_, col) => (
                    <div
                        style={{
                            background:
                                (row + col) % 2 === 0 ? "#f0d9b5" : "#b58863",
                            width: `${squareSize}px`,
                            height: `${squareSize}px`,
                        }}
                    ></div>
                ))
            )}
            <Pieces pieces={pieces} onPieceDrop={handleDrop} />
        </div>
    );
};
export default ChessBoard;
