import React, { useRef } from "react";
import Piece from "./Piece";
import squareSize from "../variables";

const Pieces = ({ pieces, onPieceDrop, isFlipped }) => {
    const ref = useRef();

    const piecesWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "absolute",
    };

    const calculateSquare = (clientX, clientY) => {
        const { left, top } = ref.current.getBoundingClientRect();
        let file = Math.floor((clientX - left) / squareSize);
        let rank = Math.floor((clientY - top) / squareSize);

        if (isFlipped) {
            file = 7 - file;
            rank = 7 - rank;
        }

        return rank * 8 + file;
    };

    const handleDrop = (event) => {
        event.preventDefault();
        const data = event.dataTransfer.getData("text/plain");
        const [fromSquare, index] = data.split(",");

        const toSquare = calculateSquare(
            event.clientX,
            event.clientY
        ).toString();

        onPieceDrop(fromSquare, toSquare, index);
    };

    // Allow drag over to enable drop
    const onDragOver = (event) => {
        event.preventDefault();
    };

    return (
        <div
            ref={ref}
            style={piecesWrapperStyles}
            onDrop={handleDrop}
            onDragOver={onDragOver}
        >
            {pieces.map((piece, index) => (
                <Piece
                    key={piece.id}
                    index={index}
                    piece={piece.piece}
                    square={piece.offset}
                    isFlipped={isFlipped}
                />
            ))}
        </div>
    );
};

export default Pieces;
