import React, { useRef } from "react";
import Piece from "./Piece";
import squareSize from "./variables";

const Pieces = ({ pieces, onPieceDrop }) => {
    const ref = useRef();

    const piecesWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "absolute",
    };

    // Calculate the board square based on mouse position
    // TODO: change is the user plays black
    const calculateSquare = (clientX, clientY) => {
        const { left, top } = ref.current.getBoundingClientRect();
        const file = Math.floor((clientX - left) / squareSize);
        const rank = Math.floor((clientY - top) / squareSize);
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
                />
            ))}
        </div>
    );
};

export default Pieces;
